package dna

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type staticDoer struct {
	status int
	body   string
}

func (d staticDoer) Do(*http.Request) (*http.Response, error) {
	return &http.Response{StatusCode: d.status, Header: http.Header{}, Body: io.NopCloser(strings.NewReader(d.body))}, nil
}

func errorOf(t *testing.T, status int, body string) map[string]any {
	t.Helper()
	c := New("r", "k", false, WithHTTPClient(staticDoer{status, body}))
	c.SetErrorReportingEnabled(false)
	res := c.GetDetails("example.com")
	if res["result"] != ResultError {
		t.Fatalf("status %d: want an error, got %v", status, res)
	}
	return res["error"].(map[string]any)
}

func TestErrorMapping(t *testing.T) {
	cases := []struct {
		name          string
		status        int
		body          string
		code, message string
		details       string
	}{
		{"empty 429", 429, "", "API_429", "HTTP 429 - Too Many Requests", "HTTP 429 - Too Many Requests"},
		{"empty unknown 4xx", 418, "", "API_418", "HTTP 418 - Client Error", "HTTP 418 - Client Error"},
		{"empty unknown 5xx", 599, "", "API_599", "HTTP 599 - Server Error", "HTTP 599 - Server Error"},
		{"top-level message", 400, `{"message":"Bad thing","details":"more"}`, "API_400", "Bad thing", "more"},
		{"nested message", 404, `{"error":{"message":"Nope","details":"gone"}}`, "API_404", "Nope", "gone"},
		{"html body", 502, "<html>bad</html>", "API_502", "<html>bad</html>", "<html>bad</html>"},
		{"301 credentials", 301, "", "API_301", "Invalid API credentials", "Authentication failed. Check your API key and reseller ID."},
		{
			"ABP boilerplate uses details", 400,
			`{"error":{"code":"X:1","message":"Exception of type 'Volo.Abp.BusinessException' was thrown.","details":"Dna Host is linked"}}`,
			"API_400", "Host is linked", "Dna Host is linked",
		},
		{
			"ABP boilerplate without details", 400,
			`{"error":{"code":"X:1","message":"Exception of type 'Volo.Abp.BusinessException' was thrown.","details":"Exception of type 'A' was thrown."}}`,
			"API_400", "The operation could not be completed by the registry (error: X:1). Please try again or contact support.", "Exception of type 'A' was thrown.",
		},
		{
			"validation errors", 400,
			`{"error":{"message":"Validation failed.","details":"x","validationErrors":[{"message":"'PostalCode' is required","members":["$.administrativeContact.postalCode","$.technicalContact.postalCode"]},{"message":"Bad","members":["$.contacts[0].tldAttributes/TRABISCITYID"]},"plain text",{"message":"'PostalCode' is required","members":["$.administrativeContact.postalCode","$.technicalContact.postalCode"]}]}}`,
			"API_400", "Validation failed.", "Admin PostalCode, Tech PostalCode - 'PostalCode' is required | Contacts TldAttributes TRABISCITYID - Bad | plain text",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			e := errorOf(t, tc.status, tc.body)
			if e["Code"] != tc.code || e["Message"] != tc.message || e["Details"] != tc.details {
				t.Errorf("got Code=%q Message=%q Details=%q\nwant Code=%q Message=%q Details=%q",
					e["Code"], e["Message"], e["Details"], tc.code, tc.message, tc.details)
			}
		})
	}
}

func TestGatewayCode(t *testing.T) {
	e := errorOf(t, 404, `{"error":{"code":"Dna.DomainService:Domain:10007","message":"Not found","details":"gone"}}`)
	if e["GatewayCode"] != "Dna.DomainService:Domain:10007" {
		t.Errorf("GatewayCode = %v", e["GatewayCode"])
	}
	for _, body := range []string{"", `{"message":"x"}`, `{"error":{"code":null,"message":"x"}}`} {
		if e := errorOf(t, 400, body); e["GatewayCode"] != nil {
			t.Errorf("body %q: unexpected GatewayCode %v", body, e["GatewayCode"])
		}
	}
	if e := errorOf(t, 302, ""); e["GatewayCode"] != nil {
		t.Errorf("credential redirect must not carry GatewayCode: %v", e)
	}
	// The message is never replaced by the raw body: code-only errors without
	// details get the synthesised text.
	e = errorOf(t, 400, `{"error":{"code":"X:1","message":"X:1","details":null}}`)
	if e["Message"] != "The operation could not be completed by the registry (error: X:1). Please try again or contact support." {
		t.Errorf("message = %v", e["Message"])
	}
}

func TestOperationFailure(t *testing.T) {
	cases := []struct {
		call func(*Client) map[string]any
		code string
	}{
		{func(c *Client) map[string]any { return c.Renew("example.com", 1) }, "RENEW_DOMAIN"},
		{func(c *Client) map[string]any { return c.Transfer("example.com", "x", 1) }, "TRANSFER_DOMAIN"},
		{func(c *Client) map[string]any { return c.RegisterWithContactInfo("example.com", 1, nil, nil) }, "REGISTER_DOMAIN"},
	}
	for _, tc := range cases {
		for body, msg := range map[string]string{
			`{"success":false,"message":"Not available"}`: "Not available",
			`{"success":false,"message":"  "}`:            "The registry did not accept the request.",
		} {
			c := New("r", "k", false, WithHTTPClient(staticDoer{200, body}))
			c.SetErrorReportingEnabled(false)
			e, _ := tc.call(c)["error"].(map[string]any)
			if e["Code"] != tc.code || e["Message"] != msg || e["Details"] != msg {
				t.Errorf("%s %s: got %v", tc.code, body, e)
			}
		}
	}
	// "success": "false" (a string) is not a refusal.
	c := New("r", "k", false, WithHTTPClient(staticDoer{200, `{"success":"false","expirationDate":"2027-01-01T00:00:00"}`}))
	if res := c.Renew("example.com", 1); res["result"] != ResultOK {
		t.Errorf("non-bool success must not fail, got %v", res)
	}
}

func TestEmptySuccessIsResponseError(t *testing.T) {
	for _, body := range []string{"", "null", "[]", "{}", "not json"} {
		c := New("r", "k", false, WithHTTPClient(staticDoer{200, body}))
		res := c.GetDetails("example.com")
		e, _ := res["error"].(map[string]any)
		if res["result"] != ResultError || e["Code"] != "RESPONSE" || e["Message"] != "Invalid response received from server! Response is empty." {
			t.Errorf("body %q: got %v", body, res)
		}
	}
}

func TestSetErrorDefaults(t *testing.T) {
	e := baseSetError("CREDENTIALS", "", " ")
	if e["Message"] != "Invalid reseller ID or API key" || e["Details"] != "The provided API credentials are invalid" {
		t.Errorf("predefined code must fall back to its texts, got %v", e)
	}
	e = baseSetError("CREDENTIALS", "real", "server said")
	if e["Message"] != "real" || e["Details"] != "server said" {
		t.Errorf("caller texts must win, got %v", e)
	}
	e = baseSetError("OTHER", "", "")
	if e["Code"] != "OTHER" || e["Message"] != "" || e["Details"] != "" {
		t.Errorf("unknown code must be kept as is, got %v", e)
	}
}

func TestHumanizeValidationMember(t *testing.T) {
	cases := map[string]string{
		"$.administrativeContact.postalCode": "Admin PostalCode",
		"$.tldAttributes.TRABISCOUNTRYID":    "TldAttributes TRABISCOUNTRYID",
		"$.contacts[2].billing/eMail":        "Contacts Billing EMail",
		"   ":                                "",
		"$.":                                 "",
	}
	for in, want := range cases {
		if got := humanizeValidationMember(in); got != want {
			t.Errorf("humanizeValidationMember(%q) = %q, want %q", in, got, want)
		}
	}
}
