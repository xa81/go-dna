// Example: Check whether a domain can be transferred
//
// Run: go run ./examples/CheckTransfer
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

	result := client.CheckTransfer("example.com", "ABC123XYZ456")
	printJSON(result)
	// {"result": "OK", "data": {"TransferAvailabilityStatus": true, "AuthCodeIsRequired": true,
	//  "AuthCodeIsValid": true, "UserTransferRequired": false, "TransferLock": false,
	//  "Message": "Domain can be transferred", "MessageKey": "TransferAvailable", "Additional": []}}
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
