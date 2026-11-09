package dna

import (
	"regexp"
	"strings"
)

type defaultError struct {
	code        string
	message     string
	description string
}

// defaultErrors are php-dna's predefined error codes (DEFAULT_ERRORS).
var defaultErrors = map[string]defaultError{
	"DOMAIN_DETAILS":          {"DOMAIN_DETAILS", "Invalid domain details! Details format is not valid", "The provided domain details are not in the expected format"},
	"CREDENTIALS":             {"CREDENTIALS", "Invalid reseller ID or API key", "The provided API credentials are invalid"},
	"DOMAIN_LIST":             {"DOMAIN_LIST", "The domain list is invalid or contains multiple entries!", "The domain list response format is incorrect or unexpected"},
	"TLD_LIST":                {"TLD_LIST", "TLD info is not a valid array or more than one TLD info has returned!", "The TLD list response is not in the expected format"},
	"RESPONSE":                {"RESPONSE", "Invalid response received from server! Response is empty.", "The API response is empty or null"},
	"RESPONSE_FORMAT":         {"RESPONSE_FORMAT", "Invalid response received from server! Response format is not valid.", "The API response format is not in the expected structure"},
	"RESPONSE_COUNT":          {"RESPONSE_COUNT", "Invalid parameters passed to function! Response data contains more than one result!", "The API response contains multiple results when only one was expected"},
	"RESPONSE_CODE":           {"RESPONSE_CODE", "Invalid parameters passed to function! Operation result or Error code not received from server", "The API response is missing required operation result or error code fields"},
	"RESPONSE_REST":           {"RESPONSE_REST", "Invalid parameters passed to function! REST return is not a valid array!", "The REST response is not in a valid array format"},
	"CONTACT_INFO":            {"CONTACT_INFO", "Invalid response received from server! Contact info is not a valid array or more than one contact info has returned!", "The contact information response is not in the expected format"},
	"CONTACT_SAVE":            {"CONTACT_SAVE", "Invalid response received from server! Contact info could not be saved!", "The contact information could not be saved on the server"},
	"DOMAIN_TRANSFER_REQUEST": {"DOMAIN_TRANSFER_REQUEST", "Invalid response received from server! Domain transfer request could not be completed!", "The domain transfer request failed to complete"},
	"DOMAIN_RENEW":            {"DOMAIN_RENEW", "Invalid response received from server! Domain renew request could not be completed!", "The domain renewal request failed to complete"},
	"DOMAIN_REGISTER":         {"DOMAIN_REGISTER", "Invalid response received from server! Domain register request could not be completed!", "The domain registration request failed to complete"},
	"DOMAIN_SYNC":             {"DOMAIN_SYNC", "Invalid response received from server! Domain sync request could not be completed!", "The domain synchronization request failed to complete"},
}

// DefaultIgnoredErrors are error messages (substrings) that are expected
// noise and never reported to telemetry.
var DefaultIgnoredErrors = []string{
	"Domain not found",
	"ERR_DOMAIN_NOT_FOUND",
	"Reseller not found",
	"Domain is not in updateable status",
	"balance is not sufficient",
	"Price definition not found",
	"TLD is not supported",
	"Invalid API credentials",
	"could not be found",
	"Insufficent reseller balance",
}

// flattenValidationErrors turns the gateway's validationErrors list into one
// readable string, e.g. "Admin PostalCode, Tech PostalCode - 'PostalCode' is
// required", entries joined by " | " and de-duplicated.
func flattenValidationErrors(validationErrors any) string {
	var parts []string
	for _, ve := range phpValues(validationErrors) {
		if s, ok := ve.(string); ok {
			if phpTrim(s) != "" {
				parts = append(parts, phpTrim(s))
			}
			continue
		}
		if !isArray(ve) {
			continue
		}
		message := ""
		if m := dig(ve, "message"); m != nil {
			message = phpTrim(phpString(m))
		}
		var members []string
		if ms := dig(ve, "members"); isArray(ms) {
			for _, member := range phpValues(ms) {
				if label := humanizeValidationMember(phpString(member)); label != "" {
					members = append(members, label)
				}
			}
		}
		switch {
		case len(members) > 0 && message != "":
			parts = append(parts, strings.Join(members, ", ")+" - "+message)
		case len(members) > 0:
			parts = append(parts, strings.Join(members, ", "))
		case message != "":
			parts = append(parts, message)
		}
	}
	seen := map[string]bool{}
	var unique []string
	for _, p := range parts {
		if !seen[p] {
			seen[p] = true
			unique = append(unique, p)
		}
	}
	return strings.Join(unique, " | ")
}

var arrayIndexRe = regexp.MustCompile(`\[\d+\]`)
var memberSplitRe = regexp.MustCompile(`[./]+`)

var validationRoles = map[string]string{
	"administrativecontact": "Admin", "administrative": "Admin",
	"technicalcontact": "Tech", "technical": "Tech",
	"registrantcontact": "Registrant", "registrant": "Registrant",
	"billingcontact": "Billing", "billing": "Billing",
}

// humanizeValidationMember turns "$.administrativeContact.postalCode" into
// "Admin PostalCode".
func humanizeValidationMember(member string) string {
	member = phpTrim(member)
	if member == "" {
		return ""
	}
	member = strings.TrimLeft(member, "$.")
	member = arrayIndexRe.ReplaceAllString(member, "")
	var out []string
	for _, seg := range memberSplitRe.Split(member, -1) {
		if seg == "" {
			continue
		}
		if r, ok := validationRoles[asciiLower(seg)]; ok {
			out = append(out, r)
		} else {
			out = append(out, ucfirst(seg))
		}
	}
	return phpTrim(strings.Join(out, " "))
}
