package dna

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestClient(t *testing.T, h http.HandlerFunc) (*Client, *httptest.Server) {
	t.Helper()
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	c := New(" reseller-id ", "\tapi-key\n", false)
	c.SetErrorReportingEnabled(false)
	c.SetServiceURL(srv.URL + "/api/v1")
	return c, srv
}

func TestRequestHeadersAndQuery(t *testing.T) {
	var got *http.Request
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		got = r.Clone(context.Background())
		_, _ = io.WriteString(w, `{"totalCount":0,"items":[]}`)
	})
	res := c.GetList(map[string]any{"DomainName": "a b~c", "Flag": true, "Off": false, "Skip": nil})
	if res["result"] != ResultOK {
		t.Fatalf("result = %v", res)
	}
	if got.Header.Get("X-Api-Key") != "api-key" {
		t.Errorf("X-API-KEY = %q (credentials must be trimmed)", got.Header.Get("X-Api-Key"))
	}
	if got.Header.Get("__reseller") != "reseller-id" {
		t.Errorf("__reseller = %q", got.Header.Get("__reseller"))
	}
	if got.Header.Get("User-Agent") != "DomainNameAPI-Go/1.0.0" {
		t.Errorf("User-Agent = %q", got.Header.Get("User-Agent"))
	}
	if got.Header.Get("Accept") != "application/json" || got.Header.Get("Content-Type") != "application/json" {
		t.Errorf("Accept/Content-Type = %q/%q", got.Header.Get("Accept"), got.Header.Get("Content-Type"))
	}
	want := "MaxResultCount=200&SkipCount=0&DomainName=a+b%7Ec&Flag=1&Off=0"
	if got.URL.RawQuery != want {
		t.Errorf("query = %q, want %q", got.URL.RawQuery, want)
	}
	if c.GetLastFunction() != "request:GET domains" {
		t.Errorf("last function = %q", c.GetLastFunction())
	}
	if rd := c.GetRequestData(); rd["method"] != "GET" || !strings.HasSuffix(rd["url"].(string), "/api/v1/domains") {
		t.Errorf("request data = %v", rd)
	}
	if h := c.GetResponseHeaders(); h["http_code"] != int64(200) {
		t.Errorf("response headers = %v", h)
	}
}

func TestPostSendsJSONBody(t *testing.T) {
	var body string
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		body = string(b)
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/domains/zones/apply" || r.URL.RawQuery != "domainName=example.com" {
			t.Errorf("got %s %s?%s", r.Method, r.URL.Path, r.URL.RawQuery)
		}
	})
	if res := c.ApplyResourceRecords("example.com"); res["result"] != ResultOK {
		t.Fatalf("result = %v", res)
	}
	if body != "[]" {
		t.Errorf("body = %q, want []", body)
	}
	if c.GetLastFunction() != "request:POST domains/zones/apply" {
		t.Errorf("last function = %q (query must not be in the operation)", c.GetLastFunction())
	}
}

func TestRedirectIsCredentialError(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/login" {
			_, _ = io.WriteString(w, `{"resellerId":"followed"}`)
			return
		}
		http.Redirect(w, r, "/login", http.StatusFound)
	})
	res := c.GetResellerDetails()
	errMap := res["error"].(map[string]any)
	if res["result"] != ResultError || errMap["Code"] != "API_302" || errMap["Message"] != "Invalid API credentials" {
		t.Fatalf("redirect must not be followed, got %v", res)
	}
	if p := c.GetParsedResponseData().(map[string]any); p["Code"] != "CREDENTIALS" {
		t.Errorf("parsed response = %v", p)
	}
}

func TestTestModeAndServiceURL(t *testing.T) {
	cases := []struct {
		mode any
		ote  bool
	}{
		{false, false}, {true, true}, {"true", true}, {" On ", true}, {"YES", true}, {"1", true},
		{"false", false}, {"off", false}, {"0", false}, {"", false}, {"no", false}, {nil, false}, {1, true}, {0, false},
	}
	for _, tc := range cases {
		c := New("r", "k", tc.mode)
		want := URLProd
		if tc.ote {
			want = URLOTE
		}
		if c.GetServiceURL() != want {
			t.Errorf("testMode %#v: service URL %s, want %s", tc.mode, c.GetServiceURL(), want)
		}
	}
	c := New("r", "k", false)
	c.SetServiceURL("http://localhost:1/api/v1")
	if c.GetServiceURL() != "http://localhost:1/api/v1" {
		t.Errorf("SetServiceURL did not apply")
	}
}

type errDoer struct{ err error }

func (d errDoer) Do(*http.Request) (*http.Response, error) { return nil, d.err }

type timeoutErr struct{}

func (timeoutErr) Error() string   { return "i/o timeout" }
func (timeoutErr) Timeout() bool   { return true }
func (timeoutErr) Temporary() bool { return true }

func TestTransportErrors(t *testing.T) {
	cases := []struct {
		err   error
		errno string
		msg   string
	}{
		{&net.OpError{Op: "dial", Err: timeoutErr{}}, "28", "Could not reach the DomainNameAPI service: the connection timed out. ("},
		{context.DeadlineExceeded, "28", "Could not reach the DomainNameAPI service: the connection timed out. ("},
		{&net.DNSError{Err: "no such host", Name: "api.example"}, "6", "Curl error during request: "},
		{errors.New("connection refused"), "7", "Curl error during request: "},
		{io.ErrUnexpectedEOF, "56", "Curl error during request: "},
	}
	for _, tc := range cases {
		c := New("r", "k", false, WithHTTPClient(errDoer{tc.err}))
		c.SetErrorReportingEnabled(false)
		res := c.EnableTheftProtectionLock("example.com")
		e := res["error"].(map[string]any)
		if res["result"] != ResultError {
			t.Errorf("%v: transport failure must be an error, got %v", tc.err, res)
			continue
		}
		if e["Code"] != "API_"+tc.errno {
			t.Errorf("%v: code %v, want API_%s", tc.err, e["Code"], tc.errno)
		}
		if !strings.HasPrefix(e["Message"].(string), tc.msg) || e["Details"] != e["Message"] {
			t.Errorf("%v: message %q / details %q", tc.err, e["Message"], e["Details"])
		}
		if p := c.GetParsedResponseData().(map[string]any); p["Code"] != "CURL_"+tc.errno {
			t.Errorf("%v: parsed code %v, want CURL_%s", tc.err, p["Code"], tc.errno)
		}
		if raw := c.GetResponseData().(map[string]any)["raw_response"]; raw != false {
			t.Errorf("%v: raw_response = %#v, want false", tc.err, raw)
		}
	}
}

func TestRealTimeout(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(300 * time.Millisecond)
	})
	c.httpClient = &http.Client{Timeout: 50 * time.Millisecond}
	res := c.GetDetails("example.com")
	e := res["error"].(map[string]any)
	if e["Code"] != "API_28" || !strings.HasPrefix(e["Message"].(string), "Could not reach the DomainNameAPI service") {
		t.Fatalf("got %v", res)
	}
}

func TestSSLVerificationToggle(t *testing.T) {
	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = io.WriteString(w, `{"resellerId":"x","usdBalance":1,"tryBalance":2}`)
	}))
	defer srv.Close()
	c := New("r", "k", false)
	c.SetErrorReportingEnabled(false)
	c.SetServiceURL(srv.URL + "/api/v1")
	res := c.GetResellerDetails()
	if res["result"] != ResultError || res["error"].(map[string]any)["Code"] != "API_60" {
		t.Fatalf("self-signed certificate must fail with verification on, got %v", res)
	}
	c.SetSSLVerification(false)
	if res := c.GetResellerDetails(); res["result"] != ResultOK {
		t.Fatalf("verification off must accept the certificate, got %v", res)
	}
}

func TestZonePacing(t *testing.T) {
	c, _ := newTestClient(t, func(w http.ResponseWriter, r *http.Request) { _, _ = io.WriteString(w, "[]") })
	now := time.Unix(1000, 0)
	var slept []time.Duration
	c.now = func() time.Time { return now }
	c.sleep = func(d time.Duration) { slept = append(slept, d); now = now.Add(d) }

	c.GetResourceRecords("example.com")
	now = now.Add(300 * time.Millisecond)
	c.GetResourceRecords("example.com")
	c.GetDetails("example.com") // not a zone call: never paced
	now = now.Add(2 * time.Second)
	c.GetResourceRecords("example.com")

	if len(slept) != 1 || slept[0] != 700*time.Millisecond {
		t.Fatalf("slept %v, want [700ms]", slept)
	}
}
