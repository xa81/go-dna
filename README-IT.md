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

# go-dna — Client DomainNameAPI per Go

Client Go per l'API REST del registrar di domini DomainNameAPI. È un porting fedele di
[php-dna](https://github.com/domainreseller/php-dna): gli stessi metodi, le stesse richieste in
rete e gli stessi valori di ritorno. La libreria standard è l'unica dipendenza.

## 📦 Download — usa sempre le Releases!

⬇️ **Scarica qui l'ultima versione testata: https://github.com/domainreseller/go-dna/releases/latest**

> ⚠️ **Non** usare il pulsante verde **Code → Download ZIP** — scarica il ramo di sviluppo grezzo. I pacchetti di release sono versionati, testati e pronti per la produzione.

## Requisiti

- Go 1.21 o superiore
- Nessun pacchetto di terze parti (solo `net/http` e `encoding/json` della libreria standard)

## Installazione

Con i moduli Go, fissa una versione rilasciata:

```bash
go get github.com/domainreseller/go-dna@v1.0.0
```

```go
import dna "github.com/domainreseller/go-dna"
```

Il nome del pacchetto è `dna`. Esempi eseguibili per ogni metodo si trovano nella cartella [examples](examples):

```bash
go run ./examples/GetDetails
```

## 🔑 Credenziali API — Reseller ID e API Key

La libreria richiede esattamente due credenziali, passate al costruttore:

| Parametro | Valore |
|---|---|
| Primo | Reseller ID — UUID nel formato `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` |
| Secondo | API Key |

> 💡 Trovi il tuo **Reseller ID** e la tua **API Key** nel pannello DomainNameAPI, alla voce **API Settings**.
> ⚠️ Si tratta di **credenziali API** — i dati di accesso al pannello **non** funzionano qui.

```go
resellerID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
apiKey := "your-api-key"

// Ambiente di produzione
client := dna.New(resellerID, apiKey, false)

// Ambiente di test (OTE)
client = dna.New(resellerID, apiKey, true)
```

Gli spazi superflui di entrambe le credenziali vengono rimossi. Il terzo argomento accetta anche le stringhe che i pannelli
memorizzano per le caselle di controllo: `"true"`, `"1"`, `"on"` e `"yes"` (senza distinzione tra maiuscole e minuscole) selezionano l'ambiente di test;
qualsiasi altro valore (`"false"`, `"0"`, `"off"`, `""`, ...) seleziona la produzione.

## Avvio rapido

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

## Valori di ritorno

I metodi restituiscono le stesse strutture dinamiche di php-dna — `map[string]any` e `[]any` — e non
modelli tipizzati, quindi le chiavi, i flag in forma di stringa (`"true"`/`"false"`) e gli importi formattati
(`"11.9900"`) sono esattamente quelli di php-dna. I numeri provenienti dal gateway sono `int64` o `float64`.

Le chiamate riuscite hanno questo aspetto:

```go
map[string]any{"result": "OK", "data": map[string]any{ /* ... */ }}
```

Le chiamate non riuscite non generano mai un panic; restituiscono l'envelope di errore:

```go
map[string]any{
	"result": "ERROR",
	"error": map[string]any{
		"Code":        "API_404", // API_<stato HTTP>, CURL_<n>, oppure un codice della libreria come RESPONSE
		"Message":     "The domain name you requested could not be found. Please check and try again.",
		"Details":     "The domain could not be found. Please check and try again.",
		"GatewayCode": "Dna.DomainService:Domain:10007", // solo con API_<stato>, quando il gateway ha inviato un proprio codice
	},
}
```

Tre metodi usano una forma diversa, esattamente come in php-dna:

| Metodo | Successo | Errore |
|---|---|---|
| `CheckAvailability` | una semplice lista `[]any` | l'envelope di errore (`map[string]any`), quindi il tipo di ritorno è `any` |
| `GetCurrentBalance` | envelope `OperationResult: "SUCCESS"` | lo stesso envelope con `OperationResult: "FAILED"` |
| `GetResellerDetails` | map piatta (`id`, `name`, `balance`, `balances`, ...) | l'envelope di errore |

### Argomenti opzionali

Go non dispone di argomenti predefiniti. I metodi mantengono l'ordine dei parametri di php-dna e ricevono ogni argomento
in modo esplicito; il valore predefinito di php-dna è indicato in ciascuna sezione seguente. Per comodità, alcuni valori zero
selezionano il valore predefinito: `GetCurrentBalance("")` significa USD, `GetTldList(0)` significa 20,
`CheckAvailability(..., "")` significa `create` e `SetForwarding(..., "")` significa `Standard`.
`RegisterWithContactInfo` riceve i suoi quattro argomenti opzionali in un `*RegisterOptions` (nil = valori predefiniti).
`AddResourceRecordValue` e `ReplaceResourceRecordValue` accettano `ttl *int`: `nil` mantiene il TTL attuale
del set.

### Contatti

I contatti sono `map[string]any` con le chiavi di php-dna (`FirstName`, `LastName`, `Company`, `EMail`,
`AddressLine1`, `AddressLine2`, `City`, `State`, `Country`, `ZipCode`, `Phone`,
`PhoneCountryCode`, `Fax`, `FaxCountryCode`). È accettata anche la forma annidata (`Address.Line1`,
`Phone.Phone.Number`, ...). Poiché le map di Go non hanno un ordine, i quattro ruoli vengono
passati come lista ordinata `dna.Contacts`:

```go
contacts := dna.Contacts{
	{Role: "Administrative", Contact: contact},
	{Role: "Billing", Contact: contact},
	{Role: "Technical", Contact: contact},
	{Role: "Registrant", Contact: contact},
}
contacts = dna.AllRoles(contact)                 // lo stesso contatto per tutti e quattro i ruoli
contacts = dna.ContactsFromMap(map[string]map[string]any{ /* ruolo => contatto */ })
```

I contatti compilati solo in parte vengono completati secondo le regole di `ValidateContact` prima dell'invio. I contatti
sono usati da `RegisterWithContactInfo` e `SaveContacts`; i trasferimenti non inviano mai contatti.

### Nomi dei metodi

I nomi dei metodi Go seguono le convenzioni di Go; il nome php-dna è lo stesso in lowerCamelCase.
`getServiceUrl`/`setServiceUrl` diventano `GetServiceURL`/`SetServiceURL` e `setSslVerification` diventa
`SetSSLVerification`.

## Utilizzo

### Registrazione Dominio

Valori predefiniti di php-dna: name server predefiniti (`ns1.domainnameapi.com`, `ns2.domainnameapi.com`), blocco EPP
attivo, privacy disattivata, nessun attributo aggiuntivo. Gli attributi aggiuntivi sono richiesti solo per domini
come `.tr`.

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
// Solo per i domini .tr:
opts.AdditionalAttributes = map[string]any{
	"TRABISDOMAINCATEGORY": 1,             // 1 = persona fisica (Owner), 0 = azienda (Corporate)
	"TRABISCITIZIENID":     "12345678901", // oppure TRABISTAXOFFICE + TRABISTAXNUMBER per le aziende
}

result := client.RegisterWithContactInfo("example.com", 1, dna.AllRoles(contact), opts)
```

Per i domini `.tr` i vecchi nomi TRABIS vengono normalizzati secondo lo schema del gateway: la categoria
diventa `Owner`/`Corporate`, `TRABISCITIZIENID` diventa `TRABISCITIZENID` e i campi paese, città,
nome e organizzazione vengono eliminati (il gateway li ricava dal contatto). Ogni valore di attributo
viene inviato come stringa.

<details>
<summary>Output di Esempio per la Registrazione del Dominio</summary>

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

L'endpoint di registrazione risponde con un record ridotto, quindi la maggior parte dei campi è vuota; chiama poi `GetDetails`
per ottenere il dominio completo. Quando il gateway rifiuta la registrazione (HTTP 200 con
`"success": false`), il risultato è un errore con codice `REGISTER_DOMAIN` e il messaggio del gateway.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Rinnovo Dominio

```go
result := client.Renew("example.com", 1)
```

<details>
<summary>Output di Esempio per il Rinnovo del Dominio</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"ExpirationDate": "2026-03-15T10:00:00"}}
```

Un rifiuto da parte del gateway (`"success": false`) restituisce il codice `RENEW_DOMAIN`; un successo senza
data di scadenza restituisce `DOMAIN_RENEW`.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Trasferimento Dominio

```go
result := client.Transfer("example.com", "Xy9#mK2$", 1)
```

I trasferimenti non inviano mai contatti: il registrar li acquisisce insieme al dominio, quindi la richiesta contiene sempre
una lista di contatti vuota e il metodo non ha un parametro per i contatti. Il risultato ha la stessa forma di quello di
`GetDetails`; un rifiuto da parte del gateway (`"success": false`) restituisce il codice `TRANSFER_DOMAIN` con il messaggio del gateway.

<details>
<summary>Output di Esempio per il Trasferimento del Dominio</summary>

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

Chiamate correlate al trasferimento:

```go
client.CancelTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Cancelled"}}
client.ApproveTransfer("example.com") // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Approved"}}
client.RejectTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Rejected"}}
```

<hr style="border: 4px solid #000; border-style: dashed;">

### Verifica Trasferimento

```go
result := client.CheckTransfer("example.com", "Ab1!xyZ9")
```

`result` è `OK` solo quando il registro indica che il dominio è trasferibile; altrimenti è `ERROR`
con codice `TRANSFER_NOT_AVAILABLE` e gli stessi `data`.

<details>
<summary>Output di Esempio per la Verifica Trasferimento</summary>

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

### Elenco Domini

```go
result := client.GetList(map[string]any{
	"OrderColumn":    "Id",
	"OrderDirection": "ASC",
	"PageNumber":     0,
	"PageSize":       1000,
})
```

`PageNumber` (a partire da 0), `PageSize`, `OrderColumn` e `OrderDirection` vengono convertiti in
`SkipCount`, `MaxResultCount` e `Sorting` del gateway; le altre chiavi (ad esempio `DomainName`) vengono
passate invariate come filtri. `nil` restituisce i primi 200 domini.

<details>
<summary>Output di Esempio per l'Elenco Domini</summary>

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

### Elenco TLD

```go
result := client.GetTldList(100) // valore predefinito di php-dna: 20
```

<details>
<summary>Output di Esempio per l'Elenco TLD</summary>

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

### Verifica Disponibilità Dominio

```go
result := client.CheckAvailability([]string{"hello", "world123x0"}, []string{"com", "net"}, 1, "create")
if list, ok := result.([]any); ok {
	// successo
	_ = list
}
```

<details>
<summary>Output di Esempio per la Verifica Disponibilità Dominio</summary>

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

`IsFee` è `true` per i domini premium.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Informazioni Dominio

```go
result := client.GetDetails("example.com")
result = client.SyncFromRegistry("example.com") // stessa forma
```

<details>
<summary>Output di Esempio per le Informazioni del Dominio</summary>

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

`Status` è l'etichetta testuale (`Active`, `TransferredOut`, `PendingDelete`, ...) anche quando il
gateway invia l'enum numerico. Le date non interpretabili diventano `""`, mai una falsa data del 1970.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Modifica Name Server

```go
result := client.ModifyNameServer("example.com", []string{"ns1.example.com", "ns2.example.com"})
```

È accettata anche una map (i suoi valori vengono usati nell'ordine delle chiavi). Prima dell'invio le voci vengono ripulite
dagli spazi, convertite in minuscolo, private dell'eventuale punto finale e deduplicate; il risultato riporta questa lista normalizzata.

<details>
<summary>Output di Esempio per la Modifica Name Server</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServers": []any{"ns1.example.com", "ns2.example.com"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Attivare il Blocco Dominio

```go
result := client.EnableTheftProtectionLock("example.com")
```

<details>
<summary>Output di Esempio per Attivare il Blocco Dominio</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Disattivare il Blocco Dominio

```go
result := client.DisableTheftProtectionLock("example.com")
```

<details>
<summary>Output di Esempio per Disattivare il Blocco Dominio</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": false}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Aggiungere Child Name Server

```go
result := client.AddChildNameServer("example.com", "ns1.example.com", "1.2.3.4")
```

<details>
<summary>Output di Esempio per Aggiungere Child Name Server</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.4"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Eliminare Child Name Server

```go
result := client.DeleteChildNameServer("example.com", "ns1.example.com")
```

<details>
<summary>Output di Esempio per Eliminare Child Name Server</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com"}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Modificare Child Name Server

```go
result := client.ModifyChildNameServer("example.com", "ns1.example.com", "1.2.3.5")
```

<details>
<summary>Output di Esempio per Modificare Child Name Server</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.5"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Modificare Protezione Privacy

```go
result := client.ModifyPrivacyProtectionStatus("example.com", true, "Owner request")
```

`status` è un bool oppure un valore memorizzato dai pannelli: `true`, `1`, `"1"`, `"true"`, `"on"` e `"yes"` attivano
la privacy, qualsiasi altro valore (`"off"`, `"0"`, `false`, ...) la disattiva. Il motivo è mantenuto solo per compatibilità
con php-dna e non viene inviato: il gateway non prevede questo campo.

<details>
<summary>Output di Esempio per Modificare Protezione Privacy</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"PrivacyProtectionStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Ottenere Contatti Dominio

```go
result := client.GetContacts("example.com")
```

I tipi di contatto del gateway vengono riconosciuti senza distinzione tra maiuscole e minuscole (`Admin`/`Administrative`, `Tech`/`Technical`,
`Registrant`/`Owner`, `Billing`); qualsiasi altro tipo viene mantenuto con il proprio nome, con l'iniziale maiuscola.

<details>
<summary>Output di Esempio per Ottenere Contatti Dominio</summary>

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
			// "Billing", "Registrant", "Technical": stessa forma
		},
	},
}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Salvare Contatti Dominio

```go
result := client.SaveContacts("example.com", dna.AllRoles(contact))
```

<details>
<summary>Output di Esempio per Salvare Contatti Dominio</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"contacts": map[string]any{}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Ottenere Dettagli Rivenditore

```go
result := client.GetResellerDetails()
```

<details>
<summary>Output di Esempio per Ottenere Dettagli Rivenditore</summary>

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

### Ottenere Saldo Attuale

```go
result := client.GetCurrentBalance("USD") // USD (predefinito), TRY, EUR o GBP
```

<details>
<summary>Output di Esempio per Ottenere Saldo Attuale</summary>

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

In caso di errore: `ErrorCode` è `API_<status>` (oppure `CURL_<n>`), `OperationMessage` contiene il messaggio di errore,
`OperationResult` è `"FAILED"` e `Balance` è `"0.00"`. L'account riporta solo i saldi in USD e TRY:
per qualsiasi altra valuta il risultato è `FAILED` con `ErrorCode: "BALANCE"` e
`No EUR balance is reported for this account. Available: USD, TRY.`

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Gestione Zona DNS

Il gateway memorizza i record come **set** (un nome + un tipo con tutti i relativi valori) e ogni scrittura
sostituisce l'intero set. I nomi dei record possono essere indicati come `@`, `""`, l'apex della zona, un'etichetta relativa
(`www`) o un nome completo (`www.example.com`); vengono ridotti all'etichetta relativa attesa dal gateway.
Tipi supportati: `A`, `AAAA`, `ALIAS`, `CAA`, `CNAME`, `DS`, `MX`, `NAPTR`, `PTR`, `SRV`,
`TLSA`, `TXT` (gli altri vengono rifiutati prima di qualsiasi chiamata). Il TTL viene limitato all'intervallo 1..86400 e il risultato riporta
il tipo in maiuscolo e il TTL effettivamente inviati. Con `apply = true` (il valore predefinito di php-dna) le modifiche vengono
applicate immediatamente; se questa applicazione non riesce, viene restituito il suo envelope di errore (la modifica resta
in attesa). Le chiamate consecutive sulla zona vengono distanziate di almeno un secondo, perché il gateway limita la frequenza
delle raffiche di richieste.

```go
// Elenco: una riga per valore, nomi completamente qualificati
client.GetResourceRecords("example.com")

// Imposta un nome+tipo esattamente su questi valori (una lista vuota rimuove il set). Predefiniti: ttl 3600, apply true
client.SetResourceRecordSet("example.com", "@", "MX", []string{"10 mx1.example.com", "20 mx2.example.com"}, 3600, true)

// Aggiunge / sostituisce / rimuove UN valore, mantenendo gli altri (ttl nil = mantiene il TTL del set)
client.AddResourceRecordValue("example.com", "www", "A", "192.0.2.10", nil, true)
client.ReplaceResourceRecordValue("example.com", "www", "A", "192.0.2.10", "192.0.2.11", nil, true)
client.RemoveResourceRecordValue("example.com", "www", "A", "192.0.2.11", true)

// Rimuove un intero nome+tipo (knownValue "" = cerca prima un valore attuale)
client.DeleteResourceRecordSet("example.com", "www", "A", true, "")

// Chiamate di basso livello (content è una string o un []string)
client.AddResourceRecord("example.com", "api", "CNAME", "example.com.", 300, false)
client.EditResourceRecord("example.com", "api", "api", "CNAME", "www.example.com.", 300, false)
client.DeleteResourceRecord("example.com", "old.example.com.", "192.0.2.99", "A", false)
client.ApplyResourceRecords("example.com")
```

<details>
<summary>Output di Esempio per la Gestione Zona DNS</summary>

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

// Valore già presente / non presente: non viene scritto nulla
map[string]any{"result": "OK", "data": map[string]any{"Name": "", "Type": "MX", "Contents": []any{"..."}, "Changed": false}}

// DeleteResourceRecordSet di un set inesistente
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

### Inoltro URL

```go
client.SetForwarding("example.com", "https://www.example.org", "Standard") // oppure "Frame" (mascherato)
client.GetForwarding("example.com")
client.DeleteForwarding("example.com")
```

<details>
<summary>Output di Esempio per l'Inoltro URL</summary>

```go
// SetForwarding
map[string]any{"result": "OK", "data": map[string]any{"RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// GetForwarding
map[string]any{"result": "OK", "data": map[string]any{"Enabled": true, "RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// GetForwarding quando non è configurato alcun inoltro
map[string]any{"result": "OK", "data": map[string]any{"Enabled": false, "RedirectAddress": "", "ForwardType": ""}}

// DeleteForwarding (Removed è false quando non era configurato nulla)
map[string]any{"result": "OK", "data": map[string]any{"DomainName": "example.com", "Removed": true}}
```

</details>

Solo la risposta "Redirect information not found" del gateway significa "non configurato"; qualsiasi altro errore,
ad esempio un dominio sconosciuto, viene restituito come errore.

<hr style="border: 4px solid #000; border-style: dashed;">

### Utilità (senza rete)

```go
client.IsTrTLD("example.com.tr") // true

client.ValidateContact(map[string]any{"FirstName": " Jo ", "Phone": "+90 (532) 374 33 88", "Country": "Turkey"})
// map[string]any{"FirstName": "Jo", "LastName": "Jo", "AddressLine1": "Addres yok", "City": "ISTANBUL",
//   "State": "ISTANBUL", "Country": "TR", "ZipCode": "34000", "Phone": "5323743388",
//   "PhoneCountryCode": "90", "Fax": "", "FaxCountryCode": ""}
```

`ValidateContact` compila i campi mancanti con valori segnaposto, rimuove cifre e punteggiatura dai nomi,
usa la città in mancanza della provincia/stato, converte i nomi di paese più comuni in codici ISO e separa i numeri di telefono e di fax
in numero e prefisso internazionale. Un numero scritto con il proprio prefisso internazionale
(`+44 20 ...`, `0044 20 ...`) mantiene il proprio prefisso anche quando ne viene indicato un altro.

<hr style="border: 4px solid #000; border-style: dashed;">

### Diagnostica

```go
client.GetRequestData()        // {"url": ..., "method": ..., "payload": ...} dell'ultima chiamata
client.GetResponseData()       // ultima risposta decodificata, oppure {"raw_response": body}
client.GetResponseHeaders()    // {"url", "http_code", "content_type", "total_time", "headers"}
client.GetParsedResponseData() // ultima risposta elaborata, oppure l'errore {Code, Message, Details}
client.GetLastFunction()       // ad es. "request:GET domains/info"
client.GetServiceURL()         // URL di base del gateway
client.SetServiceURL("https://ote.domainresellerapi.com/api/v1")
```

Un `Client` conserva questo stato relativo a ogni chiamata, quindi non condividere lo stesso `Client` tra più goroutine; creane uno
per ogni goroutine.

## Gestione degli errori

- Gli errori del gateway restituiscono `Code: "API_<HTTP status>"`. Quando il gateway ha inviato un proprio codice di errore
  (ad esempio `Dna.DomainService:Domain:10007`), questo viene aggiunto come `GatewayCode`, così è possibile distinguere
  "dominio non trovato" da "TLD non supportato". Il messaggio proviene dal gateway; un corpo di errore
  vuoto diventa `HTTP 429 - Too Many Requests` (e così via), gli errori di validazione sono elencati campo per
  campo in `Details` (`Admin PostalCode - 'PostalCode' is required`) e i messaggi generici del backend come
  `Exception of type '...' was thrown` vengono sostituiti dai dettagli del gateway, oppure da un messaggio stabile
  che indica il codice di errore — mai dal corpo grezzo della risposta.
- Una risposta 301/302 significa che le credenziali sono state rifiutate: `Message: "Invalid API credentials"`.
- Gli errori di trasporto (timeout, DNS, TLS, connessione reimpostata) sono errori, mai un successo vuoto. Un
  timeout restituisce `Code: "CURL_28"` e `Could not reach the DomainNameAPI service: the connection
  timed out. (...)`; gli altri errori usano numeri in stile curl (`CURL_6` DNS, `CURL_7` connessione,
  `CURL_35` TLS, `CURL_56` reset, `CURL_60` certificato) e `Curl error during request: ...`.
  `API_<n>` è riservato agli stati HTTP.
- Le risposte HTTP 200 con `"success": false` a registrazione, trasferimento e rinnovo sono rifiuti:
  `REGISTER_DOMAIN`, `TRANSFER_DOMAIN` o `RENEW_DOMAIN` con il messaggio del gateway.
- Un corpo di successo vuoto o `null` restituisce `Code: "RESPONSE"`.
- Timeout: 30 secondi per chiamata, 10 secondi per la connessione.

## Modalità test

Passa `true` (oppure `"true"`, `"1"`, `"on"`, `"yes"`) come terzo argomento del costruttore per usare il gateway
di test OTE `https://ote.domainresellerapi.com/api/v1`. La telemetria in modalità test è contrassegnata con
`environment: ote`.

## Telemetria (segnalazione errori)

Quando una chiamata API non riesce, la libreria può inviare una segnalazione di errore al sistema di tracciamento errori di DomainNameAPI
(Sentry). Una piccola quota casuale delle chiamate riuscite (il 4%, più ogni chiamata che dura oltre 1 secondo) viene inviata
come record di performance. La segnalazione è **attiva per impostazione predefinita**.

> ℹ️ Il progetto Sentry per `DomainNameAPI-Go` non è ancora stato creato: `DefaultSentryDSN` è vuoto, quindi
> non viene inviato nulla finché una release non lo imposta. `dna.WithSentryDSN(...)` sovrascrive il DSN.

**Cosa viene inviato:** codice e messaggio di errore; percorso dell'endpoint e metodo HTTP; durata della chiamata; nome della libreria
(`DomainNameAPI-Go`), versione della libreria e versione di Go; Reseller ID e nome di dominio; l'indirizzo IP
pubblico del server (ricavato tramite `api.ipify.org`); il nome utente del sistema operativo; la richiesta
e la risposta API, compresi i dati dei contatti. Vengono mascherate solo le credenziali: codici EPP/auth, chiavi API,
password e token diventano `***`.

**Cosa non viene mai inviato:** codici EPP/auth, chiavi API, password o token in forma leggibile; argomenti
delle chiamate di funzione negli stack trace.

**Perché e per quanto tempo:** i dati vengono usati solo per misurare l'utilizzo e diagnosticare tempestivamente i problemi e
sono conservati per 30 giorni sul server Sentry di DomainNameAPI.

**Disattivala, o invia solo il minimo:**

```go
client := dna.New(resellerID, apiKey, false)
client.SetErrorReportingEnabled(false) // non viene inviato nulla

client.SetTelemetryMinimal(true) // solo codice/messaggio di errore, endpoint, tempi, versioni e Reseller ID
```

La telemetria viene eseguita in background, non blocca mai a lungo la chiamata e non genera mai un panic nel tuo codice.

> 🔒 **Sicurezza:** la verifica dei certificati TLS è attiva. Su un host con un bundle CA obsoleto può essere
> disattivata con `client.SetSSLVerification(false)` — sconsigliato; aggiornare il bundle CA è
> la soluzione migliore.

## Opzioni avanzate

```go
client := dna.New(resellerID, apiKey, false,
	dna.WithHTTPClient(myHTTPClient),      // qualsiasi Do(*http.Request) (*http.Response, error)
	dna.WithSentryDSN("https://key@host/1"),
	dna.WithTelemetryClient(myHTTPClient),
)
```

Con un client HTTP personalizzato, configura su di esso timeout, proxy e TLS; non deve seguire i redirect
(restituisci `http.ErrUseLastResponse` da `CheckRedirect`).

## Test

Go non è necessario sull'host; basta Docker:

```bash
docker run --rm -v "$PWD":/src -w /src golang:1.23 go test ./...
docker run --rm -v "$PWD":/src -w /src golang:1.23 go vet ./...
```

La suite comprende test unitari per il livello HTTP, la mappatura degli errori, la telemetria e le utilità, oltre a un
test di conformità che riproduce tutte le 134 fixture di php-dna da `testdata/conformance` (le richieste in
rete e i valori di ritorno devono corrispondere esattamente a php-dna). Nessun test accede alla rete.

## Codici di Ritorno ed Errore con Descrizioni

| Codice | Descrizione | Dettaglio |
|--------|-------------|-----------|
| 1000 | Comando completato con successo | Operazione riuscita. |
| 1001 | Comando completato con successo; azione in sospeso | Operazione riuscita. Tuttavia, l'operazione è attualmente in coda per il completamento. |
| 2003 | Parametro richiesto mancante | Errore parametro mancante. Ad esempio: nessuna voce telefonica nelle informazioni di contatto. |
| 2105 | L'oggetto non è idoneo per il rinnovo | Lo stato del dominio non è adatto per il rinnovo, bloccato per operazioni di aggiornamento. Lo stato non dovrebbe essere "clientupdateprohibited". Potrebbe essere dovuto ad altre condizioni di stato. |
| 2200 | Errore di autenticazione | Errore di autorizzazione, il codice di sicurezza non è corretto oppure il dominio è registrato presso un altro registrar. |
| 2302 | L'oggetto esiste | Il nome di dominio o le informazioni del name server esistono già nel database. Impossibile registrare. |
| 2303 | L'oggetto non esiste | Il nome di dominio o le informazioni del name server non esistono nel database. È necessaria una nuova registrazione. |
| 2304 | Lo stato dell'oggetto impedisce l'operazione | Lo stato del dominio non è adatto per gli aggiornamenti, bloccato per operazioni di aggiornamento. Lo stato non dovrebbe essere "clientupdateprohibited". Potrebbe essere dovuto ad altre condizioni di stato. |

Codici della libreria: `RESPONSE` (risposta vuota), `CREDENTIALS` (formato di risposta non valido dall'endpoint
dell'account), `AVAILABILITY` (nessun nome di dominio fornito), `BALANCE` (valuta non riportata dall'account),
`TRANSFER_NOT_AVAILABLE`, `REGISTER_DOMAIN` / `TRANSFER_DOMAIN` / `RENEW_DOMAIN` (il gateway ha
rifiutato l'operazione), `DOMAIN_RENEW` (rinnovo senza data di scadenza), `ADD_ZONE`, `EDIT_ZONE`,
`SET_ZONE`, `DELETE_ZONE` (tipo di record non supportato o valore non trovato), `SET_FORWARD` (tipo di
inoltro non supportato o indirizzo vuoto), `CURL_<n>` (errore di trasporto), `API_<status>` (errore HTTP, con
`GatewayCode` quando il gateway ha inviato un proprio codice).

## Licenza

Apache License 2.0 — vedi [LICENSE](LICENSE).
