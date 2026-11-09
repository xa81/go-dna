// Example: Change the IP of a child nameserver
//
// Run: go run ./examples/ModifyChildNameServer
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

	result := client.ModifyChildNameServer("example.com", "ns1.example.com", "192.168.1.2")
	printJSON(result)
	// {"result": "OK", "data": {"NameServer": "ns1.example.com", "IPAdresses": ["192.168.1.2"]}}
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
