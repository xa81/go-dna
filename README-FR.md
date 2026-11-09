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

# go-dna — Client DomainNameAPI pour Go

Client Go pour l'API REST du registraire de domaines DomainNameAPI. Il s'agit d'un portage fidèle de
[php-dna](https://github.com/domainreseller/php-dna) : les mêmes méthodes, les mêmes requêtes sur le
réseau et les mêmes valeurs de retour. La bibliothèque standard est son unique dépendance.

## 📦 Téléchargement — utilisez toujours les Releases !

⬇️ **Obtenez la dernière version testée ici : https://github.com/domainreseller/go-dna/releases/latest**

> ⚠️ N'utilisez **pas** le bouton vert **Code → Download ZIP** — il télécharge la branche de développement brute. Les paquets de release sont versionnés, testés et prêts pour la production.

## Exigences

- Go 1.21 ou supérieur
- Aucun paquet tiers (uniquement `net/http` et `encoding/json` de la bibliothèque standard)

## Installation

Avec les modules Go, fixez une version publiée :

```bash
go get github.com/domainreseller/go-dna@v1.0.0
```

```go
import dna "github.com/domainreseller/go-dna"
```

Le nom du paquet est `dna`. Des exemples exécutables pour chaque méthode se trouvent dans le dossier [examples](examples) :

```bash
go run ./examples/GetDetails
```

## 🔑 Identifiants API — Reseller ID et API Key

La bibliothèque a besoin d'exactement deux identifiants, transmis au constructeur :

| Paramètre | Valeur |
|---|---|
| Premier | Reseller ID — UUID au format `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` |
| Second | API Key |

> 💡 Vous trouverez votre **Reseller ID** et votre **API Key** dans votre panneau DomainNameAPI, sous **API Settings**.
> ⚠️ Il s'agit d'**identifiants API** — vos identifiants de connexion au panneau ne fonctionneront **pas** ici.

```go
resellerID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
apiKey := "your-api-key"

// Environnement de production
client := dna.New(resellerID, apiKey, false)

// Environnement de test (OTE)
client = dna.New(resellerID, apiKey, true)
```

Les espaces superflus des deux identifiants sont supprimés. Le troisième argument accepte également les chaînes
que les panneaux enregistrent pour les cases à cocher : `"true"`, `"1"`, `"on"` et `"yes"` (quelle que soit la casse)
sélectionnent l'environnement de test ; toute autre valeur (`"false"`, `"0"`, `"off"`, `""`, ...) sélectionne la production.

## Démarrage rapide

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

## Valeurs de retour

Les méthodes renvoient les mêmes structures dynamiques que php-dna — `map[string]any` et `[]any` — et non
des modèles typés ; ainsi les clés, les indicateurs sous forme de chaîne (`"true"`/`"false"`) et les montants
formatés (`"11.9900"`) sont exactement ceux de php-dna. Les nombres provenant de la passerelle sont des `int64` ou des `float64`.

Les appels réussis ressemblent à ceci :

```go
map[string]any{"result": "OK", "data": map[string]any{ /* ... */ }}
```

Les appels en échec ne provoquent jamais de panic ; ils renvoient l'enveloppe d'erreur :

```go
map[string]any{
	"result": "ERROR",
	"error": map[string]any{
		"Code":        "API_404", // API_<statut HTTP>, CURL_<n>, ou un code de la bibliothèque tel que RESPONSE
		"Message":     "The domain name you requested could not be found. Please check and try again.",
		"Details":     "The domain could not be found. Please check and try again.",
		"GatewayCode": "Dna.DomainService:Domain:10007", // uniquement avec API_<statut>, lorsque la passerelle a envoyé son propre code
	},
}
```

Trois méthodes utilisent une forme différente, exactement comme dans php-dna :

| Méthode | Succès | Échec |
|---|---|---|
| `CheckAvailability` | une simple liste `[]any` | l'enveloppe d'erreur (`map[string]any`), le type de retour est donc `any` |
| `GetCurrentBalance` | enveloppe `OperationResult: "SUCCESS"` | la même enveloppe avec `OperationResult: "FAILED"` |
| `GetResellerDetails` | map à plat (`id`, `name`, `balance`, `balances`, ...) | l'enveloppe d'erreur |

### Arguments optionnels

Go ne possède pas d'arguments par défaut. Les méthodes conservent l'ordre des paramètres de php-dna et prennent
chaque argument explicitement ; la valeur par défaut de php-dna est indiquée dans chaque section ci-dessous. Par commodité,
quelques valeurs nulles sélectionnent la valeur par défaut : `GetCurrentBalance("")` signifie USD, `GetTldList(0)` signifie 20,
`CheckAvailability(..., "")` signifie `create` et `SetForwarding(..., "")` signifie `Standard`.
`RegisterWithContactInfo` reçoit ses quatre arguments optionnels dans un `*RegisterOptions` (nil = valeurs par défaut).
`AddResourceRecordValue` et `ReplaceResourceRecordValue` prennent `ttl *int` : `nil` conserve le TTL actuel
de l'ensemble.

### Contacts

Les contacts sont des `map[string]any` avec les clés de php-dna (`FirstName`, `LastName`, `Company`, `EMail`,
`AddressLine1`, `AddressLine2`, `City`, `State`, `Country`, `ZipCode`, `Phone`,
`PhoneCountryCode`, `Fax`, `FaxCountryCode`). La forme imbriquée (`Address.Line1`,
`Phone.Phone.Number`, ...) est également acceptée. Comme les maps Go ne sont pas ordonnées, les quatre rôles sont
transmis sous forme d'une liste ordonnée `dna.Contacts` :

```go
contacts := dna.Contacts{
	{Role: "Administrative", Contact: contact},
	{Role: "Billing", Contact: contact},
	{Role: "Technical", Contact: contact},
	{Role: "Registrant", Contact: contact},
}
contacts = dna.AllRoles(contact)                 // le même contact pour les quatre rôles
contacts = dna.ContactsFromMap(map[string]map[string]any{ /* rôle => contact */ })
```

Les contacts à moitié remplis sont complétés selon les règles de `ValidateContact` avant leur envoi. Les contacts sont
utilisés par `RegisterWithContactInfo` et `SaveContacts` ; les transferts n'envoient jamais de contacts.

### Noms des méthodes

Les noms de méthodes Go suivent les conventions de Go ; le nom php-dna est le même en lowerCamelCase.
`getServiceUrl`/`setServiceUrl` deviennent `GetServiceURL`/`SetServiceURL` et `setSslVerification` devient
`SetSSLVerification`.

## Utilisation

### Enregistrement de Domaine

Valeurs par défaut de php-dna : serveurs de noms par défaut (`ns1.domainnameapi.com`, `ns2.domainnameapi.com`), verrou EPP
activé, protection de la vie privée désactivée, aucun attribut supplémentaire. Les attributs supplémentaires ne sont requis que pour des domaines
tels que `.tr`.

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
// Uniquement pour les domaines .tr :
opts.AdditionalAttributes = map[string]any{
	"TRABISDOMAINCATEGORY": 1,             // 1 = particulier (Owner), 0 = société (Corporate)
	"TRABISCITIZIENID":     "12345678901", // ou TRABISTAXOFFICE + TRABISTAXNUMBER pour les sociétés
}

result := client.RegisterWithContactInfo("example.com", 1, dna.AllRoles(contact), opts)
```

Pour les domaines `.tr`, les anciens noms TRABIS sont normalisés selon le schéma de la passerelle : la catégorie
devient `Owner`/`Corporate`, `TRABISCITIZIENID` devient `TRABISCITIZENID`, et les champs pays, ville,
nom et organisation sont supprimés (la passerelle les reprend du contact). Chaque valeur d'attribut
est envoyée sous forme de chaîne.

<details>
<summary>Exemple de Sortie pour l'Enregistrement de Domaine</summary>

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

Le point de terminaison d'enregistrement répond avec un enregistrement abrégé, la plupart des champs sont donc vides ; appelez ensuite `GetDetails`
pour obtenir le domaine complet. Lorsque la passerelle refuse l'enregistrement (HTTP 200 avec
`"success": false`), le résultat est une erreur avec le code `REGISTER_DOMAIN` et le message de la passerelle.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Renouvellement de Domaine

```go
result := client.Renew("example.com", 1)
```

<details>
<summary>Exemple de Sortie pour le Renouvellement de Domaine</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"ExpirationDate": "2026-03-15T10:00:00"}}
```

Un refus de la passerelle (`"success": false`) renvoie le code `RENEW_DOMAIN` ; un succès sans
date d'expiration renvoie `DOMAIN_RENEW`.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Transfert de Domaine

```go
result := client.Transfer("example.com", "Xy9#mK2$", 1)
```

Les transferts n'envoient jamais de contacts : le registraire les reprend avec le domaine, la requête contient donc toujours
une liste de contacts vide et la méthode n'a pas de paramètre de contacts. Le résultat a la même forme que celui de
`GetDetails` ; un refus de la passerelle (`"success": false`) renvoie le code `TRANSFER_DOMAIN` avec le message de la passerelle.

<details>
<summary>Exemple de Sortie pour le Transfert de Domaine</summary>

```go
map[string]any{
	"result": "OK",
	"data": map[string]any{
		"ID": int64(0), "Status": "Active", "DomainName": "", "AuthCode": "",
		"LockStatus": "false", "PrivacyProtectionStatus": "false", "IsChildNameServer": "false",
		"Contacts": map[string]any{ /* Billing, Technical, Administrative, Registrant : {"ID": 0} */ },
		"Dates": map[string]any{"Start": "", "Expiration": "2027-11-09T12:00:00", "RemainingDays": int64(365)},
		"NameServers": []any{}, "Additional": []any{}, "ChildNameServers": []any{},
	},
}
```

</details>

Appels associés au transfert :

```go
client.CancelTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Cancelled"}}
client.ApproveTransfer("example.com") // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Approved"}}
client.RejectTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Rejected"}}
```

<hr style="border: 4px solid #000; border-style: dashed;">

### Vérification de Transfert

```go
result := client.CheckTransfer("example.com", "Ab1!xyZ9")
```

`result` vaut `OK` uniquement lorsque le registre indique que le domaine est transférable ; sinon il vaut `ERROR`
avec le code `TRANSFER_NOT_AVAILABLE` et les mêmes `data`.

<details>
<summary>Exemple de Sortie pour la Vérification de Transfert</summary>

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

### Liste des Domaines

```go
result := client.GetList(map[string]any{
	"OrderColumn":    "Id",
	"OrderDirection": "ASC",
	"PageNumber":     0,
	"PageSize":       1000,
})
```

`PageNumber` (à partir de 0), `PageSize`, `OrderColumn` et `OrderDirection` sont convertis en
`SkipCount`, `MaxResultCount` et `Sorting` de la passerelle ; les autres clés (par exemple `DomainName`) sont
transmises telles quelles comme filtres. `nil` renvoie les 200 premiers domaines.

<details>
<summary>Exemple de Sortie pour la Liste des Domaines</summary>

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

### Liste des TLD

```go
result := client.GetTldList(100) // valeur par défaut de php-dna : 20
```

<details>
<summary>Exemple de Sortie pour la Liste des TLD</summary>

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

### Vérification de Disponibilité de Domaine

```go
result := client.CheckAvailability([]string{"hello", "world123x0"}, []string{"com", "net"}, 1, "create")
if list, ok := result.([]any); ok {
	// succès
	_ = list
}
```

<details>
<summary>Exemple de Sortie pour la Vérification de Disponibilité de Domaine</summary>

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

`IsFee` vaut `true` pour les domaines premium.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Informations de Domaine

```go
result := client.GetDetails("example.com")
result = client.SyncFromRegistry("example.com") // même forme
```

<details>
<summary>Exemple de Sortie pour les Informations de Domaine</summary>

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

`Status` est le libellé textuel (`Active`, `TransferredOut`, `PendingDelete`, ...) même lorsque la
passerelle envoie l'énumération numérique. Les dates illisibles deviennent `""`, jamais une fausse date de 1970.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Modifier les Serveurs de Noms

```go
result := client.ModifyNameServer("example.com", []string{"ns1.example.com", "ns2.example.com"})
```

Une map est également acceptée (ses valeurs sont utilisées dans l'ordre des clés). Avant l'envoi, les entrées sont débarrassées
des espaces superflus, mises en minuscules, privées d'un éventuel point final et dédoublonnées ; le résultat indique cette liste normalisée.

<details>
<summary>Exemple de Sortie pour Modifier les Serveurs de Noms</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServers": []any{"ns1.example.com", "ns2.example.com"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Activer le Verrouillage de Domaine

```go
result := client.EnableTheftProtectionLock("example.com")
```

<details>
<summary>Exemple de Sortie pour Activer le Verrouillage de Domaine</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Désactiver le Verrouillage de Domaine

```go
result := client.DisableTheftProtectionLock("example.com")
```

<details>
<summary>Exemple de Sortie pour Désactiver le Verrouillage de Domaine</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": false}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Ajouter un Serveur de Noms Enfant

```go
result := client.AddChildNameServer("example.com", "ns1.example.com", "1.2.3.4")
```

<details>
<summary>Exemple de Sortie pour Ajouter un Serveur de Noms Enfant</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.4"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Supprimer un Serveur de Noms Enfant

```go
result := client.DeleteChildNameServer("example.com", "ns1.example.com")
```

<details>
<summary>Exemple de Sortie pour Supprimer un Serveur de Noms Enfant</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com"}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Modifier un Serveur de Noms Enfant

```go
result := client.ModifyChildNameServer("example.com", "ns1.example.com", "1.2.3.5")
```

<details>
<summary>Exemple de Sortie pour Modifier un Serveur de Noms Enfant</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.5"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Modifier la Protection de la Vie Privée

```go
result := client.ModifyPrivacyProtectionStatus("example.com", true, "Owner request")
```

`status` est un booléen ou une valeur enregistrée par les panneaux : `true`, `1`, `"1"`, `"true"`, `"on"` et `"yes"` activent
la protection de la vie privée ; toute autre valeur (`"off"`, `"0"`, `false`, ...) la désactive. Le motif n'est conservé que pour
la compatibilité avec php-dna et n'est pas envoyé : la passerelle ne possède pas ce champ.

<details>
<summary>Exemple de Sortie pour Modifier la Protection de la Vie Privée</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"PrivacyProtectionStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Obtenir les Contacts du Domaine

```go
result := client.GetContacts("example.com")
```

Les types de contact de la passerelle sont reconnus sans tenir compte de la casse (`Admin`/`Administrative`, `Tech`/`Technical`,
`Registrant`/`Owner`, `Billing`) ; tout autre type est conservé sous son propre nom, avec une majuscule initiale.

<details>
<summary>Exemple de Sortie pour Obtenir les Contacts du Domaine</summary>

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
			// "Billing", "Registrant", "Technical" : même forme
		},
	},
}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Sauvegarder les Contacts du Domaine

```go
result := client.SaveContacts("example.com", dna.AllRoles(contact))
```

<details>
<summary>Exemple de Sortie pour Sauvegarder les Contacts du Domaine</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"contacts": map[string]any{}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Obtenir les Détails du Revendeur

```go
result := client.GetResellerDetails()
```

<details>
<summary>Exemple de Sortie pour Obtenir les Détails du Revendeur</summary>

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

### Obtenir le Solde Actuel

```go
result := client.GetCurrentBalance("USD") // USD (par défaut), TRY, EUR ou GBP
```

<details>
<summary>Exemple de Sortie pour Obtenir le Solde Actuel</summary>

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

En cas d'échec : `ErrorCode` vaut `API_<status>` (ou `CURL_<n>`), `OperationMessage` contient le message d'erreur,
`OperationResult` vaut `"FAILED"` et `Balance` vaut `"0.00"`. Le compte ne renvoie que les soldes en USD et en TRY :
pour toute autre devise, le résultat est `FAILED` avec `ErrorCode: "BALANCE"` et
`No EUR balance is reported for this account. Available: USD, TRY.`

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Gestion de la Zone DNS

La passerelle stocke les enregistrements sous forme d'**ensembles** (un nom + un type avec toutes ses valeurs) et chaque écriture
remplace l'ensemble entier. Les noms d'enregistrement peuvent être indiqués sous la forme `@`, `""`, l'apex de la zone, un libellé relatif
(`www`) ou un nom complet (`www.example.com`) ; ils sont réduits au libellé relatif attendu par la passerelle.
Types pris en charge : `A`, `AAAA`, `ALIAS`, `CAA`, `CNAME`, `DS`, `MX`, `NAPTR`, `PTR`, `SRV`,
`TLSA`, `TXT` (les autres sont rejetés avant tout appel). Le TTL est borné à 1..86400 et le résultat indique
le type en majuscules et le TTL réellement envoyés. Avec `apply = true` (la valeur par défaut de php-dna), les modifications
sont appliquées immédiatement ; si cette application échoue, son enveloppe d'erreur est renvoyée (la modification
reste en attente). Les appels de zone consécutifs sont espacés d'au moins une seconde, car la passerelle limite le débit
des rafales.

```go
// Liste : une ligne par valeur, noms entièrement qualifiés
client.GetResourceRecords("example.com")

// Définir un nom+type avec exactement ces valeurs (une liste vide supprime l'ensemble). Valeurs par défaut : ttl 3600, apply true
client.SetResourceRecordSet("example.com", "@", "MX", []string{"10 mx1.example.com", "20 mx2.example.com"}, 3600, true)

// Ajouter / remplacer / supprimer UNE valeur en conservant les autres (ttl nil = conserver le TTL de l'ensemble)
client.AddResourceRecordValue("example.com", "www", "A", "192.0.2.10", nil, true)
client.ReplaceResourceRecordValue("example.com", "www", "A", "192.0.2.10", "192.0.2.11", nil, true)
client.RemoveResourceRecordValue("example.com", "www", "A", "192.0.2.11", true)

// Supprimer un nom+type entier (knownValue "" = rechercher d'abord une valeur actuelle)
client.DeleteResourceRecordSet("example.com", "www", "A", true, "")

// Appels de bas niveau (content est une string ou un []string)
client.AddResourceRecord("example.com", "api", "CNAME", "example.com.", 300, false)
client.EditResourceRecord("example.com", "api", "api", "CNAME", "www.example.com.", 300, false)
client.DeleteResourceRecord("example.com", "old.example.com.", "192.0.2.99", "A", false)
client.ApplyResourceRecords("example.com")
```

<details>
<summary>Exemples de Sortie pour la Gestion de la Zone DNS</summary>

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

// Valeur déjà présente / absente : rien n'est écrit
map[string]any{"result": "OK", "data": map[string]any{"Name": "", "Type": "MX", "Contents": []any{"..."}, "Changed": false}}

// DeleteResourceRecordSet d'un ensemble inexistant
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

### Redirection d'URL

```go
client.SetForwarding("example.com", "https://www.example.org", "Standard") // ou "Frame" (masquée)
client.GetForwarding("example.com")
client.DeleteForwarding("example.com")
```

<details>
<summary>Exemples de Sortie pour la Redirection d'URL</summary>

```go
// SetForwarding
map[string]any{"result": "OK", "data": map[string]any{"RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// GetForwarding
map[string]any{"result": "OK", "data": map[string]any{"Enabled": true, "RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// GetForwarding lorsqu'aucune redirection n'est configurée
map[string]any{"result": "OK", "data": map[string]any{"Enabled": false, "RedirectAddress": "", "ForwardType": ""}}

// DeleteForwarding (Removed vaut false lorsque rien n'était configuré)
map[string]any{"result": "OK", "data": map[string]any{"DomainName": "example.com", "Removed": true}}
```

</details>

Seule la réponse « Redirect information not found » de la passerelle signifie « non configurée » ; toute autre erreur,
par exemple un domaine inconnu, est renvoyée comme une erreur.

<hr style="border: 4px solid #000; border-style: dashed;">

### Utilitaires (sans réseau)

```go
client.IsTrTLD("example.com.tr") // true

client.ValidateContact(map[string]any{"FirstName": " Jo ", "Phone": "+90 (532) 374 33 88", "Country": "Turkey"})
// map[string]any{"FirstName": "Jo", "LastName": "Jo", "AddressLine1": "Addres yok", "City": "ISTANBUL",
//   "State": "ISTANBUL", "Country": "TR", "ZipCode": "34000", "Phone": "5323743388",
//   "PhoneCountryCode": "90", "Fax": "", "FaxCountryCode": ""}
```

`ValidateContact` remplit les champs manquants avec des valeurs de substitution, supprime les chiffres et la ponctuation des noms,
utilise la ville à défaut d'État, convertit les noms de pays courants en codes ISO et sépare les numéros de téléphone et de fax
en numéro et indicatif pays. Un numéro écrit avec son propre préfixe international
(`+44 20 ...`, `0044 20 ...`) conserve son propre indicatif pays, même si un autre est fourni.

<hr style="border: 4px solid #000; border-style: dashed;">

### Diagnostic

```go
client.GetRequestData()        // {"url": ..., "method": ..., "payload": ...} du dernier appel
client.GetResponseData()       // dernière réponse décodée, ou {"raw_response": body}
client.GetResponseHeaders()    // {"url", "http_code", "content_type", "total_time", "headers"}
client.GetParsedResponseData() // dernière réponse analysée, ou l'erreur {Code, Message, Details}
client.GetLastFunction()       // par ex. "request:GET domains/info"
client.GetServiceURL()         // URL de base de la passerelle
client.SetServiceURL("https://ote.domainresellerapi.com/api/v1")
```

Un `Client` conserve cet état propre à chaque appel ; ne partagez donc pas un même `Client` entre plusieurs goroutines : créez-en un
par goroutine.

## Gestion des erreurs

- Les erreurs de la passerelle renvoient `Code: "API_<HTTP status>"`. Lorsque la passerelle a envoyé son propre code d'erreur
  (par exemple `Dna.DomainService:Domain:10007`), celui-ci est ajouté sous `GatewayCode`, ce qui permet de distinguer
  « domaine introuvable » de « TLD non pris en charge ». Le message provient de la passerelle ; un corps d'erreur
  vide devient `HTTP 429 - Too Many Requests` (et ainsi de suite), les erreurs de validation sont listées champ par
  champ dans `Details` (`Admin PostalCode - 'PostalCode' is required`), et les messages génériques du backend tels que
  `Exception of type '...' was thrown` sont remplacés par les détails de la passerelle, ou par un message stable
  indiquant le code d'erreur — jamais par le corps brut de la réponse.
- Une réponse 301/302 signifie que les identifiants ont été refusés : `Message: "Invalid API credentials"`.
- Les échecs de transport (délai dépassé, DNS, TLS, connexion réinitialisée) sont des erreurs, jamais un succès vide. Un
  dépassement de délai renvoie `Code: "CURL_28"` et `Could not reach the DomainNameAPI service: the connection
  timed out. (...)` ; les autres échecs utilisent des numéros de style curl (`CURL_6` DNS, `CURL_7` connexion,
  `CURL_35` TLS, `CURL_56` réinitialisation, `CURL_60` certificat) et `Curl error during request: ...`.
  `API_<n>` est réservé aux statuts HTTP.
- Les réponses HTTP 200 avec `"success": false` à un enregistrement, un transfert ou un renouvellement sont des refus :
  `REGISTER_DOMAIN`, `TRANSFER_DOMAIN` ou `RENEW_DOMAIN` avec le message de la passerelle.
- Un corps de succès vide ou `null` renvoie `Code: "RESPONSE"`.
- Délais d'expiration : 30 secondes par appel, 10 secondes pour la connexion.

## Mode test

Passez `true` (ou `"true"`, `"1"`, `"on"`, `"yes"`) comme troisième argument du constructeur pour utiliser la passerelle
de test OTE `https://ote.domainresellerapi.com/api/v1`. La télémétrie en mode test est étiquetée
`environment: ote`.

## Télémétrie (rapports d'erreurs)

Lorsqu'un appel API échoue, la bibliothèque peut envoyer un rapport d'erreur au système de suivi des erreurs propre à DomainNameAPI
(Sentry). Une petite part aléatoire des appels réussis (4 %, plus tout appel de plus d'une seconde) est envoyée
sous forme d'enregistrement de performance. Les rapports sont **activés par défaut**.

> ℹ️ Le projet Sentry pour `DomainNameAPI-Go` n'est pas encore créé : `DefaultSentryDSN` est vide, donc
> rien n'est envoyé tant qu'une release ne le définit pas. `dna.WithSentryDSN(...)` remplace le DSN.

**Ce qui est envoyé :** code et message d'erreur ; chemin du point de terminaison et méthode HTTP ; durée de l'appel ; nom de la bibliothèque
(`DomainNameAPI-Go`), version de la bibliothèque et version de Go ; Reseller ID et nom de domaine ; l'adresse IP
publique du serveur (obtenue via `api.ipify.org`) ; le nom d'utilisateur du système d'exploitation ; la requête
et la réponse API, y compris les coordonnées des contacts. Seuls les identifiants sont masqués : codes EPP/auth, clés API,
mots de passe et jetons deviennent `***`.

**Ce qui n'est jamais envoyé :** codes EPP/auth, clés API, mots de passe ou jetons sous forme lisible ; arguments
des appels de fonction dans les traces de pile.

**Pourquoi et pour combien de temps :** les données servent uniquement à mesurer l'utilisation et à diagnostiquer les problèmes au plus tôt ;
elles sont conservées 30 jours sur le serveur Sentry propre à DomainNameAPI.

**Désactivez-la, ou n'envoyez que le minimum :**

```go
client := dna.New(resellerID, apiKey, false)
client.SetErrorReportingEnabled(false) // rien n'est envoyé

client.SetTelemetryMinimal(true) // uniquement code/message d'erreur, point de terminaison, durée, versions et Reseller ID
```

La télémétrie s'exécute en arrière-plan, ne bloque jamais l'appel longtemps et ne provoque jamais de panic dans votre code.

> 🔒 **Sécurité :** la vérification des certificats TLS est activée. Sur un hôte dont le bundle CA est obsolète, elle peut être
> désactivée avec `client.SetSSLVerification(false)` — non recommandé ; mettre à jour le bundle CA est
> la meilleure solution.

## Options avancées

```go
client := dna.New(resellerID, apiKey, false,
	dna.WithHTTPClient(myHTTPClient),      // tout Do(*http.Request) (*http.Response, error)
	dna.WithSentryDSN("https://key@host/1"),
	dna.WithTelemetryClient(myHTTPClient),
)
```

Avec un client HTTP personnalisé, configurez les délais d'expiration, les proxys et TLS sur celui-ci ; il ne doit pas suivre les redirections
(renvoyez `http.ErrUseLastResponse` depuis `CheckRedirect`).

## Tests

Go n'est pas requis sur l'hôte ; Docker suffit :

```bash
docker run --rm -v "$PWD":/src -w /src golang:1.23 go test ./...
docker run --rm -v "$PWD":/src -w /src golang:1.23 go vet ./...
```

La suite comprend des tests unitaires pour la couche HTTP, la correspondance des erreurs, la télémétrie et les utilitaires, ainsi qu'un
test de conformité qui rejoue les 134 fixtures php-dna de `testdata/conformance` (les requêtes sur le
réseau et les valeurs de retour doivent correspondre exactement à php-dna). Aucun test n'accède au réseau.

## Codes de Retour et d'Erreur avec Descriptions

| Code | Description | Détail |
|------|-------------|--------|
| 1000 | Commande complétée avec succès | Opération réussie. |
| 1001 | Commande complétée avec succès ; action en attente | Opération réussie. Cependant, l'opération est actuellement en file d'attente pour être finalisée. |
| 2003 | Paramètre requis manquant | Erreur de paramètre manquant. Par exemple : aucune entrée de téléphone dans les informations de contact. |
| 2105 | L'objet n'est pas éligible pour le renouvellement | Le statut du domaine n'est pas adapté pour le renouvellement, verrouillé pour les opérations de mise à jour. Le statut ne doit pas être "clientupdateprohibited". Peut être dû à d'autres conditions de statut. |
| 2200 | Erreur d'authentification | Erreur d'autorisation, le code de sécurité est incorrect ou le domaine est enregistré auprès d'un autre registraire. |
| 2302 | L'objet existe | Le nom de domaine ou les informations du serveur de noms existent déjà dans la base de données. Enregistrement impossible. |
| 2303 | L'objet n'existe pas | Le nom de domaine ou les informations du serveur de noms n'existent pas dans la base de données. Un nouvel enregistrement est nécessaire. |
| 2304 | Le statut de l'objet interdit l'opération | Le statut du domaine n'est pas adapté aux mises à jour, verrouillé pour les opérations de mise à jour. Le statut ne doit pas être "clientupdateprohibited". Peut être dû à d'autres conditions de statut. |

Codes de la bibliothèque : `RESPONSE` (réponse vide), `CREDENTIALS` (format de réponse invalide du point de terminaison
du compte), `AVAILABILITY` (aucun nom de domaine fourni), `BALANCE` (devise non renvoyée par le compte),
`TRANSFER_NOT_AVAILABLE`, `REGISTER_DOMAIN` / `TRANSFER_DOMAIN` / `RENEW_DOMAIN` (la passerelle
a refusé l'opération), `DOMAIN_RENEW` (renouvellement sans date d'expiration), `ADD_ZONE`, `EDIT_ZONE`,
`SET_ZONE`, `DELETE_ZONE` (type d'enregistrement non pris en charge ou valeur introuvable), `SET_FORWARD` (type
de redirection non pris en charge ou adresse vide), `CURL_<n>` (échec de transport), `API_<status>` (erreur HTTP, avec
`GatewayCode` lorsque la passerelle a envoyé son propre code).

## Licence

Apache License 2.0 — voir [LICENSE](LICENSE).
