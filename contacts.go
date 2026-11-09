package dna

import (
	"regexp"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

// RoleContact is one contact of a domain: the role (Administrative, Billing,
// Technical, Registrant) and the contact fields, with php-dna's keys
// (FirstName, LastName, Company, EMail, AddressLine1, AddressLine2, City,
// State, Country, ZipCode, Phone, PhoneCountryCode, Fax, FaxCountryCode,
// ...). Both the flat shape and the nested one (Address.Line1,
// Phone.Phone.Number, ...) are accepted.
type RoleContact struct {
	Role    string
	Contact map[string]any
}

// Contacts is an ordered list of role contacts. The order is the order in
// which the contacts are sent to the gateway (php-dna sends them in the
// array's order).
type Contacts []RoleContact

// ContactRoles is the canonical role order used by ContactsFromMap and
// AllRoles.
var ContactRoles = []string{"Administrative", "Billing", "Technical", "Registrant"}

// ContactsFromMap builds Contacts from a role => contact map. Go maps have no
// order, so the roles are taken in the canonical order (Administrative,
// Billing, Technical, Registrant), then any other keys alphabetically.
func ContactsFromMap(m map[string]map[string]any) Contacts {
	var out Contacts
	done := map[string]bool{}
	for _, r := range ContactRoles {
		if c, ok := m[r]; ok {
			out = append(out, RoleContact{Role: r, Contact: c})
			done[r] = true
		}
	}
	var rest []string
	for k := range m {
		if !done[k] {
			rest = append(rest, k)
		}
	}
	sort.Strings(rest)
	for _, k := range rest {
		out = append(out, RoleContact{Role: k, Contact: m[k]})
	}
	return out
}

// AllRoles uses the same contact for all four roles.
func AllRoles(contact map[string]any) Contacts {
	out := make(Contacts, 0, len(ContactRoles))
	for _, r := range ContactRoles {
		out = append(out, RoleContact{Role: r, Contact: contact})
	}
	return out
}

var contactTypeToAPI = map[string]string{
	"Administrative": "Admin",
	"Technical":      "Tech",
	"Billing":        "Billing",
	"Registrant":     "Registrant",
}

// parseContact converts one input contact to the gateway's ContactLiteDto.
func parseContact(contact any, typ string) map[string]any {
	contact = normalizeInput(contact)
	apiType := typ
	if t, ok := contactTypeToAPI[typ]; ok {
		apiType = t
	}

	address := phpString(coalesce(dig(contact, "AddressLine1"), dig(contact, "Address", "Line1"), ""))
	if line2 := coalesce(dig(contact, "AddressLine2"), dig(contact, "Address", "Line2"), ""); !phpEmpty(line2) {
		address += " " + phpString(line2)
	}

	phone := coalesce(dig(contact, "Phone"), dig(contact, "Phone", "Phone", "Number"), "")
	phoneCc := coalesce(dig(contact, "PhoneCountryCode"), dig(contact, "Phone", "Phone", "CountryCode"), "")
	fax := coalesce(dig(contact, "Fax"), dig(contact, "Phone", "Fax", "Number"), "")
	faxCc := coalesce(dig(contact, "FaxCountryCode"), dig(contact, "Phone", "Fax", "CountryCode"), "")

	if isArray(phone) {
		phoneCc = coalesce(dig(phone, "Phone", "CountryCode"), "")
		phone = coalesce(dig(phone, "Phone", "Number"), "")
	}
	if isArray(fax) {
		faxCc = coalesce(dig(fax, "Fax", "CountryCode"), "")
		fax = coalesce(dig(fax, "Fax", "Number"), "")
	}

	validated := validateContact(map[string]any{
		"FirstName":        coalesce(dig(contact, "FirstName"), ""),
		"LastName":         coalesce(dig(contact, "LastName"), ""),
		"AddressLine1":     address,
		"City":             coalesce(dig(contact, "City"), dig(contact, "Address", "City"), ""),
		"State":            coalesce(dig(contact, "State"), dig(contact, "Address", "State"), ""),
		"Country":          coalesce(dig(contact, "Country"), dig(contact, "Address", "Country"), ""),
		"ZipCode":          coalesce(dig(contact, "ZipCode"), dig(contact, "Address", "ZipCode"), ""),
		"Phone":            phone,
		"PhoneCountryCode": phoneCc,
		"Fax":              fax,
		"FaxCountryCode":   faxCc,
	})

	return map[string]any{
		"contactType":      apiType,
		"firstName":        validated["FirstName"],
		"lastName":         validated["LastName"],
		"companyName":      coalesce(dig(contact, "Company"), ""),
		"eMail":            coalesce(dig(contact, "EMail"), ""),
		"address":          validated["AddressLine1"],
		"city":             validated["City"],
		"state":            validated["State"],
		"country":          validated["Country"],
		"postalCode":       validated["ZipCode"],
		"phoneCountryCode": validated["PhoneCountryCode"],
		"phone":            validated["Phone"],
		"faxCountryCode":   validated["FaxCountryCode"],
		"fax":              validated["Fax"],
		"isHidden":         phpBool(coalesce(dig(contact, "IsHidden"), false)),
	}
}

// payloadContacts converts the caller's contacts to the request list, keeping
// their order. The role goes through ucfirst(strtolower()), as in php-dna.
func payloadContacts(contacts Contacts) []any {
	out := make([]any, 0, len(contacts))
	for _, rc := range contacts {
		out = append(out, parseContact(rc.Contact, ucfirst(asciiLower(rc.Role))))
	}
	return out
}

// parseContactInfo converts a gateway contact to php-dna's contact array.
func parseContactInfo(data any) any {
	if phpEmpty(data) {
		return []any{}
	}
	return map[string]any{
		"ID":        coalesce(dig(data, "handle"), dig(data, "id"), ""),
		"Status":    coalesce(dig(data, "status"), "Active"),
		"AuthCode":  coalesce(dig(data, "authCode"), ""),
		"FirstName": coalesce(dig(data, "firstName"), ""),
		"LastName":  coalesce(dig(data, "lastName"), ""),
		"Company":   coalesce(dig(data, "companyName"), dig(data, "organizationName"), dig(data, "company"), ""),
		"EMail":     coalesce(dig(data, "eMail"), dig(data, "emailAddress"), dig(data, "email"), ""),
		"Type":      coalesce(dig(data, "type"), ""),
		"Address": map[string]any{
			"Line1":   coalesce(dig(data, "address"), dig(data, "addressLine1"), ""),
			"Line2":   coalesce(dig(data, "addressLine2"), ""),
			"Line3":   coalesce(dig(data, "addressLine3"), ""),
			"State":   coalesce(dig(data, "state"), dig(data, "stateOrProvince"), ""),
			"City":    coalesce(dig(data, "city"), ""),
			"Country": coalesce(dig(data, "country"), dig(data, "countryCode"), ""),
			"ZipCode": coalesce(dig(data, "postalCode"), dig(data, "zipCode"), ""),
		},
		"Phone": map[string]any{
			"Phone": map[string]any{
				"Number":      coalesce(dig(data, "phone"), dig(data, "phoneNumber"), ""),
				"CountryCode": coalesce(dig(data, "phoneCountryCode"), ""),
			},
			"Fax": map[string]any{
				"Number":      coalesce(dig(data, "fax"), dig(data, "faxNumber"), ""),
				"CountryCode": coalesce(dig(data, "faxCountryCode"), ""),
			},
		},
		"Additional": coalesce(dig(data, "additionalAttributes"), []any{}),
	}
}

// contactDefaults are the placeholders validateContact fills in.
var contactDefaults = []kv{
	{"FirstName", "Isimyok"},
	{"LastName", "Isimyok"},
	{"AddressLine1", "Addres yok"},
	{"City", "ISTANBUL"},
	{"State", "ISTANBUL"},
	{"Country", "TR"},
	{"ZipCode", "34000"},
	{"Phone", "5555555555"},
	{"PhoneCountryCode", "90"},
}

// validateContact is php-dna's validateContact(): it completes a half-empty
// contact with placeholders and applies the gateway's rules up front.
func validateContact(in map[string]any) map[string]any {
	contact := make(map[string]any, len(in)+11)
	for k, v := range in {
		contact[k] = normalizeInput(v)
	}
	for _, d := range contactDefaults {
		if contact[d.k] == nil {
			contact[d.k] = d.v
		}
	}

	contact["FirstName"] = sanitizePersonName(phpString(contact["FirstName"]))
	contact["LastName"] = sanitizePersonName(phpString(contact["LastName"]))

	blank := func(k string) bool { return phpTrim(phpString(contact[k])) == "" }

	if blank("FirstName") {
		contact["FirstName"] = "Isimyok"
	}
	if blank("LastName") {
		contact["LastName"] = contact["FirstName"]
	}
	if blank("AddressLine1") {
		contact["AddressLine1"] = "Addres yok"
	}
	if blank("City") {
		contact["City"] = "ISTANBUL"
	}
	if blank("State") {
		contact["State"] = contact["City"]
	}
	contact["Country"] = normalizeCountryCode(phpString(contact["Country"]))
	if blank("Country") {
		contact["Country"] = "TR"
	}
	if blank("ZipCode") {
		contact["ZipCode"] = "34000"
	}

	phone, phoneCc := normalizePhone(phpString(contact["Phone"]), phpString(contact["PhoneCountryCode"]))
	contact["Phone"] = phone
	contact["PhoneCountryCode"] = phoneCc
	if blank("PhoneCountryCode") {
		contact["PhoneCountryCode"] = "90"
	}
	if blank("Phone") {
		contact["Phone"] = "5555555555"
	}

	fax, faxCc := normalizePhone(phpString(contact["Fax"]), phpString(contact["FaxCountryCode"]))
	contact["Fax"] = fax
	if fax == "" {
		contact["FaxCountryCode"] = ""
	} else {
		contact["FaxCountryCode"] = faxCc
	}
	return contact
}

// sanitizePersonName keeps letters (any alphabet), whitespace, hyphens and
// apostrophes, collapses whitespace and trims: the gateway rejects digits and
// punctuation in names.
func sanitizePersonName(name string) string {
	var b strings.Builder
	if utf8.ValidString(name) {
		for _, r := range name {
			if unicode.IsLetter(r) || unicode.IsSpace(r) || r == '\'' || r == '-' {
				b.WriteRune(r)
			}
		}
	} else {
		for i := 0; i < len(name); i++ {
			c := name[i]
			if (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || c == ' ' || (c >= '\t' && c <= '\r') || c == '\'' || c == '-' {
				b.WriteByte(c)
			}
		}
	}
	var out strings.Builder
	inSpace := false
	for _, r := range b.String() {
		if unicode.IsSpace(r) {
			if !inSpace {
				out.WriteByte(' ')
			}
			inSpace = true
			continue
		}
		inSpace = false
		out.WriteRune(r)
	}
	return phpTrim(out.String())
}

var intlPrefixRe = regexp.MustCompile(`^[ \t\n\r\v\f]*(\+|00)`)

func onlyDigits(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] >= '0' && s[i] <= '9' {
			b.WriteByte(s[i])
		}
	}
	return b.String()
}

// normalizePhone splits a free-form number into (subscriber number, country
// code): digits only, no international or trunk prefix.
func normalizePhone(number, countryCode string) (string, string) {
	hasIntlPrefix := intlPrefixRe.MatchString(number)
	digits := onlyDigits(number)
	cc := onlyDigits(countryCode)

	if digits == "" {
		return "", cc
	}
	if hasIntlPrefix {
		if strings.HasPrefix(digits, "00") {
			digits = digits[2:]
		}
		if cc != "" && strings.HasPrefix(digits, cc) && len(digits) > len(cc) {
			digits = digits[len(cc):]
		}
	} else if cc != "" && strings.HasPrefix(digits, cc) && len(digits) > 10 {
		digits = digits[len(cc):]
	}
	if cc == "" {
		if len(digits) == 12 && strings.HasPrefix(digits, "90") {
			digits = digits[2:]
		}
		cc = "90"
	}
	if len(digits) > 1 && digits[0] == '0' {
		digits = strings.TrimLeft(digits, "0")
	}
	return digits, cc
}

var countryNames = map[string]string{
	"turkey": "TR", "türkiye": "TR", "turkiye": "TR",
	"united states": "US", "united states of america": "US", "usa": "US",
	"united kingdom": "GB", "great britain": "GB", "england": "GB",
	"germany": "DE", "deutschland": "DE", "france": "FR", "spain": "ES",
	"italy": "IT", "netherlands": "NL", "india": "IN", "china": "CN",
	"russia": "RU", "canada": "CA", "australia": "AU", "brazil": "BR",
}

// normalizeCountryCode returns the 2-character ISO code the gateway requires
// for the few common country names it knows, upper-cases 2-character codes
// and leaves anything else untouched.
func normalizeCountryCode(country string) string {
	country = phpTrim(country)
	if country == "" {
		return ""
	}
	if len(country) == 2 {
		return asciiUpper(country)
	}
	if code, ok := countryNames[asciiLower(country)]; ok {
		return code
	}
	return country
}

// normalizeHostName trims, drops the trailing dot and lower-cases a host.
func normalizeHostName(host string) string {
	return asciiLower(phpRtrim(phpTrim(host), "."))
}

// normalizeHostNameList normalizes a nameserver list (a list, or a map whose
// values are taken in key order): non-strings, empties and duplicates are
// dropped.
func normalizeHostNameList(hosts any) []any {
	out := []any{}
	seen := map[string]bool{}
	for _, h := range phpValues(normalizeInput(hosts)) {
		s, ok := h.(string)
		if !ok {
			continue
		}
		s = normalizeHostName(s)
		if s != "" && !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}

// normalizeTrAttributes maps legacy .tr TRABIS attributes to the REST schema:
// only TRABISDOMAINCATEGORY (Owner|Corporate), TRABISCITIZENID,
// TRABISTAXOFFICE and TRABISTAXNUMBER are kept.
func normalizeTrAttributes(attrs map[string]any) map[string]any {
	keyMap := map[string]string{
		"TRABISDOMAINCATEGORY": "TRABISDOMAINCATEGORY",
		"TRABISCITIZENID":      "TRABISCITIZENID",
		"TRABISCITIZIENID":     "TRABISCITIZENID",
		"TRABISTAXOFFICE":      "TRABISTAXOFFICE",
		"TRABISTAXNUMBER":      "TRABISTAXNUMBER",
	}
	out := map[string]any{}
	for k, v := range attrs {
		target, ok := keyMap[k]
		if !ok {
			continue
		}
		// The canonical TRABISCITIZENID always wins over the typo'd variant.
		if k == "TRABISCITIZIENID" {
			if _, has := attrs["TRABISCITIZENID"]; has {
				continue
			}
		}
		out[target] = v
	}
	if v, ok := out["TRABISDOMAINCATEGORY"]; ok && v != nil {
		switch asciiLower(phpTrim(phpString(v))) {
		case "1", "owner", "personal", "individual":
			out["TRABISDOMAINCATEGORY"] = "Owner"
		case "0", "corporate", "company":
			out["TRABISDOMAINCATEGORY"] = "Corporate"
		}
	}
	return out
}

// stringifyAttributes coerces every tldAttributes value to a string: the
// gateway deserializes tldAttributes as a string dictionary.
func stringifyAttributes(attrs map[string]any) map[string]any {
	out := make(map[string]any, len(attrs))
	for k, v := range attrs {
		switch t := v.(type) {
		case bool:
			if t {
				out[k] = "true"
			} else {
				out[k] = "false"
			}
		case nil:
			out[k] = ""
		case string, int64, float64:
			out[k] = phpString(t)
		default:
			out[k] = phpJSONEncode(t)
		}
	}
	return out
}
