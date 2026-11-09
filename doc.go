// Package dna is the Go client for the DomainNameAPI domain-registrar REST
// gateway. It is a faithful port of the PHP library php-dna: the same 45
// methods, the same requests on the wire and the same return values.
//
// Methods return the dynamic structures php-dna returns (map[string]any and
// []any), not typed models. Failures are returned, never panicked: every
// network method returns the php-dna error envelope
//
//	{"result": "ERROR", "error": {"Code": "...", "Message": "...", "Details": "..."}}
//
// A Client keeps per-call diagnostics (GetRequestData, GetResponseData, ...)
// exactly like php-dna, so a single Client must not be used from several
// goroutines at the same time. Create one Client per goroutine instead.
//
//	c := dna.New("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", "your-api-key", false)
//	res := c.GetDetails("example.com")
//	if res["result"] == "OK" {
//		data := res["data"].(map[string]any)
//		fmt.Println(data["Status"])
//	}
package dna
