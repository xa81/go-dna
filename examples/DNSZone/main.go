// Example: Manage DNS zone records
//
// Run: go run ./examples/DNSZone
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

	domain := "example.com"

	// List the zone: one row per value, names fully qualified.
	printJSON(client.GetResourceRecords(domain))

	// Make a name+type hold exactly these values (ttl 3600, apply now).
	printJSON(client.SetResourceRecordSet(domain, "@", "MX", []string{"10 mx1.example.com", "20 mx2.example.com"}, 3600, true))

	// Add, replace and remove single values, keeping the others (read-merge-write).
	printJSON(client.AddResourceRecordValue(domain, "www", "A", "192.0.2.10", nil, true))
	printJSON(client.ReplaceResourceRecordValue(domain, "www", "A", "192.0.2.10", "192.0.2.11", nil, true))
	printJSON(client.RemoveResourceRecordValue(domain, "www", "A", "192.0.2.11", true))

	// Low-level calls: stage several changes, then apply once.
	client.AddResourceRecord(domain, "api", "CNAME", "example.com.", 300, false)
	client.EditResourceRecord(domain, "api", "api", "CNAME", "www.example.com.", 300, false)
	client.DeleteResourceRecord(domain, "old.example.com.", "192.0.2.99", "A", false)
	printJSON(client.ApplyResourceRecords(domain))

	// Remove a whole name+type.
	printJSON(client.DeleteResourceRecordSet(domain, "api", "CNAME", true, ""))
}

func printJSON(v any) {
	out, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(out))
}
