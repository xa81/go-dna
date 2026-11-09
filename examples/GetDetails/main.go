// Example: Get domain details
//
// Run: go run ./examples/GetDetails
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

	result := client.GetDetails("example.com")
	if result["result"] != "OK" {
		printJSON(result["error"])
		return
	}
	data := result["data"].(map[string]any)
	fmt.Println("Status:", data["Status"])
	fmt.Println("Expires:", data["Dates"].(map[string]any)["Expiration"])
	printJSON(result)
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
