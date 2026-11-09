package dna

import (
	"math"
	"net"
	"strings"
	"time"
)

var apiDateLayouts = []string{
	time.RFC3339Nano,
	"2006-01-02T15:04:05.999999999",
	"2006-01-02T15:04:05.999999999Z0700",
	"2006-01-02T15:04Z07:00",
	"2006-01-02T15:04",
	"2006-01-02 15:04:05.999999999Z07:00",
	"2006-01-02 15:04:05.999999999",
	"2006-01-02 15:04",
	"2006-01-02",
}

// parseAPIDate parses the gateway's ISO-8601 dates. Dates without an offset
// are taken as UTC; dates with one are converted to UTC.
func parseAPIDate(v any) (time.Time, bool) {
	if v == nil || v == "" {
		return time.Time{}, false
	}
	s := phpTrim(phpString(v))
	for _, layout := range apiDateLayouts {
		if t, err := time.Parse(layout, s); err == nil {
			return t.UTC(), true
		}
	}
	return time.Time{}, false
}

// formatAPIDate formats an API date as "2006-01-02T15:04:05", or "" when it
// is missing or unparseable (never a bogus 1970 date).
func formatAPIDate(v any) string {
	t, ok := parseAPIDate(v)
	if !ok {
		return ""
	}
	return t.Format("2006-01-02T15:04:05")
}

// remainingDaysFromDate counts the days until an API date, floored at 0.
func (c *Client) remainingDaysFromDate(v any) int64 {
	t, ok := parseAPIDate(v)
	if !ok {
		return 0
	}
	days := int64(math.Ceil(t.Sub(c.now()).Seconds() / 86400))
	if days < 0 {
		return 0
	}
	return days
}

var contactIDTypes = map[string]string{"Billing": "Billing", "Tech": "Technical", "Admin": "Administrative", "Registrant": "Registrant"}

func parseContactIDs(contacts any) map[string]any {
	result := map[string]any{
		"Billing":        map[string]any{"ID": int64(0)},
		"Technical":      map[string]any{"ID": int64(0)},
		"Administrative": map[string]any{"ID": int64(0)},
		"Registrant":     map[string]any{"ID": int64(0)},
	}
	for _, ct := range phpValues(contacts) {
		apiType, _ := coalesce(dig(ct, "ContactType"), dig(ct, "contactType"), "").(string)
		if key, ok := contactIDTypes[apiType]; ok {
			result[key] = map[string]any{"ID": coalesce(dig(ct, "handle"), dig(ct, "id"), int64(0))}
		}
	}
	return result
}

// parseDomainInfo converts a gateway domain to php-dna's domain array. An
// empty response is the RESPONSE error.
func (c *Client) parseDomainInfo(data any) map[string]any {
	if phpEmpty(data) {
		return map[string]any{"result": ResultError, "error": setError("RESPONSE", "", "")}
	}
	var remaining int64
	if rd := dig(data, "remainingDay"); rd != nil {
		remaining = phpInt(rd)
	} else {
		remaining = c.remainingDaysFromDate(dig(data, "expirationDate"))
	}

	nameServers := []any{}
	if ns := dig(data, "nameservers"); ns != nil {
		for _, n := range phpValues(ns) {
			nameServers = append(nameServers, phpString(n))
		}
	}

	var additional any = []any{}
	if a := dig(data, "additionalAttributes"); a != nil {
		if isArray(a) {
			additional = a
		} else {
			additional = []any{a}
		}
	}

	childNS := []any{}
	if hosts := dig(data, "hosts"); hosts != nil {
		for _, h := range phpValues(hosts) {
			// Legacy contract: the last IP as a string.
			var ip any = ""
			if ips := phpValues(dig(h, "ipAddresses")); len(ips) > 0 {
				ip = dig(ips[len(ips)-1], "ipAddress")
			}
			childNS = append(childNS, map[string]any{"ns": dig(h, "name"), "ip": ip})
		}
	}

	return map[string]any{
		"data": map[string]any{
			"ID":                      phpInt(coalesce(dig(data, "id"), int64(0))),
			"Status":                  normalizeDomainStatus(coalesce(dig(data, "status"), "")),
			"DomainName":              phpString(coalesce(dig(data, "domainName"), dig(data, "name"), "")),
			"AuthCode":                phpString(coalesce(dig(data, "authCode"), dig(data, "eppCode"), "")),
			"LockStatus":              boolString(dig(data, "lockStatus")),
			"PrivacyProtectionStatus": boolString(dig(data, "privacyProtectionStatus")),
			"IsChildNameServer":       boolString(dig(data, "hosts")),
			"Contacts":                parseContactIDs(coalesce(dig(data, "contacts"), []any{})),
			"Dates": map[string]any{
				"Start":         formatAPIDate(dig(data, "startDate")),
				"Expiration":    formatAPIDate(dig(data, "expirationDate")),
				"RemainingDays": remaining,
			},
			"NameServers":      nameServers,
			"Additional":       additional,
			"ChildNameServers": childNS,
		},
		"result": ResultOK,
	}
}

// GetDetails returns the details of a domain:
//
//	{"result": "OK", "data": {"ID": 4321, "Status": "Active",
//	 "DomainName": "example.com", "AuthCode": "...", "LockStatus": "true",
//	 "PrivacyProtectionStatus": "false", "IsChildNameServer": "false",
//	 "Contacts": {"Billing": {"ID": ...}, ...},
//	 "Dates": {"Start": "...", "Expiration": "...", "RemainingDays": 126},
//	 "NameServers": [...], "Additional": {...}, "ChildNameServers": [...]}}
func (c *Client) GetDetails(domainName string) map[string]any {
	response, err := c.request("GET", "domains/info", params{{"DomainName", domainName}})
	if err != nil {
		return c.errorEnvelope(err)
	}
	return c.parseDomainInfo(response)
}

// SyncFromRegistry refreshes a domain from the registry; same result as
// GetDetails.
func (c *Client) SyncFromRegistry(domainName string) map[string]any {
	return c.GetDetails(domainName)
}

// ModifyNameServer replaces the nameservers of a domain. nameServers is a
// []string (or any list), or a map whose values are taken in key order.
// Entries are trimmed, lower-cased and de-duplicated before sending; the
// result echoes the gateway's list, or the input as given.
func (c *Client) ModifyNameServer(domainName string, nameServers any) map[string]any {
	payload := map[string]any{"domainName": domainName, "nameServers": normalizeHostNameList(nameServers)}
	response, err := c.request("PUT", "domains/dns/name-server", payload)
	if err != nil {
		return c.errorEnvelope(err)
	}
	return map[string]any{
		"result": ResultOK,
		"data":   map[string]any{"NameServers": coalesce(dig(response, "nameServers"), nameServers)},
	}
}

// EnableTheftProtectionLock locks the domain (registrar lock):
// {"result": "OK", "data": {"LockStatus": true}}.
func (c *Client) EnableTheftProtectionLock(domainName string) map[string]any {
	if _, err := c.request("POST", "domains/lock", map[string]any{"domainName": domainName}); err != nil {
		return c.errorEnvelope(err)
	}
	return map[string]any{"result": ResultOK, "data": map[string]any{"LockStatus": true}}
}

// DisableTheftProtectionLock unlocks the domain:
// {"result": "OK", "data": {"LockStatus": false}}.
func (c *Client) DisableTheftProtectionLock(domainName string) map[string]any {
	if _, err := c.request("POST", "domains/unlock", map[string]any{"domainName": domainName}); err != nil {
		return c.errorEnvelope(err)
	}
	return map[string]any{"result": ResultOK, "data": map[string]any{"LockStatus": false}}
}

// ipVersion is "v4" for a dotted IPv4 address, otherwise "v6".
func ipVersion(ip string) string {
	if strings.Contains(ip, ":") {
		return "v6"
	}
	if p := net.ParseIP(ip); p != nil && p.To4() != nil {
		return "v4"
	}
	return "v6"
}

// AddChildNameServer creates a child nameserver (glue record):
// {"result": "OK", "data": {"NameServer": "ns1.example.com", "IPAdresses": ["192.0.2.10"]}}.
func (c *Client) AddChildNameServer(domainName, nameServer, ipAddress string) map[string]any {
	payload := map[string]any{
		"domainName": domainName,
		"hostName":   normalizeHostName(nameServer),
		"ipAddresses": []any{
			map[string]any{"ipAddress": ipAddress, "ipVersion": ipVersion(ipAddress)},
		},
	}
	if _, err := c.request("POST", "domains/dns/host", payload); err != nil {
		return c.errorEnvelope(err)
	}
	return map[string]any{
		"result": ResultOK,
		"data":   map[string]any{"NameServer": nameServer, "IPAdresses": []any{ipAddress}},
	}
}

// DeleteChildNameServer deletes a child nameserver:
// {"result": "OK", "data": {"NameServer": "ns1.example.com"}}.
func (c *Client) DeleteChildNameServer(domainName, nameServer string) map[string]any {
	p := params{{"domainName", domainName}, {"hostName", normalizeHostName(nameServer)}}
	if _, err := c.request("DELETE", "domains/dns/host", p); err != nil {
		return c.errorEnvelope(err)
	}
	return map[string]any{"result": ResultOK, "data": map[string]any{"NameServer": nameServer}}
}

// ModifyChildNameServer changes the IP address of a child nameserver:
// {"result": "OK", "data": {"NameServer": "ns1.example.com", "IPAdresses": ["192.0.2.20"]}}.
func (c *Client) ModifyChildNameServer(domainName, nameServer, ipAddress string) map[string]any {
	host := normalizeHostName(nameServer)
	payload := map[string]any{
		"domainName":  domainName,
		"hostName":    host,
		"newHostName": host,
		"ipAddresses": []any{
			map[string]any{"ipAddress": ipAddress, "ipVersion": ipVersion(ipAddress)},
		},
	}
	response, err := c.request("PUT", "domains/dns/host", payload)
	if err != nil {
		return c.errorEnvelope(err)
	}
	return map[string]any{
		"result": ResultOK,
		"data": map[string]any{
			"NameServer": coalesce(dig(response, "hostName"), host),
			"IPAdresses": coalesce(dig(response, "ipAddresses"), []any{ipAddress}),
		},
	}
}

var contactKeyFromAPI = map[string]string{"Admin": "Administrative", "Billing": "Billing", "Tech": "Technical", "Registrant": "Registrant"}

// GetContacts returns the four contacts of a domain, keyed Administrative,
// Billing, Registrant, Technical:
//
//	{"result": "OK", "data": {"contacts": {"Administrative": {"ID": "...",
//	 "Status": "Active", "FirstName": "...", "LastName": "...", "Company": "...",
//	 "EMail": "...", "Type": "Administrative", "Address": {...},
//	 "Phone": {"Phone": {...}, "Fax": {...}}, "Additional": [...]}, ...}}}
func (c *Client) GetContacts(domainName string) map[string]any {
	info, err := c.request("GET", "domains/info", params{{"DomainName", domainName}})
	if err != nil {
		return c.errorEnvelope(err)
	}
	contacts := map[string]any{}
	if list := dig(info, "contacts"); isArray(list) {
		for _, item := range phpValues(list) {
			apiType := phpString(coalesce(dig(item, "contactType"), ""))
			key, ok := contactKeyFromAPI[apiType]
			if !ok {
				key = ucfirst(asciiLower(apiType))
			}
			merged := map[string]any{}
			if m, ok := item.(map[string]any); ok {
				for k, v := range m {
					merged[k] = v
				}
			}
			merged["type"] = key
			contacts[key] = parseContactInfo(merged)
		}
	}
	ordered := map[string]any{}
	for _, k := range []string{"Administrative", "Billing", "Registrant", "Technical"} {
		if v, ok := contacts[k]; ok && v != nil {
			ordered[k] = v
		}
	}
	return map[string]any{"data": map[string]any{"contacts": ordered}, "result": ResultOK}
}

// SaveContacts updates the contacts of a domain. Half-empty contacts are
// completed by ValidateContact's rules before sending.
func (c *Client) SaveContacts(domainName string, contacts Contacts) map[string]any {
	payload := map[string]any{"domainName": domainName, "contacts": payloadContacts(contacts)}
	response, err := c.request("PUT", "domains/contacts/update", payload)
	if err != nil {
		return c.errorEnvelope(err)
	}
	parsed := map[string]any{}
	if list := dig(response, "contacts"); isArray(list) {
		for _, ct := range phpValues(list) {
			parsed[ucfirst(asciiLower(phpString(dig(ct, "type"))))] = parseContactInfo(ct)
		}
	}
	return map[string]any{"result": ResultOK, "data": map[string]any{"contacts": parsed}}
}

// Transfer starts an incoming transfer with the EPP/auth code.
//
// Like php-dna's facade (4.0.4), the contacts argument is accepted but not
// forwarded: the request always carries an empty contact list. It is part of
// the signature so the call shape matches php-dna's DNARest::transfer().
func (c *Client) Transfer(domainName, eppCode string, period int, contacts Contacts) map[string]any {
	_ = contacts
	payload := map[string]any{
		"domainName": domainName,
		"authCode":   eppCode,
		"period":     int64(period),
		"contacts":   []any{},
	}
	response, err := c.request("POST", "domains/transfer", payload)
	if err != nil {
		return c.errorEnvelope(err)
	}
	return c.parseDomainInfo(response)
}

func (c *Client) transferAction(domainName, action, status string) map[string]any {
	response, err := c.request("POST", "domains/transfers/"+action, map[string]any{"domainName": domainName})
	if err != nil {
		return c.errorEnvelope(err)
	}
	return map[string]any{
		"result": ResultOK,
		"data":   map[string]any{"DomainName": domainName, "Status": coalesce(dig(response, "status"), status)},
	}
}

// CancelTransfer cancels a pending incoming transfer.
func (c *Client) CancelTransfer(domainName string) map[string]any {
	return c.transferAction(domainName, "cancel", "Cancelled")
}

// ApproveTransfer approves a pending outgoing transfer.
func (c *Client) ApproveTransfer(domainName string) map[string]any {
	return c.transferAction(domainName, "approve", "Approved")
}

// RejectTransfer rejects a pending outgoing transfer.
func (c *Client) RejectTransfer(domainName string) map[string]any {
	return c.transferAction(domainName, "reject", "Rejected")
}

// CheckTransfer checks whether a domain can be transferred with the given
// auth code. result is OK only when the registry reports it transferable;
// otherwise the envelope also carries error TRANSFER_NOT_AVAILABLE.
func (c *Client) CheckTransfer(domainName, authCode string) map[string]any {
	response, err := c.request("POST", "domains/transfers/check", map[string]any{"domainName": domainName, "authCode": authCode})
	if err != nil {
		return c.errorEnvelope(err)
	}
	available := !phpEmpty(dig(response, "transferAvailabilityStatus"))
	var additional any = []any{}
	if a := dig(response, "additionalAttributes"); a != nil {
		if isArray(a) {
			additional = a
		} else {
			additional = []any{a}
		}
	}
	result := ResultError
	if available {
		result = ResultOK
	}
	envelope := map[string]any{
		"result": result,
		"data": map[string]any{
			"TransferAvailabilityStatus": available,
			"AuthCodeIsRequired":         !phpEmpty(dig(response, "authCodeIsRequired")),
			"AuthCodeIsValid":            !phpEmpty(dig(response, "authCodeIsValid")),
			"UserTransferRequired":       !phpEmpty(dig(response, "userTransferRequired")),
			"TransferLock":               !phpEmpty(dig(response, "transferLock")),
			"Message":                    phpString(coalesce(dig(response, "message"), "")),
			"MessageKey":                 phpString(coalesce(dig(response, "messageKey"), "")),
			"Additional":                 additional,
		},
	}
	if !available {
		envelope["error"] = setError(
			"TRANSFER_NOT_AVAILABLE",
			coalesce(dig(response, "message"), "Domain is not transferable"),
			coalesce(dig(response, "messageKey"), dig(response, "message"), ""),
		)
	}
	return envelope
}

// Renew renews a domain for period years:
// {"result": "OK", "data": {"ExpirationDate": "2027-11-09T12:00:00"}}.
func (c *Client) Renew(domainName string, period int) map[string]any {
	response, err := c.request("POST", "domains/renew", map[string]any{"domainName": domainName, "period": int64(period)})
	if err != nil {
		return c.errorEnvelope(err)
	}
	if exp := dig(response, "expirationDate"); phpBool(exp) {
		return map[string]any{"result": ResultOK, "data": map[string]any{"ExpirationDate": formatAPIDate(exp)}}
	}
	c.reportError(&apiError{message: "[DOMAIN_RENEW] " + defaultErrors["DOMAIN_RENEW"].description})
	return map[string]any{"result": ResultError, "error": setError("DOMAIN_RENEW", "", "")}
}

// RegisterOptions holds RegisterWithContactInfo's optional arguments. A nil
// *RegisterOptions means php-dna's defaults: default nameservers, EPP lock
// on, privacy off, no additional attributes. Start from
// DefaultRegisterOptions() to change only some of them.
type RegisterOptions struct {
	// NameServers is a []string (or a map, values taken in key order). When
	// it holds no valid entry, DefaultNameServers are used.
	NameServers any
	// EppLock turns the registrar (theft protection) lock on. Default true.
	EppLock bool
	// PrivacyLock turns WHOIS privacy protection on. Default false.
	PrivacyLock bool
	// AdditionalAttributes are TLD attributes, e.g. the TRABIS attributes of
	// .tr domains (legacy names are normalized).
	AdditionalAttributes map[string]any
}

// DefaultRegisterOptions returns php-dna's defaults.
func DefaultRegisterOptions() *RegisterOptions {
	return &RegisterOptions{EppLock: true}
}

// RegisterWithContactInfo registers a domain for period years with the
// given contacts. opts may be nil (php-dna's defaults). Returns the domain
// array (as GetDetails) on success.
func (c *Client) RegisterWithContactInfo(domainName string, period int, contacts Contacts, opts *RegisterOptions) map[string]any {
	if opts == nil {
		opts = DefaultRegisterOptions()
	}
	nameServers := normalizeHostNameList(opts.NameServers)
	if len(nameServers) == 0 {
		for _, ns := range DefaultNameServers {
			nameServers = append(nameServers, ns)
		}
	}

	attrs, _ := normalizeInput(opts.AdditionalAttributes).(map[string]any)
	if attrs == nil {
		attrs = map[string]any{}
	}
	if strings.HasSuffix(asciiLower(domainName), ".tr") {
		attrs = normalizeTrAttributes(attrs)
	}
	attrs = stringifyAttributes(attrs)

	payload := map[string]any{
		"domainName":     domainName,
		"period":         int64(period),
		"nameServers":    nameServers,
		"isLocked":       opts.EppLock,
		"privacyEnabled": opts.PrivacyLock,
		"contacts":       payloadContacts(contacts),
		"tldAttributes":  attrs,
	}
	response, err := c.request("POST", "domains/register-with-contacts", payload)
	if err != nil {
		return c.errorEnvelope(err)
	}
	return c.parseDomainInfo(response)
}

// ModifyPrivacyProtectionStatus turns WHOIS privacy protection on or off.
// reason is informational (empty means "Owner request").
func (c *Client) ModifyPrivacyProtectionStatus(domainName string, status bool, reason string) map[string]any {
	if reason == "" {
		reason = DefaultReason
	}
	_ = reason
	if _, err := c.request("POST", "domains/privacy", map[string]any{"domainName": domainName, "privacyStatus": status}); err != nil {
		return c.errorEnvelope(err)
	}
	return map[string]any{"data": map[string]any{"PrivacyProtectionStatus": status}, "result": ResultOK}
}

// IsTrTLD reports whether the domain ends with ".tr" (case-insensitive).
func (c *Client) IsTrTLD(domain string) bool {
	return strings.HasSuffix(asciiLower(domain), ".tr")
}

// ValidateContact completes and normalizes a contact without any network
// call: missing fields get placeholders (FirstName/LastName "Isimyok",
// AddressLine1 "Addres yok", City/State "ISTANBUL", Country "TR", ZipCode
// "34000", Phone "5555555555", PhoneCountryCode "90"), names lose digits and
// punctuation, State falls back to City, common country names become ISO
// codes and phone/fax numbers are split into number and country code.
func (c *Client) ValidateContact(contact map[string]any) map[string]any {
	return validateContact(contact)
}
