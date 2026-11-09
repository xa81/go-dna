package dna

// PHP semantics helpers.
//
// php-dna works on loosely typed PHP arrays. To return exactly what php-dna
// returns, this port reproduces the handful of PHP conversions php-dna relies
// on: (string), (int), (float), empty(), is_numeric(), number_format(),
// http_build_query(), json_decode() and a few ASCII-only string functions.

import (
	"bytes"
	"encoding/json"
	"io"
	"math"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

// kv is one ordered key/value pair; params keeps the insertion order PHP
// arrays have, so query strings come out in php-dna's order.
type kv struct {
	k string
	v any
}

type params []kv

func (p params) toMap() map[string]any {
	m := make(map[string]any, len(p))
	for _, e := range p {
		m[e.k] = e.v
	}
	return m
}

// phpWhitespace is the default character list of PHP's trim().
const phpWhitespace = " \t\n\r\x00\x0b"

func phpTrim(s string) string { return strings.Trim(s, phpWhitespace) }

func phpRtrim(s, chars string) string { return strings.TrimRight(s, chars) }

// asciiLower and asciiUpper mirror PHP 8 strtolower/strtoupper, which only
// touch ASCII letters.
func asciiLower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}

func asciiUpper(s string) string {
	b := []byte(s)
	for i, c := range b {
		if c >= 'a' && c <= 'z' {
			b[i] = c - 32
		}
	}
	return string(b)
}

// ucfirst mirrors PHP ucfirst (ASCII only).
func ucfirst(s string) string {
	if s == "" {
		return s
	}
	c := s[0]
	if c >= 'a' && c <= 'z' {
		return string(c-32) + s[1:]
	}
	return s
}

// asciiEqualFold mirrors strcasecmp(a, b) === 0.
func asciiEqualFold(a, b string) bool { return asciiLower(a) == asciiLower(b) }

// normalizeInput converts common Go container types (map[string]string,
// []string, nested maps, ...) to the map[string]any / []any shapes php-dna's
// logic works on. Scalars are normalised to string, bool, int64 or float64.
func normalizeInput(v any) any {
	switch t := v.(type) {
	case nil:
		return nil
	case string, bool, int64, float64:
		return t
	case int:
		return int64(t)
	case int8:
		return int64(t)
	case int16:
		return int64(t)
	case int32:
		return int64(t)
	case uint:
		return int64(t)
	case uint8:
		return int64(t)
	case uint16:
		return int64(t)
	case uint32:
		return int64(t)
	case uint64:
		if t > math.MaxInt64 {
			return float64(t)
		}
		return int64(t)
	case float32:
		return float64(t)
	case json.Number:
		return convertNumber(t)
	case map[string]any:
		m := make(map[string]any, len(t))
		for k, e := range t {
			m[k] = normalizeInput(e)
		}
		return m
	case []any:
		l := make([]any, len(t))
		for i, e := range t {
			l[i] = normalizeInput(e)
		}
		return l
	case params:
		return t.toMap()
	}
	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Map:
		if rv.Type().Key().Kind() != reflect.String {
			m := make(map[string]any, rv.Len())
			iter := rv.MapRange()
			for iter.Next() {
				m[phpString(normalizeInput(iter.Key().Interface()))] = normalizeInput(iter.Value().Interface())
			}
			return m
		}
		m := make(map[string]any, rv.Len())
		iter := rv.MapRange()
		for iter.Next() {
			m[iter.Key().String()] = normalizeInput(iter.Value().Interface())
		}
		return m
	case reflect.Slice, reflect.Array:
		if rv.Kind() == reflect.Slice && rv.IsNil() {
			return []any{}
		}
		l := make([]any, rv.Len())
		for i := 0; i < rv.Len(); i++ {
			l[i] = normalizeInput(rv.Index(i).Interface())
		}
		return l
	case reflect.Ptr, reflect.Interface:
		if rv.IsNil() {
			return nil
		}
		return normalizeInput(rv.Elem().Interface())
	case reflect.String:
		return rv.String()
	case reflect.Bool:
		return rv.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int()
	case reflect.Float32, reflect.Float64:
		return rv.Float()
	}
	return v
}

// sortedKeys returns the keys of m in ascending order. Go maps have no
// insertion order, so wherever php-dna iterates a caller-supplied array the
// port iterates its keys sorted.
func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// phpValues returns what `foreach ($v as $item)` visits: list items, or map
// values in key order. Anything else visits nothing.
func phpValues(v any) []any {
	switch t := v.(type) {
	case []any:
		return t
	case map[string]any:
		out := make([]any, 0, len(t))
		for _, k := range sortedKeys(t) {
			out = append(out, t[k])
		}
		return out
	}
	return nil
}

// isArray mirrors is_array(): a decoded JSON list or object.
func isArray(v any) bool {
	switch v.(type) {
	case []any, map[string]any:
		return true
	}
	return false
}

// dig follows map keys like PHP's `$a['x']['y'] ?? null`: it never fails and
// returns nil when a level is missing, null or not an object.
func dig(v any, keys ...string) any {
	for _, k := range keys {
		m, ok := v.(map[string]any)
		if !ok {
			return nil
		}
		v = m[k]
		if v == nil {
			return nil
		}
	}
	return v
}

// coalesce returns the first non-nil value, like PHP's ?? chain.
func coalesce(vals ...any) any {
	for _, v := range vals {
		if v != nil {
			return v
		}
	}
	return nil
}

// phpEmpty mirrors empty().
func phpEmpty(v any) bool {
	switch t := v.(type) {
	case nil:
		return true
	case bool:
		return !t
	case int64:
		return t == 0
	case int:
		return t == 0
	case float64:
		return t == 0
	case string:
		return t == "" || t == "0"
	case []any:
		return len(t) == 0
	case map[string]any:
		return len(t) == 0
	case []string:
		return len(t) == 0
	case map[string]string:
		return len(t) == 0
	}
	return false
}

// phpBool mirrors a (bool) cast.
func phpBool(v any) bool { return !phpEmpty(v) }

// phpString mirrors a (string) cast.
func phpString(v any) string {
	switch t := v.(type) {
	case nil:
		return ""
	case string:
		return t
	case bool:
		if t {
			return "1"
		}
		return ""
	case int64:
		return strconv.FormatInt(t, 10)
	case int:
		return strconv.Itoa(t)
	case float64:
		return phpFloatString(t, 14)
	case []any, map[string]any:
		return "Array"
	}
	return phpString(normalizeInput(v))
}

// phpFloatString renders a float the way PHP does with the given precision
// (14 for a (string) cast; 17 for serialize_precision=-1 contexts).
func phpFloatString(f float64, precision int) string {
	switch {
	case math.IsNaN(f):
		return "NAN"
	case math.IsInf(f, 1):
		return "INF"
	case math.IsInf(f, -1):
		return "-INF"
	case f == 0:
		if math.Signbit(f) {
			return "-0"
		}
		return "0"
	}
	var mant string
	var exp int
	if precision >= 17 {
		e := strconv.FormatFloat(f, 'e', -1, 64)
		i := strings.IndexByte(e, 'e')
		mant = e[:i]
		exp, _ = strconv.Atoi(e[i+1:])
	} else {
		e := strconv.FormatFloat(f, 'e', precision-1, 64)
		i := strings.IndexByte(e, 'e')
		mant = e[:i]
		exp, _ = strconv.Atoi(e[i+1:])
	}
	threshold := precision
	if precision >= 17 {
		threshold = 15
	}
	if exp < -4 || exp >= threshold {
		if strings.Contains(mant, ".") {
			mant = strings.TrimRight(mant, "0")
			mant = strings.TrimSuffix(mant, ".")
		}
		if !strings.Contains(mant, ".") {
			mant += ".0"
		}
		sign := "+"
		if exp < 0 {
			sign = "-"
			exp = -exp
		}
		return mant + "E" + sign + strconv.Itoa(exp)
	}
	r, _ := strconv.ParseFloat(mant+"e"+strconv.Itoa(exp), 64)
	return strconv.FormatFloat(r, 'f', -1, 64)
}

var numericPrefix = regexp.MustCompile(`^[ \t\n\r\v\f]*[+-]?(\d+(\.\d*)?|\.\d+)([eE][+-]?\d+)?`)
var numericFull = regexp.MustCompile(`^[ \t\n\r\v\f]*[+-]?(\d+(\.\d*)?|\.\d+)([eE][+-]?\d+)?[ \t\n\r\v\f]*$`)

// isNumeric mirrors is_numeric().
func isNumeric(v any) bool {
	switch t := v.(type) {
	case int64, int, float64:
		return true
	case string:
		return numericFull.MatchString(t)
	}
	return false
}

// phpFloat mirrors a (float) cast.
func phpFloat(v any) float64 {
	switch t := v.(type) {
	case nil:
		return 0
	case bool:
		if t {
			return 1
		}
		return 0
	case int64:
		return float64(t)
	case int:
		return float64(t)
	case float64:
		return t
	case string:
		m := numericPrefix.FindString(t)
		if m == "" {
			return 0
		}
		f, _ := strconv.ParseFloat(strings.TrimLeft(m, " \t\n\r\v\f"), 64)
		return f
	case []any:
		if len(t) == 0 {
			return 0
		}
		return 1
	case map[string]any:
		if len(t) == 0 {
			return 0
		}
		return 1
	}
	return phpFloat(normalizeInput(v))
}

// phpInt mirrors an (int) cast.
func phpInt(v any) int64 {
	switch t := v.(type) {
	case nil:
		return 0
	case bool:
		if t {
			return 1
		}
		return 0
	case int64:
		return t
	case int:
		return int64(t)
	case float64:
		if math.IsNaN(t) || math.IsInf(t, 0) || t >= 9.2233720368547758e18 || t < -9.2233720368547758e18 {
			return 0
		}
		return int64(t)
	case string:
		m := numericPrefix.FindString(t)
		if m == "" {
			return 0
		}
		m = strings.TrimLeft(m, " \t\n\r\v\f")
		if strings.ContainsAny(m, ".eE") {
			return phpInt(phpFloat(m))
		}
		n, err := strconv.ParseInt(m, 10, 64)
		if err != nil {
			if strings.HasPrefix(m, "-") {
				return math.MinInt64
			}
			return math.MaxInt64
		}
		return n
	case []any:
		if len(t) == 0 {
			return 0
		}
		return 1
	case map[string]any:
		if len(t) == 0 {
			return 0
		}
		return 1
	}
	return phpInt(normalizeInput(v))
}

// phpRound mirrors PHP's round() (half away from zero, with PHP's
// pre-rounding to 15 significant digits so 1.005 rounds the way PHP does).
func phpRound(f float64, places int) float64 {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return f
	}
	pow := math.Pow(10, float64(places))
	tmp := f * pow
	if pre, err := strconv.ParseFloat(strconv.FormatFloat(tmp, 'g', 15, 64), 64); err == nil {
		tmp = pre
	}
	r := math.Round(tmp)
	res := r / pow
	if math.IsInf(res, 0) || math.IsNaN(res) {
		return f
	}
	return res
}

// numberFormat mirrors PHP number_format with "." as decimal point and no thousands separator.
func numberFormat(f float64, decimals int) string {
	r := phpRound(f, decimals)
	s := strconv.FormatFloat(r, 'f', decimals, 64)
	if strings.HasPrefix(s, "-") && strings.Trim(s, "-0.") == "" {
		s = s[1:]
	}
	return s
}

// phpURLEncode mirrors urlencode(): letters, digits and "-_." stay, space
// becomes "+", everything else is %XX.
func phpURLEncode(s string) string {
	var b strings.Builder
	const hex = "0123456789ABCDEF"
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case c >= 'a' && c <= 'z', c >= 'A' && c <= 'Z', c >= '0' && c <= '9', c == '-', c == '_', c == '.':
			b.WriteByte(c)
		case c == ' ':
			b.WriteByte('+')
		default:
			b.WriteByte('%')
			b.WriteByte(hex[c>>4])
			b.WriteByte(hex[c&15])
		}
	}
	return b.String()
}

// httpBuildQuery mirrors http_build_query(): true is 1, false is 0, null is
// skipped and nested arrays become key[sub]=value.
func httpBuildQuery(p params) string {
	var parts []string
	for _, e := range p {
		parts = appendQuery(parts, phpURLEncode(e.k), normalizeInput(e.v))
	}
	return strings.Join(parts, "&")
}

func appendQuery(parts []string, key string, v any) []string {
	switch t := v.(type) {
	case nil:
		return parts
	case bool:
		if t {
			return append(parts, key+"=1")
		}
		return append(parts, key+"=0")
	case float64:
		return append(parts, key+"="+phpURLEncode(phpFloatString(t, 17)))
	case []any:
		for i, e := range t {
			parts = appendQuery(parts, key+"%5B"+strconv.Itoa(i)+"%5D", e)
		}
		return parts
	case map[string]any:
		for _, k := range sortedKeys(t) {
			parts = appendQuery(parts, key+"%5B"+phpURLEncode(k)+"%5D", t[k])
		}
		return parts
	}
	return append(parts, key+"="+phpURLEncode(phpString(v)))
}

// decodeJSON mirrors json_decode($s, true): objects become map[string]any,
// integers int64 (floats when they overflow), other numbers float64. It
// returns nil for empty or invalid input, like PHP.
func decodeJSON(s string) any {
	if strings.TrimSpace(s) == "" {
		return nil
	}
	dec := json.NewDecoder(strings.NewReader(s))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil
	}
	if _, err := dec.Token(); err != io.EOF {
		return nil // trailing data
	}
	return convertNumbers(v)
}

func convertNumbers(v any) any {
	switch t := v.(type) {
	case json.Number:
		return convertNumber(t)
	case map[string]any:
		for k, e := range t {
			t[k] = convertNumbers(e)
		}
		return t
	case []any:
		for i, e := range t {
			t[i] = convertNumbers(e)
		}
		return t
	}
	return v
}

func convertNumber(n json.Number) any {
	s := string(n)
	if !strings.ContainsAny(s, ".eE") {
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			return i
		}
	}
	f, _ := strconv.ParseFloat(s, 64)
	return f
}

// encodeJSON encodes a request body. HTML characters are not escaped, the
// way PHP's json_encode leaves them alone.
func encodeJSON(v any) []byte {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetEscapeHTML(false)
	if err := enc.Encode(v); err != nil {
		return []byte("null")
	}
	return bytes.TrimRight(buf.Bytes(), "\n")
}

// phpJSONEncode mirrors json_encode() with default flags: "/" is escaped as
// "\/" and non-ASCII characters as \uXXXX. Used where php-dna embeds JSON in a
// string value.
func phpJSONEncode(v any) string {
	var b strings.Builder
	writePHPJSON(&b, normalizeInput(v))
	return b.String()
}

func writePHPJSON(b *strings.Builder, v any) {
	switch t := v.(type) {
	case nil:
		b.WriteString("null")
	case bool:
		if t {
			b.WriteString("true")
		} else {
			b.WriteString("false")
		}
	case int64:
		b.WriteString(strconv.FormatInt(t, 10))
	case float64:
		s := phpFloatString(t, 17)
		if t == math.Trunc(t) && !strings.ContainsAny(s, "E") && math.Abs(t) < 1e15 {
			s = strconv.FormatFloat(t, 'f', -1, 64) + ".0"
		}
		b.WriteString(s)
	case string:
		writePHPJSONString(b, t)
	case []any:
		b.WriteByte('[')
		for i, e := range t {
			if i > 0 {
				b.WriteByte(',')
			}
			writePHPJSON(b, e)
		}
		b.WriteByte(']')
	case map[string]any:
		b.WriteByte('{')
		for i, k := range sortedKeys(t) {
			if i > 0 {
				b.WriteByte(',')
			}
			writePHPJSONString(b, k)
			b.WriteByte(':')
			writePHPJSON(b, t[k])
		}
		b.WriteByte('}')
	default:
		writePHPJSONString(b, phpString(t))
	}
}

func writePHPJSONString(b *strings.Builder, s string) {
	b.WriteByte('"')
	for i := 0; i < len(s); {
		r, size := utf8.DecodeRuneInString(s[i:])
		i += size
		switch {
		case r == '"':
			b.WriteString(`\"`)
		case r == '\\':
			b.WriteString(`\\`)
		case r == '/':
			b.WriteString(`\/`)
		case r == '\n':
			b.WriteString(`\n`)
		case r == '\r':
			b.WriteString(`\r`)
		case r == '\t':
			b.WriteString(`\t`)
		case r == '\b':
			b.WriteString(`\b`)
		case r == '\f':
			b.WriteString(`\f`)
		case r < 0x20 || r >= 0x7f:
			if r == utf8.RuneError && size == 1 {
				b.WriteString(`�`)
				continue
			}
			if r > 0xffff {
				r -= 0x10000
				hi := 0xd800 + (r>>10)&0x3ff
				lo := 0xdc00 + r&0x3ff
				b.WriteString(`\u` + hex4(hi) + `\u` + hex4(lo))
				continue
			}
			b.WriteString(`\u` + hex4(r))
		default:
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
}

func hex4(r rune) string {
	s := strconv.FormatInt(int64(r), 16)
	for len(s) < 4 {
		s = "0" + s
	}
	return s
}
