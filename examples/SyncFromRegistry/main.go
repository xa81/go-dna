// Example: Refresh a domain from the registry
//
// Run: go run ./examples/SyncFromRegistry
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

	result := client.SyncFromRegistry("example.com")
	printJSON(result) // same shape as GetDetails
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
