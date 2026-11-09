// Example: Change the nameservers of a domain
//
// Run: go run ./examples/ModifyNameServer
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

	// A list, or a map whose values are taken in key order.
	result := client.ModifyNameServer("example.com", []string{"ns1.example.com", "ns2.example.com"})
	printJSON(result)
	// {"result": "OK", "data": {"NameServers": ["ns1.example.com", "ns2.example.com"]}}
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
