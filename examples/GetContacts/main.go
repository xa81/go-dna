// Example: Get the contacts of a domain
//
// Run: go run ./examples/GetContacts
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

	result := client.GetContacts("example.com")
	printJSON(result)
	// {"result": "OK", "data": {"contacts": {"Administrative": {"ID": "...", "FirstName": "...", ...},
	//  "Billing": {...}, "Registrant": {...}, "Technical": {...}}}}
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
