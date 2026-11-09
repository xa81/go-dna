package dna

import (
	"fmt"
	"strings"
)

// GetForwarding reads the URL forwarding of a domain. "Not configured" (the
// gateway answers "Redirect information not found") is an OK envelope with
// Enabled false:
//
//	{"result": "OK", "data": {"Enabled": true,
//	 "RedirectAddress": "https://www.example.org", "ForwardType": "Frame"}}
func (c *Client) GetForwarding(domainName string) map[string]any {
	response, err := c.request("GET", "domains/forwards", params{{"domainName", domainName}})
	if err != nil {
		raw := coalesce(dig(c.lastResponse, "raw_response"), err.message)
		if strings.Contains(strings.ToLower(phpString(raw)), "not found") || strings.Contains(strings.ToLower(err.message), "not found") {
			return map[string]any{
				"result": ResultOK,
				"data":   map[string]any{"Enabled": false, "RedirectAddress": "", "ForwardType": ""},
			}
		}
		return map[string]any{
			"result": ResultError,
			"error":  setError(formatErrorCode(err.code), err.message, coalesce(dig(c.lastParsedResponse, "Details"), raw)),
		}
	}
	return map[string]any{
		"result": ResultOK,
		"data": map[string]any{
			"Enabled":         !phpEmpty(dig(response, "redirectAlias")),
			"RedirectAddress": phpString(coalesce(dig(response, "redirectAlias"), "")),
			"ForwardType":     phpString(coalesce(dig(response, "redirectType"), "")),
		},
	}
}

// SetForwarding points the domain at a URL. forwardType is "Standard" (plain
// redirect) or "Frame" (masked), case-insensitive; empty means "Standard".
func (c *Client) SetForwarding(domainName, redirectAddress, forwardType string) map[string]any {
	if forwardType == "" {
		forwardType = "Standard"
	}
	typ := ucfirst(asciiLower(phpTrim(forwardType)))
	supported := false
	for _, t := range ForwardTypes {
		if t == typ {
			supported = true
		}
	}
	if !supported {
		return errorResult("SET_FORWARD", fmt.Sprintf("Forward type \"%s\" is not supported. Supported types: %s.", forwardType, strings.Join(ForwardTypes, ", ")))
	}
	if phpTrim(redirectAddress) == "" {
		return errorResult("SET_FORWARD", "Redirect address is required.")
	}
	payload := map[string]any{"domainName": domainName, "redirectAddress": redirectAddress, "forwardType": typ}
	if _, err := c.request("POST", "domains/forwards", payload); err != nil {
		return c.errorEnvelope(err)
	}
	return map[string]any{
		"result": ResultOK,
		"data":   map[string]any{"RedirectAddress": redirectAddress, "ForwardType": typ},
	}
}

// DeleteForwarding removes the URL forwarding of a domain:
// {"result": "OK", "data": {"DomainName": "example.com"}}.
func (c *Client) DeleteForwarding(domainName string) map[string]any {
	if _, err := c.request("DELETE", "domains/forwards", params{{"DomainName", domainName}}); err != nil {
		return c.errorEnvelope(err)
	}
	return map[string]any{"result": ResultOK, "data": map[string]any{"DomainName": domainName}}
}
