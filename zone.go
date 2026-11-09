package dna

import (
	"fmt"
	"strings"
)

// relativizeRecordName reduces a record name to the relative label the
// gateway expects: "@", the apex and "" become "", "www.example.com" becomes
// "www".
func relativizeRecordName(name, domainName string) string {
	name = phpTrim(name)
	domain := phpRtrim(asciiLower(phpTrim(domainName)), ".")
	if name == "" || name == "@" {
		return ""
	}
	bare := phpRtrim(name, ".")
	if asciiEqualFold(bare, domain) {
		return ""
	}
	suffix := "." + domain
	if len(bare) > len(suffix) && asciiEqualFold(bare[len(bare)-len(suffix):], suffix) {
		return bare[:len(bare)-len(suffix)]
	}
	return bare
}

// qualifyRecordName turns a relative label into the FQDN form the zone
// listing returns ("www" => "www.example.com."); qualified names pass through.
func qualifyRecordName(name, domainName string) string {
	name = phpTrim(name)
	domain := phpRtrim(asciiLower(phpTrim(domainName)), ".")
	apex := domain + "."
	if name == "" || name == "@" {
		return apex
	}
	bare := phpRtrim(name, ".")
	if asciiEqualFold(bare, domain) || strings.HasSuffix(asciiLower(bare), "."+domain) {
		return bare + "."
	}
	return bare + "." + apex
}

// unsupportedZoneType returns an error message for a record type the zone
// API cannot store, or "" when the type is supported.
func unsupportedZoneType(typ string) string {
	typ = asciiUpper(phpTrim(typ))
	if typ == "" {
		return "Record type is required."
	}
	for _, t := range ZoneRecordTypes {
		if t == typ {
			return ""
		}
	}
	return fmt.Sprintf("Record type \"%s\" is not supported by the DNS zone API. Supported types: %s.", typ, strings.Join(ZoneRecordTypes, ", "))
}

// contentList is `array_values(array_map('strval', (array) $content))`.
func contentList(content any) []any {
	out := []any{}
	switch t := normalizeInput(content).(type) {
	case nil:
	case []any, map[string]any:
		for _, v := range phpValues(t) {
			out = append(out, phpString(v))
		}
	default:
		out = append(out, phpString(t))
	}
	return out
}

// buildZoneStruct builds the ZoneCreateDto; ttl is clamped to 1..86400 (a
// value below 1 becomes 3600).
func buildZoneStruct(name, typ string, content any, ttl int) map[string]any {
	if ttl < 1 {
		ttl = 3600
	}
	if ttl > 86400 {
		ttl = 86400
	}
	return map[string]any{
		"name":     name,
		"ttl":      int64(ttl),
		"type":     asciiUpper(typ),
		"contents": contentList(content),
	}
}

func zoneEndpoint(path string, p params) string {
	return path + "?" + httpBuildQuery(p)
}

// GetResourceRecords lists the DNS zone of a domain, one row per value:
//
//	{"result": "OK", "data": {"records": [{"Name": "www.example.com.",
//	 "Type": "A", "TTL": 300, "Content": "192.0.2.1", "Disabled": false}, ...]}}
//
// Names come back fully qualified with a trailing dot.
func (c *Client) GetResourceRecords(domainName string) map[string]any {
	response, err := c.request("GET", "domains/zones", params{{"domainName", domainName}})
	if err != nil {
		return c.errorEnvelope(err)
	}
	records := []any{}
	sets := phpValues(response)
	if response != nil && !isArray(response) {
		sets = []any{response}
	}
	for _, set := range sets {
		if !isArray(set) {
			continue
		}
		name := coalesce(dig(set, "name"), "")
		typ := coalesce(dig(set, "type"), "")
		ttl := phpInt(coalesce(dig(set, "ttl"), int64(0)))
		rows := coalesce(dig(set, "records"), []any{})
		if !isArray(rows) || phpEmpty(rows) {
			records = append(records, map[string]any{"Name": name, "Type": typ, "TTL": ttl, "Content": "", "Disabled": false})
			continue
		}
		for _, rec := range phpValues(rows) {
			records = append(records, map[string]any{
				"Name":     name,
				"Type":     typ,
				"TTL":      ttl,
				"Content":  coalesce(dig(rec, "content"), ""),
				"Disabled": !phpEmpty(dig(rec, "disabled")),
			})
		}
	}
	return map[string]any{"result": ResultOK, "data": map[string]any{"records": records}}
}

// AddResourceRecord writes a record set: content is one value (string) or a
// list ([]string). name may be "@", the apex, a relative label or an FQDN.
// ttl is clamped to 1..86400 (below 1 means 3600); php-dna's default is
// 3600. With apply the change is committed at once (php-dna's default:
// true). Note: the gateway replaces the whole name+type on every write; use
// AddResourceRecordValue to append one value.
func (c *Client) AddResourceRecord(domainName, name, recordType string, content any, ttl int, apply bool) map[string]any {
	if why := unsupportedZoneType(recordType); why != "" {
		return errorResult("ADD_ZONE", why)
	}
	name = relativizeRecordName(name, domainName)
	body := map[string]any{"zoneStruct": buildZoneStruct(name, recordType, content, ttl)}
	if _, err := c.request("POST", zoneEndpoint("domains/zones", params{{"domainName", domainName}}), body); err != nil {
		return c.errorEnvelope(err)
	}
	if apply {
		c.ApplyResourceRecords(domainName)
	}
	return map[string]any{
		"result": ResultOK,
		"data":   map[string]any{"Name": name, "Type": recordType, "TTL": int64(ttl), "Content": content},
	}
}

// EditResourceRecord rewrites the record set identified by recordName. ttl
// and apply as in AddResourceRecord.
func (c *Client) EditResourceRecord(domainName, recordName, name, recordType string, content any, ttl int, apply bool) map[string]any {
	if why := unsupportedZoneType(recordType); why != "" {
		return errorResult("EDIT_ZONE", why)
	}
	name = relativizeRecordName(name, domainName)
	body := map[string]any{"zoneStruct": buildZoneStruct(name, recordType, content, ttl)}
	endpoint := zoneEndpoint("domains/zones", params{{"domainName", domainName}, {"recordName", recordName}})
	if _, err := c.request("PUT", endpoint, body); err != nil {
		return c.errorEnvelope(err)
	}
	if apply {
		c.ApplyResourceRecords(domainName)
	}
	return map[string]any{
		"result": ResultOK,
		"data":   map[string]any{"Name": name, "Type": recordType, "TTL": int64(ttl), "Content": content},
	}
}

// DeleteResourceRecord deletes a record. Note: the gateway removes the whole
// name+type set, not only that value; use RemoveResourceRecordValue to
// remove one value. apply defaults to true in php-dna.
func (c *Client) DeleteResourceRecord(domainName, name, record, recordType string, apply bool) map[string]any {
	p := params{{"domainName", domainName}, {"Name", name}, {"Record", record}, {"RecordType", recordType}}
	if _, err := c.request("DELETE", "domains/zones", p); err != nil {
		return c.errorEnvelope(err)
	}
	if apply {
		c.ApplyResourceRecords(domainName)
	}
	return map[string]any{
		"result": ResultOK,
		"data":   map[string]any{"Name": name, "Type": recordType, "Content": record},
	}
}

// ApplyResourceRecords commits staged zone changes:
// {"result": "OK", "data": {"DomainName": "example.com"}}.
func (c *Client) ApplyResourceRecords(domainName string) map[string]any {
	if _, err := c.request("POST", zoneEndpoint("domains/zones/apply", params{{"domainName", domainName}}), []any{}); err != nil {
		return c.errorEnvelope(err)
	}
	return map[string]any{"result": ResultOK, "data": map[string]any{"DomainName": domainName}}
}

// SetResourceRecordSet makes a name+type hold exactly contents (empty values
// are dropped; an empty list removes the set). php-dna's defaults: ttl
// 3600, apply true.
//
//	{"result": "OK", "data": {"Name": "", "Type": "MX", "TTL": 3600,
//	 "Contents": ["10 mx1.example.com", "20 mx2.example.com"]}}
func (c *Client) SetResourceRecordSet(domainName, name, recordType string, contents []string, ttl int, apply bool) map[string]any {
	if why := unsupportedZoneType(recordType); why != "" {
		return errorResult("SET_ZONE", why)
	}
	name = relativizeRecordName(name, domainName)
	values := []any{}
	for _, v := range contents {
		if phpTrim(v) != "" {
			values = append(values, v)
		}
	}
	if len(values) == 0 {
		return c.DeleteResourceRecordSet(domainName, name, recordType, apply, "")
	}
	body := map[string]any{"zoneStruct": buildZoneStruct(name, recordType, values, ttl)}
	if _, err := c.request("POST", zoneEndpoint("domains/zones", params{{"domainName", domainName}}), body); err != nil {
		return c.errorEnvelope(err)
	}
	if apply {
		c.ApplyResourceRecords(domainName)
	}
	return map[string]any{
		"result": ResultOK,
		"data":   map[string]any{"Name": name, "Type": asciiUpper(recordType), "TTL": int64(ttl), "Contents": values},
	}
}

// DeleteResourceRecordSet removes every value under name+type. knownValue
// is any current value (saves a lookup); pass "" to look it up. Removing a
// set that does not exist succeeds with "Removed": false. php-dna's
// defaults: apply true, knownValue "".
func (c *Client) DeleteResourceRecordSet(domainName, name, recordType string, apply bool, knownValue string) map[string]any {
	typ := asciiUpper(phpTrim(recordType))
	value := knownValue
	if value == "" {
		current := c.GetResourceRecords(domainName)
		if current["result"] != ResultOK {
			return current
		}
		fqdn := qualifyRecordName(name, domainName)
		for _, rec := range phpValues(dig(current, "data", "records")) {
			if asciiEqualFold(phpString(dig(rec, "Name")), fqdn) && asciiEqualFold(phpString(dig(rec, "Type")), typ) {
				value = phpString(dig(rec, "Content"))
				break
			}
		}
		if value == "" {
			return map[string]any{"result": ResultOK, "data": map[string]any{"Name": name, "Type": typ, "Removed": false}}
		}
	}
	return c.DeleteResourceRecord(domainName, qualifyRecordName(name, domainName), value, typ, apply)
}

// recordSet is one name+type read from the zone.
type recordSet struct {
	ttl      int64
	contents []any
}

// readResourceRecordSet reads one set; an absent set is empty. On a read
// failure it returns the error envelope.
func (c *Client) readResourceRecordSet(domainName, name, typ string) (*recordSet, map[string]any) {
	current := c.GetResourceRecords(domainName)
	if current["result"] != ResultOK {
		return nil, current
	}
	fqdn := qualifyRecordName(name, domainName)
	typ = asciiUpper(phpTrim(typ))
	set := &recordSet{contents: []any{}}
	for _, rec := range phpValues(dig(current, "data", "records")) {
		if !asciiEqualFold(phpString(dig(rec, "Name")), fqdn) || !asciiEqualFold(phpString(dig(rec, "Type")), typ) {
			continue
		}
		if v := phpString(dig(rec, "Content")); v != "" {
			set.contents = append(set.contents, v)
		}
		set.ttl = phpInt(dig(rec, "TTL"))
	}
	return set, nil
}

// sameRecordValue compares record values the way the gateway means them:
// TXT without surrounding quotes, host targets case-insensitively and
// without the trailing dot.
func sameRecordValue(typ, a, b string) bool {
	norm := func(v string) string {
		v = phpTrim(v)
		if asciiUpper(typ) == "TXT" {
			if len(v) > 1 && v[0] == '"' && v[len(v)-1] == '"' {
				v = v[1 : len(v)-1]
			}
			return v
		}
		return asciiLower(phpRtrim(v, "."))
	}
	return norm(a) == norm(b)
}

func setTTL(ttl *int, set *recordSet) int {
	if ttl != nil {
		return *ttl
	}
	if set.ttl != 0 {
		return int(set.ttl)
	}
	return 3600
}

func stringsOf(values []any) []string {
	out := make([]string, len(values))
	for i, v := range values {
		out[i] = phpString(v)
	}
	return out
}

// AddResourceRecordValue adds one value to a name+type, keeping the values
// already there (read-merge-write). ttl nil keeps the set's TTL (3600 for a
// new set). Adding a value that is already present changes nothing and
// returns "Changed": false. apply defaults to true in php-dna.
func (c *Client) AddResourceRecordValue(domainName, name, recordType, value string, ttl *int, apply bool) map[string]any {
	if why := unsupportedZoneType(recordType); why != "" {
		return errorResult("ADD_ZONE", why)
	}
	name = relativizeRecordName(name, domainName)
	set, errEnv := c.readResourceRecordSet(domainName, name, recordType)
	if errEnv != nil {
		return errEnv
	}
	for _, existing := range set.contents {
		if sameRecordValue(recordType, phpString(existing), value) {
			return map[string]any{
				"result": ResultOK,
				"data":   map[string]any{"Name": name, "Type": asciiUpper(recordType), "Contents": set.contents, "Changed": false},
			}
		}
	}
	contents := append(stringsOf(set.contents), value)
	return c.SetResourceRecordSet(domainName, name, recordType, contents, setTTL(ttl, set), apply)
}

// RemoveResourceRecordValue removes one value from a name+type, leaving the
// others; the set is deleted only when that was its last value. Removing a
// value that is not there returns "Changed": false.
func (c *Client) RemoveResourceRecordValue(domainName, name, recordType, value string, apply bool) map[string]any {
	if why := unsupportedZoneType(recordType); why != "" {
		return errorResult("DELETE_ZONE", why)
	}
	name = relativizeRecordName(name, domainName)
	set, errEnv := c.readResourceRecordSet(domainName, name, recordType)
	if errEnv != nil {
		return errEnv
	}
	var remaining []string
	found := false
	for _, existing := range set.contents {
		if !found && sameRecordValue(recordType, phpString(existing), value) {
			found = true
			continue
		}
		remaining = append(remaining, phpString(existing))
	}
	if !found {
		return map[string]any{
			"result": ResultOK,
			"data":   map[string]any{"Name": name, "Type": asciiUpper(recordType), "Contents": set.contents, "Changed": false},
		}
	}
	if len(remaining) == 0 {
		return c.DeleteResourceRecordSet(domainName, name, recordType, apply, value)
	}
	return c.SetResourceRecordSet(domainName, name, recordType, remaining, setTTL(nil, set), apply)
}

// ReplaceResourceRecordValue swaps one value of a name+type for another,
// keeping the others. ttl nil keeps the set's TTL. An oldValue that is not
// in the set is an EDIT_ZONE error.
func (c *Client) ReplaceResourceRecordValue(domainName, name, recordType, oldValue, newValue string, ttl *int, apply bool) map[string]any {
	if why := unsupportedZoneType(recordType); why != "" {
		return errorResult("EDIT_ZONE", why)
	}
	name = relativizeRecordName(name, domainName)
	set, errEnv := c.readResourceRecordSet(domainName, name, recordType)
	if errEnv != nil {
		return errEnv
	}
	var contents []string
	found := false
	for _, existing := range set.contents {
		if !found && sameRecordValue(recordType, phpString(existing), oldValue) {
			found = true
			contents = append(contents, newValue)
			continue
		}
		contents = append(contents, phpString(existing))
	}
	if !found {
		why := fmt.Sprintf("No %s record \"%s\" with value \"%s\" to edit.", asciiUpper(recordType), name, oldValue)
		return errorResult("EDIT_ZONE", why)
	}
	return c.SetResourceRecordSet(domainName, name, recordType, contents, setTTL(ttl, set), apply)
}
