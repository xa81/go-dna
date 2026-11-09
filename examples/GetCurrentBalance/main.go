// Example: Get the balance in one currency
//
// Run: go run ./examples/GetCurrentBalance
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

	result := client.GetCurrentBalance("USD") // USD, TRY, EUR or GBP
	printJSON(result)
	// {"ErrorCode": 0, "OperationMessage": "Command completed succesfully.", "OperationResult": "SUCCESS",
	//  "Balance": "601.64", "CurrencyId": 2, "CurrencyInfo": null, "CurrencyName": "USD", "CurrencySymbol": "$"}
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
