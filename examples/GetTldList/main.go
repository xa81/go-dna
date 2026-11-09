// Example: Get the TLD list with prices
//
// Run: go run ./examples/GetTldList
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

	result := client.GetTldList(2)
	printJSON(result)
	// {"result": "OK", "data": [{"id": 1, "status": "Active", "tld": "com",
	//  "pricing": {"registration": {"1": "10.8100"}, ...}, "currencies": {"registration": "USD", ...}, ...}]}
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
