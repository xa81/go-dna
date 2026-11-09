// Example: Complete and normalize a contact (no network call)
//
// Run: go run ./examples/ValidateContact
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

	contact := client.ValidateContact(map[string]any{
		"FirstName": " Jo ",
		"Phone":     "+90 (532) 374 33 88",
		"Country":   "Turkey",
	})
	printJSON(contact)
	// {"FirstName": "Jo", "LastName": "Jo", "AddressLine1": "Addres yok", "City": "ISTANBUL",
	//  "State": "ISTANBUL", "Country": "TR", "ZipCode": "34000", "Phone": "5323743388",
	//  "PhoneCountryCode": "90", "Fax": "", "FaxCountryCode": ""}
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
