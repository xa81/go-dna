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

# go-dna — Go için DomainNameAPI istemcisi

DomainNameAPI alan adı kayıt (registrar) REST API'si için Go istemcisi. [php-dna](https://github.com/domainreseller/php-dna)
kütüphanesinin birebir bir uyarlamasıdır: aynı metotlar, ağa giden aynı istekler ve aynı dönüş
değerleri. Tek bağımlılığı standart kütüphanedir.

## 📦 İndirme — her zaman Releases kullanın!

⬇️ **En güncel test edilmiş sürümü buradan indirin: https://github.com/domainreseller/go-dna/releases/latest**

> ⚠️ Yeşil **Code → Download ZIP** butonunu **kullanmayın** — bu buton ham geliştirme dalını (development branch) indirir. Release paketleri sürümlenmiş, test edilmiş ve canlı kullanıma hazırdır.

## Gereksinimler

- Go 1.21 veya daha üstü
- Üçüncü taraf paket yok (yalnızca standart kütüphaneden `net/http` ve `encoding/json`)

## Kurulum

Go modülleriyle, yayınlanmış bir sürümü sabitleyin:

```bash
go get github.com/domainreseller/go-dna@v1.0.0
```

```go
import dna "github.com/domainreseller/go-dna"
```

Paket adı `dna`'dır. Her metot için çalıştırılabilir örnekler [examples](examples) klasöründe yer alır:

```bash
go run ./examples/GetDetails
```

## 🔑 API Bilgileri — Reseller ID ve API Key

Kütüphane, constructor'a verilen tam olarak iki bilgiye ihtiyaç duyar:

| Parametre | Değer |
|---|---|
| Birinci | Reseller ID — `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` biçiminde UUID |
| İkinci | API Key |

> 💡 **Reseller ID** ve **API Key** bilgilerinizi DomainNameAPI panelinizdeki **API Ayarları** bölümünde bulabilirsiniz.
> ⚠️ Bunlar **API bilgileridir** — panel giriş bilgileriniz burada **ÇALIŞMAZ**.

```go
resellerID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
apiKey := "your-api-key"

// Canlı ortam
client := dna.New(resellerID, apiKey, false)

// Test (OTE) ortamı
client = dna.New(resellerID, apiKey, true)
```

Her iki bilginin de baştaki ve sondaki boşlukları temizlenir. Üçüncü argüman, panellerin onay kutuları
için sakladığı metin değerlerini de kabul eder: `"true"`, `"1"`, `"on"` ve `"yes"` (büyük/küçük harf
fark etmez) test ortamını seçer; diğer her şey (`"false"`, `"0"`, `"off"`, `""`, ...) canlı ortamı seçer.

## Hızlı başlangıç

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

## Dönüş değerleri

Metotlar, php-dna'nın döndürdüğü dinamik yapıların aynısını döndürür — tipli modeller değil,
`map[string]any` ve `[]any` — böylece anahtarlar, metin bayrakları (`"true"`/`"false"`) ve
biçimlendirilmiş tutar metinleri (`"11.9900"`) php-dna'dakiyle birebir aynıdır. Gateway'den gelen
sayılar `int64` veya `float64` türündedir.

Başarılı çağrılar şöyle görünür:

```go
map[string]any{"result": "OK", "data": map[string]any{ /* ... */ }}
```

Başarısız çağrılar asla panic oluşturmaz; hata zarfını (error envelope) döndürür:

```go
map[string]any{
	"result": "ERROR",
	"error": map[string]any{
		"Code":        "API_404", // API_<HTTP durum kodu>, CURL_<n> veya RESPONSE gibi bir kütüphane kodu
		"Message":     "The domain name you requested could not be found. Please check and try again.",
		"Details":     "The domain could not be found. Please check and try again.",
		"GatewayCode": "Dna.DomainService:Domain:10007", // yalnızca API_<status> ile, gateway kendi kodunu gönderdiğinde
	},
}
```

Üç metot, tıpkı php-dna'da olduğu gibi farklı bir yapı kullanır:

| Metot | Başarılı | Başarısız |
|---|---|---|
| `CheckAvailability` | düz bir `[]any` listesi | hata zarfı (`map[string]any`); bu yüzden dönüş tipi `any`'dir |
| `GetCurrentBalance` | `OperationResult: "SUCCESS"` zarfı | `OperationResult: "FAILED"` içeren aynı zarf |
| `GetResellerDetails` | düz map (`id`, `name`, `balance`, `balances`, ...) | hata zarfı |

### İsteğe bağlı argümanlar

Go'da varsayılan argüman yoktur. Metotlar php-dna'nın parametre sırasını korur ve her argümanı açıkça
alır; php-dna'daki varsayılan değer aşağıdaki her bölümde gösterilmiştir. Kolaylık olsun diye bazı sıfır
değerler varsayılanı seçer: `GetCurrentBalance("")` USD, `GetTldList(0)` 20,
`CheckAvailability(..., "")` `create` ve `SetForwarding(..., "")` `Standard` anlamına gelir.
`RegisterWithContactInfo`, dört isteğe bağlı argümanını bir `*RegisterOptions` içinde alır (nil = varsayılanlar).
`AddResourceRecordValue` ve `ReplaceResourceRecordValue` `ttl *int` alır: `nil`, kayıt kümesinin
mevcut TTL değerini korur.

### Kişiler (Contacts)

Kişiler, php-dna'nın anahtarlarını kullanan `map[string]any` değerleridir (`FirstName`, `LastName`,
`Company`, `EMail`, `AddressLine1`, `AddressLine2`, `City`, `State`, `Country`, `ZipCode`, `Phone`,
`PhoneCountryCode`, `Fax`, `FaxCountryCode`). İç içe yapı (`Address.Line1`,
`Phone.Phone.Number`, ...) da kabul edilir. Go map'lerinde sıra bulunmadığından dört rol, sıralı bir
`dna.Contacts` listesi olarak verilir:

```go
contacts := dna.Contacts{
	{Role: "Administrative", Contact: contact},
	{Role: "Billing", Contact: contact},
	{Role: "Technical", Contact: contact},
	{Role: "Registrant", Contact: contact},
}
contacts = dna.AllRoles(contact)                 // dört rolün tamamı için aynı kişi
contacts = dna.ContactsFromMap(map[string]map[string]any{ /* rol => kişi */ })
```

Yarı boş kişi bilgileri gönderilmeden önce `ValidateContact` kurallarıyla tamamlanır. Kişi bilgileri
`RegisterWithContactInfo` ve `SaveContacts` tarafından kullanılır; transferler asla kişi bilgisi göndermez.

### Metot adları

Go metot adları Go kurallarına uyar; php-dna'daki ad, lowerCamelCase biçiminde aynıdır.
`getServiceUrl`/`setServiceUrl`, `GetServiceURL`/`SetServiceURL` olur; `setSslVerification` ise
`SetSSLVerification` olur.

## Kullanım

### Domain Kayıt işlemleri

php-dna varsayılanları: varsayılan nameserver'lar (`ns1.domainnameapi.com`, `ns2.domainnameapi.com`),
EPP kilidi açık, gizlilik kapalı, ek özellik (additional attributes) yok. Ek özellikler yalnızca
`.tr` gibi alan adları için gereklidir.

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
// Yalnızca .tr alan adları için:
opts.AdditionalAttributes = map[string]any{
	"TRABISDOMAINCATEGORY": 1,             // 1 = bireysel (Owner), 0 = kurumsal (Corporate)
	"TRABISCITIZIENID":     "12345678901", // kurumlar için TRABISTAXOFFICE + TRABISTAXNUMBER
}

result := client.RegisterWithContactInfo("example.com", 1, dna.AllRoles(contact), opts)
```

`.tr` alan adlarında eski (legacy) TRABIS adları gateway şemasına göre normalize edilir: kategori
`Owner`/`Corporate` olur, `TRABISCITIZIENID` `TRABISCITIZENID` olur ve ülke, şehir, ad ve kurum
alanları çıkarılır (gateway bunları kişi bilgisinden alır). Her özellik değeri metin (string) olarak
gönderilir.

<details>
<summary>Domain Kayıt işlemleri için Örnek Çıktı</summary>

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

Kayıt uç noktası kısa bir kayıtla yanıt verdiği için alanların çoğu boştur; alan adının tüm
bilgileri için ardından `GetDetails` çağırın. Gateway kaydı reddettiğinde (`"success": false` içeren
HTTP 200 yanıtı) sonuç, `REGISTER_DOMAIN` kodlu ve gateway'in mesajını içeren bir hatadır.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain Yenileme

```go
result := client.Renew("example.com", 1)
```

<details>
<summary>Domain Yenileme Örnek Çıktı</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"ExpirationDate": "2026-03-15T10:00:00"}}
```

Gateway'in reddi (`"success": false`) `RENEW_DOMAIN` kodunu döndürür; bitiş tarihi içermeyen başarılı
bir yanıt ise `DOMAIN_RENEW` döndürür.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain Transfer

```go
result := client.Transfer("example.com", "Xy9#mK2$", 1)
```

Transferler asla kişi bilgisi göndermez: kayıt firması kişi bilgilerini alan adıyla birlikte devralır;
bu nedenle istek her zaman boş bir kişi listesi taşır ve metodun kişi parametresi yoktur. Sonuç,
`GetDetails` ile aynı yapıdadır; gateway'in reddi (`"success": false`) gateway'in mesajıyla birlikte
`TRANSFER_DOMAIN` kodunu döndürür.

<details>
<summary>Domain Transfer Örnek Çıktı</summary>

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

İlgili transfer çağrıları:

```go
client.CancelTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Cancelled"}}
client.ApproveTransfer("example.com") // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Approved"}}
client.RejectTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Rejected"}}
```

<hr style="border: 4px solid #000; border-style: dashed;">

### Transfer Kontrolü

```go
result := client.CheckTransfer("example.com", "Ab1!xyZ9")
```

`result` yalnızca kayıt otoritesi (registry) alan adının transfer edilebilir olduğunu bildirdiğinde
`OK` olur; aksi halde `TRANSFER_NOT_AVAILABLE` koduyla ve aynı `data` ile `ERROR` olur.

<details>
<summary>Transfer Kontrolü Örnek Çıktı</summary>

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

### Alan adları listesi

```go
result := client.GetList(map[string]any{
	"OrderColumn":    "Id",
	"OrderDirection": "ASC",
	"PageNumber":     0,
	"PageSize":       1000,
})
```

`PageNumber` (0'dan başlar), `PageSize`, `OrderColumn` ve `OrderDirection`, gateway'in `SkipCount`,
`MaxResultCount` ve `Sorting` parametrelerine çevrilir; diğer anahtarlar (örneğin `DomainName`) filtre
olarak olduğu gibi iletilir. `nil` verildiğinde ilk 200 alan adı döner.

<details>
<summary>Alan adları listesi Örnek Çıktı</summary>

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

### TLD Listesi

```go
result := client.GetTldList(100) // php-dna varsayılanı: 20
```

<details>
<summary>TLD Listesi Örnek Çıktı</summary>

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

### Alan adı uygunluğu kontrolü

```go
result := client.CheckAvailability([]string{"hello", "world123x0"}, []string{"com", "net"}, 1, "create")
if list, ok := result.([]any); ok {
	// başarılı
	_ = list
}
```

<details>
<summary>Alan adı uygunluğu kontrolü Örnek Çıktı</summary>

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

Premium alan adlarında `IsFee` değeri `true` olur.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Alan adı detayları

```go
result := client.GetDetails("example.com")
result = client.SyncFromRegistry("example.com") // aynı yapı
```

<details>
<summary>Alan adı detayları Örnek Çıktı</summary>

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

`Status`, gateway sayısal enum gönderse bile metin etiketidir (`Active`, `TransferredOut`,
`PendingDelete`, ...). Ayrıştırılamayan tarihler `""` olur; asla anlamsız bir 1970 tarihi dönmez.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Nameserver Düzenlemesi

```go
result := client.ModifyNameServer("example.com", []string{"ns1.example.com", "ns2.example.com"})
```

Map de kabul edilir (değerleri anahtar sırasına göre kullanılır). Girdiler gönderilmeden önce
boşluklardan arındırılır, küçük harfe çevrilir, sondaki nokta kaldırılır ve tekrarlar ayıklanır;
sonuçta bu normalize edilmiş liste döner.

<details>
<summary>Nameserver Düzenlemesi Örnek Çıktı</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServers": []any{"ns1.example.com", "ns2.example.com"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain Kilidi aktifleştirme

```go
result := client.EnableTheftProtectionLock("example.com")
```

<details>
<summary>Domain Kilidi aktifleştirme Örnek Çıktı</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain Kilidi kaldırma

```go
result := client.DisableTheftProtectionLock("example.com")
```

<details>
<summary>Domain Kilidi kaldırma Örnek Çıktı</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": false}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domaine ChildNS ekleme

```go
result := client.AddChildNameServer("example.com", "ns1.example.com", "1.2.3.4")
```

<details>
<summary>Domaine ChildNS ekleme Örnek Çıktı</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.4"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domaine ChildNS silme

```go
result := client.DeleteChildNameServer("example.com", "ns1.example.com")
```

<details>
<summary>Domaine ChildNS silme Örnek Çıktı</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com"}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Child NS Güncelleme

```go
result := client.ModifyChildNameServer("example.com", "ns1.example.com", "1.2.3.5")
```

<details>
<summary>Child NS Güncelleme Örnek Çıktı</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.5"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain gizliliği değiştirme

```go
result := client.ModifyPrivacyProtectionStatus("example.com", true, "Owner request")
```

`status` bir bool ya da panellerin sakladığı bir değerdir: `true`, `1`, `"1"`, `"true"`, `"on"` ve
`"yes"` gizliliği açar; diğer her şey (`"off"`, `"0"`, `false`, ...) kapatır. Gerekçe (reason)
yalnızca php-dna uyumluluğu için korunur ve gönderilmez: gateway'de böyle bir alan yoktur.

<details>
<summary>Domain gizliliği değiştirme Örnek Çıktı</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"PrivacyProtectionStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain contact bilgilerini alma

```go
result := client.GetContacts("example.com")
```

Gateway kişi türleri büyük/küçük harf duyarsız eşleştirilir (`Admin`/`Administrative`,
`Tech`/`Technical`, `Registrant`/`Owner`, `Billing`); diğer türler kendi baş harfi büyük adıyla korunur.

<details>
<summary>Domain contact bilgilerini alma Örnek Çıktı</summary>

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
			// "Billing", "Registrant", "Technical": aynı yapı
		},
	},
}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Domain contact kaydetme

```go
result := client.SaveContacts("example.com", dna.AllRoles(contact))
```

<details>
<summary>Domain contact kaydetme Örnek Çıktı</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"contacts": map[string]any{}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### Reseller bilgilerini alma

```go
result := client.GetResellerDetails()
```

<details>
<summary>Reseller bilgilerini alma Örnek Çıktı</summary>

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

### Güncel bakiyeyi alma

```go
result := client.GetCurrentBalance("USD") // USD (varsayılan), TRY, EUR veya GBP
```

<details>
<summary>Güncel bakiyeyi alma Örnek Çıktı</summary>

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

Hata durumunda: `ErrorCode` `API_<status>` (veya `CURL_<n>`), `OperationMessage` hata mesajı,
`OperationResult` `"FAILED"` ve `Balance` `"0.00"` olur. Hesap yalnızca USD ve TRY bakiyelerini
bildirir: başka herhangi bir para birimi için sonuç `ErrorCode: "BALANCE"` ve
`No EUR balance is reported for this account. Available: USD, TRY.` ile `FAILED` olur.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### DNS Zone Yönetimi

Gateway kayıtları **kümeler** (set) halinde saklar (tüm değerleriyle birlikte bir ad + tür) ve her
yazma işlemi kümenin tamamını değiştirir. Kayıt adları `@`, `""`, zone kökü (apex), göreli bir etiket
(`www`) veya tam ad (`www.example.com`) olarak verilebilir; hepsi gateway'in beklediği göreli etikete
indirgenir. Desteklenen türler: `A`, `AAAA`, `ALIAS`, `CAA`, `CNAME`, `DS`, `MX`, `NAPTR`, `PTR`,
`SRV`, `TLSA`, `TXT` (diğerleri herhangi bir çağrı yapılmadan reddedilir). TTL 1..86400 aralığına
sınırlanır ve sonuç, gerçekte gönderilen büyük harfli türü ve TTL değerini bildirir. `apply = true` ile
(php-dna varsayılanı) değişiklikler hemen uygulanır; bu uygulama (apply) başarısız olursa onun hata
zarfı döndürülür (değişiklik beklemede kalır). Gateway ani istek yığınlarına hız sınırı uyguladığı için
ardışık zone çağrıları arasında en az bir saniye bırakılır.

```go
// Listeleme: her değer için bir satır, adlar tam nitelikli (FQDN)
client.GetResourceRecords("example.com")

// Bir ad+türü tam olarak bu değerlere ayarla (boş liste kümeyi siler). Varsayılanlar: ttl 3600, apply true
client.SetResourceRecordSet("example.com", "@", "MX", []string{"10 mx1.example.com", "20 mx2.example.com"}, 3600, true)

// Diğerlerini koruyarak TEK bir değer ekle / değiştir / kaldır (ttl nil = kümenin TTL'ini koru)
client.AddResourceRecordValue("example.com", "www", "A", "192.0.2.10", nil, true)
client.ReplaceResourceRecordValue("example.com", "www", "A", "192.0.2.10", "192.0.2.11", nil, true)
client.RemoveResourceRecordValue("example.com", "www", "A", "192.0.2.11", true)

// Bir ad+türü tamamen kaldır (knownValue "" = önce mevcut bir değeri sorgula)
client.DeleteResourceRecordSet("example.com", "www", "A", true, "")

// Düşük seviyeli çağrılar (content bir string veya []string)
client.AddResourceRecord("example.com", "api", "CNAME", "example.com.", 300, false)
client.EditResourceRecord("example.com", "api", "api", "CNAME", "www.example.com.", 300, false)
client.DeleteResourceRecord("example.com", "old.example.com.", "192.0.2.99", "A", false)
client.ApplyResourceRecords("example.com")
```

<details>
<summary>DNS Zone Yönetimi Örnek Çıktılar</summary>

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

// Değer zaten mevcut / mevcut değil: hiçbir şey yazılmaz
map[string]any{"result": "OK", "data": map[string]any{"Name": "", "Type": "MX", "Contents": []any{"..."}, "Changed": false}}

// Var olmayan bir küme için DeleteResourceRecordSet
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

### URL Yönlendirme

```go
client.SetForwarding("example.com", "https://www.example.org", "Standard") // veya "Frame" (maskeli)
client.GetForwarding("example.com")
client.DeleteForwarding("example.com")
```

<details>
<summary>URL Yönlendirme Örnek Çıktılar</summary>

```go
// SetForwarding
map[string]any{"result": "OK", "data": map[string]any{"RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// GetForwarding
map[string]any{"result": "OK", "data": map[string]any{"Enabled": true, "RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// Yönlendirme tanımlı değilken GetForwarding
map[string]any{"result": "OK", "data": map[string]any{"Enabled": false, "RedirectAddress": "", "ForwardType": ""}}

// DeleteForwarding (hiçbir şey tanımlı değilse Removed false olur)
map[string]any{"result": "OK", "data": map[string]any{"DomainName": "example.com", "Removed": true}}
```

</details>

Yalnızca gateway'in "Redirect information not found" yanıtı "tanımlı değil" anlamına gelir; diğer tüm
hatalar, örneğin bilinmeyen bir alan adı, hata olarak döndürülür.

<hr style="border: 4px solid #000; border-style: dashed;">

### Yardımcı araçlar (ağ erişimi yok)

```go
client.IsTrTLD("example.com.tr") // true

client.ValidateContact(map[string]any{"FirstName": " Jo ", "Phone": "+90 (532) 374 33 88", "Country": "Turkey"})
// map[string]any{"FirstName": "Jo", "LastName": "Jo", "AddressLine1": "Addres yok", "City": "ISTANBUL",
//   "State": "ISTANBUL", "Country": "TR", "ZipCode": "34000", "Phone": "5323743388",
//   "PhoneCountryCode": "90", "Fax": "", "FaxCountryCode": ""}
```

`ValidateContact` eksik alanları yer tutucu değerlerle doldurur, adlardan rakamları ve noktalama
işaretlerini temizler, State boşsa City değerini kullanır, yaygın ülke adlarını ISO kodlarına çevirir
ve telefon ile faks numaralarını numara ve ülke kodu olarak ayırır. Kendi uluslararası önekiyle
yazılmış bir numara (`+44 20 ...`, `0044 20 ...`), başka bir ülke kodu verilse bile kendi ülke kodunu korur.

<hr style="border: 4px solid #000; border-style: dashed;">

### Tanılama (Diagnostics)

```go
client.GetRequestData()        // son çağrının {"url": ..., "method": ..., "payload": ...} bilgisi
client.GetResponseData()       // son çözümlenmiş yanıt veya {"raw_response": body}
client.GetResponseHeaders()    // {"url", "http_code", "content_type", "total_time", "headers"}
client.GetParsedResponseData() // son ayrıştırılmış yanıt veya {Code, Message, Details} hatası
client.GetLastFunction()       // örn. "request:GET domains/info"
client.GetServiceURL()         // gateway temel URL'i
client.SetServiceURL("https://ote.domainresellerapi.com/api/v1")
```

Bir `Client` bu çağrıya özgü durumu kendi içinde tuttuğu için tek bir `Client`'ı goroutine'ler arasında
paylaşmayın; her goroutine için ayrı bir tane oluşturun.

## Hata yönetimi

- Gateway hataları `Code: "API_<HTTP status>"` döndürür. Gateway kendi hata kodunu gönderdiyse
  (örneğin `Dna.DomainService:Domain:10007`), bu kod `GatewayCode` olarak eklenir; böylece çağıran
  taraf "alan adı bulunamadı" ile "TLD desteklenmiyor" durumlarını ayırt edebilir. Mesaj gateway'den
  gelir; boş bir hata gövdesi `HTTP 429 - Too Many Requests` (ve benzerleri) olur, doğrulama hataları
  `Details` içinde alan alan listelenir (`Admin PostalCode - 'PostalCode' is required`) ve
  `Exception of type '...' was thrown` gibi arka uç kalıp metinleri gateway'in ayrıntılarıyla ya da
  hata kodunu belirten sabit bir mesajla değiştirilir — asla ham yanıt gövdesiyle değil.
- 301/302 yanıtı, kimlik bilgilerinin reddedildiği anlamına gelir: `Message: "Invalid API credentials"`.
- Aktarım hataları (zaman aşımı, DNS, TLS, bağlantı sıfırlanması) her zaman hatadır, asla boş bir
  başarılı yanıt değildir. Zaman aşımı `Code: "CURL_28"` ve `Could not reach the DomainNameAPI service: the connection
  timed out. (...)` döndürür; diğer hatalar curl tarzı numaralar (`CURL_6` DNS, `CURL_7` bağlantı,
  `CURL_35` TLS, `CURL_56` sıfırlama, `CURL_60` sertifika) ve `Curl error during request: ...` kullanır.
  `API_<n>` yalnızca HTTP durum kodlarına ayrılmıştır.
- Kayıt, transfer ve yenileme işlemlerinde `"success": false` içeren HTTP 200 yanıtları birer rettir:
  gateway'in mesajıyla birlikte `REGISTER_DOMAIN`, `TRANSFER_DOMAIN` veya `RENEW_DOMAIN` döner.
- Boş veya `null` bir başarılı yanıt gövdesi `Code: "RESPONSE"` döndürür.
- Zaman aşımları: çağrı başına 30 saniye, bağlantı kurmak için 10 saniye.

## Test modu

OTE test gateway'ini (`https://ote.domainresellerapi.com/api/v1`) kullanmak için constructor'ın üçüncü
argümanı olarak `true` (veya `"true"`, `"1"`, `"on"`, `"yes"`) verin. Test modundaki telemetri
`environment: ote` etiketiyle işaretlenir.

## Telemetri (hata raporlama)

Bir API çağrısı başarısız olduğunda kütüphane, DomainNameAPI'nin kendi hata takip sistemine (Sentry)
bir hata raporu gönderebilir. Başarılı çağrıların rastgele seçilen küçük bir kısmı (%4, ayrıca 1
saniyeden uzun süren her çağrı) performans kaydı olarak gönderilir. Raporlama **varsayılan olarak açıktır**.

> ℹ️ `DomainNameAPI-Go` için Sentry projesi henüz oluşturulmadı: `DefaultSentryDSN` boştur; bu nedenle
> bir sürüm bu değeri ayarlayana kadar hiçbir şey gönderilmez. `dna.WithSentryDSN(...)` DSN'i geçersiz kılar.

**Gönderilenler:** hata kodu ve mesajı; uç nokta yolu ve HTTP yöntemi; çağrı süresi; kütüphane adı
(`DomainNameAPI-Go`), kütüphane sürümü ve Go sürümü; bayi kimliği (Reseller ID) ve alan adı; sunucunun
genel (public) IP adresi (`api.ipify.org` üzerinden öğrenilir); işletim sistemi kullanıcı adı; kişi
bilgileri dahil API isteği ve yanıtı. Yalnızca kimlik bilgileri maskelenir: EPP/transfer kodları, API
anahtarları, parolalar ve token'lar `***` ile değiştirilir.

**Asla gönderilmeyenler:** okunabilir biçimde EPP/transfer kodları, API anahtarları, parolalar veya
token'lar; hata izlerindeki (stack trace) fonksiyon çağrısı argümanları.

**Neden ve ne kadar süre:** veriler yalnızca kullanım ölçümü ve sorunların erken teşhisi için kullanılır
ve DomainNameAPI'nin kendi Sentry sunucusunda 30 gün saklanır.

**Kapatma veya yalnızca asgari bilgiyi gönderme:**

```go
client := dna.New(resellerID, apiKey, false)
client.SetErrorReportingEnabled(false) // hiçbir şey gönderilmez

client.SetTelemetryMinimal(true) // yalnızca hata kodu/mesajı, uç nokta, süre, sürümler ve Reseller ID
```

Telemetri arka planda çalışır, çağrıyı asla uzun süre bekletmez ve kodunuza asla panic yansıtmaz.

> 🔒 **Güvenlik:** TLS sertifika doğrulaması açıktır. CA paketi güncel olmayan bir sunucuda
> `client.SetSSLVerification(false)` ile kapatılabilir — önerilmez; asıl çözüm CA paketini güncellemektir.

## Gelişmiş seçenekler

```go
client := dna.New(resellerID, apiKey, false,
	dna.WithHTTPClient(myHTTPClient),      // herhangi bir Do(*http.Request) (*http.Response, error)
	dna.WithSentryDSN("https://key@host/1"),
	dna.WithTelemetryClient(myHTTPClient),
)
```

Özel bir HTTP istemcisi kullanıyorsanız zaman aşımlarını, proxy'leri ve TLS ayarlarını o istemci
üzerinde yapılandırın; istemci yönlendirmeleri takip etmemelidir (`CheckRedirect` içinden
`http.ErrUseLastResponse` döndürün).

## Testler

Sunucuda Go kurulu olması gerekmez; Docker yeterlidir:

```bash
docker run --rm -v "$PWD":/src -w /src golang:1.23 go test ./...
docker run --rm -v "$PWD":/src -w /src golang:1.23 go vet ./...
```

Test paketi; HTTP katmanı, hata eşleme, telemetri ve yardımcı araçlar için birim testleri ile
`testdata/conformance` altındaki 134 php-dna fixture'ının tamamını yeniden oynatan bir uyumluluk
(conformance) testi içerir (ağa giden istekler ve dönüş değerleri php-dna ile birebir eşleşmelidir).
Hiçbir test ağa erişmez.

## Dönüş ve Hata Kodları ile Açıklamaları

| Kod  | Açıklama                                       | Detay                                                                                                                                                          |
|------|------------------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------|
| 1000 | Komut başarıyla tamamlandı                     | İşlem başarılı.                                                                                                                                                |
| 1001 | Komut başarıyla tamamlandı; işlem beklemede    | İşlem başarılı. Fakat işlem şu an tamamlanmak için kuyruğa alındı.                                                                                             |
| 2003 | Zorunlu parametre eksik                        | Parametre eksik hatası. Örneğin; kontak bilgisinde telefon girişi yapılmaması.                                                                                 |
| 2105 | Nesne yenilemeye uygun değil                   | Domain durumu yenilemeye müsait değil, güncelleme işlemlerine kilitlenmiştir. Durum "clientupdateprohibited" olmamalı. Diğer durum koşullarından kaynaklanabilir. |
| 2200 | Kimlik doğrulama hatası                        | Yetki hatası, güvenlik kodu hatalı veya domain başka bir kayıt firmasında bulunuyor.                                                                           |
| 2302 | Nesne zaten mevcut                             | Domain adı veya name server bilgisi veritabanında mevcut. Kayıt edilemez.                                                                                      |
| 2303 | Nesne mevcut değil                             | Domain adı veya name server bilgisi veritabanında mevcut değil. Yeni kayıt oluşturulmalı.                                                                      |
| 2304 | Nesne durumu işleme izin vermiyor              | Domain durumu güncellemeye müsait değildir, güncelleme işlemlerine kilitlenmiştir. Durum "clientupdateprohibited" olmamalı. Diğer durum koşullarından kaynaklanabilir. |

Kütüphane kodları: `RESPONSE` (boş yanıt), `CREDENTIALS` (hesap uç noktasından geçersiz yanıt biçimi),
`AVAILABILITY` (alan adı verilmedi), `BALANCE` (para birimi hesap tarafından bildirilmiyor),
`TRANSFER_NOT_AVAILABLE`, `REGISTER_DOMAIN` / `TRANSFER_DOMAIN` / `RENEW_DOMAIN` (gateway işlemi
reddetti), `DOMAIN_RENEW` (bitiş tarihi olmadan yenileme), `ADD_ZONE`, `EDIT_ZONE`, `SET_ZONE`,
`DELETE_ZONE` (desteklenmeyen kayıt türü veya değer bulunamadı), `SET_FORWARD` (desteklenmeyen
yönlendirme türü veya boş adres), `CURL_<n>` (aktarım hatası), `API_<status>` (HTTP hatası; gateway
kendi kodunu gönderdiyse `GatewayCode` ile birlikte).

## Lisans

Apache License 2.0 — bkz. [LICENSE](LICENSE).
