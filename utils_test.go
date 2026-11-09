package dna

import (
	"math"
	"reflect"
	"testing"
)

func TestPHPString(t *testing.T) {
	cases := []struct {
		in   any
		want string
	}{
		{nil, ""}, {true, "1"}, {false, ""}, {int64(42), "42"}, {601.64, "601.64"}, {1250.5, "1250.5"},
		{99.999, "99.999"}, {0.1 + 0.2, "0.3"}, {1e25, "1.0E+25"}, {0.00001, "1.0E-5"}, {-2.5, "-2.5"},
		{float64(3), "3"}, {math.Inf(1), "INF"}, {"x", "x"}, {[]any{1}, "Array"}, {7, "7"},
	}
	for _, tc := range cases {
		if got := phpString(tc.in); got != tc.want {
			t.Errorf("phpString(%#v) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestPHPIntAndNumeric(t *testing.T) {
	ints := []struct {
		in   any
		want int64
	}{
		{"12abc", 12}, {" 7", 7}, {"1e3", 1000}, {"abc", 0}, {3.9, 3}, {-3.9, -3}, {true, 1}, {nil, 0}, {"", 0},
	}
	for _, tc := range ints {
		if got := phpInt(tc.in); got != tc.want {
			t.Errorf("phpInt(%#v) = %d, want %d", tc.in, got, tc.want)
		}
	}
	numeric := map[any]bool{"3": true, " 3 ": true, "1.5e3": true, "0x1A": false, "abc": false, "": false, int64(3): true, true: false}
	for in, want := range numeric {
		if got := isNumeric(in); got != want {
			t.Errorf("isNumeric(%#v) = %v, want %v", in, got, want)
		}
	}
}

func TestNumberFormat(t *testing.T) {
	cases := []struct {
		f    float64
		d    int
		want string
	}{
		{11.99, 4, "11.9900"}, {601.64, 2, "601.64"}, {1.005, 2, "1.01"}, {2.675, 2, "2.68"},
		{0, 2, "0.00"}, {-0.001, 2, "0.00"}, {1250.5, 2, "1250.50"}, {2500, 4, "2500.0000"},
	}
	for _, tc := range cases {
		if got := numberFormat(tc.f, tc.d); got != tc.want {
			t.Errorf("numberFormat(%v, %d) = %q, want %q", tc.f, tc.d, got, tc.want)
		}
	}
}

func TestHTTPBuildQuery(t *testing.T) {
	got := httpBuildQuery(params{
		{"a", "x y/z"}, {"t", true}, {"f", false}, {"n", nil}, {"i", int64(5)},
		{"list", []any{"p", "q"}}, {"m", map[string]any{"k": "v"}},
	})
	want := "a=x+y%2Fz&t=1&f=0&i=5&list%5B0%5D=p&list%5B1%5D=q&m%5Bk%5D=v"
	if got != want {
		t.Errorf("httpBuildQuery = %q, want %q", got, want)
	}
}

func TestDecodeJSON(t *testing.T) {
	v := decodeJSON(`{"i":1,"f":1.5,"big":12345678901234567890,"l":[true,null]}`).(map[string]any)
	if v["i"] != int64(1) || v["f"] != 1.5 || v["big"] != 12345678901234567890.0 {
		t.Errorf("numbers decoded as %#v", v)
	}
	for _, s := range []string{"", "  ", "nope", "{} trailing"} {
		if decodeJSON(s) != nil {
			t.Errorf("decodeJSON(%q) must be nil", s)
		}
	}
}

func TestPHPJSONEncode(t *testing.T) {
	got := phpJSONEncode(map[string]any{"a": "ü/x", "b": []any{int64(1), 2.0, nil}})
	want := "{\"a\":\"\\u00fc\\/x\",\"b\":[1,2.0,null]}"
	if got != want {
		t.Errorf("phpJSONEncode = %s, want %s", got, want)
	}
}

func TestNormalizePhone(t *testing.T) {
	cases := []struct{ number, cc, wantN, wantCC string }{
		{"+90 (532) 374 33 88", "", "5323743388", "90"},
		{"0532 374 33 88", "90", "5323743388", "90"},
		{"905323743388", "90", "5323743388", "90"},
		{"0090 212 123 45 67", "", "2121234567", "90"},
		{"+1 (202) 555-0100", "1", "2025550100", "1"},
		{"+1 (202) 555-0100", "90", "2025550100", "1"},
		{"+44 20 7946 0958", "90", "2079460958", "44"},
		{"+358 40 123 4567", "", "401234567", "358"},
		{"", "44", "", "44"},
		{"0", "", "0", "90"},
	}
	for _, tc := range cases {
		n, cc := normalizePhone(tc.number, tc.cc)
		if n != tc.wantN || cc != tc.wantCC {
			t.Errorf("normalizePhone(%q, %q) = %q, %q; want %q, %q", tc.number, tc.cc, n, cc, tc.wantN, tc.wantCC)
		}
	}
}

func TestCallingCodePrefix(t *testing.T) {
	cases := map[string]string{"12025550100": "1", "79161234567": "7", "905321112233": "90", "358401234567": "358", "0123": "", "": "", "35": ""}
	for in, want := range cases {
		if got := callingCodePrefix(in); got != want {
			t.Errorf("callingCodePrefix(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestParseBoolFlag(t *testing.T) {
	cases := map[any]bool{true: true, false: false, 1: true, 0: false, 2: false, "1": true, " On ": true, "TRUE": true, "yes": true, "off": false, "": false, nil: false}
	for in, want := range cases {
		if got := parseBoolFlag(in); got != want {
			t.Errorf("parseBoolFlag(%#v) = %v, want %v", in, got, want)
		}
	}
}

func TestSanitizePersonName(t *testing.T) {
	cases := map[string]string{
		"  Jane 3rd ":    "Jane rd",
		"Ömer  Çağlar":   "Ömer Çağlar",
		"O'Neil-Smith!":  "O'Neil-Smith",
		"12345":          "",
		"\xffBad\xfe Na": "Bad Na",
	}
	for in, want := range cases {
		if got := sanitizePersonName(in); got != want {
			t.Errorf("sanitizePersonName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestValidateContact(t *testing.T) {
	c := New("r", "k", false)
	got := c.ValidateContact(map[string]any{
		"FirstName": "007", "LastName": "", "City": "Ankara", "State": "", "Country": "germany",
		"ZipCode": 6100, "Phone": "0312 555 44 33", "PhoneCountryCode": 90, "Fax": "", "FaxCountryCode": "90",
	})
	want := map[string]any{
		"FirstName": "Isimyok", "LastName": "Isimyok", "AddressLine1": "Addres yok", "City": "Ankara", "State": "Ankara",
		"Country": "DE", "ZipCode": int64(6100), "Phone": "3125554433", "PhoneCountryCode": "90", "Fax": "", "FaxCountryCode": "",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("ValidateContact =\n%#v\nwant\n%#v", got, want)
	}
}

func TestCountryCodes(t *testing.T) {
	cases := map[string]string{"tr": "TR", " Turkey ": "TR", "Türkiye": "TR", "USA": "US", "Narnia": "Narnia", "": ""}
	for in, want := range cases {
		if got := normalizeCountryCode(in); got != want {
			t.Errorf("normalizeCountryCode(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIsTrTLD(t *testing.T) {
	c := New("r", "k", false)
	cases := map[string]bool{"example.com.tr": true, "EXAMPLE.TR": true, "example.trade": false, "tr": false, "x.com": false}
	for in, want := range cases {
		if got := c.IsTrTLD(in); got != want {
			t.Errorf("IsTrTLD(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestRecordNames(t *testing.T) {
	rel := map[[2]string]string{
		{"@", "example.com"}: "", {"", "example.com"}: "", {"Example.COM.", "example.com"}: "",
		{"www.example.com", "example.com"}: "www", {"a.b.EXAMPLE.com.", "example.com"}: "a.b", {"www", "example.com"}: "www",
	}
	for in, want := range rel {
		if got := relativizeRecordName(in[0], in[1]); got != want {
			t.Errorf("relativizeRecordName(%q, %q) = %q, want %q", in[0], in[1], got, want)
		}
	}
	q := map[[2]string]string{
		{"@", "Example.com"}: "example.com.", {"www", "example.com"}: "www.example.com.",
		{"www.example.com.", "example.com"}: "www.example.com.", {"example.com", "example.com"}: "example.com.",
	}
	for in, want := range q {
		if got := qualifyRecordName(in[0], in[1]); got != want {
			t.Errorf("qualifyRecordName(%q, %q) = %q, want %q", in[0], in[1], got, want)
		}
	}
	if !sameRecordValue("TXT", `"v=spf1 -all"`, "v=spf1 -all") || !sameRecordValue("MX", "10 MX1.example.com.", "10 mx1.example.com") {
		t.Error("sameRecordValue must ignore TXT quotes, trailing dots and case")
	}
}

func TestTrAttributes(t *testing.T) {
	got := stringifyAttributes(normalizeTrAttributes(map[string]any{
		"TRABISDOMAINCATEGORY": int64(1), "TRABISCITIZIENID": "1", "TRABISCITIZENID": "2", "TRABISCITYID": "34",
	}))
	want := map[string]any{"TRABISDOMAINCATEGORY": "Owner", "TRABISCITIZENID": "2"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	got = stringifyAttributes(map[string]any{"a": true, "b": nil, "c": 2.5, "d": []any{"x"}})
	want = map[string]any{"a": "true", "b": "", "c": "2.5", "d": `["x"]`}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestContactsHelpers(t *testing.T) {
	c := ContactsFromMap(map[string]map[string]any{"Registrant": {}, "Zeta": {}, "Administrative": {}, "Billing": {}, "Technical": {}})
	var roles []string
	for _, rc := range c {
		roles = append(roles, rc.Role)
	}
	if !reflect.DeepEqual(roles, []string{"Administrative", "Billing", "Technical", "Registrant", "Zeta"}) {
		t.Errorf("roles = %v", roles)
	}
	if len(AllRoles(map[string]any{"FirstName": "A"})) != 4 {
		t.Error("AllRoles must give four roles")
	}
}

func TestNameServerList(t *testing.T) {
	got := normalizeHostNameList(map[string]string{"ns2": "NS2.Example.NET.", "ns1": " ns1.example.net ", "ns3": "", "ns4": "ns1.example.net"})
	want := []any{"ns1.example.net", "ns2.example.net"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
	if got := normalizeHostNameList("not a list"); len(got) != 0 {
		t.Errorf("non-list must give an empty list, got %v", got)
	}
}

func TestFormatAPIDate(t *testing.T) {
	cases := map[any]string{
		"2027-03-15T10:20:30":         "2027-03-15T10:20:30",
		"2027-03-15T10:20:30.1234567": "2027-03-15T10:20:30",
		"2027-03-15T10:20:30Z":        "2027-03-15T10:20:30",
		"2027-03-15T12:20:30+02:00":   "2027-03-15T10:20:30",
		"2027-03-15":                  "2027-03-15T00:00:00",
		"0001-01-01T00:00:00":         "0001-01-01T00:00:00",
		"not-a-date":                  "",
		"":                            "",
		nil:                           "",
	}
	for in, want := range cases {
		if got := formatAPIDate(in); got != want {
			t.Errorf("formatAPIDate(%#v) = %q, want %q", in, got, want)
		}
	}
}
