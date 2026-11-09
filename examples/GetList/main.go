// Example: List the domains in the account
//
// Run: go run ./examples/GetList
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

	result := client.GetList(map[string]any{
		"OrderColumn":    "ExpirationDate",
		"OrderDirection": "ASC",
		"PageNumber":     0,
		"PageSize":       100,
	})
	printJSON(result)
	// {"result": "OK", "TotalCount": 1, "data": {"Domains": [{"ID": 564346, "Status": "Active",
	//  "DomainName": "example.com", ...}]}}
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
