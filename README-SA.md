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
# go-dna — عميل DomainNameAPI للغة Go

عميل بلغة Go لواجهة REST API الخاصة بمسجّل النطاقات DomainNameAPI. وهو نقل أمين لمكتبة
[php-dna](https://github.com/domainreseller/php-dna): الدوال نفسها، والطلبات نفسها المُرسَلة عبر
الشبكة، والقيم المُعادة نفسها. المكتبة القياسية هي الاعتمادية الوحيدة.

## 📦 التنزيل — استخدم دائمًا صفحة الإصدارات (Releases)!

⬇️ **احصل على أحدث نسخة مختبَرة من هنا: https://github.com/domainreseller/go-dna/releases/latest**

> ⚠️ **لا** تستخدم الزر الأخضر **Code → Download ZIP** — فهو ينزّل فرع التطوير الخام. حزم الإصدارات مرقّمة الإصدار ومختبَرة وجاهزة لبيئة الإنتاج.

## المتطلبات

- Go 1.21 أو أعلى
- لا توجد حزم من أطراف ثالثة (فقط `net/http` و `encoding/json` من المكتبة القياسية)

## التثبيت

باستخدام Go modules، ثبّت إصدارًا منشورًا محددًا:

```bash
go get github.com/domainreseller/go-dna@v1.0.0
```

```go
import dna "github.com/domainreseller/go-dna"
```

اسم الحزمة هو `dna`. توجد أمثلة قابلة للتشغيل لكل دالة في مجلد [examples](examples):

```bash
go run ./examples/GetDetails
```

## 🔑 بيانات اعتماد API — Reseller ID و API Key

تحتاج المكتبة إلى بيانات اعتماد اثنتين بالضبط تُمرَّران إلى الدالة الإنشائية (constructor):

| المُعامِل | القيمة |
|---|---|
| الأول | Reseller ID — بصيغة UUID مثل `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` |
| الثاني | API Key |

> 💡 ستجد **Reseller ID** و **API Key** في لوحة DomainNameAPI الخاصة بك ضمن **إعدادات API**.
> ⚠️ هذه **بيانات اعتماد API** — بيانات تسجيل الدخول إلى اللوحة **لن تعمل** هنا.

```go
resellerID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
apiKey := "your-api-key"

// البيئة الحية
client := dna.New(resellerID, apiKey, false)

// بيئة الاختبار (OTE)
client = dna.New(resellerID, apiKey, true)
```

تُزال المسافات الزائدة من بداية بيانات الاعتماد ونهايتها. يقبل الوسيط الثالث أيضًا السلاسل النصية التي
تخزّنها اللوحات لمربعات الاختيار: `"true"` و `"1"` و `"on"` و `"yes"` (بأي حالة أحرف) تختار بيئة
الاختبار، وكل ما عدا ذلك (`"false"` و `"0"` و `"off"` و `""` ...) يختار البيئة الحية.

## البدء السريع

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

## القيم المُعادة

تُعيد الدوال البنى الديناميكية نفسها التي تُعيدها php-dna — أي `map[string]any` و `[]any` — وليس
نماذج ذات أنواع محددة، لذا فإن المفاتيح والأعلام النصية (`"true"`/`"false"`) والمبالغ المالية المنسّقة
كسلاسل نصية (`"11.9900"`) مطابقة تمامًا لما في php-dna. الأرقام القادمة من البوابة تكون من النوع
`int64` أو `float64`.

تبدو الاستدعاءات الناجحة على هذا النحو:

```go
map[string]any{"result": "OK", "data": map[string]any{ /* ... */ }}
```

الاستدعاءات الفاشلة لا تُسبب panic أبدًا؛ بل تُعيد غلاف الخطأ:

```go
map[string]any{
	"result": "ERROR",
	"error": map[string]any{
		"Code":    "API_404",      // API_<حالة HTTP>، أو رمز من المكتبة مثل RESPONSE
		"Message": "The domain name you requested could not be found. Please check and try again.",
		"Details": "The domain could not be found. Please check and try again.",
	},
}
```

تستخدم ثلاث دوال شكلًا مختلفًا، تمامًا كما في php-dna:

| الدالة | النجاح | الفشل |
|---|---|---|
| `CheckAvailability` | قائمة `[]any` مجرّدة | غلاف الخطأ (`map[string]any`)، لذا فإن نوع القيمة المُعادة هو `any` |
| `GetCurrentBalance` | غلاف يحتوي `OperationResult: "SUCCESS"` | الغلاف نفسه مع `OperationResult: "FAILED"` |
| `GetResellerDetails` | خريطة مسطّحة (`id` و `name` و `balance` و `balances` ...) | غلاف الخطأ |

### الوسائط الاختيارية

لا تدعم Go الوسائط الافتراضية. تحافظ الدوال على ترتيب المعاملات في php-dna وتأخذ كل وسيط بشكل
صريح؛ وتُعرض القيمة الافتراضية في php-dna في كل قسم أدناه. ولتسهيل الاستخدام، تختار بعض القيم الصفرية
القيمة الافتراضية: `GetCurrentBalance("")` تعني USD، و `GetTldList(0)` تعني 20،
و `CheckAvailability(..., "")` تعني `create`، و `SetForwarding(..., "")` تعني `Standard`،
و `ModifyPrivacyProtectionStatus(..., "")` تعني `Owner request`.
تأخذ `RegisterWithContactInfo` وسائطها الاختيارية الأربعة ضمن `*RegisterOptions` (القيمة nil = القيم الافتراضية).
تأخذ `AddResourceRecordValue` و `ReplaceResourceRecordValue` المعامل `ttl *int`: القيمة `nil` تُبقي
على TTL الحالي للمجموعة.

### جهات الاتصال

جهات الاتصال هي `map[string]any` بمفاتيح php-dna (`FirstName` و `LastName` و `Company` و `EMail`
و `AddressLine1` و `AddressLine2` و `City` و `State` و `Country` و `ZipCode` و `Phone`
و `PhoneCountryCode` و `Fax` و `FaxCountryCode`). كما يُقبل الشكل المتداخل (`Address.Line1`
و `Phone.Phone.Number` ...). ولأن خرائط Go ليس لها ترتيب، تُمرَّر الأدوار الأربعة كقائمة مرتّبة
من النوع `dna.Contacts`:

```go
contacts := dna.Contacts{
	{Role: "Administrative", Contact: contact},
	{Role: "Billing", Contact: contact},
	{Role: "Technical", Contact: contact},
	{Role: "Registrant", Contact: contact},
}
contacts = dna.AllRoles(contact)                 // جهة الاتصال نفسها للأدوار الأربعة
contacts = dna.ContactsFromMap(map[string]map[string]any{ /* role => contact */ })
```

تُستكمل جهات الاتصال الناقصة وفق قواعد `ValidateContact` قبل إرسالها.

### أسماء الدوال

تتبع أسماء الدوال في Go أعراف لغة Go؛ واسم الدالة في php-dna هو نفسه بصيغة lowerCamelCase.
`getServiceUrl`/`setServiceUrl` أصبحتا `GetServiceURL`/`SetServiceURL`، و `setSslVerification` أصبحت
`SetSSLVerification`.

## الاستخدام

### تسجيل النطاق

القيم الافتراضية في php-dna: خوادم الأسماء الافتراضية (`ns1.domainnameapi.com` و `ns2.domainnameapi.com`)،
قفل EPP مفعّل، الخصوصية معطّلة، ولا توجد سمات إضافية. السمات الإضافية مطلوبة فقط لنطاقات
مثل `.tr`.

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
// لنطاقات .tr فقط:
opts.AdditionalAttributes = map[string]any{
	"TRABISDOMAINCATEGORY": 1,             // 1 = فرد (Owner)، 0 = شركة (Corporate)
	"TRABISCITIZIENID":     "12345678901", // أو TRABISTAXOFFICE + TRABISTAXNUMBER للشركات
}

result := client.RegisterWithContactInfo("example.com", 1, dna.AllRoles(contact), opts)
```

بالنسبة لنطاقات `.tr`، تُوحَّد أسماء TRABIS القديمة وفق مخطط البوابة: تصبح الفئة
`Owner`/`Corporate`، ويصبح `TRABISCITIZIENID` هو `TRABISCITIZENID`، وتُحذف حقول الدولة والمدينة
والاسم والمؤسسة (إذ تأخذها البوابة من جهة الاتصال). تُرسَل قيمة كل سمة كسلسلة نصية.

<details>
<summary>مثال على مخرجات تسجيل النطاق</summary>

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

تُجيب نقطة نهاية التسجيل بسجل مختصر، لذا تكون معظم الحقول فارغة؛ استدعِ `GetDetails`
بعد ذلك للحصول على بيانات النطاق كاملة.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### تجديد النطاق

```go
result := client.Renew("example.com", 1)
```

<details>
<summary>مثال على مخرجات تجديد النطاق</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"ExpirationDate": "2026-03-15T10:00:00"}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### نقل النطاق

```go
result := client.Transfer("example.com", "Xy9#mK2$", 1, nil)
```

كما في الواجهة (facade) الخاصة بـ php-dna، لا يُمرَّر وسيط جهات الاتصال: يحمل طلب النقل دائمًا
قائمة جهات اتصال فارغة. للنتيجة الشكل نفسه الذي تُعيده `GetDetails`.

<details>
<summary>مثال على مخرجات نقل النطاق</summary>

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

استدعاءات النقل ذات الصلة:

```go
client.CancelTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Cancelled"}}
client.ApproveTransfer("example.com") // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Approved"}}
client.RejectTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Rejected"}}
```

<hr style="border: 4px solid #000; border-style: dashed;">

### التحقق من إمكانية النقل

```go
result := client.CheckTransfer("example.com", "Ab1!xyZ9")
```

تكون قيمة `result` هي `OK` فقط عندما يُفيد السجل (registry) بأن النطاق قابل للنقل؛ وإلا فهي `ERROR`
مع الرمز `TRANSFER_NOT_AVAILABLE` ومع قيمة `data` نفسها.

<details>
<summary>مثال على مخرجات التحقق من إمكانية النقل</summary>

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

### قائمة النطاقات

```go
result := client.GetList(map[string]any{
	"OrderColumn":    "Id",
	"OrderDirection": "ASC",
	"PageNumber":     0,
	"PageSize":       1000,
})
```

تُترجَم `PageNumber` (تبدأ من 0) و `PageSize` و `OrderColumn` و `OrderDirection` إلى
`SkipCount` و `MaxResultCount` و `Sorting` الخاصة بالبوابة؛ أما المفاتيح الأخرى (مثل `DomainName`) فتُمرَّر
كما هي كمرشّحات. القيمة `nil` تُعيد أول 200 نطاق.

<details>
<summary>مثال على مخرجات قائمة النطاقات</summary>

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

### قائمة امتدادات النطاقات (TLD)

```go
result := client.GetTldList(100) // القيمة الافتراضية في php-dna: 20
```

<details>
<summary>مثال على مخرجات قائمة امتدادات النطاقات</summary>

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

### التحقق من توفر النطاق

```go
result := client.CheckAvailability([]string{"hello", "world123x0"}, []string{"com", "net"}, 1, "create")
if list, ok := result.([]any); ok {
	// نجاح
	_ = list
}
```

<details>
<summary>مثال على مخرجات التحقق من توفر النطاق</summary>

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

تكون قيمة `IsFee` هي `true` للنطاقات المميزة (premium).

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### تفاصيل النطاق

```go
result := client.GetDetails("example.com")
result = client.SyncFromRegistry("example.com") // الشكل نفسه
```

<details>
<summary>مثال على مخرجات تفاصيل النطاق</summary>

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

تكون قيمة `Status` هي التسمية النصية (`Active` و `TransferredOut` و `PendingDelete` ...) حتى عندما
ترسل البوابة القيمة الرقمية للتعداد (enum). التواريخ غير القابلة للتحليل تصبح `""`، ولا تصبح أبدًا تاريخًا زائفًا من عام 1970.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### تعديل خوادم الأسماء

```go
result := client.ModifyNameServer("example.com", []string{"ns1.example.com", "ns2.example.com"})
```

تُقبل الخريطة (map) أيضًا (تُستخدم قيمها بترتيب المفاتيح). قبل الإرسال، تُزال المسافات الزائدة من
المُدخلات وتُحوَّل إلى أحرف صغيرة وتُحذف النقطة الأخيرة منها وتُزال القيم المكرّرة.

<details>
<summary>مثال على مخرجات تعديل خوادم الأسماء</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServers": []string{"ns1.example.com", "ns2.example.com"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### تفعيل قفل النطاق

```go
result := client.EnableTheftProtectionLock("example.com")
```

<details>
<summary>مثال على مخرجات تفعيل قفل النطاق</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### إلغاء قفل النطاق

```go
result := client.DisableTheftProtectionLock("example.com")
```

<details>
<summary>مثال على مخرجات إلغاء قفل النطاق</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": false}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### إضافة خادم أسماء فرعي

```go
result := client.AddChildNameServer("example.com", "ns1.example.com", "1.2.3.4")
```

<details>
<summary>مثال على مخرجات إضافة خادم أسماء فرعي</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.4"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### حذف خادم أسماء فرعي

```go
result := client.DeleteChildNameServer("example.com", "ns1.example.com")
```

<details>
<summary>مثال على مخرجات حذف خادم أسماء فرعي</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com"}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### تحديث خادم الأسماء الفرعي

```go
result := client.ModifyChildNameServer("example.com", "ns1.example.com", "1.2.3.5")
```

<details>
<summary>مثال على مخرجات تحديث خادم الأسماء الفرعي</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.5"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### تعديل حالة حماية الخصوصية

```go
result := client.ModifyPrivacyProtectionStatus("example.com", true, "Owner request")
```

<details>
<summary>مثال على مخرجات تعديل حالة حماية الخصوصية</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"PrivacyProtectionStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### الحصول على معلومات الاتصال للنطاق

```go
result := client.GetContacts("example.com")
```

<details>
<summary>مثال على مخرجات الحصول على معلومات الاتصال</summary>

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
			// "Billing" و "Registrant" و "Technical": الشكل نفسه
		},
	},
}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### حفظ معلومات الاتصال للنطاق

```go
result := client.SaveContacts("example.com", dna.AllRoles(contact))
```

<details>
<summary>مثال على مخرجات حفظ معلومات الاتصال</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"contacts": map[string]any{}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### الحصول على معلومات الموزع

```go
result := client.GetResellerDetails()
```

<details>
<summary>مثال على مخرجات معلومات الموزع</summary>

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

### الحصول على الرصيد الحالي

```go
result := client.GetCurrentBalance("USD") // USD (افتراضي)، TRY أو EUR أو GBP
```

<details>
<summary>مثال على مخرجات الرصيد الحالي</summary>

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

عند الفشل: تكون `ErrorCode` هي `API_<status>`، و `OperationMessage` هي رسالة الخطأ،
و `OperationResult` هي `"FAILED"`، و `Balance` هي `"0.00"`.

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### إدارة منطقة DNS

تخزّن البوابة السجلات على شكل **مجموعات** (اسم + نوع واحد مع جميع قيمه)، وكل عملية كتابة
تستبدل المجموعة كاملة. يمكن إعطاء أسماء السجلات بالصيغة `@` أو `""` أو جذر المنطقة (apex) أو تسمية نسبية
(`www`) أو اسم كامل (`www.example.com`)؛ وتُختزل جميعها إلى التسمية النسبية التي تتوقعها البوابة.
الأنواع المدعومة: `A` و `AAAA` و `ALIAS` و `CAA` و `CNAME` و `DS` و `MX` و `NAPTR` و `PTR` و `SRV`
و `TLSA` و `TXT` (تُرفض الأنواع الأخرى قبل إجراء أي استدعاء). تُحصر قيمة TTL ضمن النطاق 1..86400. مع `apply = true`
(القيمة الافتراضية في php-dna) تُعتمد التغييرات فورًا؛ ويُفصل بين استدعاءات المنطقة المتتالية بثانية
واحدة على الأقل لأن البوابة تحدّ من معدل الطلبات المتلاحقة.

```go
// القائمة: صف واحد لكل قيمة، والأسماء مؤهّلة بالكامل
client.GetResourceRecords("example.com")

// تعيين اسم+نوع إلى هذه القيم بالضبط (القائمة الفارغة تحذف المجموعة). الافتراضي: ttl 3600، apply true
client.SetResourceRecordSet("example.com", "@", "MX", []string{"10 mx1.example.com", "20 mx2.example.com"}, 3600, true)

// إضافة / استبدال / حذف قيمة واحدة مع الإبقاء على البقية (ttl nil = الإبقاء على TTL المجموعة)
client.AddResourceRecordValue("example.com", "www", "A", "192.0.2.10", nil, true)
client.ReplaceResourceRecordValue("example.com", "www", "A", "192.0.2.10", "192.0.2.11", nil, true)
client.RemoveResourceRecordValue("example.com", "www", "A", "192.0.2.11", true)

// حذف اسم+نوع بالكامل (knownValue "" = البحث عن قيمة حالية أولًا)
client.DeleteResourceRecordSet("example.com", "www", "A", true, "")

// استدعاءات منخفضة المستوى (المحتوى سلسلة نصية أو []string)
client.AddResourceRecord("example.com", "api", "CNAME", "example.com.", 300, false)
client.EditResourceRecord("example.com", "api", "api", "CNAME", "www.example.com.", 300, false)
client.DeleteResourceRecord("example.com", "old.example.com.", "192.0.2.99", "A", false)
client.ApplyResourceRecords("example.com")
```

<details>
<summary>أمثلة على مخرجات إدارة منطقة DNS</summary>

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

// القيمة موجودة مسبقًا / غير موجودة: لا يُكتب شيء
map[string]any{"result": "OK", "data": map[string]any{"Name": "", "Type": "MX", "Contents": []any{"..."}, "Changed": false}}

// DeleteResourceRecordSet لمجموعة غير موجودة
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

### إعادة توجيه URL

```go
client.SetForwarding("example.com", "https://www.example.org", "Standard") // أو "Frame" (مُقنَّع)
client.GetForwarding("example.com")
client.DeleteForwarding("example.com")
```

<details>
<summary>أمثلة على مخرجات إعادة توجيه URL</summary>

```go
// SetForwarding
map[string]any{"result": "OK", "data": map[string]any{"RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// GetForwarding (في حال عدم إعداد إعادة توجيه: Enabled تساوي false، والنتيجة تبقى OK)
map[string]any{"result": "OK", "data": map[string]any{"Enabled": true, "RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// DeleteForwarding
map[string]any{"result": "OK", "data": map[string]any{"DomainName": "example.com"}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### الأدوات المساعدة (دون اتصال بالشبكة)

```go
client.IsTrTLD("example.com.tr") // true

client.ValidateContact(map[string]any{"FirstName": " Jo ", "Phone": "+90 (532) 374 33 88", "Country": "Turkey"})
// map[string]any{"FirstName": "Jo", "LastName": "Jo", "AddressLine1": "Addres yok", "City": "ISTANBUL",
//   "State": "ISTANBUL", "Country": "TR", "ZipCode": "34000", "Phone": "5323743388",
//   "PhoneCountryCode": "90", "Fax": "", "FaxCountryCode": ""}
```

تملأ `ValidateContact` الحقول الناقصة بقيم بديلة، وتزيل الأرقام وعلامات الترقيم من الأسماء،
وتستخدم المدينة عند غياب الولاية، وتحوّل أسماء الدول الشائعة إلى رموز ISO، وتفصل أرقام الهاتف
والفاكس إلى الرقم ورمز الدولة.

<hr style="border: 4px solid #000; border-style: dashed;">

### التشخيص

```go
client.GetRequestData()        // {"url": ..., "method": ..., "payload": ...} لآخر استدعاء
client.GetResponseData()       // آخر استجابة مفكوكة، أو {"raw_response": body}
client.GetResponseHeaders()    // {"url", "http_code", "content_type", "total_time", "headers"}
client.GetParsedResponseData() // آخر استجابة محلَّلة، أو الخطأ {Code, Message, Details}
client.GetLastFunction()       // مثلًا "request:GET domains/info"
client.GetServiceURL()         // عنوان URL الأساسي للبوابة
client.SetServiceURL("https://ote.domainresellerapi.com/api/v1")
```

يحتفظ `Client` بهذه الحالة الخاصة بكل استدعاء، لذا لا تشارك `Client` واحدًا بين عدة goroutines؛ بل أنشئ
واحدًا لكل goroutine.

## معالجة الأخطاء

- تُعيد أخطاء البوابة `Code: "API_<HTTP status>"`. تأتي الرسالة من البوابة؛ ويصبح جسم الخطأ الفارغ
  `HTTP 429 - Too Many Requests` (وهكذا)، وتُسرد أخطاء التحقق حقلًا تلو الآخر
  في `Details` (`Admin PostalCode - 'PostalCode' is required`)، وتُستبدل العبارات النمطية للواجهة الخلفية مثل
  `Exception of type '...' was thrown` بالتفاصيل المفيدة.
- الاستجابة 301/302 تعني أن بيانات الاعتماد قد رُفضت: `Message: "Invalid API credentials"`.
- أعطال النقل (انتهاء المهلة، DNS، TLS، إعادة تعيين الاتصال) تُعدّ أخطاء، ولا تكون أبدًا نجاحًا فارغًا.
  عند انتهاء المهلة تُعاد `Code: "API_28"` والرسالة `Could not reach the DomainNameAPI service: the connection
  timed out. (...)`؛ وتستخدم الأعطال الأخرى أرقامًا على نمط curl (`API_6` لـ DNS، و `API_7` للاتصال، و `API_35`
  لـ TLS، و `API_56` لإعادة التعيين، و `API_60` للشهادة) مع الرسالة `Curl error during request: ...`.
  تحتوي `GetParsedResponseData()["Code"]` على `CURL_<n>`.
- جسم النجاح الفارغ أو `null` يُعيد `Code: "RESPONSE"`.
- المهل الزمنية: 30 ثانية لكل استدعاء، و 10 ثوانٍ لإنشاء الاتصال.

## وضع الاختبار

مرّر `true` (أو `"true"` أو `"1"` أو `"on"` أو `"yes"`) كوسيط ثالث للدالة الإنشائية لاستخدام بوابة
الاختبار OTE ‏`https://ote.domainresellerapi.com/api/v1`. يُوسَم القياس عن بُعد في وضع الاختبار بـ
`environment: ote`.

## القياس عن بُعد (تقارير الأخطاء)

عندما يفشل استدعاء API، يمكن للمكتبة إرسال تقرير خطأ إلى نظام تتبع الأخطاء الخاص بـ DomainNameAPI
(Sentry). كما يُرسَل جزء صغير عشوائي من الاستدعاءات الناجحة (4%، بالإضافة إلى كل استدعاء تتجاوز مدته ثانية واحدة)
كسجل أداء. تقارير الأخطاء **مفعّلة افتراضيًا**.

> ℹ️ لم يُنشأ مشروع Sentry الخاص بـ `DomainNameAPI-Go` بعد: `DefaultSentryDSN` فارغ، لذا
> لا يُرسَل شيء حتى يقوم أحد الإصدارات بتعيينه. تتجاوز `dna.WithSentryDSN(...)` قيمة DSN.

**ما يتم إرساله:** رمز الخطأ ورسالته؛ مسار نقطة النهاية وطريقة HTTP؛ مدة الاستدعاء؛ اسم المكتبة
(`DomainNameAPI-Go`) وإصدار المكتبة وإصدار Go؛ معرّف الموزع (Reseller ID) واسم النطاق؛ عنوان IP
العام للخادم (يتم الحصول عليه عبر `api.ipify.org`)؛ اسم مستخدم نظام التشغيل؛ طلب API
والاستجابة بما في ذلك بيانات الاتصال. يتم إخفاء بيانات الاعتماد فقط: تُستبدل رموز EPP/auth ومفاتيح API
وكلمات المرور والرموز المميزة (tokens) بـ `***`.

**ما لا يتم إرساله أبدًا:** رموز EPP/auth أو مفاتيح API أو كلمات المرور أو الرموز المميزة (tokens) بصيغة مقروءة؛ وسائط
استدعاءات الدوال في تتبع المكدس (stack traces).

**لماذا وإلى متى:** تُستخدم البيانات فقط لقياس الاستخدام والتشخيص المبكر للمشكلات،
وتُحفظ لمدة 30 يومًا على خادم Sentry الخاص بـ DomainNameAPI.

**الإيقاف، أو إرسال الحد الأدنى:**

```go
client := dna.New(resellerID, apiKey, false)
client.SetErrorReportingEnabled(false) // لا يُرسَل أي شيء

client.SetTelemetryMinimal(true) // فقط رمز الخطأ/رسالته، ونقطة النهاية، والتوقيت، والإصدارات، ومعرّف الموزع
```

يعمل القياس عن بُعد في الخلفية، ولا يحجب الاستدعاء لمدة طويلة أبدًا، ولا يُسبب panic في شيفرتك أبدًا.

> 🔒 **الأمان:** التحقق من شهادة TLS مفعّل. على خادم تكون حزمة شهادات CA فيه قديمة، يمكن
> إيقافه باستخدام `client.SetSSLVerification(false)` — غير مستحسن؛ تحديث حزمة CA هو
> الحل الأفضل.

## الخيارات المتقدمة

```go
client := dna.New(resellerID, apiKey, false,
	dna.WithHTTPClient(myHTTPClient),      // أي نوع يملك Do(*http.Request) (*http.Response, error)
	dna.WithSentryDSN("https://key@host/1"),
	dna.WithTelemetryClient(myHTTPClient),
)
```

عند استخدام عميل HTTP مخصص، اضبط المهل الزمنية والوكلاء (proxies) و TLS عليه؛ ويجب ألّا يتبع عمليات
إعادة التوجيه (أعِد `http.ErrUseLastResponse` من `CheckRedirect`).

## الاختبارات

لا يلزم تثبيت Go على الجهاز المضيف؛ يكفي Docker:

```bash
docker run --rm -v "$PWD":/src -w /src golang:1.23 go test ./...
docker run --rm -v "$PWD":/src -w /src golang:1.23 go vet ./...
```

تتضمن مجموعة الاختبارات اختبارات وحدات لطبقة HTTP وربط الأخطاء والقياس عن بُعد والأدوات المساعدة،
واختبار توافق يعيد تشغيل جميع حالات php-dna البالغ عددها 126 من `testdata/conformance` (يجب أن تتطابق
الطلبات المُرسَلة عبر الشبكة والقيم المُعادة مع php-dna تمامًا). لا يتصل أي اختبار بالشبكة.

## رموز العودة والأخطاء مع تفسيراتها

| الرمز | الوصف | التفاصيل |
|------|------|------|
| 1000 | تم تنفيذ الأمر بنجاح | العملية ناجحة. |
| 1001 | تم تنفيذ الأمر بنجاح؛ الإجراء معلق | العملية ناجحة. ولكن العملية حاليًا في قائمة الانتظار للإكمال. |
| 2003 | معلمة مطلوبة مفقودة | خطأ في المعلمات المفقودة. على سبيل المثال، عدم إدخال رقم الهاتف في معلومات الاتصال. |
| 2105 | الكائن غير مؤهل للتجديد | حالة النطاق غير مناسبة للتجديد، وهو مقفل لعمليات التحديث. يجب ألّا تكون الحالة "clientupdateprohibited". قد يكون ذلك بسبب شروط حالة أخرى. |
| 2200 | خطأ في المصادقة | خطأ في التفويض، رمز الأمان غير صحيح أو النطاق مسجل لدى شركة تسجيل أخرى. |
| 2302 | الكائن موجود | اسم النطاق أو معلومات خادم الأسماء موجودة بالفعل في قاعدة البيانات. لا يمكن التسجيل. |
| 2303 | الكائن غير موجود | اسم النطاق أو معلومات خادم الأسماء غير موجودة في قاعدة البيانات. يلزم تسجيل جديد. |
| 2304 | حالة الكائن تمنع العملية | حالة النطاق غير مناسبة للتحديث، وهو مقفل لعمليات التحديث. يجب ألّا تكون الحالة "clientupdateprohibited". قد يكون ذلك بسبب شروط حالة أخرى. |

رموز المكتبة: `RESPONSE` (استجابة فارغة)، و `CREDENTIALS` (صيغة استجابة غير صالحة من نقطة نهاية
الحساب)، و `AVAILABILITY` (لم تُعطَ أسماء نطاقات)، و `TRANSFER_NOT_AVAILABLE`، و `DOMAIN_RENEW`
(تجديد دون تاريخ انتهاء)، و `ADD_ZONE` و `EDIT_ZONE` و `SET_ZONE` و `DELETE_ZONE` (نوع سجل غير
مدعوم أو قيمة غير موجودة)، و `SET_FORWARD` (نوع إعادة توجيه غير مدعوم أو عنوان فارغ).

## الترخيص

Apache License 2.0 — راجع [LICENSE](LICENSE).
