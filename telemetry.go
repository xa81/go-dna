package dna

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"io"
	"net"
	"net/http"
	"net/url"
	"os/user"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

// DefaultSentryDSN is where error reports and sampled performance records
// go. Empty: telemetry is a no-op until then.
// It is filled in when the DomainNameAPI-Go Sentry project exists.
const DefaultSentryDSN = ""

// PerformanceSampleRate is the share of successful calls (per 1000) sent as
// performance records; calls slower than one second are always sent.
const PerformanceSampleRate = 40

const ipCacheTTL = 512 * time.Second

var (
	ipCacheMu    sync.Mutex
	ipCacheValue string
	ipCacheAt    time.Time
)

var telemetryMaskRe = regexp.MustCompile(`(?i)auth|epp|password|passwd|secret|token|api[_-]?key`)

// telemetryMask masks credentials (auth/EPP codes, API keys, passwords,
// tokens) as "***"; everything else is kept.
func telemetryMask(data any, depth int) any {
	if depth > 8 {
		return "[…]"
	}
	switch t := data.(type) {
	case map[string]any:
		out := make(map[string]any, len(t))
		for k, v := range t {
			if telemetryMaskRe.MatchString(k) {
				if v == nil || v == "" {
					out[k] = v
				} else {
					out[k] = "***"
				}
			} else {
				out[k] = telemetryMask(v, depth+1)
			}
		}
		return out
	case []any:
		out := make([]any, len(t))
		for i, v := range t {
			out[i] = telemetryMask(v, depth+1)
		}
		return out
	}
	return data
}

func randomHex(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

type sentryDSN struct {
	apiURL    string
	publicKey string
	secretKey string
}

func parseDSN(dsn string) (sentryDSN, bool) {
	if dsn == "" {
		return sentryDSN{}, false
	}
	u, err := url.Parse(dsn)
	if err != nil || u.Host == "" || u.User == nil || u.User.Username() == "" || strings.Trim(u.Path, "/") == "" {
		return sentryDSN{}, false
	}
	secret, _ := u.User.Password()
	return sentryDSN{
		apiURL:    "https://" + u.Host + "/api/" + strings.TrimLeft(u.Path, "/") + "/store/",
		publicKey: u.User.Username(),
		secretKey: secret,
	}, true
}

func (d sentryDSN) authHeader(client string) string {
	parts := []string{"sentry_version=7", "sentry_client=" + client, "sentry_key=" + d.publicKey}
	if d.secretKey != "" {
		parts = append(parts, "sentry_secret="+d.secretKey)
	}
	return "Sentry " + strings.Join(parts, ", ")
}

func (c *Client) environment() string {
	if strings.Contains(strings.ToLower(c.serviceURL), "ote.") {
		return "ote"
	}
	return "production"
}

func (c *Client) traceID() string {
	if c.sessionTraceID == "" {
		c.sessionTraceID = randomHex(16)
	}
	return c.sessionTraceID
}

var apiCodeRe = regexp.MustCompile(`API_(-?\d+)_ERROR`)

// commonTags are the tags shared by error and performance events.
func (c *Client) commonTags(e *apiError) map[string]any {
	tags := map[string]any{}
	if c.resellerID != "" {
		tags["reseller"] = c.resellerID
	}
	if d, ok := dig(c.lastRequest, "payload", "domainName").(string); ok && d != "" && !c.telemetryMinimal {
		tags["domain"] = d
	}
	if c.lastFunction != "" {
		tags["operation"] = c.lastFunction
	}
	tags["api_type"] = "REST"
	if e != nil {
		if m := apiCodeRe.FindStringSubmatch(e.message); m != nil {
			tags["api_code"] = "API_" + m[1]
		} else if e.code != 0 {
			tags["api_code"] = "API_" + strconv.Itoa(e.code)
		}
		if e.code != 0 {
			tags["status_code"] = strconv.Itoa(e.code)
		}
	}
	return tags
}

var (
	ipInMessageRe  = regexp.MustCompile(`\b(\d{1,3}(?:\.\d{1,3}){3})\b`)
	apiErrPrefixRe = regexp.MustCompile(`(?i)^\[API_ERROR\]:\s*`)
	apiFailedRe    = regexp.MustCompile(`(?i)API_-?\d+_ERROR\s*-\s*Failed\s*-\s*`)
	resellerIDRe   = regexp.MustCompile(`(?i)\s*\(Reseller Id[^)]*\)\s*`)
	ipRe           = regexp.MustCompile(`\b\d{1,3}(?:\.\d{1,3}){3}\b`)
	quotedRe       = regexp.MustCompile(`'[^']{0,80}'|"[^"]{0,80}"`)
	numbersRe      = regexp.MustCompile(`\b\d{2,}\b`)
	spacesRe       = regexp.MustCompile(`\s+`)
)

// classify builds a stable title (volatile parts replaced by placeholders),
// the tags and the fingerprint for an error event.
func (c *Client) classify(e *apiError) (string, map[string]any, []any) {
	msg := e.message
	tags := c.commonTags(e)
	if m := ipInMessageRe.FindStringSubmatch(msg); m != nil {
		tags["ip_value"] = m[1]
	}
	apiCode, _ := tags["api_code"].(string)

	n := apiErrPrefixRe.ReplaceAllString(msg, "")
	n = apiFailedRe.ReplaceAllString(n, "")
	n = resellerIDRe.ReplaceAllString(n, "")
	n = ipRe.ReplaceAllString(n, "{ip}")
	n = quotedRe.ReplaceAllString(n, "{v}")
	n = numbersRe.ReplaceAllString(n, "{n}")
	n = strings.TrimSpace(spacesRe.ReplaceAllString(n, " "))

	if n == "" {
		if apiCode != "" {
			n = apiCode
		} else {
			n = "Exception"
		}
	} else if apiCode != "" && !strings.HasPrefix(n, apiCode) {
		n = apiCode + ": " + n
	}
	if len(n) > 160 {
		n = n[:157] + "..."
	}
	op, _ := tags["operation"].(string)
	if op == "" {
		op = "unknown"
	}
	return n, tags, []any{"dna", ApplicationName, op, n}
}

// stackFrames returns the caller's frames (file, line, function; no
// arguments), oldest first as Sentry expects.
func stackFrames(minimal bool) []any {
	pcs := make([]uintptr, 32)
	n := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	var list []any
	for {
		f, more := frames.Next()
		file := f.File
		if minimal {
			file = filepath.Base(file)
		}
		list = append(list, map[string]any{"filename": file, "lineno": f.Line, "function": f.Function})
		if !more {
			break
		}
	}
	for i, j := 0, len(list)-1; i < j; i, j = i+1, j-1 {
		list[i], list[j] = list[j], list[i]
	}
	return list
}

func vhostUser() string {
	if u, err := user.Current(); err == nil {
		return u.Username
	}
	return ""
}

func ignoredError(msg string) bool {
	for _, ev := range DefaultIgnoredErrors {
		if strings.Contains(msg, ev) {
			return true
		}
	}
	return false
}

// buildErrorEvent builds the Sentry event for an error (ipAddress is filled
// in by the sender).
func (c *Client) buildErrorEvent(e *apiError) map[string]any {
	value, tags, fingerprint := c.classify(e)
	elapsed := c.now().Sub(c.startAt).Seconds()

	baseTags := map[string]any{
		"handled":         "yes",
		"level":           "error",
		"url":             "NA",
		"transaction":     "NA",
		"trace_id":        randomHex(8),
		"runtime_name":    "Go",
		"runtime_version": runtime.Version(),
		"elapsed_time":    strconv.FormatFloat(elapsed, 'f', 4, 64),
		"application":     ApplicationName,
		"library_version": Version,
		"ip_address":      "",
		"vhost_user":      "",
	}
	if !c.telemetryMinimal {
		baseTags["vhost_user"] = vhostUser()
	}
	for k, v := range tags {
		baseTags[k] = v
	}

	event := map[string]any{
		"event_id":  randomHex(16),
		"timestamp": time.Now().UTC().Format("2006-01-02T15:04:05Z"),
		"level":     "error",
		"logger":    "go",
		"platform":  "go",
		"culprit":   "dna.Client via github.com/domainreseller/go-dna",
		"message":   map[string]any{"formatted": value},
		"exception": map[string]any{
			"values": []any{map[string]any{
				"type":       ApplicationName + " Exception",
				"value":      value,
				"stacktrace": map[string]any{"frames": stackFrames(c.telemetryMinimal)},
			}},
		},
		"tags": baseTags,
		"extra": map[string]any{
			"request_data":     telemetryMask(normalizeInput(c.lastRequest), 0),
			"response_data":    telemetryMask(normalizeInput(c.lastResponse), 0),
			"original_message": e.message,
			"go_version":       runtime.Version(),
		},
		"release":     Version,
		"environment": c.environment(),
		"fingerprint": fingerprint,
	}

	if c.telemetryMinimal {
		delete(baseTags, "url")
		delete(baseTags, "ip_address")
		delete(baseTags, "vhost_user")
		delete(baseTags, "ip_value")
		endpoint := ""
		if u, err := url.Parse(phpString(dig(c.lastRequest, "url"))); err == nil {
			endpoint = u.Path
		}
		extra := event["extra"].(map[string]any)
		extra["request_data"] = map[string]any{
			"method":   phpString(dig(c.lastRequest, "method")),
			"endpoint": endpoint,
		}
		extra["response_data"] = map[string]any{
			"http_code": phpInt(dig(c.lastResponseHeaders, "http_code")),
		}
		event["culprit"] = "dna.Client via go-dna"
	}
	return event
}

// reportError sends an error event in the background. It never blocks the
// caller for long and never panics into it.
func (c *Client) reportError(e *apiError) {
	defer func() { _ = recover() }()
	if !c.errorReportingEnabled || ignoredError(e.message) {
		return
	}
	dsn, ok := parseDSN(c.errorReportingDSN)
	if !ok {
		return
	}
	event := c.buildErrorEvent(e)
	minimal := c.telemetryMinimal
	client := c.telemetryDoer()
	go func() {
		defer func() { _ = recover() }()
		if !minimal {
			if tags, ok := event["tags"].(map[string]any); ok {
				tags["ip_address"] = serverIP(client)
			}
		}
		postEvent(client, dsn, event)
	}()
}

// reportPerformance sends a sampled transaction for a successful call.
func (c *Client) reportPerformance(duration time.Duration, start, end time.Time) {
	defer func() { _ = recover() }()
	if !c.errorReportingEnabled {
		return
	}
	dsn, ok := parseDSN(c.errorReportingDSN)
	if !ok {
		return
	}
	event := c.buildPerformanceEvent(duration, start, end)
	minimal := c.telemetryMinimal
	client := c.telemetryDoer()
	go func() {
		defer func() { _ = recover() }()
		if !minimal {
			if tags, ok := event["tags"].(map[string]any); ok {
				tags["ip_address"] = serverIP(client)
			}
		}
		postEvent(client, dsn, event)
	}()
}

func sentryTime(t time.Time) string {
	return t.UTC().Format("2006-01-02T15:04:05.000Z")
}

func (c *Client) buildPerformanceEvent(duration time.Duration, start, end time.Time) map[string]any {
	traceID := c.traceID()
	rootSpan := randomHex(8)
	childSpan := randomHex(8)
	status := "ok"

	tags := map[string]any{
		"application":     ApplicationName,
		"library_version": Version,
		"ip_address":      "",
		"vhost_user":      "",
	}
	if !c.telemetryMinimal {
		tags["vhost_user"] = vhostUser()
	}
	for k, v := range c.commonTags(nil) {
		tags[k] = v
	}
	if code := phpInt(dig(c.lastResponseHeaders, "http_code")); code != 0 {
		tags["status_code"] = strconv.FormatInt(code, 10)
	}

	method := phpString(dig(c.lastRequest, "method"))
	if method == "" {
		method = "GET"
	}
	callURL := phpString(dig(c.lastRequest, "url"))
	host := ""
	if u, err := url.Parse(callURL); err == nil {
		host = u.Host
	}
	spanData := map[string]any{"http.request.method": method}
	if host != "" {
		spanData["server.address"] = host
	}
	if callURL != "" && !c.telemetryMinimal {
		if i := strings.IndexByte(callURL, '?'); i >= 0 {
			callURL = callURL[:i]
		}
		spanData["url.full"] = callURL
	}
	if sc, ok := tags["status_code"].(string); ok {
		n, _ := strconv.Atoi(sc)
		spanData["http.response.status_code"] = n
	}
	description := c.lastFunction
	if host != "" {
		description = method + " https://" + host
	}

	return map[string]any{
		"event_id":         randomHex(16),
		"platform":         "go",
		"level":            "info",
		"type":             "transaction",
		"transaction":      "API." + c.lastFunction,
		"transaction_info": map[string]any{"source": "custom"},
		"release":          Version,
		"environment":      c.environment(),
		"timestamp":        sentryTime(end),
		"start_timestamp":  sentryTime(start),
		"contexts": map[string]any{
			"trace":   map[string]any{"trace_id": traceID, "span_id": rootSpan, "op": "http.client", "status": status},
			"runtime": map[string]any{"name": "go", "version": runtime.Version()},
		},
		"tags": tags,
		"measurements": map[string]any{
			"duration": map[string]any{"value": float64(duration) / float64(time.Millisecond), "unit": "millisecond"},
		},
		"spans": []any{map[string]any{
			"span_id":         childSpan,
			"parent_span_id":  rootSpan,
			"trace_id":        traceID,
			"op":              "http.client",
			"description":     description,
			"status":          status,
			"data":            spanData,
			"start_timestamp": sentryTime(start),
			"timestamp":       sentryTime(end),
		}},
	}
}

func (c *Client) telemetryDoer() Doer {
	if c.telemetryClient != nil {
		return c.telemetryClient
	}
	return &http.Client{
		Timeout: 2 * time.Second,
		Transport: &http.Transport{
			Proxy:       http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{Timeout: time.Second}).DialContext,
		},
	}
}

// postEvent posts one event to Sentry's store endpoint, ignoring the result.
func postEvent(client Doer, dsn sentryDSN, event map[string]any) {
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, dsn.apiURL, bytes.NewReader(encodeJSON(event)))
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("X-Sentry-Auth", dsn.authHeader(ApplicationName+"/"+Version))
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	_, _ = io.Copy(io.Discard, resp.Body)
	resp.Body.Close()
}

// serverIP looks up the public IP (api.ipify.org), cached for 512 seconds.
func serverIP(client Doer) string {
	ipCacheMu.Lock()
	if ipCacheValue != "" && time.Since(ipCacheAt) < ipCacheTTL {
		v := ipCacheValue
		ipCacheMu.Unlock()
		return v
	}
	ipCacheMu.Unlock()

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.ipify.org", nil)
	if err != nil {
		return "unknown"
	}
	req.Header.Set("User-Agent", UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		return "unknown"
	}
	defer resp.Body.Close()
	b, _ := io.ReadAll(io.LimitReader(resp.Body, 64))
	ip := strings.TrimSpace(string(b))
	if net.ParseIP(ip) == nil {
		return "unknown"
	}
	ipCacheMu.Lock()
	ipCacheValue, ipCacheAt = ip, time.Now()
	ipCacheMu.Unlock()
	return ip
}
