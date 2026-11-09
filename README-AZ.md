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

# go-dna — Go üçün DomainNameAPI klienti

DomainNameAPI domen qeydiyyatçısının REST API-si üçün Go klienti. Bu,
[php-dna](https://github.com/domainreseller/php-dna) kitabxanasının dəqiq portudur: eyni metodlar,
şəbəkə üzərindən eyni sorğular və eyni qaytarılan dəyərlər. Yeganə asılılıq standart kitabxanadır.

## 📦 Yükləmə — həmişə Releases istifadə edin!

⬇️ **Ən son test edilmiş versiyanı buradan əldə edin: https://github.com/domainreseller/go-dna/releases/latest**

> ⚠️ Yaşıl **Code → Download ZIP** düyməsindən istifadə **etməyin** — o, xam inkişaf budağını yükləyir. Buraxılış paketləri versiyalanmış, test edilmiş və istehsala hazırdır.

## Tələblər

- Go 1.21 və ya daha yüksək
- Üçüncü tərəf paketləri yoxdur (yalnız standart kitabxanadan `net/http` və `encoding/json`)

## Quraşdırma

Go modules ilə buraxılmış bir versiyanı sabitləyin:

```bash
go get github.com/domainreseller/go-dna@v1.0.0
```

```go
import dna "github.com/domainreseller/go-dna"
```

Paketin adı `dna`-dır. Hər metod üçün işlədilə bilən nümunələr [examples](examples) qovluğundadır:

```bash
go run ./examples/GetDetails
```

## 🔑 API Məlumatları — Reseller ID və API Key

Kitabxanaya konstruktora ötürülən dəqiq iki giriş məlumatı lazımdır:

| Parametr | Dəyər |
|---|---|
| Birinci | Reseller ID — `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` kimi UUID |
| İkinci | API Key |

> 💡 **Reseller ID** və **API Key** məlumatlarınızı DomainNameAPI panelinizdə **API Parametrləri** bölməsində tapa bilərsiniz.
> ⚠️ Bunlar **API giriş məlumatlarıdır** — panelə giriş məlumatlarınız burada **işləməyəcək**.

```go
resellerID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
apiKey := "your-api-key"

// Canlı mühit
client := dna.New(resellerID, apiKey, false)

// Test (OTE) mühiti
client = dna.New(resellerID, apiKey, true)
```

Hər iki giriş məlumatının əvvəlindəki və sonundakı boşluqlar silinir. Üçüncü arqument panellərin
onay qutuları üçün saxladığı sətirləri də qəbul edir: `"true"`, `"1"`, `"on"` və `"yes"` (böyük-kiçik
hərfdən asılı olmayaraq) test mühitini seçir, qalan hər şey (`"false"`, `"0"`, `"off"`, `""`, ...)
canlı mühiti seçir.

## Sürətli başlanğıc

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

## Qaytarılan dəyərlər

Metodlar tipli modellər deyil, php-dna-nın qaytardığı eyni dinamik strukturları — `map[string]any`
və `[]any` — qaytarır; beləliklə açarlar, sətir bayraqları (`"true"`/`"false"`) və formatlanmış pul
sətirləri (`"11.9900"`) php-dna-dakı ilə tam eynidir. Gateway-dən gələn ədədlər `int64` və ya
`float64` tipindədir.

Uğurlu çağırışlar belə görünür:

```go
map[string]any{"result": "OK", "data": map[string]any{ /* ... */ }}
```

Uğursuz çağırışlar heç vaxt panic yaratmır; onlar xəta zərfini qaytarır:

```go
map[string]any{
	"result": "ERROR",
	"error": map[string]any{
		"Code":    "API_404",      // API_<HTTP statusu> və ya RESPONSE kimi kitabxana kodu
		"Message": "The domain name you requested could not be found. Please check and try again.",
		"Details": "The domain could not be found. Please check and try again.",
	},
}
```

Üç metod php-dna-da olduğu kimi fərqli formadan istifadə edir:

| Metod | Uğur | Uğursuzluq |
|---|---|---|
| `CheckAvailability` | sadə `[]any` siyahısı | xəta zərfi (`map[string]any`), buna görə qaytarılan tip `any`-dir |
| `GetCurrentBalance` | `OperationResult: "SUCCESS"` zərfi | `OperationResult: "FAILED"` olan eyni zərf |
| `GetResellerDetails` | düz map (`id`, `name`, `balance`, `balances`, ...) | xəta zərfi |

### İstəyə bağlı arqumentlər

Go-da standart arqumentlər yoxdur. Metodlar php-dna-nın parametr sırasını saxlayır və hər arqumenti
açıq şəkildə qəbul edir; php-dna-nın standart dəyəri aşağıdakı hər bölmədə göstərilib. Rahatlıq üçün
bəzi sıfır dəyərlər standart dəyəri seçir: `GetCurrentBalance("")` USD, `GetTldList(0)` 20,
`CheckAvailability(..., "")` `create`, `SetForwarding(..., "")` `Standard`,
`ModifyPrivacyProtectionStatus(..., "")` isə `Owner request` deməkdir.
`RegisterWithContactInfo` dörd istəyə bağlı arqumentini `*RegisterOptions` daxilində qəbul edir
(nil = standart dəyərlər). `AddResourceRecordValue` və `ReplaceResourceRecordValue` `ttl *int` qəbul
edir: `nil` dəstin cari TTL-ini saxlayır.

### Əlaqələr

Əlaqələr php-dna açarları (`FirstName`, `LastName`, `Company`, `EMail`, `AddressLine1`,
`AddressLine2`, `City`, `State`, `Country`, `ZipCode`, `Phone`, `PhoneCountryCode`, `Fax`,
`FaxCountryCode`) olan `map[string]any`-dir. İç-içə forma (`Address.Line1`,
`Phone.Phone.Number`, ...) da qəbul edilir. Go map-lərində sıra olmadığı üçün dörd rol sıralı
`dna.Contacts` siyahısı kimi ötürülür:

```go
contacts := dna.Contacts{
	{Role: "Administrative", Contact: contact},
	{Role: "Billing", Contact: contact},
	{Role: "Technical", Contact: contact},
	{Role: "Registrant", Contact: contact},
}
contacts = dna.AllRoles(contact)                 // dörd rolun hamısı üçün eyni əlaqə
contacts = dna.ContactsFromMap(map[string]map[string]any{ /* rol => əlaqə */ })
```

Yarımçıq doldurulmuş əlaqələr göndərilməzdən əvvəl `ValidateContact` qaydalarına əsasən tamamlanır.

### Metod adları

Go metod adları Go konvensiyalarına uyğundur; php-dna adı lowerCamelCase formatında eynidir.
`getServiceUrl`/`setServiceUrl` `GetServiceURL`/`SetServiceURL`, `setSslVerification` isə
`SetSSLVerification` adlanır.

## İstifadə

### Domen Qeydiyyatı

php-dna standart dəyərləri: standart nameserver-lər (`ns1.domainnameapi.com`,
`ns2.domainnameapi.com`), EPP kilidi aktiv, məxfilik deaktiv, əlavə atributlar yoxdur. Əlavə
atributlar yalnız `.tr` kimi domenlər üçün tələb olunur.

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
// Yalnız .tr domenləri üçün:
opts.AdditionalAttributes = map[string]any{
	"TRABISDOMAINCATEGORY": 1,             // 1 = fərdi şəxs (Owner), 0 = şirkət (Corporate)
	"TRABISCITIZIENID":     "12345678901", // və ya şirkətlər üçün TRABISTAXOFFICE + TRABISTAXNUMBER
}

result := client.RegisterWithContactInfo("example.com", 1, dna.AllRoles(contact), opts)
```

`.tr` domenləri üçün köhnə TRABIS adları gateway sxeminə uyğun normallaşdırılır: kateqoriya
`Owner`/`Corporate` olur, `TRABISCITIZIENID` `TRABISCITIZENID`-ə çevrilir, ölkə, şəhər, ad və təşkilat
sahələri isə çıxarılır (gateway onları əlaqədən götürür). Hər atribut dəyəri sətir kimi göndərilir.

<details>
<summary>Domen Qeydiyyatı üçün Nümunə Çıxış</summary>

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

Qeydiyyat endpoint-i qısa bir qeydlə cavab verir, buna görə sahələrin çoxu boşdur; tam domen
məlumatları üçün sonra `GetDetails` çağırın.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domen Yenilənməsi

```go
result := client.Renew("example.com", 1)
```

<details>
<summary>Domen Yenilənməsi üçün Nümunə Çıxış</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"ExpirationDate": "2026-03-15T10:00:00"}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domen Transferi

```go
result := client.Transfer("example.com", "Xy9#mK2$", 1, nil)
```

php-dna-nın fasadında olduğu kimi, əlaqələr arqumenti ötürülmür: transfer sorğusu həmişə boş əlaqə
siyahısı daşıyır. Nəticə `GetDetails` ilə eyni formadadır.

<details>
<summary>Domen Transferi üçün Nümunə Çıxış</summary>

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

Əlaqəli transfer çağırışları:

```go
client.CancelTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Cancelled"}}
client.ApproveTransfer("example.com") // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Approved"}}
client.RejectTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Rejected"}}
```

<hr style="border: 4px solid #000; border-style: dashed;">

### Transfer Yoxlaması

```go
result := client.CheckTransfer("example.com", "Ab1!xyZ9")
```

`result` yalnız registry domenin transfer edilə bilən olduğunu bildirdikdə `OK` olur; əks halda
`TRANSFER_NOT_AVAILABLE` kodu və eyni `data` ilə `ERROR` olur.

<details>
<summary>Transfer Yoxlaması üçün Nümunə Çıxış</summary>

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

### Domen Siyahısı

```go
result := client.GetList(map[string]any{
	"OrderColumn":    "Id",
	"OrderDirection": "ASC",
	"PageNumber":     0,
	"PageSize":       1000,
})
```

`PageNumber` (0-dan başlayır), `PageSize`, `OrderColumn` və `OrderDirection` gateway-in `SkipCount`,
`MaxResultCount` və `Sorting` parametrlərinə çevrilir; digər açarlar (məsələn, `DomainName`) filtr
kimi ötürülür. `nil` ilk 200 domeni qaytarır.

<details>
<summary>Domen Siyahısı üçün Nümunə Çıxış</summary>

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

### TLD Siyahısı

```go
result := client.GetTldList(100) // php-dna standart dəyəri: 20
```

<details>
<summary>TLD Siyahısı üçün Nümunə Çıxış</summary>

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

### Domen Mövcudluğunun Yoxlanması

```go
result := client.CheckAvailability([]string{"hello", "world123x0"}, []string{"com", "net"}, 1, "create")
if list, ok := result.([]any); ok {
	// uğurlu
	_ = list
}
```

<details>
<summary>Domen Mövcudluğunun Yoxlanması üçün Nümunə Çıxış</summary>

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

Premium domenlər üçün `IsFee` `true` olur.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domen Məlumatları

```go
result := client.GetDetails("example.com")
result = client.SyncFromRegistry("example.com") // eyni forma
```

<details>
<summary>Domen Məlumatları üçün Nümunə Çıxış</summary>

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

Gateway rəqəmsal enum göndərsə belə, `Status` mətn etiketidir (`Active`, `TransferredOut`,
`PendingDelete`, ...). Təhlil edilə bilməyən tarixlər `""` olur, heç vaxt yanlış 1970 tarixi olmur.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Nameserver Dəyişdirmə

```go
result := client.ModifyNameServer("example.com", []string{"ns1.example.com", "ns2.example.com"})
```

Map də qəbul edilir (onun dəyərləri açar sırası ilə istifadə olunur). Göndərilməzdən əvvəl elementlərin
boşluqları silinir, kiçik hərfə çevrilir, sondakı nöqtə götürülür və təkrarlar aradan qaldırılır.

<details>
<summary>Nameserver Dəyişdirmə üçün Nümunə Çıxış</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServers": []string{"ns1.example.com", "ns2.example.com"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domen Kilidini Aktivləşdirmə

```go
result := client.EnableTheftProtectionLock("example.com")
```

<details>
<summary>Domen Kilidini Aktivləşdirmə üçün Nümunə Çıxış</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domen Kilidini Deaktivləşdirmə

```go
result := client.DisableTheftProtectionLock("example.com")
```

<details>
<summary>Domen Kilidini Deaktivləşdirmə üçün Nümunə Çıxış</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": false}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Child Nameserver Əlavə Etmək

```go
result := client.AddChildNameServer("example.com", "ns1.example.com", "1.2.3.4")
```

<details>
<summary>Child Nameserver Əlavə Etmək üçün Nümunə Çıxış</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.4"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Child Nameserver Silmə

```go
result := client.DeleteChildNameServer("example.com", "ns1.example.com")
```

<details>
<summary>Child Nameserver Silmə üçün Nümunə Çıxış</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com"}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Child Nameserver Yeniləmə

```go
result := client.ModifyChildNameServer("example.com", "ns1.example.com", "1.2.3.5")
```

<details>
<summary>Child Nameserver Yeniləmə üçün Nümunə Çıxış</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.5"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domen Məxfiliyini Dəyişdirmə

```go
result := client.ModifyPrivacyProtectionStatus("example.com", true, "Owner request")
```

<details>
<summary>Domen Məxfiliyini Dəyişdirmə üçün Nümunə Çıxış</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"PrivacyProtectionStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domen Əlaqələrini Əldə Etmək

```go
result := client.GetContacts("example.com")
```

<details>
<summary>Domen Əlaqələrini Əldə Etmək üçün Nümunə Çıxış</summary>

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
			// "Billing", "Registrant", "Technical": eyni forma
		},
	},
}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domen Əlaqələrini Saxlamaq

```go
result := client.SaveContacts("example.com", dna.AllRoles(contact))
```

<details>
<summary>Domen Əlaqələrini Saxlamaq üçün Nümunə Çıxış</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"contacts": map[string]any{}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Satıcı Təfərrüatlarını Əldə Etmək

```go
result := client.GetResellerDetails()
```

<details>
<summary>Satıcı Təfərrüatlarını Əldə Etmək üçün Nümunə Çıxış</summary>

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

### Cari Balansı Əldə Etmək

```go
result := client.GetCurrentBalance("USD") // USD (standart), TRY, EUR və ya GBP
```

<details>
<summary>Cari Balansı Əldə Etmək üçün Nümunə Çıxış</summary>

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

Uğursuzluq halında: `ErrorCode` `API_<status>`, `OperationMessage` xəta mesajı,
`OperationResult` `"FAILED"`, `Balance` isə `"0.00"` olur.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### DNS Zonasının İdarə Edilməsi

Gateway qeydləri **dəstlər** şəklində saxlayır (bütün dəyərləri ilə birlikdə bir ad + tip) və hər
yazma əməliyyatı bütün dəsti əvəz edir. Qeyd adları `@`, `""`, zona apeksi, nisbi etiket (`www`) və ya
tam ad (`www.example.com`) kimi verilə bilər; onlar gateway-in gözlədiyi nisbi etiketə endirilir.
Dəstəklənən tiplər: `A`, `AAAA`, `ALIAS`, `CAA`, `CNAME`, `DS`, `MX`, `NAPTR`, `PTR`, `SRV`, `TLSA`,
`TXT` (digərləri heç bir çağırış edilməzdən əvvəl rədd edilir). TTL 1..86400 aralığına məhdudlaşdırılır.
`apply = true` (php-dna-nın standart dəyəri) ilə dəyişikliklər dərhal tətbiq olunur; gateway ardıcıl
sorğu partlayışlarına sürət limiti tətbiq etdiyi üçün ardıcıl zona çağırışları arasında ən azı bir
saniyə fasilə saxlanılır.

```go
// Siyahı: hər dəyər üçün bir sətir, adlar tam ixtisaslı
client.GetResourceRecords("example.com")

// Ad+tipi dəqiq bu dəyərlərə təyin edin (boş siyahı dəsti silir). Standart: ttl 3600, apply true
client.SetResourceRecordSet("example.com", "@", "MX", []string{"10 mx1.example.com", "20 mx2.example.com"}, 3600, true)

// BİR dəyəri əlavə edin / əvəz edin / silin, digərlərini saxlayın (ttl nil = dəstin TTL-ini saxla)
client.AddResourceRecordValue("example.com", "www", "A", "192.0.2.10", nil, true)
client.ReplaceResourceRecordValue("example.com", "www", "A", "192.0.2.10", "192.0.2.11", nil, true)
client.RemoveResourceRecordValue("example.com", "www", "A", "192.0.2.11", true)

// Bütöv ad+tipi silin (knownValue "" = əvvəlcə cari dəyəri axtar)
client.DeleteResourceRecordSet("example.com", "www", "A", true, "")

// Aşağı səviyyəli çağırışlar (content sətir və ya []string-dir)
client.AddResourceRecord("example.com", "api", "CNAME", "example.com.", 300, false)
client.EditResourceRecord("example.com", "api", "api", "CNAME", "www.example.com.", 300, false)
client.DeleteResourceRecord("example.com", "old.example.com.", "192.0.2.99", "A", false)
client.ApplyResourceRecords("example.com")
```

<details>
<summary>DNS Zonasının İdarə Edilməsi üçün Nümunə Çıxışlar</summary>

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

// Dəyər artıq mövcuddur / mövcud deyil: heç nə yazılmır
map[string]any{"result": "OK", "data": map[string]any{"Name": "", "Type": "MX", "Contents": []any{"..."}, "Changed": false}}

// Mövcud olmayan dəst üçün DeleteResourceRecordSet
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

### URL Yönləndirməsi

```go
client.SetForwarding("example.com", "https://www.example.org", "Standard") // və ya "Frame" (maskalanmış)
client.GetForwarding("example.com")
client.DeleteForwarding("example.com")
```

<details>
<summary>URL Yönləndirməsi üçün Nümunə Çıxışlar</summary>

```go
// SetForwarding
map[string]any{"result": "OK", "data": map[string]any{"RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// GetForwarding (yönləndirmə qurulmayıb: Enabled false, yenə də OK)
map[string]any{"result": "OK", "data": map[string]any{"Enabled": true, "RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// DeleteForwarding
map[string]any{"result": "OK", "data": map[string]any{"DomainName": "example.com"}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Köməkçi funksiyalar (şəbəkəsiz)

```go
client.IsTrTLD("example.com.tr") // true

client.ValidateContact(map[string]any{"FirstName": " Jo ", "Phone": "+90 (532) 374 33 88", "Country": "Turkey"})
// map[string]any{"FirstName": "Jo", "LastName": "Jo", "AddressLine1": "Addres yok", "City": "ISTANBUL",
//   "State": "ISTANBUL", "Country": "TR", "ZipCode": "34000", "Phone": "5323743388",
//   "PhoneCountryCode": "90", "Fax": "", "FaxCountryCode": ""}
```

`ValidateContact` çatışmayan sahələri yer tutucularla doldurur, adlardan rəqəmləri və durğu işarələrini
silir, State boş olduqda City-dən istifadə edir, geniş yayılmış ölkə adlarını ISO kodlarına çevirir
və telefon və faks nömrələrini nömrə və ölkə koduna ayırır.

<hr style="border: 4px solid #000; border-style: dashed;">

### Diaqnostika

```go
client.GetRequestData()        // son çağırışın {"url": ..., "method": ..., "payload": ...} məlumatı
client.GetResponseData()       // son deşifrə edilmiş cavab və ya {"raw_response": body}
client.GetResponseHeaders()    // {"url", "http_code", "content_type", "total_time", "headers"}
client.GetParsedResponseData() // son təhlil edilmiş cavab və ya {Code, Message, Details} xətası
client.GetLastFunction()       // məsələn, "request:GET domains/info"
client.GetServiceURL()         // gateway-in əsas URL-i
client.SetServiceURL("https://ote.domainresellerapi.com/api/v1")
```

`Client` bu vəziyyəti hər çağırış üçün saxlayır, buna görə bir `Client`-i goroutine-lər arasında
paylaşmayın; hər goroutine üçün ayrıca bir dənə yaradın.

## Xətaların idarə edilməsi

- Gateway xətaları `Code: "API_<HTTP status>"` qaytarır. Mesaj gateway-dən gəlir; boş xəta gövdəsi
  `HTTP 429 - Too Many Requests` (və s.) olur, validasiya xətaları `Details` daxilində sahə-sahə
  sadalanır (`Admin PostalCode - 'PostalCode' is required`), `Exception of type '...' was thrown`
  kimi backend şablon mətnləri isə faydalı təfərrüatlarla əvəz olunur.
- 301/302 cavabı giriş məlumatlarının rədd edildiyini bildirir: `Message: "Invalid API credentials"`.
- Nəqliyyat xətaları (vaxt aşımı, DNS, TLS, bağlantının qırılması) heç vaxt boş uğur deyil, xətadır.
  Vaxt aşımı `Code: "API_28"` və `Could not reach the DomainNameAPI service: the connection
  timed out. (...)` qaytarır; digər xətalar curl üslublu nömrələrdən (`API_6` DNS, `API_7` bağlantı,
  `API_35` TLS, `API_56` bağlantının qırılması, `API_60` sertifikat) və `Curl error during request: ...`
  mətnindən istifadə edir. `GetParsedResponseData()["Code"]` `CURL_<n>` dəyərini saxlayır.
- Boş və ya `null` uğurlu cavab gövdəsi `Code: "RESPONSE"` qaytarır.
- Vaxt aşımları: hər çağırış üçün 30 saniyə, bağlantı qurmaq üçün 10 saniyə.

## Test rejimi

OTE test gateway-indən (`https://ote.domainresellerapi.com/api/v1`) istifadə etmək üçün üçüncü
konstruktor arqumenti kimi `true` (və ya `"true"`, `"1"`, `"on"`, `"yes"`) ötürün. Test rejimindəki
telemetriya `environment: ote` ilə etiketlənir.

## Telemetriya (xəta hesabatları)

API çağırışı uğursuz olduqda kitabxana DomainNameAPI-nin öz xəta izləyicisinə (Sentry) xəta hesabatı
göndərə bilər. Uğurlu çağırışların kiçik təsadüfi bir hissəsi (4%, üstəlik 1 saniyədən uzun çəkən hər
çağırış) performans qeydi kimi göndərilir. Hesabat **standart olaraq aktivdir**.

> ℹ️ `DomainNameAPI-Go` üçün Sentry layihəsi hələ yaradılmayıb: `DefaultSentryDSN` boşdur, buna görə
> bir buraxılış onu təyin edənə qədər heç nə göndərilmir. `dna.WithSentryDSN(...)` DSN-i əvəz edir.

**Nə göndərilir:** xəta kodu və mesajı; endpoint yolu və HTTP metodu; çağırışın müddəti; kitabxananın
adı (`DomainNameAPI-Go`), kitabxana versiyası və Go versiyası; reseller ID və domen adı; serverin
ictimai IP ünvanı (`api.ipify.org` vasitəsilə müəyyən edilir); əməliyyat sisteminin istifadəçi adı;
əlaqə məlumatları daxil olmaqla API sorğusu və cavabı. Yalnız giriş məlumatları maskalanır:
EPP/avtorizasiya kodları, API açarları, şifrələr və tokenlər `***` olur.

**Heç vaxt nə göndərilmir:** oxuna bilən formada EPP/avtorizasiya kodları, API açarları, şifrələr və
ya tokenlər; stack trace-lərdə funksiya çağırışlarının arqumentləri.

**Niyə və nə qədər müddətə:** məlumatlar yalnız istifadənin ölçülməsi və problemlərin erkən
diaqnostikası üçün istifadə olunur və DomainNameAPI-nin öz Sentry serverində 30 gün saxlanılır.

**Söndürmək və ya yalnız minimumu göndərmək:**

```go
client := dna.New(resellerID, apiKey, false)
client.SetErrorReportingEnabled(false) // heç nə göndərilmir

client.SetTelemetryMinimal(true) // yalnız xəta kodu/mesajı, endpoint, müddət, versiyalar və reseller ID
```

Telemetriya arxa planda işləyir, çağırışı heç vaxt uzun müddət bloklamır və kodunuzda heç vaxt panic
yaratmır.

> 🔒 **Təhlükəsizlik:** TLS sertifikat yoxlaması aktivdir. Köhnəlmiş CA paketi olan hostda onu
> `client.SetSSLVerification(false)` ilə söndürmək olar — tövsiyə edilmir; CA paketini yeniləmək daha
> yaxşı həlldir.

## Qabaqcıl seçimlər

```go
client := dna.New(resellerID, apiKey, false,
	dna.WithHTTPClient(myHTTPClient),      // istənilən Do(*http.Request) (*http.Response, error)
	dna.WithSentryDSN("https://key@host/1"),
	dna.WithTelemetryClient(myHTTPClient),
)
```

Xüsusi HTTP klienti istifadə etdikdə vaxt aşımlarını, proxy-ləri və TLS-i həmin klientdə konfiqurasiya
edin; o, yönləndirmələri izləməməlidir (`CheckRedirect`-dən `http.ErrUseLastResponse` qaytarın).

## Testlər

Hostda Go quraşdırılmalı deyil; Docker kifayətdir:

```bash
docker run --rm -v "$PWD":/src -w /src golang:1.23 go test ./...
docker run --rm -v "$PWD":/src -w /src golang:1.23 go vet ./...
```

Test dəsti HTTP təbəqəsi, xətaların uyğunlaşdırılması, telemetriya və köməkçi funksiyalar üçün vahid
testlərdən, həmçinin `testdata/conformance` qovluğundakı bütün 126 php-dna fixture-unu yenidən
icra edən uyğunluq testindən ibarətdir (şəbəkə üzərindəki sorğular və qaytarılan dəyərlər php-dna ilə
tam eyni olmalıdır). Heç bir test şəbəkəyə müraciət etmir.

## Geri Dönüş və Xəta Kodları ilə Açıqlamaları

| Kod  | Açıqlama                                       | Təfərrüat                                                                                                                                                      |
|------|------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 1000 | Əmr uğurla tamamlandı                          | Əməliyyat uğurlu oldu.                                                                                                                                         |
| 1001 | Əmr uğurla tamamlandı; hərəkət gözləmədədir    | Əməliyyat uğurlu oldu. Bununla belə, əməliyyat hazırda tamamlanmaq üçün növbəyə alınıb.                                                                        |
| 2003 | Tələb olunan parametr çatışmır                 | Parametr çatışmazlığı xətası. Məsələn, əlaqə məlumatlarında telefon nömrəsi yoxdur.                                                                            |
| 2105 | Obyekt yenilənmə üçün uyğun deyil              | Domen statusu yenilənmə üçün uyğun deyil, yeniləmə əməliyyatları üçün kilidlənib. Status "clientupdateprohibited" olmamalıdır. Digər status şərtlərindən qaynaqlana bilər. |
| 2200 | Kimlik doğrulama xətası                        | İcazə xətası, təhlükəsizlik kodu səhvdir və ya domen başqa qeydiyyatçıda qeydiyyatdadır.                                                                       |
| 2302 | Obyekt mövcuddur                               | Domen adı və ya nameserver məlumatları verilənlər bazasında artıq mövcuddur. Qeydiyyatdan keçirilə bilməz.                                                     |
| 2303 | Obyekt mövcud deyil                            | Domen adı və ya nameserver məlumatları verilənlər bazasında mövcud deyil. Yeni qeydiyyat tələb olunur.                                                         |
| 2304 | Obyekt statusu əməliyyatı qadağan edir         | Domen statusu yeniləmələr üçün uyğun deyil, yeniləmə əməliyyatları üçün kilidlənib. Status "clientupdateprohibited" olmamalıdır. Digər status şərtlərindən qaynaqlana bilər. |

Kitabxana kodları: `RESPONSE` (boş cavab), `CREDENTIALS` (hesab endpoint-indən yanlış cavab formatı),
`AVAILABILITY` (domen adı verilməyib), `TRANSFER_NOT_AVAILABLE`, `DOMAIN_RENEW` (bitmə tarixi olmadan
yenilənmə), `ADD_ZONE`, `EDIT_ZONE`, `SET_ZONE`, `DELETE_ZONE` (dəstəklənməyən qeyd tipi və ya dəyər
tapılmadı), `SET_FORWARD` (dəstəklənməyən yönləndirmə tipi və ya boş ünvan).

## Lisenziya

Apache License 2.0 — bax: [LICENSE](LICENSE).
