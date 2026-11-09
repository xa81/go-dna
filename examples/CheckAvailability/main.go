// Example: Check domain availability
//
// Run: go run ./examples/CheckAvailability
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

	result := client.CheckAvailability([]string{"example", "testdomain"}, []string{"com", "net"}, 1, "create")

	// Success: a list of results. Failure: the error envelope (a map).
	if list, ok := result.([]any); ok {
		for _, item := range list {
			row := item.(map[string]any)
			fmt.Printf("%s.%s: %s (%v %v)\n", row["DomainName"], row["TLD"], row["Status"], row["Price"], row["Currency"])
		}
		return
	}
	printJSON(result)
	// [{"TLD": "com", "DomainName": "example", "Status": "notavailable", "Command": "create",
	//   "Period": 1, "IsFee": false, "Price": "11.9900", "Currency": "USD", "Reason": "Domain exists"}, ...]
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
