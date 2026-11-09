package dna

import (
	"bytes"
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// apiError is the Go counterpart of the Exception php-dna's request() throws:
// a message plus a numeric code (the HTTP status, or a curl-style errno for
// transport failures).
type apiError struct {
	message string
	code    int
}

func (e *apiError) Error() string { return e.message }

var httpReasons = map[int]string{
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	408: "Request Timeout",
	409: "Conflict",
	422: "Unprocessable Entity",
	429: "Too Many Requests",
	500: "Internal Server Error",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
}

// httpStatusReason returns the reason phrase used to synthesise a message
// when an error status comes with an empty body.
func httpStatusReason(status int) string {
	if r, ok := httpReasons[status]; ok {
		return r
	}
	if status >= 500 {
		return "Server Error"
	}
	if status >= 400 {
		return "Client Error"
	}
	return "Unexpected Response"
}

// formatErrorCode formats a numeric error code: CURL_<errno> after a
// transport failure, API_<status> otherwise.
func (c *Client) formatErrorCode(code int) string {
	if c.transportFailed {
		return "CURL_" + strconv.Itoa(code)
	}
	return "API_" + strconv.Itoa(code)
}

// setError builds the {Code, Message, Details} error array. For the codes
// php-dna predefines, the generic texts are used only when no message or
// details were given.
func baseSetError(code any, message, details any) map[string]any {
	if k, ok := code.(string); ok {
		if def, ok := defaultErrors[k]; ok {
			msg := message
			if phpTrim(phpString(message)) == "" {
				msg = def.message
			}
			det := details
			if phpTrim(phpString(details)) == "" {
				det = def.description
			}
			return map[string]any{"Code": def.code, "Message": msg, "Details": det}
		}
	}
	return map[string]any{"Code": code, "Message": message, "Details": details}
}

// setError is baseSetError plus the gateway's own error code (GatewayCode)
// next to an API_<status> code, when the last failed request carried one.
func (c *Client) setError(code any, message, details any) map[string]any {
	res := baseSetError(code, message, details)
	if s, ok := res["Code"].(string); ok && apiStatusCodeRe.MatchString(s) && c.lastGatewayCode != "" {
		res["GatewayCode"] = c.lastGatewayCode
	}
	return res
}

var apiStatusCodeRe = regexp.MustCompile(`^API_\d+$`)

// errorEnvelope is the shared catch block of every php-dna method:
// {result: ERROR, error: {Code: API_<status>, Message, Details}}.
func (c *Client) errorEnvelope(e *apiError) map[string]any {
	details := coalesce(dig(c.lastParsedResponse, "Details"), dig(c.lastResponse, "raw_response"), e.message)
	return map[string]any{
		"result": ResultError,
		"error":  c.setError(c.formatErrorCode(e.code), e.message, details),
	}
}

func (c *Client) errorResult(code, message string) map[string]any {
	return map[string]any{"result": ResultError, "error": c.setError(code, message, message)}
}

// paceZoneRequest keeps consecutive domains/zones* calls at least
// zoneMinInterval apart.
func (c *Client) paceZoneRequest(endpoint string) {
	if !strings.HasPrefix(strings.TrimLeft(endpoint, "/"), "domains/zones") {
		return
	}
	if !c.lastZoneRequestAt.IsZero() && c.zoneMinInterval > 0 {
		wait := c.zoneMinInterval - c.now().Sub(c.lastZoneRequestAt)
		if wait > 0 {
			c.sleep(wait)
		}
	}
	c.lastZoneRequestAt = c.now()
}

func (c *Client) get(endpoint string, p params) (any, *apiError) {
	return c.request(http.MethodGet, endpoint, p)
}

func (c *Client) post(endpoint string, data any) (any, *apiError) {
	return c.request(http.MethodPost, endpoint, data)
}

func (c *Client) put(endpoint string, data any) (any, *apiError) {
	return c.request(http.MethodPut, endpoint, data)
}

func (c *Client) del(endpoint string, p params) (any, *apiError) {
	return c.request(http.MethodDelete, endpoint, p)
}

var boilerplateRe = regexp.MustCompile(`(?i)Exception of type '[^']*' was thrown`)
var dnaPrefixRe = regexp.MustCompile(`^Dna\s+`)

// request performs one gateway call, following DNARest::request(). GET and
// DELETE send data (params) as a query string; POST and PUT send it as a
// JSON body. A non-nil *apiError means the call failed; the error array is
// then in lastParsedResponse.
func (c *Client) request(method, endpoint string, data any) (any, *apiError) {
	c.paceZoneRequest(endpoint)
	c.transportFailed = false
	c.lastGatewayCode = ""
	c.startAt = c.now()

	opName := endpoint
	if i := strings.IndexByte(opName, '?'); i >= 0 {
		opName = opName[:i]
	}
	c.lastFunction = "request:" + method + " " + opName

	url := c.serviceURL + "/" + strings.TrimLeft(endpoint, "/")

	payload := data
	if p, ok := data.(params); ok {
		payload = p.toMap()
	}
	if payload == nil {
		payload = []any{}
	}
	c.lastRequest = map[string]any{"url": url, "method": method, "payload": payload}

	var body io.Reader
	if method == http.MethodGet || method == http.MethodDelete {
		if p, ok := data.(params); ok && len(p) > 0 {
			if q := httpBuildQuery(p); q != "" {
				url += "?" + q
			}
		}
	} else {
		if data == nil {
			data = []any{}
		}
		if p, ok := data.(params); ok {
			data = p.toMap()
		}
		body = bytes.NewReader(encodeJSON(data))
	}

	status, respBody, headers, err := c.send(method, url, body)

	c.lastResponseHeaders = map[string]any{
		"url":          url,
		"http_code":    int64(status),
		"content_type": headers.Get("Content-Type"),
		"total_time":   c.now().Sub(c.startAt).Seconds(),
		"headers":      flattenHeaders(headers),
	}
	if decoded := decodeJSON(respBody); isArray(decoded) {
		c.lastResponse = decoded
	} else if err != nil {
		c.lastResponse = map[string]any{"raw_response": false}
	} else {
		c.lastResponse = map[string]any{"raw_response": respBody}
	}

	if err != nil {
		// Transport failure (timeout, DNS, TLS, reset): an error, never an
		// empty success.
		errno := transportErrno(err)
		c.transportFailed = true
		msg := "Curl error during request: " + err.Error()
		if errno == 28 {
			msg = "Could not reach the DomainNameAPI service: the connection timed out. (" + err.Error() + ")"
		}
		c.lastParsedResponse = c.setError("CURL_"+strconv.Itoa(errno), msg, msg)
		e := &apiError{message: msg, code: errno}
		c.reportError(e)
		return nil, e
	}

	if status >= 200 && status <= 299 {
		parsed := decodeJSON(respBody)
		c.lastParsedResponse = parsed
		duration := c.now().Sub(c.startAt)
		if duration > time.Second || rand.Intn(1000)+1 <= PerformanceSampleRate {
			c.reportPerformance(duration, c.startAt, c.now())
		}
		return parsed, nil
	}

	parsed := decodeJSON(respBody)
	if status == 302 || status == 301 {
		parsed = map[string]any{
			"message": "Invalid API credentials",
			"code":    "CREDENTIALS",
			"details": "Authentication failed. Check your API key and reseller ID.",
		}
	}

	errorMessage := coalesce(dig(parsed, "message"), dig(parsed, "error", "message"), respBody)
	errorCode := coalesce(dig(parsed, "code"), dig(parsed, "error", "code"), "HTTP_"+strconv.Itoa(status))
	apiDetails := coalesce(dig(parsed, "details"), dig(parsed, "error", "details"))
	errorDetails := coalesce(apiDetails, respBody)
	// The gateway's own code (e.g. Dna.DomainService:Domain:10007) is kept as
	// GatewayCode next to the API_<status> code.
	if s, ok := errorCode.(string); ok && s != "" && !strings.HasPrefix(s, "HTTP_") && s != "CREDENTIALS" {
		c.lastGatewayCode = s
	}

	if phpTrim(phpString(errorMessage)) == "" {
		errorMessage = "HTTP " + strconv.Itoa(status) + " - " + httpStatusReason(status)
		if phpTrim(phpString(errorDetails)) == "" {
			errorDetails = errorMessage
		}
	}

	// ABP boilerplate or the bare error code as the message tells the
	// reseller nothing: prefer the details, else a stable message (BUG-10124).
	msgStr := phpString(errorMessage)
	codeStr := phpString(errorCode)
	if boilerplateRe.MatchString(msgStr) || msgStr == codeStr {
		// Only the gateway's details field may replace the message; the raw
		// response body is never user-facing text.
		detailsClean := phpTrim(dnaPrefixRe.ReplaceAllString(phpString(apiDetails), ""))
		if detailsClean != "" && detailsClean != msgStr && detailsClean != codeStr && !boilerplateRe.MatchString(detailsClean) {
			errorMessage = detailsClean
		} else {
			errorMessage = "The operation could not be completed by the registry (error: " + codeStr + "). Please try again or contact support."
		}
	}

	validationErrors := coalesce(dig(parsed, "error", "validationErrors"), dig(parsed, "validationErrors"))
	if isArray(validationErrors) && !phpEmpty(validationErrors) {
		if flat := flattenValidationErrors(validationErrors); flat != "" {
			errorDetails = flat
		}
	}

	c.lastParsedResponse = c.setError(errorCode, errorMessage, errorDetails)
	e := &apiError{message: phpString(errorMessage), code: status}
	c.reportError(e)
	return nil, e
}

// send performs the HTTP exchange. A transport failure (including a failure
// while reading the body) is returned as err.
func (c *Client) send(method, url string, body io.Reader) (int, string, http.Header, error) {
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return 0, "", http.Header{}, err
	}
	req.Header["Content-Type"] = []string{"application/json"}
	req.Header["Accept"] = []string{"application/json"}
	req.Header["X-API-KEY"] = []string{c.apiKey}
	req.Header["__reseller"] = []string{c.resellerID}
	req.Header["User-Agent"] = []string{UserAgent}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, "", http.Header{}, err
	}
	defer resp.Body.Close()
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, "", resp.Header, err
	}
	return resp.StatusCode, string(raw), resp.Header, nil
}

func flattenHeaders(h http.Header) map[string]any {
	out := make(map[string]any, len(h))
	for k, v := range h {
		out[k] = strings.Join(v, ", ")
	}
	return out
}

// transportErrno maps a Go transport error to the curl error number php-dna
// would report: 28 timeout, 6 DNS, 7 connection refused or other connect
// failure, 35 TLS handshake, 60 certificate, 56 connection reset.
func transportErrno(err error) int {
	var netErr net.Error
	switch {
	case errors.Is(err, context.DeadlineExceeded), errors.Is(err, os.ErrDeadlineExceeded):
		return 28
	case errors.As(err, &netErr) && netErr.Timeout():
		return 28
	}
	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		return 6
	}
	var certErr *tls.CertificateVerificationError
	var unknownAuth x509.UnknownAuthorityError
	var hostErr x509.HostnameError
	var invalidErr x509.CertificateInvalidError
	if errors.As(err, &certErr) || errors.As(err, &unknownAuth) || errors.As(err, &hostErr) || errors.As(err, &invalidErr) {
		return 60
	}
	var recErr tls.RecordHeaderError
	if errors.As(err, &recErr) || strings.Contains(err.Error(), "tls:") {
		return 35
	}
	if errors.Is(err, syscall.ECONNRESET) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.EOF) || errors.Is(err, syscall.EPIPE) {
		return 56
	}
	return 7
}
