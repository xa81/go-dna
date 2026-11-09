// Example: Enable the registrar (theft protection) lock
//
// Run: go run ./examples/EnableTheftProtectionLock
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

	result := client.EnableTheftProtectionLock("example.com")
	printJSON(result)
	// {"result": "OK", "data": {"LockStatus": true}}
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
