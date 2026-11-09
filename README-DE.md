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

# go-dna — DomainNameAPI-Client für Go

Go-Client für die REST-API des Domain-Registrars DomainNameAPI. Er ist eine originalgetreue Portierung
von [php-dna](https://github.com/domainreseller/php-dna): dieselben Methoden, dieselben Anfragen auf
der Leitung und dieselben Rückgabewerte. Die einzige Abhängigkeit ist die Standardbibliothek.

## 📦 Download — immer die Releases verwenden!

⬇️ **Die neueste getestete Version erhalten Sie hier: https://github.com/domainreseller/go-dna/releases/latest**

> ⚠️ Verwenden Sie **nicht** den grünen Button **Code → Download ZIP** — damit laden Sie den rohen Entwicklungszweig herunter. Release-Pakete sind versioniert, getestet und produktionsreif.

## Voraussetzungen

- Go 1.21 oder höher
- Keine Pakete von Drittanbietern (nur `net/http` und `encoding/json` aus der Standardbibliothek)

## Installation

Pinnen Sie mit Go-Modulen eine veröffentlichte Version:

```bash
go get github.com/domainreseller/go-dna@v1.0.0
```

```go
import dna "github.com/domainreseller/go-dna"
```

Der Paketname lautet `dna`. Ausführbare Beispiele für jede Methode finden Sie im Ordner [examples](examples):

```bash
go run ./examples/GetDetails
```

## 🔑 API-Zugangsdaten — Reseller-ID & API-Key

Die Bibliothek benötigt genau zwei Zugangsdaten, die dem Konstruktor übergeben werden:

| Parameter | Wert |
|---|---|
| Erster | Reseller-ID — UUID wie `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` |
| Zweiter | API-Key |

> 💡 Ihre **Reseller-ID** und Ihren **API-Key** finden Sie in Ihrem DomainNameAPI-Panel unter **API-Einstellungen**.
> ⚠️ Dies sind **API-Zugangsdaten** — Ihre Panel-Anmeldedaten funktionieren hier **nicht**.

```go
resellerID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
apiKey := "your-api-key"

// Live-Umgebung
client := dna.New(resellerID, apiKey, false)

// Test-Umgebung (OTE)
client = dna.New(resellerID, apiKey, true)
```

Bei beiden Zugangsdaten werden führende und nachfolgende Leerzeichen entfernt. Das dritte Argument
akzeptiert auch die Zeichenketten, die Panels für Checkboxen speichern: `"true"`, `"1"`, `"on"` und
`"yes"` (unabhängig von Groß-/Kleinschreibung) wählen die Test-Umgebung, alles andere (`"false"`,
`"0"`, `"off"`, `""`, ...) wählt die Live-Umgebung.

## Schnellstart

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

## Rückgabewerte

Die Methoden liefern dieselben dynamischen Strukturen wie php-dna — `map[string]any` und `[]any` —
und keine typisierten Modelle. Schlüssel, String-Flags (`"true"`/`"false"`) und formatierte
Betragsstrings (`"11.9900"`) entsprechen daher exakt denen von php-dna. Zahlen vom Gateway sind
`int64` oder `float64`.

Erfolgreiche Aufrufe sehen so aus:

```go
map[string]any{"result": "OK", "data": map[string]any{ /* ... */ }}
```

Fehlgeschlagene Aufrufe lösen nie eine Panic aus; sie liefern den Fehler-Umschlag (Error Envelope):

```go
map[string]any{
	"result": "ERROR",
	"error": map[string]any{
		"Code":    "API_404",      // API_<HTTP-Status> oder ein Bibliothekscode wie RESPONSE
		"Message": "The domain name you requested could not be found. Please check and try again.",
		"Details": "The domain could not be found. Please check and try again.",
	},
}
```

Drei Methoden verwenden – genau wie in php-dna – eine andere Struktur:

| Methode | Erfolg | Fehler |
|---|---|---|
| `CheckAvailability` | eine einfache `[]any`-Liste | der Fehler-Umschlag (`map[string]any`), daher ist der Rückgabetyp `any` |
| `GetCurrentBalance` | Umschlag mit `OperationResult: "SUCCESS"` | derselbe Umschlag mit `OperationResult: "FAILED"` |
| `GetResellerDetails` | flache Map (`id`, `name`, `balance`, `balances`, ...) | der Fehler-Umschlag |

### Optionale Argumente

Go kennt keine Standardargumente. Die Methoden behalten die Parameterreihenfolge von php-dna bei und
erwarten jedes Argument explizit; der Standardwert von php-dna ist in jedem der folgenden Abschnitte
angegeben. Der Einfachheit halber wählen einige Nullwerte den Standard: `GetCurrentBalance("")`
bedeutet USD, `GetTldList(0)` bedeutet 20, `CheckAvailability(..., "")` bedeutet `create`,
`SetForwarding(..., "")` bedeutet `Standard` und `ModifyPrivacyProtectionStatus(..., "")` bedeutet
`Owner request`.
`RegisterWithContactInfo` erhält seine vier optionalen Argumente in einem `*RegisterOptions` (nil = Standardwerte).
`AddResourceRecordValue` und `ReplaceResourceRecordValue` erwarten `ttl *int`: `nil` behält die
aktuelle TTL des Record-Sets bei.

### Kontakte

Kontakte sind `map[string]any` mit den Schlüsseln von php-dna (`FirstName`, `LastName`, `Company`,
`EMail`, `AddressLine1`, `AddressLine2`, `City`, `State`, `Country`, `ZipCode`, `Phone`,
`PhoneCountryCode`, `Fax`, `FaxCountryCode`). Die verschachtelte Struktur (`Address.Line1`,
`Phone.Phone.Number`, ...) wird ebenfalls akzeptiert. Da Go-Maps keine Reihenfolge haben, werden die
vier Rollen als geordnete `dna.Contacts`-Liste übergeben:

```go
contacts := dna.Contacts{
	{Role: "Administrative", Contact: contact},
	{Role: "Billing", Contact: contact},
	{Role: "Technical", Contact: contact},
	{Role: "Registrant", Contact: contact},
}
contacts = dna.AllRoles(contact)                 // derselbe Kontakt für alle vier Rollen
contacts = dna.ContactsFromMap(map[string]map[string]any{ /* Rolle => Kontakt */ })
```

Unvollständige Kontakte werden vor dem Senden nach den Regeln von `ValidateContact` ergänzt.

### Methodennamen

Die Go-Methodennamen folgen den Go-Konventionen; der Name in php-dna ist derselbe in lowerCamelCase.
Aus `getServiceUrl`/`setServiceUrl` werden `GetServiceURL`/`SetServiceURL`, aus `setSslVerification`
wird `SetSSLVerification`.

## Verwendung

### Domain-Registrierung

Standardwerte von php-dna: Standard-Nameserver (`ns1.domainnameapi.com`, `ns2.domainnameapi.com`),
EPP-Sperre aktiviert, Datenschutz deaktiviert, keine zusätzlichen Attribute. Zusätzliche Attribute
sind nur für Domains wie `.tr` erforderlich.

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
// Nur für .tr-Domains:
opts.AdditionalAttributes = map[string]any{
	"TRABISDOMAINCATEGORY": 1,             // 1 = Privatperson (Owner), 0 = Unternehmen (Corporate)
	"TRABISCITIZIENID":     "12345678901", // oder TRABISTAXOFFICE + TRABISTAXNUMBER für Unternehmen
}

result := client.RegisterWithContactInfo("example.com", 1, dna.AllRoles(contact), opts)
```

Bei `.tr`-Domains werden die alten TRABIS-Namen auf das Gateway-Schema normalisiert: Die Kategorie
wird zu `Owner`/`Corporate`, `TRABISCITIZIENID` wird zu `TRABISCITIZENID`, und die Felder für Land,
Stadt, Name und Organisation entfallen (das Gateway übernimmt sie aus dem Kontakt). Jeder
Attributwert wird als String gesendet.

<details>
<summary>Beispielausgabe für Domain-Registrierung</summary>

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

Der Registrierungs-Endpunkt antwortet mit einem kurzen Datensatz, daher sind die meisten Felder leer;
rufen Sie anschließend `GetDetails` auf, um die vollständigen Domain-Daten zu erhalten.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain-Verlängerung

```go
result := client.Renew("example.com", 1)
```

<details>
<summary>Beispielausgabe für Domain-Verlängerung</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"ExpirationDate": "2026-03-15T10:00:00"}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain-Transfer

```go
result := client.Transfer("example.com", "Xy9#mK2$", 1, nil)
```

Wie in der Facade von php-dna wird das Kontakt-Argument nicht weitergereicht: Die Transfer-Anfrage
enthält immer eine leere Kontaktliste. Das Ergebnis hat dieselbe Struktur wie `GetDetails`.

<details>
<summary>Beispielausgabe für Domain-Transfer</summary>

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

Zugehörige Transfer-Aufrufe:

```go
client.CancelTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Cancelled"}}
client.ApproveTransfer("example.com") // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Approved"}}
client.RejectTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Rejected"}}
```

<hr style="border: 4px solid #000; border-style: dashed;">

### Transfer-Prüfung

```go
result := client.CheckTransfer("example.com", "Ab1!xyZ9")
```

`result` ist nur dann `OK`, wenn die Registry die Domain als transferierbar meldet; andernfalls ist es
`ERROR` mit dem Code `TRANSFER_NOT_AVAILABLE` und denselben `data`.

<details>
<summary>Beispielausgabe für Transfer-Prüfung</summary>

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

### Domain-Liste

```go
result := client.GetList(map[string]any{
	"OrderColumn":    "Id",
	"OrderDirection": "ASC",
	"PageNumber":     0,
	"PageSize":       1000,
})
```

`PageNumber` (beginnend bei 0), `PageSize`, `OrderColumn` und `OrderDirection` werden in die
Gateway-Parameter `SkipCount`, `MaxResultCount` und `Sorting` übersetzt; andere Schlüssel (zum
Beispiel `DomainName`) werden als Filter durchgereicht. `nil` liefert die ersten 200 Domains.

<details>
<summary>Beispielausgabe für Domain-Liste</summary>

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

### TLD-Liste

```go
result := client.GetTldList(100) // Standardwert in php-dna: 20
```

<details>
<summary>Beispielausgabe für TLD-Liste</summary>

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

### Domain-Verfügbarkeitsprüfung

```go
result := client.CheckAvailability([]string{"hello", "world123x0"}, []string{"com", "net"}, 1, "create")
if list, ok := result.([]any); ok {
	// Erfolg
	_ = list
}
```

<details>
<summary>Beispielausgabe für Domain-Verfügbarkeitsprüfung</summary>

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

Bei Premium-Domains ist `IsFee` `true`.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain-Details

```go
result := client.GetDetails("example.com")
result = client.SyncFromRegistry("example.com") // gleiche Struktur
```

<details>
<summary>Beispielausgabe für Domain-Details</summary>

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

`Status` ist die Textbezeichnung (`Active`, `TransferredOut`, `PendingDelete`, ...), auch wenn das
Gateway den numerischen Enum-Wert sendet. Nicht parsebare Datumsangaben werden zu `""`, niemals zu
einem unsinnigen Datum aus dem Jahr 1970.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Nameserver-Änderung

```go
result := client.ModifyNameServer("example.com", []string{"ns1.example.com", "ns2.example.com"})
```

Eine Map wird ebenfalls akzeptiert (ihre Werte werden in Schlüsselreihenfolge verwendet). Vor dem
Senden werden die Einträge von Leerzeichen befreit, in Kleinbuchstaben umgewandelt, von einem
abschließenden Punkt befreit und dedupliziert.

<details>
<summary>Beispielausgabe für Nameserver-Änderung</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServers": []string{"ns1.example.com", "ns2.example.com"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain-Sperre aktivieren

```go
result := client.EnableTheftProtectionLock("example.com")
```

<details>
<summary>Beispielausgabe für Domain-Sperre aktivieren</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain-Sperre deaktivieren

```go
result := client.DisableTheftProtectionLock("example.com")
```

<details>
<summary>Beispielausgabe für Domain-Sperre deaktivieren</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": false}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Child-Nameserver hinzufügen

```go
result := client.AddChildNameServer("example.com", "ns1.example.com", "1.2.3.4")
```

<details>
<summary>Beispielausgabe für Child-Nameserver hinzufügen</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.4"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Child-Nameserver löschen

```go
result := client.DeleteChildNameServer("example.com", "ns1.example.com")
```

<details>
<summary>Beispielausgabe für Child-Nameserver löschen</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com"}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Child-Nameserver aktualisieren

```go
result := client.ModifyChildNameServer("example.com", "ns1.example.com", "1.2.3.5")
```

<details>
<summary>Beispielausgabe für Child-Nameserver aktualisieren</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.5"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain-Datenschutz ändern

```go
result := client.ModifyPrivacyProtectionStatus("example.com", true, "Owner request")
```

<details>
<summary>Beispielausgabe für Domain-Datenschutz ändern</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"PrivacyProtectionStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain-Kontakte abrufen

```go
result := client.GetContacts("example.com")
```

<details>
<summary>Beispielausgabe für Domain-Kontakte abrufen</summary>

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
			// "Billing", "Registrant", "Technical": gleiche Struktur
		},
	},
}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain-Kontakte speichern

```go
result := client.SaveContacts("example.com", dna.AllRoles(contact))
```

<details>
<summary>Beispielausgabe für Domain-Kontakte speichern</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"contacts": map[string]any{}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Reseller-Details abrufen

```go
result := client.GetResellerDetails()
```

<details>
<summary>Beispielausgabe für Reseller-Details abrufen</summary>

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
		map[string]any{"balance": "1250.5", "currency": "TL", "symbol": "TL"},
	},
}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Aktuelles Guthaben abrufen

```go
result := client.GetCurrentBalance("USD") // USD (Standard), TRY, EUR oder GBP
```

<details>
<summary>Beispielausgabe für Aktuelles Guthaben abrufen</summary>

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

Im Fehlerfall: `ErrorCode` ist `API_<status>`, `OperationMessage` die Fehlermeldung,
`OperationResult` `"FAILED"` und `Balance` `"0.00"`.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### DNS-Zonenverwaltung

Das Gateway speichert Records als **Sets** (ein Name + Typ mit allen zugehörigen Werten), und jeder
Schreibvorgang ersetzt das gesamte Set. Record-Namen können als `@`, `""`, Zonen-Apex, relatives
Label (`www`) oder vollständiger Name (`www.example.com`) angegeben werden; sie werden auf das relative
Label reduziert, das das Gateway erwartet. Unterstützte Typen: `A`, `AAAA`, `ALIAS`, `CAA`, `CNAME`,
`DS`, `MX`, `NAPTR`, `PTR`, `SRV`, `TLSA`, `TXT` (andere werden vor jedem Aufruf abgewiesen). Die TTL
wird auf 1..86400 begrenzt. Mit `apply = true` (Standard in php-dna) werden Änderungen sofort
übernommen; aufeinanderfolgende Zonen-Aufrufe liegen mindestens eine Sekunde auseinander, weil das
Gateway Anfrage-Spitzen per Rate-Limit drosselt.

```go
// Auflisten: eine Zeile pro Wert, Namen vollqualifiziert
client.GetResourceRecords("example.com")

// Name+Typ auf genau diese Werte setzen (leere Liste entfernt das Set). Standardwerte: ttl 3600, apply true
client.SetResourceRecordSet("example.com", "@", "MX", []string{"10 mx1.example.com", "20 mx2.example.com"}, 3600, true)

// EINEN Wert hinzufügen / ersetzen / entfernen, die anderen bleiben erhalten (ttl nil = TTL des Sets beibehalten)
client.AddResourceRecordValue("example.com", "www", "A", "192.0.2.10", nil, true)
client.ReplaceResourceRecordValue("example.com", "www", "A", "192.0.2.10", "192.0.2.11", nil, true)
client.RemoveResourceRecordValue("example.com", "www", "A", "192.0.2.11", true)

// Einen kompletten Name+Typ entfernen (knownValue "" = zuerst einen aktuellen Wert abfragen)
client.DeleteResourceRecordSet("example.com", "www", "A", true, "")

// Low-Level-Aufrufe (content ist ein String oder []string)
client.AddResourceRecord("example.com", "api", "CNAME", "example.com.", 300, false)
client.EditResourceRecord("example.com", "api", "api", "CNAME", "www.example.com.", 300, false)
client.DeleteResourceRecord("example.com", "old.example.com.", "192.0.2.99", "A", false)
client.ApplyResourceRecords("example.com")
```

<details>
<summary>Beispielausgaben für DNS-Zonenverwaltung</summary>

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

// Wert bereits vorhanden / nicht vorhanden: es wird nichts geschrieben
map[string]any{"result": "OK", "data": map[string]any{"Name": "", "Type": "MX", "Contents": []any{"..."}, "Changed": false}}

// DeleteResourceRecordSet für ein nicht existierendes Set
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

### URL-Weiterleitung

```go
client.SetForwarding("example.com", "https://www.example.org", "Standard") // oder "Frame" (maskiert)
client.GetForwarding("example.com")
client.DeleteForwarding("example.com")
```

<details>
<summary>Beispielausgaben für URL-Weiterleitung</summary>

```go
// SetForwarding
map[string]any{"result": "OK", "data": map[string]any{"RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// GetForwarding (keine Weiterleitung konfiguriert: Enabled false, trotzdem OK)
map[string]any{"result": "OK", "data": map[string]any{"Enabled": true, "RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// DeleteForwarding
map[string]any{"result": "OK", "data": map[string]any{"DomainName": "example.com"}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Hilfsfunktionen (ohne Netzwerk)

```go
client.IsTrTLD("example.com.tr") // true

client.ValidateContact(map[string]any{"FirstName": " Jo ", "Phone": "+90 (532) 374 33 88", "Country": "Turkey"})
// map[string]any{"FirstName": "Jo", "LastName": "Jo", "AddressLine1": "Addres yok", "City": "ISTANBUL",
//   "State": "ISTANBUL", "Country": "TR", "ZipCode": "34000", "Phone": "5323743388",
//   "PhoneCountryCode": "90", "Fax": "", "FaxCountryCode": ""}
```

`ValidateContact` füllt fehlende Felder mit Platzhaltern, entfernt Ziffern und Satzzeichen aus Namen,
greift bei fehlendem State auf City zurück, ordnet gängige Ländernamen ISO-Codes zu und zerlegt
Telefon- und Faxnummern in Rufnummer und Ländervorwahl.

<hr style="border: 4px solid #000; border-style: dashed;">

### Diagnose

```go
client.GetRequestData()        // {"url": ..., "method": ..., "payload": ...} des letzten Aufrufs
client.GetResponseData()       // letzte dekodierte Antwort oder {"raw_response": body}
client.GetResponseHeaders()    // {"url", "http_code", "content_type", "total_time", "headers"}
client.GetParsedResponseData() // letzte geparste Antwort oder der Fehler {Code, Message, Details}
client.GetLastFunction()       // z. B. "request:GET domains/info"
client.GetServiceURL()         // Basis-URL des Gateways
client.SetServiceURL("https://ote.domainresellerapi.com/api/v1")
```

Ein `Client` speichert diesen aufrufbezogenen Zustand; teilen Sie daher einen `Client` nicht zwischen
Goroutinen, sondern erstellen Sie pro Goroutine einen eigenen.

## Fehlerbehandlung

- Gateway-Fehler liefern `Code: "API_<HTTP status>"`. Die Meldung stammt vom Gateway; ein leerer
  Fehler-Body wird zu `HTTP 429 - Too Many Requests` (usw.), Validierungsfehler werden in `Details`
  Feld für Feld aufgelistet (`Admin PostalCode - 'PostalCode' is required`), und Standardtexte des
  Backends wie `Exception of type '...' was thrown` werden durch die aussagekräftigen Details ersetzt.
- Eine 301/302-Antwort bedeutet, dass die Zugangsdaten abgelehnt wurden: `Message: "Invalid API credentials"`.
- Transportfehler (Timeout, DNS, TLS, Verbindungsabbruch) sind Fehler, niemals ein leerer Erfolg. Ein
  Timeout liefert `Code: "API_28"` und `Could not reach the DomainNameAPI service: the connection
  timed out. (...)`; andere Fehler verwenden Nummern im curl-Stil (`API_6` DNS, `API_7` Verbindung,
  `API_35` TLS, `API_56` Reset, `API_60` Zertifikat) und `Curl error during request: ...`.
  `GetParsedResponseData()["Code"]` enthält `CURL_<n>`.
- Ein leerer oder `null`-Erfolgs-Body liefert `Code: "RESPONSE"`.
- Timeouts: 30 Sekunden pro Aufruf, 10 Sekunden für den Verbindungsaufbau.

## Testmodus

Übergeben Sie `true` (oder `"true"`, `"1"`, `"on"`, `"yes"`) als drittes Konstruktor-Argument, um das
OTE-Test-Gateway `https://ote.domainresellerapi.com/api/v1` zu verwenden. Telemetrie im Testmodus wird
mit `environment: ote` markiert.

## Telemetrie (Fehlerberichte)

Wenn ein API-Aufruf fehlschlägt, kann die Bibliothek einen Fehlerbericht an den eigenen Fehler-Tracker
von DomainNameAPI (Sentry) senden. Ein kleiner, zufällig gewählter Anteil erfolgreicher Aufrufe (4 %,
dazu jeder Aufruf, der länger als 1 Sekunde dauert) wird als Performance-Datensatz gesendet.
Fehlerberichte sind **standardmäßig aktiviert**.

> ℹ️ Das Sentry-Projekt für `DomainNameAPI-Go` ist noch nicht angelegt: `DefaultSentryDSN` ist leer,
> daher wird nichts gesendet, bis ein Release diesen Wert setzt. `dna.WithSentryDSN(...)` überschreibt den DSN.

**Was gesendet wird:** Fehlercode und -meldung; Endpunktpfad und HTTP-Methode; Dauer des Aufrufs;
Bibliotheksname (`DomainNameAPI-Go`), Bibliotheksversion und Go-Version; Reseller-ID und Domainname;
die öffentliche IP-Adresse des Servers (ermittelt über `api.ipify.org`); der Benutzername des
Betriebssystems; API-Anfrage und -Antwort einschließlich Kontaktdaten. Nur Zugangsdaten werden
maskiert: EPP-/Auth-Codes, API-Keys, Passwörter und Tokens werden durch `***` ersetzt.

**Was nie gesendet wird:** EPP-/Auth-Codes, API-Keys, Passwörter oder Tokens in lesbarer Form;
Funktionsargumente in Stacktraces.

**Wozu und wie lange:** Die Daten werden ausschließlich zur Nutzungsmessung und zur frühzeitigen
Erkennung von Problemen verwendet und 30 Tage lang auf dem eigenen Sentry-Server von DomainNameAPI
gespeichert.

**Abschalten oder nur das Minimum senden:**

```go
client := dna.New(resellerID, apiKey, false)
client.SetErrorReportingEnabled(false) // es wird nichts gesendet

client.SetTelemetryMinimal(true) // nur Fehlercode/-meldung, Endpunkt, Zeitmessung, Versionen und Reseller-ID
```

Die Telemetrie läuft im Hintergrund, blockiert den Aufruf nie lange und lässt niemals eine Panic in
Ihren Code durchschlagen.

> 🔒 **Sicherheit:** Die Prüfung von TLS-Zertifikaten ist aktiviert. Auf einem Host mit veraltetem
> CA-Bundle kann sie mit `client.SetSSLVerification(false)` abgeschaltet werden — nicht empfohlen; die
> bessere Lösung ist, das CA-Bundle zu aktualisieren.

## Erweiterte Optionen

```go
client := dna.New(resellerID, apiKey, false,
	dna.WithHTTPClient(myHTTPClient),      // beliebiges Do(*http.Request) (*http.Response, error)
	dna.WithSentryDSN("https://key@host/1"),
	dna.WithTelemetryClient(myHTTPClient),
)
```

Bei einem eigenen HTTP-Client konfigurieren Sie Timeouts, Proxys und TLS direkt an diesem Client; er
darf Weiterleitungen nicht folgen (geben Sie in `CheckRedirect` `http.ErrUseLastResponse` zurück).

## Tests

Go muss auf dem Host nicht installiert sein; Docker genügt:

```bash
docker run --rm -v "$PWD":/src -w /src golang:1.23 go test ./...
docker run --rm -v "$PWD":/src -w /src golang:1.23 go vet ./...
```

Die Testsuite enthält Unit-Tests für die HTTP-Schicht, die Fehlerzuordnung, die Telemetrie und die
Hilfsfunktionen sowie einen Konformitätstest, der alle 126 php-dna-Fixtures aus `testdata/conformance`
erneut abspielt (Anfragen auf der Leitung und Rückgabewerte müssen exakt mit php-dna übereinstimmen).
Kein Test greift auf das Netzwerk zu.

## Rückgabe- und Fehlercodes mit Erklärungen

| Code | Beschreibung                                   | Details                                                                                                                                                        |
|------|------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 1000 | Befehl erfolgreich abgeschlossen               | Operation erfolgreich.                                                                                                                                         |
| 1001 | Befehl erfolgreich abgeschlossen; Aktion ausstehend | Operation erfolgreich. Der Vorgang wurde jedoch zur Fertigstellung in die Warteschlange eingereiht.                                                       |
| 2003 | Erforderlicher Parameter fehlt                 | Parameterfehler. Zum Beispiel: Keine Telefonnummer in den Kontaktinformationen angegeben.                                                                      |
| 2105 | Objekt ist nicht für Verlängerung geeignet     | Domain-Status erlaubt keine Verlängerung, ist für Updates gesperrt. Status darf nicht "clientupdateprohibited" sein. Kann auch an anderen Statusbedingungen liegen. |
| 2200 | Authentifizierungsfehler                       | Berechtigungsfehler, Sicherheitscode falsch oder Domain bei einem anderen Registrar registriert.                                                               |
| 2302 | Objekt existiert bereits                       | Domain-Name oder Nameserver-Information existiert bereits in der Datenbank. Kann nicht registriert werden.                                                     |
| 2303 | Objekt existiert nicht                         | Domain-Name oder Nameserver-Information existiert nicht in der Datenbank. Neuregistrierung erforderlich.                                                       |
| 2304 | Objektstatus verhindert Operation              | Domain-Status erlaubt keine Aktualisierung, ist für Updates gesperrt. Status darf nicht "clientupdateprohibited" sein. Kann auch an anderen Statusbedingungen liegen. |

Bibliothekscodes: `RESPONSE` (leere Antwort), `CREDENTIALS` (ungültiges Antwortformat vom
Account-Endpunkt), `AVAILABILITY` (keine Domainnamen angegeben), `TRANSFER_NOT_AVAILABLE`,
`DOMAIN_RENEW` (Verlängerung ohne Ablaufdatum), `ADD_ZONE`, `EDIT_ZONE`, `SET_ZONE`, `DELETE_ZONE`
(nicht unterstützter Record-Typ oder Wert nicht gefunden), `SET_FORWARD` (nicht unterstützter
Weiterleitungstyp oder leere Adresse).

## Lizenz

Apache License 2.0 — siehe [LICENSE](LICENSE).
