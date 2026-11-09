package dna

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

// captureDoer records telemetry posts and answers the IP lookup.
type captureDoer struct {
	mu     sync.Mutex
	events []map[string]any
	auth   []string
	urls   []string
	got    chan struct{}
}

func newCaptureDoer() *captureDoer { return &captureDoer{got: make(chan struct{}, 16)} }

func (d *captureDoer) Do(req *http.Request) (*http.Response, error) {
	if req.URL.Host == "api.ipify.org" {
		return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("203.0.113.7"))}, nil
	}
	b, _ := io.ReadAll(req.Body)
	var ev map[string]any
	_ = json.Unmarshal(b, &ev)
	d.mu.Lock()
	d.events = append(d.events, ev)
	d.auth = append(d.auth, req.Header.Get("X-Sentry-Auth"))
	d.urls = append(d.urls, req.URL.String())
	d.mu.Unlock()
	d.got <- struct{}{}
	return &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader("{}"))}, nil
}

func (d *captureDoer) wait(t *testing.T) map[string]any {
	t.Helper()
	select {
	case <-d.got:
	case <-time.After(2 * time.Second):
		t.Fatal("no telemetry event was posted")
	}
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.events[len(d.events)-1]
}

func (d *captureDoer) none(t *testing.T) {
	t.Helper()
	select {
	case <-d.got:
		t.Fatal("unexpected telemetry event")
	case <-time.After(150 * time.Millisecond):
	}
}

const testDSN = "https://publickey@sentry.example.test/42"

func TestTelemetryMask(t *testing.T) {
	in := map[string]any{
		"url": "https://x", "payload": map[string]any{
			"domainName": "example.com", "authCode": "Ab1!", "eppCode": "", "X-API-KEY": "k", "api_key": "k2",
			"password": "p", "nested": []any{map[string]any{"token": "t", "keep": "v"}},
		},
	}
	got := telemetryMask(in, 0).(map[string]any)
	p := got["payload"].(map[string]any)
	if p["authCode"] != "***" || p["X-API-KEY"] != "***" || p["api_key"] != "***" || p["password"] != "***" {
		t.Errorf("credentials not masked: %v", p)
	}
	if p["eppCode"] != "" || p["domainName"] != "example.com" || got["url"] != "https://x" {
		t.Errorf("non-credentials or empty values changed: %v", got)
	}
	n := p["nested"].([]any)[0].(map[string]any)
	if n["token"] != "***" || n["keep"] != "v" {
		t.Errorf("nested masking: %v", n)
	}
}

func TestParseDSN(t *testing.T) {
	d, ok := parseDSN("https://pub:sec@sentry.example.test/9")
	if !ok || d.apiURL != "https://sentry.example.test/api/9/store/" || d.publicKey != "pub" || d.secretKey != "sec" {
		t.Errorf("parseDSN = %+v, %v", d, ok)
	}
	if !strings.Contains(d.authHeader("x/1"), "sentry_secret=sec") {
		t.Error("auth header must carry the secret")
	}
	for _, bad := range []string{"", "not a url", "https://sentry.example.test/9", "https://pub@sentry.example.test/"} {
		if _, ok := parseDSN(bad); ok {
			t.Errorf("parseDSN(%q) must fail", bad)
		}
	}
	if DefaultSentryDSN != "" {
		t.Error("DefaultSentryDSN must stay empty until the Sentry project exists")
	}
}

func TestErrorEventIsSent(t *testing.T) {
	tel := newCaptureDoer()
	c := New("reseller-1", "secret-key", true,
		WithHTTPClient(staticDoer{500, `{"error":{"message":"Registry exploded at 10.1.2.3 for 'example.com' code 12345"}}`}),
		WithSentryDSN(testDSN), WithTelemetryClient(tel))
	c.Renew("example.com", 2)

	ev := tel.wait(t)
	tags := ev["tags"].(map[string]any)
	if tags["application"] != "DomainNameAPI-Go" || tags["library_version"] != Version || ev["release"] != Version {
		t.Errorf("application/version tags: %v / %v", tags, ev["release"])
	}
	if ev["environment"] != "ote" || tags["reseller"] != "reseller-1" || tags["domain"] != "example.com" {
		t.Errorf("environment/reseller/domain: %v %v", ev["environment"], tags)
	}
	if tags["api_code"] != "API_500" || tags["status_code"] != "500" || tags["operation"] != "request:POST domains/renew" {
		t.Errorf("code/operation tags: %v", tags)
	}
	if tags["ip_value"] != "10.1.2.3" || tags["ip_address"] != "203.0.113.7" {
		t.Errorf("ip tags: %v", tags)
	}
	title := ev["message"].(map[string]any)["formatted"]
	if title != "API_500: Registry exploded at {ip} for {v} code {n}" {
		t.Errorf("title = %v", title)
	}
	fp := ev["fingerprint"].([]any)
	if len(fp) != 4 || fp[0] != "dna" || fp[1] != "DomainNameAPI-Go" || fp[3] != title {
		t.Errorf("fingerprint = %v", fp)
	}
	extra := ev["extra"].(map[string]any)
	if !strings.Contains(extra["original_message"].(string), "10.1.2.3") {
		t.Errorf("original message lost: %v", extra)
	}
	raw, _ := json.Marshal(ev)
	if strings.Contains(string(raw), "secret-key") {
		t.Error("API key leaked into telemetry")
	}
	if !strings.Contains(tel.auth[0], "sentry_key=publickey") || tel.urls[0] != "https://sentry.example.test/api/42/store/" {
		t.Errorf("auth/url: %v %v", tel.auth, tel.urls)
	}
}

func TestMinimalTelemetry(t *testing.T) {
	tel := newCaptureDoer()
	c := New("reseller-1", "k", false, WithHTTPClient(staticDoer{503, ""}), WithSentryDSN(testDSN), WithTelemetryClient(tel))
	c.SetTelemetryMinimal(true)
	c.EnableTheftProtectionLock("example.com")

	ev := tel.wait(t)
	tags := ev["tags"].(map[string]any)
	for _, k := range []string{"url", "ip_address", "vhost_user", "domain"} {
		if _, ok := tags[k]; ok {
			t.Errorf("minimal event must not carry tag %s: %v", k, tags)
		}
	}
	extra := ev["extra"].(map[string]any)
	req := extra["request_data"].(map[string]any)
	if req["method"] != "POST" || req["endpoint"] != "/api/v1/domains/lock" || len(req) != 2 {
		t.Errorf("request_data = %v", req)
	}
	if resp := extra["response_data"].(map[string]any); resp["http_code"] != 503.0 || len(resp) != 1 {
		t.Errorf("response_data = %v", resp)
	}
	if ev["message"].(map[string]any)["formatted"] != "API_503: HTTP {n} - Service Unavailable" {
		t.Errorf("title = %v", ev["message"])
	}
}

func TestTelemetryOffAndIgnored(t *testing.T) {
	tel := newCaptureDoer()
	c := New("r", "k", false, WithHTTPClient(staticDoer{500, `{"message":"boom"}`}), WithSentryDSN(testDSN), WithTelemetryClient(tel))
	c.SetErrorReportingEnabled(false)
	c.GetDetails("example.com")
	tel.none(t)

	c = New("r", "k", false, WithHTTPClient(staticDoer{404, `{"message":"Domain not found"}`}), WithSentryDSN(testDSN), WithTelemetryClient(tel))
	c.GetDetails("example.com")
	tel.none(t)

	c = New("r", "k", false, WithHTTPClient(staticDoer{500, `{"message":"boom"}`}), WithTelemetryClient(tel))
	c.GetDetails("example.com") // empty DSN: no-op
	tel.none(t)
}

func TestSlowSuccessSendsTransaction(t *testing.T) {
	tel := newCaptureDoer()
	c := New("r", "k", false, WithHTTPClient(staticDoer{200, `{"resellerId":"x"}`}), WithSentryDSN(testDSN), WithTelemetryClient(tel))
	base := time.Unix(1000, 0)
	calls := 0
	c.now = func() time.Time { calls++; return base.Add(time.Duration(calls) * 600 * time.Millisecond) }
	c.GetResellerDetails()

	ev := tel.wait(t)
	if ev["type"] != "transaction" || ev["transaction"] != "API.request:GET deposit/accounts/me" {
		t.Errorf("transaction event = %v", ev)
	}
	trace := ev["contexts"].(map[string]any)["trace"].(map[string]any)
	if trace["status"] != "ok" || trace["op"] != "http.client" {
		t.Errorf("trace = %v", trace)
	}
	span := ev["spans"].([]any)[0].(map[string]any)
	if span["description"] != "GET https://api.domainresellerapi.com" {
		t.Errorf("span = %v", span)
	}
	if ev["tags"].(map[string]any)["application"] != "DomainNameAPI-Go" {
		t.Errorf("tags = %v", ev["tags"])
	}
}
