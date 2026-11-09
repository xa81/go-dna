package dna

import (
	"crypto/tls"
	"net"
	"net/http"
	"strings"
	"time"
)

// Gateway base URLs.
const (
	URLProd = "https://api.domainresellerapi.com/api/v1"
	URLOTE  = "https://ote.domainresellerapi.com/api/v1"
)

// Envelope result values.
const (
	ResultOK      = "OK"
	ResultError   = "ERROR"
	ResultSuccess = "SUCCESS"
)

// Timeouts, as in php-dna (DEFAULT_TIMEOUT / DEFAULT_CONNECT_TIMEOUT).
const (
	DefaultTimeout        = 30 * time.Second
	DefaultConnectTimeout = 10 * time.Second
)

// DefaultReason is the privacy-change reason used when none is given.
const DefaultReason = "Owner request"

// ZoneMinInterval is the minimum gap between two calls to domains/zones*.
// The gateway answers HTTP 429 to a burst of four or five zone calls.
const ZoneMinInterval = time.Second

// DefaultNameServers are used by RegisterWithContactInfo when no valid
// nameserver is given.
var DefaultNameServers = []string{"ns1.domainnameapi.com", "ns2.domainnameapi.com"}

// ZoneRecordTypes are the record types the DNS zone API accepts.
var ZoneRecordTypes = []string{"A", "AAAA", "ALIAS", "CAA", "CNAME", "DS", "MX", "NAPTR", "PTR", "SRV", "TLSA", "TXT"}

// ForwardTypes are the redirect types accepted by domains/forwards:
// Standard (plain redirect) and Frame (masked redirect).
var ForwardTypes = []string{"Standard", "Frame"}

// Doer sends an HTTP request. *http.Client implements it; tests and callers
// can inject their own transport with WithHTTPClient.
type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

// Option configures a Client at construction time.
type Option func(*Client)

// WithHTTPClient replaces the HTTP transport used for gateway calls. With a
// custom Doer, SetSSLVerification and the built-in timeouts no longer apply:
// configure them on the Doer.
func WithHTTPClient(d Doer) Option {
	return func(c *Client) { c.httpClient = d }
}

// WithSentryDSN overrides the telemetry DSN (DefaultSentryDSN by default).
func WithSentryDSN(dsn string) Option {
	return func(c *Client) { c.errorReportingDSN = dsn }
}

// WithTelemetryClient replaces the HTTP transport used to post telemetry
// events (and the public-IP lookup).
func WithTelemetryClient(d Doer) Option {
	return func(c *Client) { c.telemetryClient = d }
}

// Client is the DomainNameAPI client: the Go counterpart of php-dna's
// DomainNameAPI_PHPLibrary facade. Create it with New. A Client is not safe
// for concurrent use.
type Client struct {
	resellerID string
	apiKey     string
	serviceURL string

	httpClient    Doer
	customHTTP    bool
	sslVerify     bool
	defaultClient *http.Client

	errorReportingEnabled bool
	telemetryMinimal      bool
	errorReportingDSN     string
	telemetryClient       Doer
	sessionTraceID        string

	lastRequest         map[string]any
	lastResponse        any
	lastResponseHeaders map[string]any
	lastParsedResponse  any
	lastFunction        string
	startAt             time.Time

	zoneMinInterval   time.Duration
	lastZoneRequestAt time.Time
	sleep             func(time.Duration)
	now               func() time.Time
}

// New creates a client. resellerID is the Reseller ID (UUID) and apiKey the
// API Key from the DomainNameAPI panel (API Settings); panel login details do
// not work. Both are trimmed.
//
// testMode selects the OTE test gateway. It accepts a bool, or a string as
// panels store checkbox settings: "true", "1", "on" and "yes" (any case) are
// true, everything else is false. Numbers are true when non-zero; nil is
// false.
func New(resellerID, apiKey string, testMode any, opts ...Option) *Client {
	c := &Client{
		resellerID:            phpTrim(resellerID),
		apiKey:                phpTrim(apiKey),
		serviceURL:            URLProd,
		sslVerify:             true,
		errorReportingEnabled: true,
		errorReportingDSN:     DefaultSentryDSN,
		lastRequest:           map[string]any{},
		lastResponse:          map[string]any{},
		lastResponseHeaders:   map[string]any{},
		lastParsedResponse:    map[string]any{},
		startAt:               time.Now(),
		zoneMinInterval:       ZoneMinInterval,
		sleep:                 time.Sleep,
		now:                   time.Now,
	}
	if parseTestMode(testMode) {
		c.serviceURL = URLOTE
	}
	for _, o := range opts {
		o(c)
	}
	if c.httpClient != nil {
		c.customHTTP = true
	} else {
		c.defaultClient = newHTTPClient(true)
		c.httpClient = c.defaultClient
	}
	return c
}

// parseTestMode mirrors php-dna's FILTER_VALIDATE_BOOLEAN handling of the
// test-mode flag (BUG-10284).
func parseTestMode(v any) bool {
	switch t := v.(type) {
	case nil:
		return false
	case bool:
		return t
	case string:
		switch strings.ToLower(strings.Trim(t, phpWhitespace)) {
		case "1", "true", "on", "yes":
			return true
		}
		return false
	}
	return phpBool(normalizeInput(v))
}

func newHTTPClient(verify bool) *http.Client {
	tr := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           (&net.Dialer{Timeout: DefaultConnectTimeout, KeepAlive: 30 * time.Second}).DialContext,
		TLSHandshakeTimeout:   DefaultConnectTimeout,
		ForceAttemptHTTP2:     true,
		MaxIdleConns:          10,
		IdleConnTimeout:       90 * time.Second,
		ExpectContinueTimeout: time.Second,
	}
	if !verify {
		tr.TLSClientConfig = &tls.Config{InsecureSkipVerify: true} //nolint:gosec // explicit opt-out, see SetSSLVerification
	}
	return &http.Client{
		Transport: tr,
		Timeout:   DefaultTimeout,
		// curl does not follow redirects for php-dna either: a 301/302 from
		// the gateway means the credentials were rejected.
		CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
	}
}

// SetSSLVerification turns TLS certificate verification for gateway calls on
// or off (on by default). Only switch it off on a host whose CA bundle is
// too old to validate the certificate. (php-dna: setSslVerification)
func (c *Client) SetSSLVerification(enabled bool) {
	c.sslVerify = enabled
	if !c.customHTTP {
		c.defaultClient = newHTTPClient(enabled)
		c.httpClient = c.defaultClient
	}
}

// SetErrorReportingEnabled turns telemetry on or off (on by default).
func (c *Client) SetErrorReportingEnabled(enabled bool) { c.errorReportingEnabled = enabled }

// SetTelemetryMinimal sends only the error code and message, endpoint,
// timing, versions and reseller ID: no request/response data, domain or
// server IP.
func (c *Client) SetTelemetryMinimal(minimal bool) { c.telemetryMinimal = minimal }

// GetRequestData returns the last request: url, method and payload.
func (c *Client) GetRequestData() map[string]any { return c.lastRequest }

// GetResponseData returns the last decoded response, or
// {"raw_response": body} when the body was not JSON.
func (c *Client) GetResponseData() any { return c.lastResponse }

// GetResponseHeaders returns transfer information about the last call (the
// counterpart of php-dna's curl_getinfo array): url, http_code,
// content_type, total_time and the response headers.
func (c *Client) GetResponseHeaders() map[string]any { return c.lastResponseHeaders }

// GetParsedResponseData returns the last parsed response, or the error
// array ({Code, Message, Details}) when the last call failed.
func (c *Client) GetParsedResponseData() any { return c.lastParsedResponse }

// GetLastFunction returns the last operation, e.g. "request:GET domains/info".
func (c *Client) GetLastFunction() string { return c.lastFunction }

// GetServiceURL returns the gateway base URL. (php-dna: getServiceUrl)
func (c *Client) GetServiceURL() string { return c.serviceURL }

// SetServiceURL overrides the gateway base URL. (php-dna: setServiceUrl)
func (c *Client) SetServiceURL(url string) { c.serviceURL = url }
