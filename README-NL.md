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

# go-dna — DomainNameAPI-client voor Go

Go-client voor de REST API van domeinregistrar DomainNameAPI. Het is een getrouwe port van
[php-dna](https://github.com/domainreseller/php-dna): dezelfde methoden, dezelfde verzoeken over de
lijn en dezelfde retourwaarden. De standaardbibliotheek is de enige afhankelijkheid.

## 📦 Downloaden — gebruik altijd de Releases!

⬇️ **Download hier de nieuwste geteste versie: https://github.com/domainreseller/go-dna/releases/latest**

> ⚠️ Gebruik **niet** de groene knop **Code → Download ZIP** — daarmee downloadt u de ruwe ontwikkelbranch. Releasepakketten zijn voorzien van een versienummer, getest en klaar voor productie.

## Vereisten

- Go 1.21 of hoger
- Geen pakketten van derden (alleen `net/http` en `encoding/json` uit de standaardbibliotheek)

## Installatie

Leg met Go modules een uitgebrachte versie vast:

```bash
go get github.com/domainreseller/go-dna@v1.0.0
```

```go
import dna "github.com/domainreseller/go-dna"
```

De pakketnaam is `dna`. Uitvoerbare voorbeelden voor elke methode staan in de map [examples](examples):

```bash
go run ./examples/GetDetails
```

## 🔑 API-inloggegevens — Reseller ID en API Key

De bibliotheek heeft precies twee inloggegevens nodig, die aan de constructor worden meegegeven:

| Parameter | Waarde |
|---|---|
| Eerste | Reseller ID — UUID zoals `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` |
| Tweede | API Key |

> 💡 U vindt uw **Reseller ID** en **API Key** in uw DomainNameAPI-paneel onder **API-instellingen**.
> ⚠️ Dit zijn **API-inloggegevens** — de inloggegevens van uw paneel werken hier **niet**.

```go
resellerID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
apiKey := "your-api-key"

// Live-omgeving
client := dna.New(resellerID, apiKey, false)

// Test-omgeving (OTE)
client = dna.New(resellerID, apiKey, true)
```

Van beide inloggegevens worden spaties aan begin en eind verwijderd. Het derde argument accepteert
ook de strings die panelen voor selectievakjes opslaan: `"true"`, `"1"`, `"on"` en `"yes"`
(hoofdletterongevoelig) kiezen de testomgeving, al het andere (`"false"`, `"0"`, `"off"`, `""`, ...)
kiest de live-omgeving.

## Snel aan de slag

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

## Retourwaarden

Methoden retourneren dezelfde dynamische structuren als php-dna — `map[string]any` en `[]any` — en
geen getypeerde modellen, zodat de sleutels, string-vlaggen (`"true"`/`"false"`) en opgemaakte
geldbedragen (`"11.9900"`) exact overeenkomen met die van php-dna. Getallen van de gateway zijn
`int64` of `float64`.

Geslaagde aanroepen zien er zo uit:

```go
map[string]any{"result": "OK", "data": map[string]any{ /* ... */ }}
```

Mislukte aanroepen veroorzaken nooit een panic; ze retourneren de foutenvelop:

```go
map[string]any{
	"result": "ERROR",
	"error": map[string]any{
		"Code":        "API_404", // API_<HTTP-status>, CURL_<n>, of een bibliotheekcode zoals RESPONSE
		"Message":     "The domain name you requested could not be found. Please check and try again.",
		"Details":     "The domain could not be found. Please check and try again.",
		"GatewayCode": "Dna.DomainService:Domain:10007", // alleen bij API_<status>, wanneer de gateway een eigen code stuurde
	},
}
```

Drie methoden gebruiken een andere vorm, precies zoals in php-dna:

| Methode | Succes | Fout |
|---|---|---|
| `CheckAvailability` | een kale `[]any`-lijst | de foutenvelop (`map[string]any`), dus het retourtype is `any` |
| `GetCurrentBalance` | envelop met `OperationResult: "SUCCESS"` | dezelfde envelop met `OperationResult: "FAILED"` |
| `GetResellerDetails` | platte map (`id`, `name`, `balance`, `balances`, ...) | de foutenvelop |

### Optionele argumenten

Go kent geen standaardargumenten. De methoden behouden de parametervolgorde van php-dna en nemen
elk argument expliciet aan; de standaardwaarde van php-dna staat in elke sectie hieronder vermeld.
Voor het gemak selecteren enkele nulwaarden de standaardwaarde: `GetCurrentBalance("")` betekent USD,
`GetTldList(0)` betekent 20, `CheckAvailability(..., "")` betekent `create` en
`SetForwarding(..., "")` betekent `Standard`.
`RegisterWithContactInfo` neemt zijn vier optionele argumenten aan in een `*RegisterOptions`
(nil = standaardwaarden). `AddResourceRecordValue` en `ReplaceResourceRecordValue` nemen `ttl *int`
aan: `nil` behoudt de huidige TTL van de set.

### Contacten

Contacten zijn `map[string]any` met de sleutels van php-dna (`FirstName`, `LastName`, `Company`,
`EMail`, `AddressLine1`, `AddressLine2`, `City`, `State`, `Country`, `ZipCode`, `Phone`,
`PhoneCountryCode`, `Fax`, `FaxCountryCode`). De geneste vorm (`Address.Line1`,
`Phone.Phone.Number`, ...) wordt ook geaccepteerd. Omdat Go-maps geen volgorde hebben, worden de
vier rollen doorgegeven als een geordende `dna.Contacts`-lijst:

```go
contacts := dna.Contacts{
	{Role: "Administrative", Contact: contact},
	{Role: "Billing", Contact: contact},
	{Role: "Technical", Contact: contact},
	{Role: "Registrant", Contact: contact},
}
contacts = dna.AllRoles(contact)                 // hetzelfde contact voor alle vier rollen
contacts = dna.ContactsFromMap(map[string]map[string]any{ /* rol => contact */ })
```

Half ingevulde contacten worden vóór verzending aangevuld volgens de regels van `ValidateContact`.
Contacten worden gebruikt door `RegisterWithContactInfo` en `SaveContacts`; bij transfers worden nooit
contacten verzonden.

### Methodenamen

Go-methodenamen volgen de Go-conventies; de php-dna-naam is dezelfde in lowerCamelCase.
`getServiceUrl`/`setServiceUrl` heten `GetServiceURL`/`SetServiceURL` en `setSslVerification` heet
`SetSSLVerification`.

## Gebruik

### Domeinregistratie

Standaardwaarden van php-dna: standaard nameservers (`ns1.domainnameapi.com`,
`ns2.domainnameapi.com`), EPP-vergrendeling aan, privacy uit, geen aanvullende attributen.
Aanvullende attributen zijn alleen vereist voor domeinen zoals `.tr`.

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
// Alleen voor .tr-domeinen:
opts.AdditionalAttributes = map[string]any{
	"TRABISDOMAINCATEGORY": 1,             // 1 = particulier (Owner), 0 = bedrijf (Corporate)
	"TRABISCITIZIENID":     "12345678901", // of TRABISTAXOFFICE + TRABISTAXNUMBER voor bedrijven
}

result := client.RegisterWithContactInfo("example.com", 1, dna.AllRoles(contact), opts)
```

Voor `.tr`-domeinen worden de verouderde TRABIS-namen genormaliseerd naar het gatewayschema: de
categorie wordt `Owner`/`Corporate`, `TRABISCITIZIENID` wordt `TRABISCITIZENID`, en de velden voor
land, stad, naam en organisatie vervallen (de gateway haalt ze uit het contact). Elke attribuutwaarde
wordt als string verzonden.

<details>
<summary>Voorbeelduitvoer voor Domeinregistratie</summary>

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

Het registratie-endpoint antwoordt met een beknopt record, waardoor de meeste velden leeg zijn; roep
daarna `GetDetails` aan voor het volledige domein. Wanneer de gateway de registratie weigert (HTTP 200
met `"success": false`), is het resultaat een fout met code `REGISTER_DOMAIN` en de melding van de gateway.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domeinverlenging

```go
result := client.Renew("example.com", 1)
```

<details>
<summary>Voorbeelduitvoer voor Domeinverlenging</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"ExpirationDate": "2026-03-15T10:00:00"}}
```

Een weigering door de gateway (`"success": false`) retourneert code `RENEW_DOMAIN`; een succes zonder
vervaldatum retourneert `DOMAIN_RENEW`.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domeintransfer

```go
result := client.Transfer("example.com", "Xy9#mK2$", 1)
```

Bij transfers worden nooit contacten verzonden: de registrar neemt ze samen met het domein over,
dus het verzoek bevat altijd een lege contactenlijst en de methode heeft geen contactenparameter. Het
resultaat heeft dezelfde vorm als `GetDetails`; een weigering door de gateway (`"success": false`)
retourneert code `TRANSFER_DOMAIN` met de melding van de gateway.

<details>
<summary>Voorbeelduitvoer voor Domeintransfer</summary>

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

Gerelateerde transferaanroepen:

```go
client.CancelTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Cancelled"}}
client.ApproveTransfer("example.com") // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Approved"}}
client.RejectTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Rejected"}}
```

<hr style="border: 4px solid #000; border-style: dashed;">

### Transfercontrole

```go
result := client.CheckTransfer("example.com", "Ab1!xyZ9")
```

`result` is alleen `OK` wanneer de registry meldt dat het domein overdraagbaar is; anders is het
`ERROR` met code `TRANSFER_NOT_AVAILABLE` en dezelfde `data`.

<details>
<summary>Voorbeelduitvoer voor Transfercontrole</summary>

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

### Domeinlijst

```go
result := client.GetList(map[string]any{
	"OrderColumn":    "Id",
	"OrderDirection": "ASC",
	"PageNumber":     0,
	"PageSize":       1000,
})
```

`PageNumber` (vanaf 0), `PageSize`, `OrderColumn` en `OrderDirection` worden vertaald naar
`SkipCount`, `MaxResultCount` en `Sorting` van de gateway; andere sleutels (bijvoorbeeld
`DomainName`) worden als filters doorgegeven. `nil` retourneert de eerste 200 domeinen.

<details>
<summary>Voorbeelduitvoer voor Domeinlijst</summary>

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

### TLD-lijst

```go
result := client.GetTldList(100) // standaardwaarde php-dna: 20
```

<details>
<summary>Voorbeelduitvoer voor TLD-lijst</summary>

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

### Domeinbeschikbaarheid Controleren

```go
result := client.CheckAvailability([]string{"hello", "world123x0"}, []string{"com", "net"}, 1, "create")
if list, ok := result.([]any); ok {
	// geslaagd
	_ = list
}
```

<details>
<summary>Voorbeelduitvoer voor Domeinbeschikbaarheid Controleren</summary>

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

`IsFee` is `true` voor premium domeinen.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domeininformatie

```go
result := client.GetDetails("example.com")
result = client.SyncFromRegistry("example.com") // zelfde vorm
```

<details>
<summary>Voorbeelduitvoer voor Domeininformatie</summary>

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

`Status` is het tekstuele label (`Active`, `TransferredOut`, `PendingDelete`, ...), ook wanneer de
gateway de numerieke enum stuurt. Onleesbare datums worden `""`, nooit een onjuiste datum uit 1970.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Nameserver Wijzigen

```go
result := client.ModifyNameServer("example.com", []string{"ns1.example.com", "ns2.example.com"})
```

Een map wordt ook geaccepteerd (de waarden worden in sleutelvolgorde gebruikt). Vóór verzending
worden de items ontdaan van spaties, naar kleine letters omgezet, van een afsluitende punt ontdaan
en ontdubbeld; het resultaat bevat die genormaliseerde lijst.

<details>
<summary>Voorbeelduitvoer voor Nameserver Wijzigen</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServers": []any{"ns1.example.com", "ns2.example.com"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domeinvergrendeling Inschakelen

```go
result := client.EnableTheftProtectionLock("example.com")
```

<details>
<summary>Voorbeelduitvoer voor Domeinvergrendeling Inschakelen</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domeinvergrendeling Uitschakelen

```go
result := client.DisableTheftProtectionLock("example.com")
```

<details>
<summary>Voorbeelduitvoer voor Domeinvergrendeling Uitschakelen</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": false}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Child Nameserver Toevoegen

```go
result := client.AddChildNameServer("example.com", "ns1.example.com", "1.2.3.4")
```

<details>
<summary>Voorbeelduitvoer voor Child Nameserver Toevoegen</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.4"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Child Nameserver Verwijderen

```go
result := client.DeleteChildNameServer("example.com", "ns1.example.com")
```

<details>
<summary>Voorbeelduitvoer voor Child Nameserver Verwijderen</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com"}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Child Nameserver Bijwerken

```go
result := client.ModifyChildNameServer("example.com", "ns1.example.com", "1.2.3.5")
```

<details>
<summary>Voorbeelduitvoer voor Child Nameserver Bijwerken</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.5"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domeinprivacy Wijzigen

```go
result := client.ModifyPrivacyProtectionStatus("example.com", true, "Owner request")
```

`status` is een bool of een waarde die panelen opslaan: `true`, `1`, `"1"`, `"true"`, `"on"` en
`"yes"` schakelen privacy in, al het andere (`"off"`, `"0"`, `false`, ...) schakelt het uit. De reden
wordt alleen behouden voor compatibiliteit met php-dna en wordt niet verzonden: de gateway kent zo'n
veld niet.

<details>
<summary>Voorbeelduitvoer voor Domeinprivacy Wijzigen</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"PrivacyProtectionStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domeincontacten Ophalen

```go
result := client.GetContacts("example.com")
```

Contacttypen van de gateway worden hoofdletterongevoelig herkend (`Admin`/`Administrative`,
`Tech`/`Technical`, `Registrant`/`Owner`, `Billing`); elk ander type wordt bewaard onder zijn eigen
naam met een hoofdletter.

<details>
<summary>Voorbeelduitvoer voor Domeincontacten Ophalen</summary>

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
			// "Billing", "Registrant", "Technical": zelfde vorm
		},
	},
}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domeincontacten Opslaan

```go
result := client.SaveContacts("example.com", dna.AllRoles(contact))
```

<details>
<summary>Voorbeelduitvoer voor Domeincontacten Opslaan</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"contacts": map[string]any{}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Resellergegevens Ophalen

```go
result := client.GetResellerDetails()
```

<details>
<summary>Voorbeelduitvoer voor Resellergegevens Ophalen</summary>

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

### Huidig Saldo Ophalen

```go
result := client.GetCurrentBalance("USD") // USD (standaard), TRY, EUR of GBP
```

<details>
<summary>Voorbeelduitvoer voor Huidig Saldo Ophalen</summary>

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

Bij een fout: `ErrorCode` is `API_<status>` (of `CURL_<n>`), `OperationMessage` de foutmelding,
`OperationResult` `"FAILED"` en `Balance` `"0.00"`. Het account rapporteert alleen saldi in USD en TRY:
voor elke andere valuta is het resultaat `FAILED` met `ErrorCode: "BALANCE"` en
`No EUR balance is reported for this account. Available: USD, TRY.`

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### DNS-zonebeheer

De gateway slaat records op als **sets** (één naam + type met al zijn waarden) en elke schrijfactie
vervangt de hele set. Recordnamen kunnen worden opgegeven als `@`, `""`, de zone-apex, een relatief
label (`www`) of een volledige naam (`www.example.com`); ze worden teruggebracht tot het relatieve
label dat de gateway verwacht. Ondersteunde typen: `A`, `AAAA`, `ALIAS`, `CAA`, `CNAME`, `DS`, `MX`,
`NAPTR`, `PTR`, `SRV`, `TLSA`, `TXT` (andere worden geweigerd vóórdat er een aanroep plaatsvindt).
De TTL wordt begrensd tot 1..86400 en het resultaat vermeldt het type in hoofdletters en de TTL die
daadwerkelijk zijn verzonden. Met `apply = true` (de standaardwaarde van php-dna) worden wijzigingen
direct doorgevoerd; als die apply mislukt, wordt de foutenvelop daarvan geretourneerd (de wijziging
blijft klaargezet). Opeenvolgende zone-aanroepen liggen minstens één seconde uit elkaar, omdat de
gateway snelle reeksen verzoeken aan een rate limit onderwerpt.

```go
// Lijst: één rij per waarde, namen volledig gekwalificeerd
client.GetResourceRecords("example.com")

// Stel een naam+type in op precies deze waarden (lege lijst verwijdert de set). Standaard: ttl 3600, apply true
client.SetResourceRecordSet("example.com", "@", "MX", []string{"10 mx1.example.com", "20 mx2.example.com"}, 3600, true)

// Voeg ÉÉN waarde toe / vervang of verwijder deze en behoud de andere (ttl nil = TTL van de set behouden)
client.AddResourceRecordValue("example.com", "www", "A", "192.0.2.10", nil, true)
client.ReplaceResourceRecordValue("example.com", "www", "A", "192.0.2.10", "192.0.2.11", nil, true)
client.RemoveResourceRecordValue("example.com", "www", "A", "192.0.2.11", true)

// Verwijder een volledige naam+type (knownValue "" = eerst een huidige waarde opzoeken)
client.DeleteResourceRecordSet("example.com", "www", "A", true, "")

// Low-level aanroepen (content is een string of []string)
client.AddResourceRecord("example.com", "api", "CNAME", "example.com.", 300, false)
client.EditResourceRecord("example.com", "api", "api", "CNAME", "www.example.com.", 300, false)
client.DeleteResourceRecord("example.com", "old.example.com.", "192.0.2.99", "A", false)
client.ApplyResourceRecords("example.com")
```

<details>
<summary>Voorbeelduitvoer voor DNS-zonebeheer</summary>

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

// Waarde al aanwezig / niet aanwezig: er wordt niets geschreven
map[string]any{"result": "OK", "data": map[string]any{"Name": "", "Type": "MX", "Contents": []any{"..."}, "Changed": false}}

// DeleteResourceRecordSet van een set die niet bestaat
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

### URL-doorsturing

```go
client.SetForwarding("example.com", "https://www.example.org", "Standard") // of "Frame" (gemaskeerd)
client.GetForwarding("example.com")
client.DeleteForwarding("example.com")
```

<details>
<summary>Voorbeelduitvoer voor URL-doorsturing</summary>

```go
// SetForwarding
map[string]any{"result": "OK", "data": map[string]any{"RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// GetForwarding
map[string]any{"result": "OK", "data": map[string]any{"Enabled": true, "RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// GetForwarding wanneer er geen doorsturing is ingesteld
map[string]any{"result": "OK", "data": map[string]any{"Enabled": false, "RedirectAddress": "", "ForwardType": ""}}

// DeleteForwarding (Removed is false wanneer er niets was ingesteld)
map[string]any{"result": "OK", "data": map[string]any{"DomainName": "example.com", "Removed": true}}
```

</details>

Alleen het antwoord "Redirect information not found" van de gateway betekent "niet ingesteld"; elke
andere fout, bijvoorbeeld een onbekend domein, wordt als fout geretourneerd.

<hr style="border: 4px solid #000; border-style: dashed;">

### Hulpfuncties (zonder netwerk)

```go
client.IsTrTLD("example.com.tr") // true

client.ValidateContact(map[string]any{"FirstName": " Jo ", "Phone": "+90 (532) 374 33 88", "Country": "Turkey"})
// map[string]any{"FirstName": "Jo", "LastName": "Jo", "AddressLine1": "Addres yok", "City": "ISTANBUL",
//   "State": "ISTANBUL", "Country": "TR", "ZipCode": "34000", "Phone": "5323743388",
//   "PhoneCountryCode": "90", "Fax": "", "FaxCountryCode": ""}
```

`ValidateContact` vult ontbrekende velden aan met plaatsaanduidingen, verwijdert cijfers en
leestekens uit namen, valt terug van State op City, zet gangbare landnamen om naar ISO-codes en
splitst telefoon- en faxnummers op in nummer en landcode. Een nummer dat met een eigen
internationaal voorvoegsel is geschreven (`+44 20 ...`, `0044 20 ...`) behoudt zijn eigen landcode,
ook wanneer er een andere is opgegeven.

<hr style="border: 4px solid #000; border-style: dashed;">

### Diagnostiek

```go
client.GetRequestData()        // {"url": ..., "method": ..., "payload": ...} van de laatste aanroep
client.GetResponseData()       // laatste gedecodeerde respons, of {"raw_response": body}
client.GetResponseHeaders()    // {"url", "http_code", "content_type", "total_time", "headers"}
client.GetParsedResponseData() // laatste geparste respons, of de fout {Code, Message, Details}
client.GetLastFunction()       // bijv. "request:GET domains/info"
client.GetServiceURL()         // basis-URL van de gateway
client.SetServiceURL("https://ote.domainresellerapi.com/api/v1")
```

Een `Client` houdt deze status per aanroep bij, dus deel één `Client` niet tussen goroutines; maak er
één per goroutine aan.

## Foutafhandeling

- Gatewayfouten retourneren `Code: "API_<HTTP status>"`. Wanneer de gateway een eigen foutcode stuurde
  (bijvoorbeeld `Dna.DomainService:Domain:10007`), wordt die toegevoegd als `GatewayCode`, zodat
  aanroepers "domein niet gevonden" kunnen onderscheiden van "TLD niet ondersteund". De melding komt
  van de gateway; een lege foutbody wordt `HTTP 429 - Too Many Requests` (enzovoort),
  validatiefouten worden per veld vermeld in `Details` (`Admin PostalCode - 'PostalCode' is required`),
  en standaardteksten van de backend zoals `Exception of type '...' was thrown` worden vervangen door
  de details van de gateway, of door een vaste melding met de foutcode — nooit door de ruwe responsbody.
- Een 301/302-antwoord betekent dat de inloggegevens zijn geweigerd: `Message: "Invalid API credentials"`.
- Transportfouten (time-out, DNS, TLS, verbinding verbroken) zijn fouten, nooit een leeg succes. Een
  time-out retourneert `Code: "CURL_28"` en `Could not reach the DomainNameAPI service: the connection
  timed out. (...)`; andere fouten gebruiken nummers in curl-stijl (`CURL_6` DNS, `CURL_7` verbinden,
  `CURL_35` TLS, `CURL_56` verbinding verbroken, `CURL_60` certificaat) en `Curl error during request: ...`.
  `API_<n>` is gereserveerd voor HTTP-statussen.
- Antwoorden op registratie, transfer en verlenging met HTTP 200 en `"success": false` zijn
  weigeringen: `REGISTER_DOMAIN`, `TRANSFER_DOMAIN` of `RENEW_DOMAIN` met de melding van de gateway.
- Een lege of `null` succesbody retourneert `Code: "RESPONSE"`.
- Time-outs: 30 seconden per aanroep, 10 seconden voor het opzetten van de verbinding.

## Testmodus

Geef `true` (of `"true"`, `"1"`, `"on"`, `"yes"`) mee als derde constructorargument om de
OTE-testgateway `https://ote.domainresellerapi.com/api/v1` te gebruiken. Telemetrie in testmodus wordt
gelabeld met `environment: ote`.

## Telemetrie (foutrapportage)

Wanneer een API-aanroep mislukt, kan de bibliotheek een foutrapport sturen naar de eigen
fouttracker van DomainNameAPI (Sentry). Een klein willekeurig deel van de geslaagde aanroepen (4%,
plus elke aanroep die langer duurt dan 1 seconde) wordt als prestatierecord verzonden. Rapportage
staat **standaard aan**.

> ℹ️ Het Sentry-project voor `DomainNameAPI-Go` is nog niet aangemaakt: `DefaultSentryDSN` is leeg,
> dus er wordt niets verzonden totdat een release deze instelt. `dna.WithSentryDSN(...)` overschrijft de DSN.

**Wat wordt verzonden:** foutcode en -melding; endpointpad en HTTP-methode; duur van de aanroep;
bibliotheeknaam (`DomainNameAPI-Go`), bibliotheekversie en Go-versie; reseller ID en domeinnaam; het
publieke IP-adres van de server (opgevraagd via `api.ipify.org`); de gebruikersnaam van het
besturingssysteem; het API-verzoek en de API-respons, inclusief contactgegevens. Alleen
inloggegevens worden gemaskeerd: EPP-/autorisatiecodes, API-sleutels, wachtwoorden en tokens worden `***`.

**Wat wordt nooit verzonden:** EPP-/autorisatiecodes, API-sleutels, wachtwoorden of tokens in
leesbare vorm; argumenten van functieaanroepen in stacktraces.

**Waarom en hoe lang:** de gegevens worden uitsluitend gebruikt voor gebruiksmeting en vroege
diagnose van problemen en worden 30 dagen bewaard op de eigen Sentry-server van DomainNameAPI.

**Uitschakelen, of alleen het minimum verzenden:**

```go
client := dna.New(resellerID, apiKey, false)
client.SetErrorReportingEnabled(false) // er wordt niets verzonden

client.SetTelemetryMinimal(true) // alleen foutcode/-melding, endpoint, timing, versies en reseller ID
```

Telemetrie draait op de achtergrond, blokkeert de aanroep nooit lang en veroorzaakt nooit een panic
in uw code.

> 🔒 **Beveiliging:** verificatie van TLS-certificaten staat aan. Op een host met een verouderde
> CA-bundel kan deze worden uitgeschakeld met `client.SetSSLVerification(false)` — niet aanbevolen;
> de CA-bundel bijwerken is de betere oplossing.

## Geavanceerde opties

```go
client := dna.New(resellerID, apiKey, false,
	dna.WithHTTPClient(myHTTPClient),      // elke Do(*http.Request) (*http.Response, error)
	dna.WithSentryDSN("https://key@host/1"),
	dna.WithTelemetryClient(myHTTPClient),
)
```

Configureer bij een eigen HTTP-client de time-outs, proxy's en TLS op die client; deze mag geen
redirects volgen (retourneer `http.ErrUseLastResponse` vanuit `CheckRedirect`).

## Testen

Go hoeft niet op de host geïnstalleerd te zijn; Docker is voldoende:

```bash
docker run --rm -v "$PWD":/src -w /src golang:1.23 go test ./...
docker run --rm -v "$PWD":/src -w /src golang:1.23 go vet ./...
```

De testsuite bevat unittests voor de HTTP-laag, foutmapping, telemetrie en hulpfuncties, en een
conformiteitstest die alle 134 php-dna-fixtures uit `testdata/conformance` opnieuw afspeelt
(verzoeken over de lijn en retourwaarden moeten exact overeenkomen met php-dna). Geen enkele test
maakt gebruik van het netwerk.

## Return- en Foutcodes met Beschrijvingen

| Code | Beschrijving                                   | Detail                                                                                                                                                         |
|------|------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 1000 | Commando succesvol voltooid                    | Operatie succesvol.                                                                                                                                            |
| 1001 | Commando succesvol voltooid; actie in behandeling | Operatie succesvol. De operatie staat momenteel echter in de wachtrij voor voltooiing.                                                                      |
| 2003 | Vereiste parameter ontbreekt                   | Fout door ontbrekende parameter. Bijvoorbeeld: geen telefoonnummer in de contactgegevens.                                                                      |
| 2105 | Object komt niet in aanmerking voor verlenging | Domeinstatus is niet geschikt voor verlenging, vergrendeld voor updatebewerkingen. Status mag niet "clientupdateprohibited" zijn. Kan ook door andere statusvoorwaarden komen. |
| 2200 | Authenticatiefout                              | Autorisatiefout, beveiligingscode is onjuist of het domein is geregistreerd bij een andere registrar.                                                          |
| 2302 | Object bestaat al                              | Domeinnaam of nameservergegevens bestaan al in de database. Kan niet worden geregistreerd.                                                                     |
| 2303 | Object bestaat niet                            | Domeinnaam of nameservergegevens bestaan niet in de database. Nieuwe registratie vereist.                                                                      |
| 2304 | Objectstatus verbiedt bewerking                | Domeinstatus is niet geschikt voor updates, vergrendeld voor updatebewerkingen. Status mag niet "clientupdateprohibited" zijn. Kan ook door andere statusvoorwaarden komen. |

Bibliotheekcodes: `RESPONSE` (lege respons), `CREDENTIALS` (ongeldig responsformaat van het
account-endpoint), `AVAILABILITY` (geen domeinnamen opgegeven), `BALANCE` (valuta niet door het
account gerapporteerd), `TRANSFER_NOT_AVAILABLE`, `REGISTER_DOMAIN` / `TRANSFER_DOMAIN` /
`RENEW_DOMAIN` (de gateway heeft de bewerking geweigerd), `DOMAIN_RENEW` (verlenging zonder
vervaldatum), `ADD_ZONE`, `EDIT_ZONE`, `SET_ZONE`, `DELETE_ZONE` (niet-ondersteund recordtype of
waarde niet gevonden), `SET_FORWARD` (niet-ondersteund doorstuurtype of leeg adres), `CURL_<n>`
(transportfout), `API_<status>` (HTTP-fout, met `GatewayCode` wanneer de gateway een eigen code stuurde).

## Licentie

Apache License 2.0 — zie [LICENSE](LICENSE).
