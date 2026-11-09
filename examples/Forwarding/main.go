// Example: Manage URL forwarding
//
// Run: go run ./examples/Forwarding
package main

import (
	"encoding/json"
	"fmt"

	dna "github.com/domainreseller/go-dna"
)

func main() {
	// Credentials: your Reseller ID (UUID) + API Key from Panel -> API Settings.
	// Panel login credentials will NOT work.
	resellerID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx" // Reseller ID (UUID)
	apiKey := "your-api-key"                             // API Key

	client := dna.New(resellerID, apiKey, false) // true = test (OTE) environment

	domain := "example.com"

	// "Standard" (plain redirect) or "Frame" (masked).
	printJSON(client.SetForwarding(domain, "https://www.example.org", "Standard"))

	printJSON(client.GetForwarding(domain))
	// {"result": "OK", "data": {"Enabled": true, "RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

	printJSON(client.DeleteForwarding(domain))
	// {"result": "OK", "data": {"DomainName": "example.com"}}
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
