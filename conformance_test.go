package dna

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// Conformance runner: replays every php-dna fixture from testdata/conformance
// (see php-dna/conformance/README.md) against the Go client and compares the
// requests on the wire and the return value with php-dna's.

const (
	confReseller = "00000000-0000-4000-8000-000000000001"
	confAPIKey   = "conformance-api-key"
	confBase     = "http://conformance.test/api/v1"
)

type fixtureResponse struct {
	Status int    `json:"status"`
	Body   string `json:"body"`
}

type fixtureRequest struct {
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Query   json.RawMessage   `json:"query"`
	Body    json.RawMessage   `json:"body"`
	Headers map[string]string `json:"headers"`
}

type fixture struct {
	ID        string            `json:"id"`
	Call      string            `json:"call"`
	Args      []json.RawMessage `json:"args"`
	Responses []fixtureResponse `json:"responses"`
	Requests  []fixtureRequest  `json:"requests"`
	Output    json.RawMessage   `json:"output"`
	Ignore    []string          `json:"ignore"`
}

type recordedRequest struct {
	method  string
	path    string
	query   map[string]any
	body    any
	headers map[string]string
}

// fakeGateway replays the fixture's responses in order and records requests.
type fakeGateway struct {
	t         *testing.T
	responses []fixtureResponse
	requests  []recordedRequest
}

func (f *fakeGateway) Do(req *http.Request) (*http.Response, error) {
	rec := recordedRequest{
		method:  req.Method,
		path:    strings.TrimPrefix(req.URL.Path, "/api/v1/"),
		query:   map[string]any{},
		headers: map[string]string{},
	}
	for k, v := range req.URL.Query() {
		rec.query[k] = v[len(v)-1]
	}
	if req.Body != nil {
		raw, _ := io.ReadAll(req.Body)
		if len(raw) > 0 {
			rec.body = canonical(raw)
		}
	}
	if v, ok := req.Header["X-API-KEY"]; ok {
		rec.headers["X-API-KEY"] = v[0]
	}
	if v, ok := req.Header["__reseller"]; ok {
		rec.headers["__reseller"] = v[0]
	}
	if ua := req.Header.Get("User-Agent"); ua != UserAgent {
		f.t.Errorf("User-Agent = %q, want %q", ua, UserAgent)
	}
	f.requests = append(f.requests, rec)

	if len(f.requests) > len(f.responses) {
		f.t.Errorf("request #%d (%s %s) has no queued response", len(f.requests), req.Method, rec.path)
		return newResponse(500, `{"error":{"message":"conformance: no queued response"}}`), nil
	}
	r := f.responses[len(f.requests)-1]
	return newResponse(r.Status, r.Body), nil
}

func newResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

// canonical decodes JSON keeping numbers as json.Number.
func canonical(raw []byte) any {
	dec := json.NewDecoder(bytes.NewReader(raw))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return fmt.Sprintf("<invalid json: %s>", raw)
	}
	return v
}

func roundTrip(v any) any { return canonical(encodeJSON(v)) }

// equalPHP compares two canonical JSON values with the fixture rules: an
// empty list equals an empty object, numbers compare by value, key order is
// irrelevant.
func equalPHP(a, b any, path string, diffs *[]string) {
	isEmpty := func(v any) bool {
		switch t := v.(type) {
		case []any:
			return len(t) == 0
		case map[string]any:
			return len(t) == 0
		}
		return false
	}
	if isEmpty(a) && isEmpty(b) {
		return
	}
	switch x := a.(type) {
	case json.Number:
		y, ok := b.(json.Number)
		if !ok {
			*diffs = append(*diffs, fmt.Sprintf("%s: got %#v, want number %s", path, b, x))
			return
		}
		fx, _ := x.Float64()
		fy, _ := y.Float64()
		if fx != fy {
			*diffs = append(*diffs, fmt.Sprintf("%s: got %s, want %s", path, y, x))
		}
	case map[string]any:
		y, ok := b.(map[string]any)
		if !ok {
			*diffs = append(*diffs, fmt.Sprintf("%s: got %#v, want object", path, b))
			return
		}
		keys := map[string]bool{}
		for k := range x {
			keys[k] = true
		}
		for k := range y {
			keys[k] = true
		}
		var sorted []string
		for k := range keys {
			sorted = append(sorted, k)
		}
		sort.Strings(sorted)
		for _, k := range sorted {
			xv, xok := x[k]
			yv, yok := y[k]
			switch {
			case !xok:
				*diffs = append(*diffs, fmt.Sprintf("%s.%s: unexpected key (value %#v)", path, k, yv))
			case !yok:
				*diffs = append(*diffs, fmt.Sprintf("%s.%s: missing (want %#v)", path, k, xv))
			default:
				equalPHP(xv, yv, path+"."+k, diffs)
			}
		}
	case []any:
		y, ok := b.([]any)
		if !ok {
			*diffs = append(*diffs, fmt.Sprintf("%s: got %#v, want list", path, b))
			return
		}
		if len(x) != len(y) {
			*diffs = append(*diffs, fmt.Sprintf("%s: got %d items, want %d", path, len(y), len(x)))
			return
		}
		for i := range x {
			equalPHP(x[i], y[i], fmt.Sprintf("%s[%d]", path, i), diffs)
		}
	default:
		if a != b {
			*diffs = append(*diffs, fmt.Sprintf("%s: got %#v, want %#v", path, b, a))
		}
	}
}

// takeIgnored removes a dot path from both values; the actual value must be
// a non-negative integer.
func takeIgnored(t *testing.T, want, got any, path string) {
	parts := strings.Split(path, ".")
	wm, gm := want, got
	for _, p := range parts[:len(parts)-1] {
		wm = wm.(map[string]any)[p]
		gm = gm.(map[string]any)[p]
	}
	last := parts[len(parts)-1]
	gmap, ok := gm.(map[string]any)
	if !ok {
		t.Errorf("ignored path %s: no object in result", path)
		return
	}
	n, ok := gmap[last].(json.Number)
	if !ok {
		t.Errorf("ignored path %s: got %#v, want a non-negative integer", path, gmap[last])
	} else if i, err := n.Int64(); err != nil || i < 0 {
		t.Errorf("ignored path %s: got %s, want a non-negative integer", path, n)
	}
	delete(gmap, last)
	delete(wm.(map[string]any), last)
}

// args decodes positional arguments with defaults.
type args struct {
	t   *testing.T
	raw []json.RawMessage
}

func (a args) has(i int) bool { return i < len(a.raw) }

func (a args) any(i int) any {
	if !a.has(i) {
		return nil
	}
	return decodeJSON(string(a.raw[i]))
}

func (a args) str(i int, def string) string {
	if !a.has(i) {
		return def
	}
	return phpString(a.any(i))
}

func (a args) num(i int, def int) int {
	if !a.has(i) {
		return def
	}
	return int(phpInt(a.any(i)))
}

func (a args) boolean(i int, def bool) bool {
	if !a.has(i) {
		return def
	}
	return phpBool(a.any(i))
}

func (a args) ttlPtr(i int) *int {
	if !a.has(i) || a.any(i) == nil {
		return nil
	}
	n := int(phpInt(a.any(i)))
	return &n
}

func (a args) strings(i int) []string {
	var out []string
	for _, v := range phpValues(a.any(i)) {
		out = append(out, phpString(v))
	}
	return out
}

func (a args) object(i int) map[string]any {
	m, _ := a.any(i).(map[string]any)
	return m
}

// contacts decodes a role => contact object keeping the key order.
func (a args) contacts(i int) Contacts {
	if !a.has(i) {
		return nil
	}
	dec := json.NewDecoder(bytes.NewReader(a.raw[i]))
	tok, err := dec.Token()
	if err != nil {
		a.t.Fatalf("contacts arg: %v", err)
	}
	if d, ok := tok.(json.Delim); !ok || d != '{' {
		return nil // [] = no contacts
	}
	var out Contacts
	for dec.More() {
		k, _ := dec.Token()
		var raw json.RawMessage
		if err := dec.Decode(&raw); err != nil {
			a.t.Fatalf("contacts arg: %v", err)
		}
		m, _ := decodeJSON(string(raw)).(map[string]any)
		out = append(out, RoleContact{Role: k.(string), Contact: m})
	}
	return out
}

// call maps a php-dna call to the Go method, filling php-dna's defaults.
func call(c *Client, name string, a args) any {
	switch name {
	case "getResellerDetails":
		return c.GetResellerDetails()
	case "getCurrentBalance":
		return c.GetCurrentBalance(a.str(0, "USD"))
	case "checkAvailability":
		return c.CheckAvailability(a.strings(0), a.strings(1), a.num(2, 1), a.str(3, "create"))
	case "getList":
		return c.GetList(a.object(0))
	case "getTldList":
		return c.GetTldList(a.num(0, 20))
	case "getDetails":
		return c.GetDetails(a.str(0, ""))
	case "syncFromRegistry":
		return c.SyncFromRegistry(a.str(0, ""))
	case "getContacts":
		return c.GetContacts(a.str(0, ""))
	case "saveContacts":
		return c.SaveContacts(a.str(0, ""), a.contacts(1))
	case "modifyNameServer":
		return c.ModifyNameServer(a.str(0, ""), a.any(1))
	case "enableTheftProtectionLock":
		return c.EnableTheftProtectionLock(a.str(0, ""))
	case "disableTheftProtectionLock":
		return c.DisableTheftProtectionLock(a.str(0, ""))
	case "modifyPrivacyProtectionStatus":
		return c.ModifyPrivacyProtectionStatus(a.str(0, ""), a.boolean(1, false), a.str(2, DefaultReason))
	case "renew":
		return c.Renew(a.str(0, ""), a.num(1, 1))
	case "registerWithContactInfo":
		opts := DefaultRegisterOptions()
		if a.has(3) {
			opts.NameServers = a.any(3)
		}
		opts.EppLock = a.boolean(4, true)
		opts.PrivacyLock = a.boolean(5, false)
		opts.AdditionalAttributes = a.object(6)
		return c.RegisterWithContactInfo(a.str(0, ""), a.num(1, 1), a.contacts(2), opts)
	case "addChildNameServer":
		return c.AddChildNameServer(a.str(0, ""), a.str(1, ""), a.str(2, ""))
	case "deleteChildNameServer":
		return c.DeleteChildNameServer(a.str(0, ""), a.str(1, ""))
	case "modifyChildNameServer":
		return c.ModifyChildNameServer(a.str(0, ""), a.str(1, ""), a.str(2, ""))
	case "transfer":
		return c.Transfer(a.str(0, ""), a.str(1, ""), a.num(2, 1), a.contacts(3))
	case "cancelTransfer":
		return c.CancelTransfer(a.str(0, ""))
	case "approveTransfer":
		return c.ApproveTransfer(a.str(0, ""))
	case "rejectTransfer":
		return c.RejectTransfer(a.str(0, ""))
	case "checkTransfer":
		return c.CheckTransfer(a.str(0, ""), a.str(1, ""))
	case "getResourceRecords":
		return c.GetResourceRecords(a.str(0, ""))
	case "addResourceRecord":
		return c.AddResourceRecord(a.str(0, ""), a.str(1, ""), a.str(2, ""), a.any(3), a.num(4, 3600), a.boolean(5, true))
	case "editResourceRecord":
		return c.EditResourceRecord(a.str(0, ""), a.str(1, ""), a.str(2, ""), a.str(3, ""), a.any(4), a.num(5, 3600), a.boolean(6, true))
	case "deleteResourceRecord":
		return c.DeleteResourceRecord(a.str(0, ""), a.str(1, ""), a.str(2, ""), a.str(3, ""), a.boolean(4, true))
	case "applyResourceRecords":
		return c.ApplyResourceRecords(a.str(0, ""))
	case "setResourceRecordSet":
		return c.SetResourceRecordSet(a.str(0, ""), a.str(1, ""), a.str(2, ""), a.strings(3), a.num(4, 3600), a.boolean(5, true))
	case "deleteResourceRecordSet":
		return c.DeleteResourceRecordSet(a.str(0, ""), a.str(1, ""), a.str(2, ""), a.boolean(3, true), a.str(4, ""))
	case "addResourceRecordValue":
		return c.AddResourceRecordValue(a.str(0, ""), a.str(1, ""), a.str(2, ""), a.str(3, ""), a.ttlPtr(4), a.boolean(5, true))
	case "removeResourceRecordValue":
		return c.RemoveResourceRecordValue(a.str(0, ""), a.str(1, ""), a.str(2, ""), a.str(3, ""), a.boolean(4, true))
	case "replaceResourceRecordValue":
		return c.ReplaceResourceRecordValue(a.str(0, ""), a.str(1, ""), a.str(2, ""), a.str(3, ""), a.str(4, ""), a.ttlPtr(5), a.boolean(6, true))
	case "getForwarding":
		return c.GetForwarding(a.str(0, ""))
	case "setForwarding":
		return c.SetForwarding(a.str(0, ""), a.str(1, ""), a.str(2, "Standard"))
	case "deleteForwarding":
		return c.DeleteForwarding(a.str(0, ""))
	case "isTrTLD":
		return c.IsTrTLD(a.str(0, ""))
	case "validateContact":
		m := a.object(0)
		if m == nil {
			m = map[string]any{}
		}
		return c.ValidateContact(m)
	}
	return fmt.Sprintf("<unmapped call %s>", name)
}

func TestConformance(t *testing.T) {
	files, err := filepath.Glob(filepath.Join("testdata", "conformance", "*.json"))
	if err != nil || len(files) == 0 {
		t.Fatalf("no fixtures found: %v", err)
	}
	if len(files) != 126 {
		t.Errorf("found %d fixtures, want 126", len(files))
	}
	calls := map[string]bool{}
	for _, file := range files {
		raw, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		var fx fixture
		if err := json.Unmarshal(raw, &fx); err != nil {
			t.Fatalf("%s: %v", file, err)
		}
		calls[fx.Call] = true
		t.Run(fx.ID, func(t *testing.T) { runFixture(t, fx) })
	}
	t.Logf("%d fixtures over %d php-dna methods", len(files), len(calls))
}

func runFixture(t *testing.T, fx fixture) {
	gw := &fakeGateway{t: t, responses: fx.Responses}
	c := New(confReseller, confAPIKey, false, WithHTTPClient(gw))
	c.SetErrorReportingEnabled(false)
	c.SetServiceURL(confBase)
	c.zoneMinInterval = 0

	result := call(c, fx.Call, args{t: t, raw: fx.Args})
	if s, ok := result.(string); ok && strings.HasPrefix(s, "<unmapped") {
		t.Fatalf("%s", s)
	}
	got := roundTrip(result)
	want := canonical(fx.Output)

	if len(gw.requests) != len(fx.Requests) {
		t.Errorf("sent %d requests, want %d", len(gw.requests), len(fx.Requests))
	}
	for i := 0; i < len(gw.requests) && i < len(fx.Requests); i++ {
		g, w := gw.requests[i], fx.Requests[i]
		var diffs []string
		if g.method != w.Method {
			diffs = append(diffs, fmt.Sprintf("method: got %s, want %s", g.method, w.Method))
		}
		if g.path != w.Path {
			diffs = append(diffs, fmt.Sprintf("path: got %s, want %s", g.path, w.Path))
		}
		equalPHP(canonical(w.Query), roundTrip(g.query), "query", &diffs)
		var wantBody any
		if len(w.Body) > 0 && string(w.Body) != "null" {
			wantBody = canonical(w.Body)
		}
		equalPHP(wantBody, g.body, "body", &diffs)
		for _, h := range []string{"X-API-KEY", "__reseller"} {
			if g.headers[h] != w.Headers[h] {
				diffs = append(diffs, fmt.Sprintf("header %s: got %q, want %q", h, g.headers[h], w.Headers[h]))
			}
		}
		for _, d := range diffs {
			t.Errorf("request #%d %s", i+1, d)
		}
	}

	for _, p := range fx.Ignore {
		takeIgnored(t, want, got, p)
	}
	var diffs []string
	equalPHP(want, got, "output", &diffs)
	for _, d := range diffs {
		t.Errorf("%s", d)
	}
}
