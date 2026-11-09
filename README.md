<div align="center">
  <a href="README-TR.md"   >   TR <img style="padding-top: 8px" src="https://raw.githubusercontent.com/yammadev/flag-icons/master/png/TR.png" alt="TR" height="20" /></a>
  <a href="README.md"> | EN <img style="padding-top: 8px" src="https://raw.githubusercontent.com/yammadev/flag-icons/master/png/US.png" alt="EN" height="20" /></a>
   <a href="README-DE.md"> | DE <img style="padding-top: 8px" src="https://raw.githubusercontent.com/yammadev/flag-icons/master/png/DE.png" alt="DE" height="20" /></a>
  <a href="README-SA.md"> | SA <img style="padding-top: 8px" src="https://raw.githubusercontent.com/yammadev/flag-icons/master/png/SA.png" alt="AR" height="20" /></a>
  <a href="README-NL.md"> | NL <img style="padding-top: 8px" src="https://raw.githubusercontent.com/yammadev/flag-icons/master/png/NL.png" alt="NL" height="20" /></a>
  <a href="README-AZ.md"> | AZ <img style="padding-top: 8px" src="https://raw.githubusercontent.com/yammadev/flag-icons/master/png/AZ.png" alt="AZ" height="20" /></a>
  <a href="README-CN.md"> | CN <img style="padding-top: 8px" src="https://raw.githubusercontent.com/yammadev/flag-icons/master/png/CN.png" alt="CN" height="20" /></a>
  <a href="README-FR.md"> | FR <img style="padding-top: 8px" src="https://raw.githubusercontent.com/yammadev/flag-icons/master/png/FR.png" alt="FR" height="20" /></a>
  <a href="README-IT.md"> | IT <img style="padding-top: 8px" src="https://raw.githubusercontent.com/yammadev/flag-icons/master/png/IT.png" alt="IT" height="20" /></a>
</div>

# go-dna — DomainNameAPI client for Go

Go client for the DomainNameAPI domain-registrar REST API. It is a faithful port of
[php-dna](https://github.com/domainreseller/php-dna): the same methods, the same requests on the
wire and the same return values. The standard library is the only dependency.

## 📦 Download — always use Releases!

⬇️ **Get the latest tested version here: https://github.com/domainreseller/go-dna/releases/latest**

> ⚠️ Do **not** use the green **Code → Download ZIP** button — that downloads the raw development branch. Release packages are versioned, tested and production-ready.

## Requirements

- Go 1.21 or higher
- No third-party packages (only `net/http` and `encoding/json` from the standard library)

## Installation

With Go modules, pin a released version:

```bash
go get github.com/domainreseller/go-dna@v1.0.0
```

```go
import dna "github.com/domainreseller/go-dna"
```

The package name is `dna`. Runnable examples for every method are in the [examples](examples) folder:

```bash
go run ./examples/GetDetails
```

## 🔑 API Credentials — Reseller ID & API Key

The library needs exactly two credentials, passed to the constructor:

| Parameter | Value |
|---|---|
| First | Reseller ID — UUID like `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` |
| Second | API Key |

> 💡 Find your **Reseller ID** and **API Key** in your DomainNameAPI panel under **API Settings**.
> ⚠️ These are **API credentials** — your panel login details will **not** work here.

```go
resellerID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
apiKey := "your-api-key"

// Live environment
client := dna.New(resellerID, apiKey, false)

// Test (OTE) environment
client = dna.New(resellerID, apiKey, true)
```

Both credentials are trimmed. The third argument also accepts the strings panels store for
checkboxes: `"true"`, `"1"`, `"on"` and `"yes"` (any case) select the test environment, everything
else (`"false"`, `"0"`, `"off"`, `""`, ...) selects live.

## Quick start

```go
package main

import (
	"fmt"

	dna "github.com/domainreseller/go-dna"
)

func main() {
	client := dna.New("xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx", "your-api-key", false)

	result := client.GetDetails("example.com")
	if result["result"] != "OK" {
		e := result["error"].(map[string]any)
		fmt.Println("Error:", e["Code"], e["Message"], e["Details"])
		return
	}
	data := result["data"].(map[string]any)
	fmt.Println(data["DomainName"], data["Status"], data["Dates"].(map[string]any)["Expiration"])
}
```

## Return values

Methods return the same dynamic structures php-dna returns — `map[string]any` and `[]any` — not
typed models, so the keys, string flags (`"true"`/`"false"`) and formatted money strings
(`"11.9900"`) are exactly php-dna's. Numbers from the gateway are `int64` or `float64`.

Successful calls look like this:

```go
map[string]any{"result": "OK", "data": map[string]any{ /* ... */ }}
```

Failed calls never panic; they return the error envelope:

```go
map[string]any{
	"result": "ERROR",
	"error": map[string]any{
		"Code":        "API_404", // API_<HTTP status>, CURL_<n>, or a library code such as RESPONSE
		"Message":     "The domain name you requested could not be found. Please check and try again.",
		"Details":     "The domain could not be found. Please check and try again.",
		"GatewayCode": "Dna.DomainService:Domain:10007", // only with API_<status>, when the gateway sent its own code
	},
}
```

Three methods use a different shape, exactly as in php-dna:

| Method | Success | Failure |
|---|---|---|
| `CheckAvailability` | a bare `[]any` list | the error envelope (`map[string]any`), so the return type is `any` |
| `GetCurrentBalance` | `OperationResult: "SUCCESS"` envelope | the same envelope with `OperationResult: "FAILED"` |
| `GetResellerDetails` | flat map (`id`, `name`, `balance`, `balances`, ...) | the error envelope |

### Optional arguments

Go has no default arguments. The methods keep php-dna's parameter order and take every argument
explicitly; the php-dna default is shown in each section below. A few zero values select the
default for convenience: `GetCurrentBalance("")` means USD, `GetTldList(0)` means 20,
`CheckAvailability(..., "")` means `create` and `SetForwarding(..., "")` means `Standard`.
`RegisterWithContactInfo` takes its four optional arguments in a `*RegisterOptions` (nil = defaults).
`AddResourceRecordValue` and `ReplaceResourceRecordValue` take `ttl *int`: `nil` keeps the set's
current TTL.

### Contacts

Contacts are `map[string]any` with php-dna's keys (`FirstName`, `LastName`, `Company`, `EMail`,
`AddressLine1`, `AddressLine2`, `City`, `State`, `Country`, `ZipCode`, `Phone`,
`PhoneCountryCode`, `Fax`, `FaxCountryCode`). The nested shape (`Address.Line1`,
`Phone.Phone.Number`, ...) is accepted too. Because Go maps have no order, the four roles are
passed as an ordered `dna.Contacts` list:

```go
contacts := dna.Contacts{
	{Role: "Administrative", Contact: contact},
	{Role: "Billing", Contact: contact},
	{Role: "Technical", Contact: contact},
	{Role: "Registrant", Contact: contact},
}
contacts = dna.AllRoles(contact)                 // the same contact for all four roles
contacts = dna.ContactsFromMap(map[string]map[string]any{ /* role => contact */ })
```

Half-empty contacts are completed by `ValidateContact`'s rules before they are sent. Contacts are
used by `RegisterWithContactInfo` and `SaveContacts`; transfers never send contacts.

### Method names

Go method names follow Go conventions; the php-dna name is the same in lowerCamelCase.
`getServiceUrl`/`setServiceUrl` are `GetServiceURL`/`SetServiceURL` and `setSslVerification` is
`SetSSLVerification`.

## Usage

### Domain Registration

php-dna defaults: default nameservers (`ns1.domainnameapi.com`, `ns2.domainnameapi.com`), EPP lock
on, privacy off, no additional attributes. Additional attributes are only required for domains such
as `.tr`.

```go
contact := map[string]any{
	"FirstName":        "John",
	"LastName":         "Doe",
	"Company":          "Example Corp",
	"EMail":            "john.doe@example.com",
	"AddressLine1":     "123 Lorem Street",
	"AddressLine2":     "Suite 456",
	"City":             "Springfield",
	"Country":          "US",
	"Fax":              "1234567890",
	"FaxCountryCode":   "1",
	"Phone":            "9876543210",
	"PhoneCountryCode": 1,
	"ZipCode":          "12345",
	"State":            "IL",
}

opts := dna.DefaultRegisterOptions()
opts.NameServers = []string{"ns1.example.com", "ns2.example.com"}
opts.EppLock = true
opts.PrivacyLock = false
// Only for .tr domains:
opts.AdditionalAttributes = map[string]any{
	"TRABISDOMAINCATEGORY": 1,             // 1 = individual (Owner), 0 = company (Corporate)
	"TRABISCITIZIENID":     "12345678901", // or TRABISTAXOFFICE + TRABISTAXNUMBER for companies
}

result := client.RegisterWithContactInfo("example.com", 1, dna.AllRoles(contact), opts)
```

For `.tr` domains the legacy TRABIS names are normalized to the gateway schema: the category
becomes `Owner`/`Corporate`, `TRABISCITIZIENID` becomes `TRABISCITIZENID`, and the country, city,
name and organization fields are dropped (the gateway takes them from the contact). Every attribute
value is sent as a string.

<details>
<summary>Sample Output for Domain Registration</summary>

```go
map[string]any{
	"result": "OK",
	"data": map[string]any{
		"ID":                      int64(0),
		"Status":                  "Active",
		"DomainName":              "",
		"AuthCode":                "",
		"LockStatus":              "false",
		"PrivacyProtectionStatus": "false",
		"IsChildNameServer":       "false",
		"Contacts": map[string]any{
			"Billing":        map[string]any{"ID": int64(0)},
			"Technical":      map[string]any{"ID": int64(0)},
			"Administrative": map[string]any{"ID": int64(0)},
			"Registrant":     map[string]any{"ID": int64(0)},
		},
		"Dates": map[string]any{
			"Start":         "",
			"Expiration":    "2027-11-09T12:00:00",
			"RemainingDays": int64(365),
		},
		"NameServers":      []any{},
		"Additional":       []any{},
		"ChildNameServers": []any{},
	},
}
```

The register endpoint answers with a short record, so most fields are empty; call `GetDetails`
afterwards for the full domain. When the gateway refuses the registration (HTTP 200 with
`"success": false`) the result is an error with code `REGISTER_DOMAIN` and the gateway's message.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain Renewal

```go
result := client.Renew("example.com", 1)
```

<details>
<summary>Sample Output for Domain Renewal</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"ExpirationDate": "2026-03-15T10:00:00"}}
```

A refusal from the gateway (`"success": false`) returns code `RENEW_DOMAIN`; a success without an
expiration date returns `DOMAIN_RENEW`.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain Transfer

```go
result := client.Transfer("example.com", "Xy9#mK2$", 1)
```

Transfers never send contacts: the registrar takes them over with the domain, so the request always
carries an empty contact list and the method has no contacts parameter. The result has the same
shape as `GetDetails`; a refusal from the gateway (`"success": false`) returns code
`TRANSFER_DOMAIN` with the gateway's message.

<details>
<summary>Sample Output for Domain Transfer</summary>

```go
map[string]any{
	"result": "OK",
	"data": map[string]any{
		"ID": int64(0), "Status": "Active", "DomainName": "", "AuthCode": "",
		"LockStatus": "false", "PrivacyProtectionStatus": "false", "IsChildNameServer": "false",
		"Contacts": map[string]any{ /* Billing, Technical, Administrative, Registrant: {"ID": 0} */ },
		"Dates": map[string]any{"Start": "", "Expiration": "2027-11-09T12:00:00", "RemainingDays": int64(365)},
		"NameServers": []any{}, "Additional": []any{}, "ChildNameServers": []any{},
	},
}
```

</details>

Related transfer calls:

```go
client.CancelTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Cancelled"}}
client.ApproveTransfer("example.com") // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Approved"}}
client.RejectTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Rejected"}}
```

<hr style="border: 4px solid #000; border-style: dashed;">

### Transfer Check

```go
result := client.CheckTransfer("example.com", "Ab1!xyZ9")
```

`result` is `OK` only when the registry reports the domain transferable; otherwise it is `ERROR`
with code `TRANSFER_NOT_AVAILABLE` and the same `data`.

<details>
<summary>Sample Output for Transfer Check</summary>

```go
map[string]any{
	"result": "OK",
	"data": map[string]any{
		"TransferAvailabilityStatus": true,
		"AuthCodeIsRequired":         true,
		"AuthCodeIsValid":            true,
		"UserTransferRequired":       false,
		"TransferLock":               false,
		"Message":                    "Domain can be transferred",
		"MessageKey":                 "TransferAvailable",
		"Additional":                 []any{},
	},
}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain List

```go
result := client.GetList(map[string]any{
	"OrderColumn":    "Id",
	"OrderDirection": "ASC",
	"PageNumber":     0,
	"PageSize":       1000,
})
```

`PageNumber` (0-based), `PageSize`, `OrderColumn` and `OrderDirection` are translated to the
gateway's `SkipCount`, `MaxResultCount` and `Sorting`; other keys (for example `DomainName`) are
passed through as filters. `nil` returns the first 200 domains.

<details>
<summary>Sample Output for Domain List</summary>

```go
map[string]any{
	"result":     "OK",
	"TotalCount": int64(1),
	"data": map[string]any{
		"Domains": []any{
			map[string]any{
				"ID":                      int64(564346),
				"Status":                  "Active",
				"DomainName":              "example.com",
				"AuthCode":                "",
				"LockStatus":              "true",
				"PrivacyProtectionStatus": "false",
				"IsChildNameServer":       "false",
				"Contacts": map[string]any{
					"Billing":        map[string]any{"ID": ""},
					"Technical":      map[string]any{"ID": ""},
					"Administrative": map[string]any{"ID": ""},
					"Registrant":     map[string]any{"ID": ""},
				},
				"Dates": map[string]any{
					"Start":         "2024-03-15T10:20:30",
					"Expiration":    "2027-03-15T10:20:30",
					"RemainingDays": int64(126),
				},
				"NameServers":      []any{"ns1.example.com", "ns2.example.com"},
				"Additional":       []any{},
				"ChildNameServers": []any{},
			},
		},
	},
}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### TLD List

```go
result := client.GetTldList(100) // php-dna default: 20
```

<details>
<summary>Sample Output for TLD List</summary>

```go
map[string]any{
	"result": "OK",
	"data": []any{
		map[string]any{
			"id": int64(1), "status": "Active", "tld": "com",
			"maxchar": int64(63), "maxperiod": int64(10), "minchar": int64(1), "minperiod": int64(1),
			"pricing": map[string]any{
				"registration": map[string]any{"1": "10.8100", "2": "21.6200"},
				"renew":        map[string]any{"1": "11.0100"},
				"transfer":     map[string]any{"1": "10.6100"},
				"restore":      map[string]any{"1": "99.9000"},
			},
			"currencies": map[string]any{
				"registration": "USD", "renew": "USD", "transfer": "USD", "restore": "USD",
			},
		},
	},
}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain Availability Check

```go
result := client.CheckAvailability([]string{"hello", "world123x0"}, []string{"com", "net"}, 1, "create")
if list, ok := result.([]any); ok {
	// success
	_ = list
}
```

<details>
<summary>Sample Output for Domain Availability Check</summary>

```go
[]any{
	map[string]any{
		"TLD": "com", "DomainName": "hello", "Status": "notavailable", "Command": "create",
		"Period": int64(1), "IsFee": false, "Price": "10.8100", "Currency": "USD", "Reason": "Domain exists",
	},
	map[string]any{
		"TLD": "com", "DomainName": "world123x0", "Status": "available", "Command": "create",
		"Period": int64(1), "IsFee": false, "Price": "10.8100", "Currency": "USD", "Reason": nil,
	},
}
```

`IsFee` is `true` for premium domains.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain Details

```go
result := client.GetDetails("example.com")
result = client.SyncFromRegistry("example.com") // same shape
```

<details>
<summary>Sample Output for Domain Details</summary>

```go
map[string]any{
	"result": "OK",
	"data": map[string]any{
		"ID":                      int64(564346),
		"Status":                  "Active",
		"DomainName":              "example.com",
		"AuthCode":                "DHQ!K52",
		"LockStatus":              "true",
		"PrivacyProtectionStatus": "false",
		"IsChildNameServer":       "true",
		"Contacts": map[string]any{
			"Billing":        map[string]any{"ID": "DNA-B-4"},
			"Technical":      map[string]any{"ID": "DNA-T-3"},
			"Administrative": map[string]any{"ID": "DNA-A-2"},
			"Registrant":     map[string]any{"ID": "DNA-R-1"},
		},
		"Dates": map[string]any{
			"Start":         "2024-03-15T10:20:30",
			"Expiration":    "2027-03-15T10:20:30",
			"RemainingDays": int64(126),
		},
		"NameServers":      []any{"ns1.example.com", "ns2.example.com"},
		"Additional":       map[string]any{"TRABISDOMAINCATEGORY": "1"},
		"ChildNameServers": []any{map[string]any{"ns": "ns1.example.com", "ip": "1.2.3.4"}},
	},
}
```

`Status` is the textual label (`Active`, `TransferredOut`, `PendingDelete`, ...) even when the
gateway sends the numeric enum. Unparseable dates become `""`, never a bogus 1970 date.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Nameserver Modification

```go
result := client.ModifyNameServer("example.com", []string{"ns1.example.com", "ns2.example.com"})
```

A map is accepted too (its values are used in key order). Entries are trimmed, lower-cased, stripped
of a trailing dot and de-duplicated before sending; the result reports that normalized list.

<details>
<summary>Sample Output for Nameserver Modification</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServers": []any{"ns1.example.com", "ns2.example.com"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Enable Domain Lock

```go
result := client.EnableTheftProtectionLock("example.com")
```

<details>
<summary>Sample Output for Enable Domain Lock</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Disable Domain Lock

```go
result := client.DisableTheftProtectionLock("example.com")
```

<details>
<summary>Sample Output for Disable Domain Lock</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": false}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Add Child Nameserver

```go
result := client.AddChildNameServer("example.com", "ns1.example.com", "1.2.3.4")
```

<details>
<summary>Sample Output for Add Child Nameserver</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.4"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Delete Child Nameserver

```go
result := client.DeleteChildNameServer("example.com", "ns1.example.com")
```

<details>
<summary>Sample Output for Delete Child Nameserver</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com"}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Update Child Nameserver

```go
result := client.ModifyChildNameServer("example.com", "ns1.example.com", "1.2.3.5")
```

<details>
<summary>Sample Output for Update Child Nameserver</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.5"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Modify Domain Privacy

```go
result := client.ModifyPrivacyProtectionStatus("example.com", true, "Owner request")
```

`status` is a bool or a value panels store: `true`, `1`, `"1"`, `"true"`, `"on"` and `"yes"` enable
privacy, anything else (`"off"`, `"0"`, `false`, ...) disables it. The reason is kept for php-dna
compatibility only and is not sent: the gateway has no such field.

<details>
<summary>Sample Output for Modify Domain Privacy</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"PrivacyProtectionStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Get Domain Contacts

```go
result := client.GetContacts("example.com")
```

Gateway contact types are matched case-insensitively (`Admin`/`Administrative`, `Tech`/`Technical`,
`Registrant`/`Owner`, `Billing`); any other type is kept under its own capitalized name.

<details>
<summary>Sample Output for Get Domain Contacts</summary>

```go
map[string]any{
	"result": "OK",
	"data": map[string]any{
		"contacts": map[string]any{
			"Administrative": map[string]any{
				"ID": "DNA-A-2", "Status": "Active", "AuthCode": "",
				"FirstName": "Grace", "LastName": "Hopper", "Company": "", "EMail": "grace@example.com",
				"Type": "Administrative",
				"Address": map[string]any{
					"Line1": "Example Street 2", "Line2": "", "Line3": "",
					"State": "DC", "City": "Washington", "Country": "US", "ZipCode": "20001",
				},
				"Phone": map[string]any{
					"Phone": map[string]any{"Number": "2025550100", "CountryCode": "1"},
					"Fax":   map[string]any{"Number": "2025550101", "CountryCode": "1"},
				},
				"Additional": []any{},
			},
			// "Billing", "Registrant", "Technical": same shape
		},
	},
}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Save Domain Contacts

```go
result := client.SaveContacts("example.com", dna.AllRoles(contact))
```

<details>
<summary>Sample Output for Save Domain Contacts</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"contacts": map[string]any{}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Get Reseller Details

```go
result := client.GetResellerDetails()
```

<details>
<summary>Sample Output for Get Reseller Details</summary>

```go
map[string]any{
	"result":   "OK",
	"id":       "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
	"active":   true,
	"name":     "Example Reseller",
	"balance":  "601.64",
	"currency": "USD",
	"symbol":   "$",
	"balances": []any{
		map[string]any{"balance": "601.64", "currency": "USD", "symbol": "$"},
		map[string]any{"balance": "1250.50", "currency": "TL", "symbol": "TL"},
	},
}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Get Current Balance

```go
result := client.GetCurrentBalance("USD") // USD (default), TRY, EUR or GBP
```

<details>
<summary>Sample Output for Get Current Balance</summary>

```go
map[string]any{
	"ErrorCode":        int64(0),
	"OperationMessage": "Command completed succesfully.",
	"OperationResult":  "SUCCESS",
	"Balance":          "601.64",
	"CurrencyId":       int64(2),
	"CurrencyInfo":     nil,
	"CurrencyName":     "USD",
	"CurrencySymbol":   "$",
}
```

On failure: `ErrorCode` is `API_<status>` (or `CURL_<n>`), `OperationMessage` the error message,
`OperationResult` `"FAILED"` and `Balance` `"0.00"`. The account only reports USD and TRY balances:
for any other currency the result is `FAILED` with `ErrorCode: "BALANCE"` and
`No EUR balance is reported for this account. Available: USD, TRY.`

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### DNS Zone Management

The gateway stores records as **sets** (one name + type with all its values) and every write
replaces the whole set. Record names may be given as `@`, `""`, the zone apex, a relative label
(`www`) or a full name (`www.example.com`); they are reduced to the relative label the gateway
expects. Supported types: `A`, `AAAA`, `ALIAS`, `CAA`, `CNAME`, `DS`, `MX`, `NAPTR`, `PTR`, `SRV`,
`TLSA`, `TXT` (others are rejected before any call). TTL is clamped to 1..86400 and the result reports
the upper-cased type and the TTL that were actually sent. With `apply = true` (php-dna's default)
changes are committed immediately; when that apply fails, its error envelope is returned (the change
stays staged). Consecutive zone calls are kept at least one second apart because the gateway
rate-limits bursts.

```go
// List: one row per value, names fully qualified
client.GetResourceRecords("example.com")

// Set a name+type to exactly these values (empty list removes the set). Defaults: ttl 3600, apply true
client.SetResourceRecordSet("example.com", "@", "MX", []string{"10 mx1.example.com", "20 mx2.example.com"}, 3600, true)

// Add / replace / remove ONE value, keeping the others (ttl nil = keep the set's TTL)
client.AddResourceRecordValue("example.com", "www", "A", "192.0.2.10", nil, true)
client.ReplaceResourceRecordValue("example.com", "www", "A", "192.0.2.10", "192.0.2.11", nil, true)
client.RemoveResourceRecordValue("example.com", "www", "A", "192.0.2.11", true)

// Remove a whole name+type (knownValue "" = look up a current value first)
client.DeleteResourceRecordSet("example.com", "www", "A", true, "")

// Low-level calls (content is a string or []string)
client.AddResourceRecord("example.com", "api", "CNAME", "example.com.", 300, false)
client.EditResourceRecord("example.com", "api", "api", "CNAME", "www.example.com.", 300, false)
client.DeleteResourceRecord("example.com", "old.example.com.", "192.0.2.99", "A", false)
client.ApplyResourceRecords("example.com")
```

<details>
<summary>Sample Outputs for DNS Zone Management</summary>

```go
// GetResourceRecords
map[string]any{"result": "OK", "data": map[string]any{"records": []any{
	map[string]any{"Name": "example.com.", "Type": "A", "TTL": int64(3600), "Content": "192.0.2.1", "Disabled": false},
	map[string]any{"Name": "www.example.com.", "Type": "CNAME", "TTL": int64(300), "Content": "example.com.", "Disabled": false},
}}}

// SetResourceRecordSet, AddResourceRecordValue, ReplaceResourceRecordValue, RemoveResourceRecordValue
map[string]any{"result": "OK", "data": map[string]any{
	"Name": "", "Type": "MX", "TTL": int64(3600), "Contents": []any{"10 mx1.example.com", "20 mx2.example.com"},
}}

// Value already present / not present: nothing is written
map[string]any{"result": "OK", "data": map[string]any{"Name": "", "Type": "MX", "Contents": []any{"..."}, "Changed": false}}

// DeleteResourceRecordSet of a set that does not exist
map[string]any{"result": "OK", "data": map[string]any{"Name": "nothing", "Type": "A", "Removed": false}}

// AddResourceRecord / EditResourceRecord
map[string]any{"result": "OK", "data": map[string]any{"Name": "api", "Type": "CNAME", "TTL": int64(300), "Content": "example.com."}}

// DeleteResourceRecord
map[string]any{"result": "OK", "data": map[string]any{"Name": "old.example.com.", "Type": "A", "Content": "192.0.2.99"}}

// ApplyResourceRecords
map[string]any{"result": "OK", "data": map[string]any{"DomainName": "example.com"}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### URL Forwarding

```go
client.SetForwarding("example.com", "https://www.example.org", "Standard") // or "Frame" (masked)
client.GetForwarding("example.com")
client.DeleteForwarding("example.com")
```

<details>
<summary>Sample Outputs for URL Forwarding</summary>

```go
// SetForwarding
map[string]any{"result": "OK", "data": map[string]any{"RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// GetForwarding
map[string]any{"result": "OK", "data": map[string]any{"Enabled": true, "RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// GetForwarding when no forwarding is configured
map[string]any{"result": "OK", "data": map[string]any{"Enabled": false, "RedirectAddress": "", "ForwardType": ""}}

// DeleteForwarding (Removed is false when nothing was configured)
map[string]any{"result": "OK", "data": map[string]any{"DomainName": "example.com", "Removed": true}}
```

</details>

Only the gateway's "Redirect information not found" answer means "not configured"; any other error,
for example an unknown domain, is returned as an error.

<hr style="border: 4px solid #000; border-style: dashed;">

### Utilities (no network)

```go
client.IsTrTLD("example.com.tr") // true

client.ValidateContact(map[string]any{"FirstName": " Jo ", "Phone": "+90 (532) 374 33 88", "Country": "Turkey"})
// map[string]any{"FirstName": "Jo", "LastName": "Jo", "AddressLine1": "Addres yok", "City": "ISTANBUL",
//   "State": "ISTANBUL", "Country": "TR", "ZipCode": "34000", "Phone": "5323743388",
//   "PhoneCountryCode": "90", "Fax": "", "FaxCountryCode": ""}
```

`ValidateContact` fills missing fields with placeholders, removes digits and punctuation from names,
falls back from State to City, maps common country names to ISO codes and splits phone and fax
numbers into number and country code. A number written with its own international prefix
(`+44 20 ...`, `0044 20 ...`) keeps its own country code even when another one is given.

<hr style="border: 4px solid #000; border-style: dashed;">

### Diagnostics

```go
client.GetRequestData()        // {"url": ..., "method": ..., "payload": ...} of the last call
client.GetResponseData()       // last decoded response, or {"raw_response": body}
client.GetResponseHeaders()    // {"url", "http_code", "content_type", "total_time", "headers"}
client.GetParsedResponseData() // last parsed response, or the {Code, Message, Details} error
client.GetLastFunction()       // e.g. "request:GET domains/info"
client.GetServiceURL()         // gateway base URL
client.SetServiceURL("https://ote.domainresellerapi.com/api/v1")
```

A `Client` keeps this per-call state, so do not share one `Client` between goroutines; create one
per goroutine.

## Error handling

- Gateway errors return `Code: "API_<HTTP status>"`. When the gateway sent its own error code
  (for example `Dna.DomainService:Domain:10007`), it is added as `GatewayCode`, so callers can tell
  "domain not found" from "TLD not supported". The message comes from the gateway; an empty error
  body becomes `HTTP 429 - Too Many Requests` (and so on), validation errors are listed field by
  field in `Details` (`Admin PostalCode - 'PostalCode' is required`), and backend boilerplate such as
  `Exception of type '...' was thrown` is replaced by the gateway's details, or by a stable message
  naming the error code — never by the raw response body.
- A 301/302 answer means the credentials were rejected: `Message: "Invalid API credentials"`.
- Transport failures (timeout, DNS, TLS, connection reset) are errors, never an empty success. A
  timeout returns `Code: "CURL_28"` and `Could not reach the DomainNameAPI service: the connection
  timed out. (...)`; other failures use curl-style numbers (`CURL_6` DNS, `CURL_7` connect,
  `CURL_35` TLS, `CURL_56` reset, `CURL_60` certificate) and `Curl error during request: ...`.
  `API_<n>` is reserved for HTTP statuses.
- Register, transfer and renew answers of HTTP 200 with `"success": false` are refusals:
  `REGISTER_DOMAIN`, `TRANSFER_DOMAIN` or `RENEW_DOMAIN` with the gateway's message.
- An empty or `null` success body returns `Code: "RESPONSE"`.
- Timeouts: 30 seconds per call, 10 seconds to connect.

## Test mode

Pass `true` (or `"true"`, `"1"`, `"on"`, `"yes"`) as the third constructor argument to use the OTE
test gateway `https://ote.domainresellerapi.com/api/v1`. Test-mode telemetry is tagged
`environment: ote`.

## Telemetry (error reporting)

When an API call fails, the library can send an error report to DomainNameAPI's own error tracker
(Sentry). A small random share of successful calls (4%, plus every call slower than 1 second) is sent
as a performance record. Reporting is **on by default**.

> ℹ️ The Sentry project for `DomainNameAPI-Go` is not created yet: `DefaultSentryDSN` is empty, so
> nothing is sent until a release sets it. `dna.WithSentryDSN(...)` overrides the DSN.

**What is sent:** error code and message; endpoint path and HTTP method; call duration; library name
(`DomainNameAPI-Go`), library version and Go version; reseller ID and domain name; the server's
public IP address (looked up via `api.ipify.org`); the operating-system user name; the API request
and response including contact details. Only credentials are masked: EPP/auth codes, API keys,
passwords and tokens become `***`.

**What is never sent:** EPP/auth codes, API keys, passwords or tokens in readable form; function call
arguments in stack traces.

**Why and how long:** the data is used only for usage measurement and early diagnosis of problems and
is kept for 30 days on DomainNameAPI's own Sentry server.

**Turn it off, or send the minimum:**

```go
client := dna.New(resellerID, apiKey, false)
client.SetErrorReportingEnabled(false) // nothing is sent

client.SetTelemetryMinimal(true) // only error code/message, endpoint, timing, versions and reseller ID
```

Telemetry runs in the background, never blocks the call for long and never panics into your code.

> 🔒 **Security:** TLS certificate verification is on. On a host with an outdated CA bundle it can be
> switched off with `client.SetSSLVerification(false)` — not recommended; updating the CA bundle is
> the better fix.

## Advanced options

```go
client := dna.New(resellerID, apiKey, false,
	dna.WithHTTPClient(myHTTPClient),      // any Do(*http.Request) (*http.Response, error)
	dna.WithSentryDSN("https://key@host/1"),
	dna.WithTelemetryClient(myHTTPClient),
)
```

With a custom HTTP client, configure timeouts, proxies and TLS on it; it must not follow redirects
(return `http.ErrUseLastResponse` from `CheckRedirect`).

## Testing

Go is not required on the host; Docker is enough:

```bash
docker run --rm -v "$PWD":/src -w /src golang:1.23 go test ./...
docker run --rm -v "$PWD":/src -w /src golang:1.23 go vet ./...
```

The suite has unit tests for the HTTP layer, error mapping, telemetry and utilities, and a
conformance test that replays all 134 php-dna fixtures from `testdata/conformance` (requests on the
wire and return values must match php-dna exactly). No test touches the network.

## Return and Error Codes with Descriptions

| Code | Description                                    | Detail                                                                                                                                                         |
|------|------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 1000 | Command completed successfully                 | Operation successful.                                                                                                                                          |
| 1001 | Command completed successfully; action pending | Operation successful. However, the operation is currently queued for completion.                                                                               |
| 2003 | Required parameter missing                     | Parameter missing error. For example, no phone entry in contact information.                                                                                   |
| 2105 | Object is not eligible for renewal             | Domain status is not suitable for renewal, locked for update operations. Status should not be "clientupdateprohibited". May be due to other status conditions. |
| 2200 | Authentication error                           | Authorization error, security code is incorrect or domain is registered with another registrar.                                                                |
| 2302 | Object exists                                  | Domain name or nameserver information already exists in database. Cannot be registered.                                                                        |
| 2303 | Object does not exist                          | Domain name or nameserver information does not exist in database. New registration required.                                                                   |
| 2304 | Object status prohibits operation              | Domain status is not suitable for updates, locked for update operations. Status should not be "clientupdateprohibited". May be due to other status conditions. |

Library codes: `RESPONSE` (empty response), `CREDENTIALS` (invalid response format from the account
endpoint), `AVAILABILITY` (no domain names given), `BALANCE` (currency not reported by the account),
`TRANSFER_NOT_AVAILABLE`, `REGISTER_DOMAIN` / `TRANSFER_DOMAIN` / `RENEW_DOMAIN` (the gateway
refused the operation), `DOMAIN_RENEW` (renewal without expiration date), `ADD_ZONE`, `EDIT_ZONE`,
`SET_ZONE`, `DELETE_ZONE` (unsupported record type or value not found), `SET_FORWARD` (unsupported
forward type or empty address), `CURL_<n>` (transport failure), `API_<status>` (HTTP error, with
`GatewayCode` when the gateway sent its own code).

## License

Apache License 2.0 — see [LICENSE](LICENSE).
