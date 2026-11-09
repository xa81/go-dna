// Example: Get the reseller account and balances
//
// Run: go run ./examples/GetResellerDetails
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

	result := client.GetResellerDetails()
	printJSON(result)
	// {"result": "OK", "id": "...", "active": true, "name": "...", "balance": "500", "currency": "USD",
	//  "symbol": "$", "balances": [{"balance": "500", "currency": "USD", "symbol": "$"}, ...]}
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
