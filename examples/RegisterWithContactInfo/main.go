// Example: Register a domain
//
// Run: go run ./examples/RegisterWithContactInfo
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

	opts := dna.DefaultRegisterOptions() // nameservers default, EPP lock on, privacy off
	opts.NameServers = []string{"ns1.example.com", "ns2.example.com"}
	opts.EppLock = true
	opts.PrivacyLock = false

	result := client.RegisterWithContactInfo("example.com", 1, dna.AllRoles(contact), opts)
	printJSON(result)
	// {"result": "OK", "data": {"ID": 0, "Status": "Active", "DomainName": "example.com", ...,
	//  "Dates": {"Start": "", "Expiration": "2027-11-09T12:00:00", "RemainingDays": 365}, ...}}
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
