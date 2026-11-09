// Example: Register a .tr domain with TRABIS attributes
//
// Run: go run ./examples/RegisterWithContactInfo-TR
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

	individual := map[string]any{
		"TRABISDOMAINCATEGORY": 1,
		"TRABISCITIZIENID":     "11111111111",
		"TRABISNAMESURNAME":    "Ahmet Yilmaz",
		"TRABISCOUNTRYID":      "215",
		"TRABISCITYID":         "34",
		"TRABISCOUNTRYNAME":    "Türkiye",
		"TRABISCITYNAME":       "Istanbul",
	}

	commercial := map[string]any{
		"TRABISDOMAINCATEGORY": 0,
		"TRABISORGANIZATION":   "Example Ltd. Sti.",
		"TRABISTAXOFFICE":      "Istanbul VD",
		"TRABISTAXNUMBER":      "1234567890",
		"TRABISCOUNTRYID":      "215",
		"TRABISCITYID":         "34",
		"TRABISCOUNTRYNAME":    "Türkiye",
		"TRABISCITYNAME":       "Istanbul",
	}
	_ = commercial // use commercial instead of individual for a company owner

	// Legacy TRABIS names are normalized: category 1/0 becomes Owner/Corporate,
	// TRABISCITIZIENID becomes TRABISCITIZENID and the country/city/name fields
	// are dropped (the gateway takes them from the contact).
	opts := dna.DefaultRegisterOptions()
	opts.NameServers = []string{"ns1.example.com", "ns2.example.com"}
	opts.AdditionalAttributes = individual

	contacts := dna.Contacts{
		{Role: "Administrative", Contact: contact},
		{Role: "Billing", Contact: contact},
		{Role: "Technical", Contact: contact},
		{Role: "Registrant", Contact: contact},
	}
	result := client.RegisterWithContactInfo("example.com.tr", 1, contacts, opts)
	printJSON(result)
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
