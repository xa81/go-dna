// Example: Update the contacts of a domain
//
// Run: go run ./examples/SaveContacts
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

	contact := map[string]any{
		"FirstName":        "Ahmet",
		"LastName":         "Yilmaz",
		"Company":          "",
		"EMail":            "ahmet.yilmaz@example.com",
		"AddressLine1":     "Ornek Sokak No: 123",
		"AddressLine2":     "Daire 4",
		"AddressLine3":     "",
		"City":             "Istanbul",
		"Country":          "TR",
		"Fax":              "2121234567",
		"FaxCountryCode":   "90",
		"Phone":            "5551234567",
		"PhoneCountryCode": 90,
		"Type":             "Contact",
		"ZipCode":          "34000",
		"State":            "Istanbul",
	}

	result := client.SaveContacts("example.com", dna.AllRoles(contact))
	printJSON(result)
	// {"result": "OK", "data": {"contacts": []}}
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
