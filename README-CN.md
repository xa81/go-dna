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
# go-dna — 适用于 Go 的 DomainNameAPI 客户端

DomainNameAPI 域名注册商 REST API 的 Go 客户端。它是
[php-dna](https://github.com/domainreseller/php-dna) 的忠实移植：相同的方法、相同的线上请求以及相同的返回值。
唯一的依赖是标准库。

## 📦 下载 — 请务必使用 Releases！

⬇️ **在此获取最新的经过测试的版本：https://github.com/domainreseller/go-dna/releases/latest**

> ⚠️ 请**不要**使用绿色的 **Code → Download ZIP** 按钮 — 它下载的是原始开发分支。Release 包经过版本管理和测试，可直接用于生产环境。

## 系统要求

- Go 1.21 或更高版本
- 无第三方包（仅使用标准库中的 `net/http` 和 `encoding/json`）

## 安装

使用 Go modules 时，请固定到一个已发布的版本：

```bash
go get github.com/domainreseller/go-dna@v1.0.0
```

```go
import dna "github.com/domainreseller/go-dna"
```

包名为 `dna`。每个方法的可运行示例都位于 [examples](examples) 目录中：

```bash
go run ./examples/GetDetails
```

## 🔑 API 凭据 — Reseller ID 与 API Key

本库只需要两项凭据，传入构造函数即可：

| 参数 | 值 |
|---|---|
| 第一个 | Reseller ID — 形如 `xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx` 的 UUID |
| 第二个 | API Key |

> 💡 您可以在 DomainNameAPI 面板的 **API 设置** 中找到您的 **Reseller ID** 和 **API Key**。
> ⚠️ 这些是 **API 凭据** — 面板登录信息在此**无法使用**。

```go
resellerID := "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
apiKey := "your-api-key"

// 正式环境
client := dna.New(resellerID, apiKey, false)

// 测试（OTE）环境
client = dna.New(resellerID, apiKey, true)
```

两项凭据都会去除首尾空白。第三个参数也接受面板为复选框存储的字符串：`"true"`、`"1"`、`"on"`
和 `"yes"`（不区分大小写）选择测试环境，其他任何值（`"false"`、`"0"`、`"off"`、`""` 等）选择正式环境。

## 快速开始

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

## 返回值

各方法返回与 php-dna 相同的动态结构 — `map[string]any` 和 `[]any` — 而不是强类型模型，因此键名、
字符串标志（`"true"`/`"false"`）以及格式化的金额字符串（`"11.9900"`）都与 php-dna 完全一致。
来自网关的数字为 `int64` 或 `float64`。

成功的调用如下所示：

```go
map[string]any{"result": "OK", "data": map[string]any{ /* ... */ }}
```

失败的调用永远不会 panic，而是返回错误信封：

```go
map[string]any{
	"result": "ERROR",
	"error": map[string]any{
		"Code":        "API_404", // API_<HTTP 状态码>、CURL_<n>，或库自身的代码，例如 RESPONSE
		"Message":     "The domain name you requested could not be found. Please check and try again.",
		"Details":     "The domain could not be found. Please check and try again.",
		"GatewayCode": "Dna.DomainService:Domain:10007", // 仅与 API_<status> 一起出现，当网关发送了自己的代码时
	},
}
```

有三个方法使用不同的结构，与 php-dna 完全一致：

| 方法 | 成功 | 失败 |
|---|---|---|
| `CheckAvailability` | 一个裸 `[]any` 列表 | 错误信封（`map[string]any`），因此返回类型为 `any` |
| `GetCurrentBalance` | `OperationResult: "SUCCESS"` 信封 | 相同的信封，但 `OperationResult: "FAILED"` |
| `GetResellerDetails` | 扁平 map（`id`、`name`、`balance`、`balances` 等） | 错误信封 |

### 可选参数

Go 没有默认参数。各方法保持 php-dna 的参数顺序，并显式接收每一个参数；php-dna 的默认值在下文各节中给出。
为方便起见，部分零值会选择默认值：`GetCurrentBalance("")` 表示 USD，`GetTldList(0)` 表示 20，
`CheckAvailability(..., "")` 表示 `create`，`SetForwarding(..., "")` 表示 `Standard`。
`RegisterWithContactInfo` 通过 `*RegisterOptions` 接收其四个可选参数（nil = 默认值）。
`AddResourceRecordValue` 和 `ReplaceResourceRecordValue` 接收 `ttl *int`：`nil` 表示保留记录集当前的 TTL。

### 联系人

联系人是使用 php-dna 键名的 `map[string]any`（`FirstName`、`LastName`、`Company`、`EMail`、
`AddressLine1`、`AddressLine2`、`City`、`State`、`Country`、`ZipCode`、`Phone`、
`PhoneCountryCode`、`Fax`、`FaxCountryCode`）。嵌套结构（`Address.Line1`、
`Phone.Phone.Number` 等）同样被接受。由于 Go 的 map 没有顺序，四个角色以有序的 `dna.Contacts` 列表传入：

```go
contacts := dna.Contacts{
	{Role: "Administrative", Contact: contact},
	{Role: "Billing", Contact: contact},
	{Role: "Technical", Contact: contact},
	{Role: "Registrant", Contact: contact},
}
contacts = dna.AllRoles(contact)                 // 四个角色使用同一个联系人
contacts = dna.ContactsFromMap(map[string]map[string]any{ /* role => contact */ })
```

信息不完整的联系人会在发送前按照 `ValidateContact` 的规则补全。联系人由 `RegisterWithContactInfo`
和 `SaveContacts` 使用；转移操作从不发送联系人。

### 方法名称

Go 方法名遵循 Go 的命名惯例；对应的 php-dna 名称就是其 lowerCamelCase 形式。
`getServiceUrl`/`setServiceUrl` 对应 `GetServiceURL`/`SetServiceURL`，`setSslVerification` 对应
`SetSSLVerification`。

## 使用方法

### 域名注册

php-dna 默认值：默认名称服务器（`ns1.domainnameapi.com`、`ns2.domainnameapi.com`），EPP 锁定开启，
隐私保护关闭，无附加属性。仅 `.tr` 等域名需要附加属性。

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
// 仅适用于 .tr 域名：
opts.AdditionalAttributes = map[string]any{
	"TRABISDOMAINCATEGORY": 1,             // 1 = 个人（Owner），0 = 公司（Corporate）
	"TRABISCITIZIENID":     "12345678901", // 公司则使用 TRABISTAXOFFICE + TRABISTAXNUMBER
}

result := client.RegisterWithContactInfo("example.com", 1, dna.AllRoles(contact), opts)
```

对于 `.tr` 域名，旧版 TRABIS 名称会被规范化为网关的数据结构：类别变为 `Owner`/`Corporate`，
`TRABISCITIZIENID` 变为 `TRABISCITIZENID`，国家、城市、姓名和组织字段会被删除（网关会从联系人中获取这些信息）。
所有属性值都以字符串形式发送。

<details>
<summary>域名注册示例输出</summary>

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

注册端点只返回一条简短记录，因此大多数字段为空；如需完整的域名信息，请随后调用 `GetDetails`。当网关拒绝注册时
（HTTP 200 且 `"success": false`），结果为错误，代码为 `REGISTER_DOMAIN`，并附带网关的消息。

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### 域名续费

```go
result := client.Renew("example.com", 1)
```

<details>
<summary>域名续费示例输出</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"ExpirationDate": "2026-03-15T10:00:00"}}
```

网关拒绝（`"success": false`）时返回代码 `RENEW_DOMAIN`；成功但没有到期日期时返回 `DOMAIN_RENEW`。

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### 域名转移

```go
result := client.Transfer("example.com", "Xy9#mK2$", 1)
```

转移操作从不发送联系人：注册商会随域名一并接管联系人，因此请求始终携带一个空的联系人列表，
该方法也没有联系人参数。返回结果与 `GetDetails` 的结构相同；网关拒绝（`"success": false`）时
返回代码 `TRANSFER_DOMAIN` 并附带网关的消息。

<details>
<summary>域名转移示例输出</summary>

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

相关的转移调用：

```go
client.CancelTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Cancelled"}}
client.ApproveTransfer("example.com") // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Approved"}}
client.RejectTransfer("example.com")  // {"result": "OK", "data": {"DomainName": "example.com", "Status": "Rejected"}}
```

<hr style="border: 4px solid #000; border-style: dashed;">

### 转移检查

```go
result := client.CheckTransfer("example.com", "Ab1!xyZ9")
```

只有当注册局报告该域名可转移时，`result` 才为 `OK`；否则为 `ERROR`，代码为 `TRANSFER_NOT_AVAILABLE`，
并附带相同的 `data`。

<details>
<summary>转移检查示例输出</summary>

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

### 域名列表

```go
result := client.GetList(map[string]any{
	"OrderColumn":    "Id",
	"OrderDirection": "ASC",
	"PageNumber":     0,
	"PageSize":       1000,
})
```

`PageNumber`（从 0 开始）、`PageSize`、`OrderColumn` 和 `OrderDirection` 会被转换为网关的
`SkipCount`、`MaxResultCount` 和 `Sorting`；其他键（例如 `DomainName`）作为筛选条件原样传递。
传入 `nil` 时返回前 200 个域名。

<details>
<summary>域名列表示例输出</summary>

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

### TLD 列表

```go
result := client.GetTldList(100) // php-dna 默认值：20
```

<details>
<summary>TLD 列表示例输出</summary>

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

### 域名可用性检查

```go
result := client.CheckAvailability([]string{"hello", "world123x0"}, []string{"com", "net"}, 1, "create")
if list, ok := result.([]any); ok {
	// 成功
	_ = list
}
```

<details>
<summary>域名可用性检查示例输出</summary>

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

对于溢价（premium）域名，`IsFee` 为 `true`。

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### 域名详情

```go
result := client.GetDetails("example.com")
result = client.SyncFromRegistry("example.com") // 结构相同
```

<details>
<summary>域名详情示例输出</summary>

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

即使网关发送的是数字枚举值，`Status` 也是文本标签（`Active`、`TransferredOut`、`PendingDelete` 等）。
无法解析的日期会变为 `""`，绝不会变成虚假的 1970 年日期。

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### 修改名称服务器

```go
result := client.ModifyNameServer("example.com", []string{"ns1.example.com", "ns2.example.com"})
```

也接受 map（按键的顺序使用其值）。发送前，各条目会被去除首尾空白、转换为小写、去掉末尾的点并去重；结果中返回的就是规范化后的列表。

<details>
<summary>修改名称服务器示例输出</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServers": []any{"ns1.example.com", "ns2.example.com"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### 启用域名锁定

```go
result := client.EnableTheftProtectionLock("example.com")
```

<details>
<summary>启用域名锁定示例输出</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### 禁用域名锁定

```go
result := client.DisableTheftProtectionLock("example.com")
```

<details>
<summary>禁用域名锁定示例输出</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"LockStatus": false}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### 添加子名称服务器

```go
result := client.AddChildNameServer("example.com", "ns1.example.com", "1.2.3.4")
```

<details>
<summary>添加子名称服务器示例输出</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.4"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### 删除子名称服务器

```go
result := client.DeleteChildNameServer("example.com", "ns1.example.com")
```

<details>
<summary>删除子名称服务器示例输出</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com"}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### 修改子名称服务器

```go
result := client.ModifyChildNameServer("example.com", "ns1.example.com", "1.2.3.5")
```

<details>
<summary>修改子名称服务器示例输出</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"NameServer": "ns1.example.com", "IPAdresses": []any{"1.2.3.5"}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### 修改隐私保护状态

```go
result := client.ModifyPrivacyProtectionStatus("example.com", true, "Owner request")
```

`status` 可以是 bool，也可以是面板存储的值：`true`、`1`、`"1"`、`"true"`、`"on"` 和 `"yes"` 启用
隐私保护，其他任何值（`"off"`、`"0"`、`false` 等）则禁用。原因参数仅为兼容 php-dna 而保留，不会被发送：
网关没有这样的字段。

<details>
<summary>修改隐私保护状态示例输出</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"PrivacyProtectionStatus": true}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### 获取域名联系人

```go
result := client.GetContacts("example.com")
```

网关返回的联系人类型按不区分大小写的方式匹配（`Admin`/`Administrative`、`Tech`/`Technical`、
`Registrant`/`Owner`、`Billing`）；其他任何类型都会以其首字母大写的名称保留。

<details>
<summary>获取域名联系人示例输出</summary>

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
			// "Billing"、"Registrant"、"Technical"：结构相同
		},
	},
}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### 保存域名联系人

```go
result := client.SaveContacts("example.com", dna.AllRoles(contact))
```

<details>
<summary>保存域名联系人示例输出</summary>

```go
map[string]any{"result": "OK", "data": map[string]any{"contacts": map[string]any{}}}
```

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### 获取经销商详细信息

```go
result := client.GetResellerDetails()
```

<details>
<summary>获取经销商详细信息示例输出</summary>

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

### 获取当前余额

```go
result := client.GetCurrentBalance("USD") // USD（默认）、TRY、EUR 或 GBP
```

<details>
<summary>获取当前余额示例输出</summary>

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

失败时：`ErrorCode` 为 `API_<status>`（或 `CURL_<n>`），`OperationMessage` 为错误消息，
`OperationResult` 为 `"FAILED"`，`Balance` 为 `"0.00"`。账户只报告 USD 和 TRY 余额：
对于其他任何货币，结果为 `FAILED`，`ErrorCode: "BALANCE"`，消息为
`No EUR balance is reported for this account. Available: USD, TRY.`

</details>

<hr style="border: 4px solid #000; border-style: dashed;">

### DNS 区域管理

网关以**记录集**的形式存储记录（一个名称 + 类型及其所有值），每次写入都会替换整个记录集。
记录名称可以写作 `@`、`""`、区域顶点（apex）、相对标签（`www`）或完整名称（`www.example.com`）；
它们都会被归约为网关所期望的相对标签。支持的类型：`A`、`AAAA`、`ALIAS`、`CAA`、`CNAME`、`DS`、`MX`、
`NAPTR`、`PTR`、`SRV`、`TLSA`、`TXT`（其他类型会在发出任何调用之前被拒绝）。TTL 会被限制在 1..86400 范围内，结果中返回实际发送的大写类型和 TTL。
当 `apply = true`（php-dna 的默认值）时，更改会立即提交；若该提交失败，则返回其错误信封（更改仍处于暂存状态）。
由于网关会对突发请求进行限流，连续的区域调用之间至少间隔一秒。

```go
// 列表：每个值一行，名称为完全限定名
client.GetResourceRecords("example.com")

// 将某个名称+类型设置为恰好这些值（空列表会删除该记录集）。默认值：ttl 3600，apply true
client.SetResourceRecordSet("example.com", "@", "MX", []string{"10 mx1.example.com", "20 mx2.example.com"}, 3600, true)

// 添加 / 替换 / 删除单个值，保留其他值（ttl nil = 保留记录集的 TTL）
client.AddResourceRecordValue("example.com", "www", "A", "192.0.2.10", nil, true)
client.ReplaceResourceRecordValue("example.com", "www", "A", "192.0.2.10", "192.0.2.11", nil, true)
client.RemoveResourceRecordValue("example.com", "www", "A", "192.0.2.11", true)

// 删除整个名称+类型（knownValue "" = 先查询一个当前值）
client.DeleteResourceRecordSet("example.com", "www", "A", true, "")

// 底层调用（content 为字符串或 []string）
client.AddResourceRecord("example.com", "api", "CNAME", "example.com.", 300, false)
client.EditResourceRecord("example.com", "api", "api", "CNAME", "www.example.com.", 300, false)
client.DeleteResourceRecord("example.com", "old.example.com.", "192.0.2.99", "A", false)
client.ApplyResourceRecords("example.com")
```

<details>
<summary>DNS 区域管理示例输出</summary>

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

// 值已存在 / 不存在：不写入任何内容
map[string]any{"result": "OK", "data": map[string]any{"Name": "", "Type": "MX", "Contents": []any{"..."}, "Changed": false}}

// 对不存在的记录集调用 DeleteResourceRecordSet
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

### URL 转发

```go
client.SetForwarding("example.com", "https://www.example.org", "Standard") // 或 "Frame"（隐藏式）
client.GetForwarding("example.com")
client.DeleteForwarding("example.com")
```

<details>
<summary>URL 转发示例输出</summary>

```go
// SetForwarding
map[string]any{"result": "OK", "data": map[string]any{"RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// GetForwarding
map[string]any{"result": "OK", "data": map[string]any{"Enabled": true, "RedirectAddress": "https://www.example.org", "ForwardType": "Standard"}}

// 未配置任何转发时的 GetForwarding
map[string]any{"result": "OK", "data": map[string]any{"Enabled": false, "RedirectAddress": "", "ForwardType": ""}}

// DeleteForwarding（未配置任何转发时 Removed 为 false）
map[string]any{"result": "OK", "data": map[string]any{"DomainName": "example.com", "Removed": true}}
```

</details>

只有网关返回 "Redirect information not found" 时才表示"未配置"；其他任何错误（例如未知域名）都会作为错误返回。

<hr style="border: 4px solid #000; border-style: dashed;">

### 实用工具（无网络访问）

```go
client.IsTrTLD("example.com.tr") // true

client.ValidateContact(map[string]any{"FirstName": " Jo ", "Phone": "+90 (532) 374 33 88", "Country": "Turkey"})
// map[string]any{"FirstName": "Jo", "LastName": "Jo", "AddressLine1": "Addres yok", "City": "ISTANBUL",
//   "State": "ISTANBUL", "Country": "TR", "ZipCode": "34000", "Phone": "5323743388",
//   "PhoneCountryCode": "90", "Fax": "", "FaxCountryCode": ""}
```

`ValidateContact` 会用占位值填充缺失字段，删除姓名中的数字和标点符号，在缺少 State 时回退使用 City，
将常见的国家名称映射为 ISO 代码，并将电话和传真号码拆分为号码和国家代码。带有自身国际前缀的号码（`+44 20 ...`、`0044 20 ...`）
即使另外提供了其他国家代码，也会保留其自身的国家代码。

<hr style="border: 4px solid #000; border-style: dashed;">

### 诊断

```go
client.GetRequestData()        // 上一次调用的 {"url": ..., "method": ..., "payload": ...}
client.GetResponseData()       // 上一次解码的响应，或 {"raw_response": body}
client.GetResponseHeaders()    // {"url", "http_code", "content_type", "total_time", "headers"}
client.GetParsedResponseData() // 上一次解析的响应，或 {Code, Message, Details} 错误
client.GetLastFunction()       // 例如 "request:GET domains/info"
client.GetServiceURL()         // 网关基础 URL
client.SetServiceURL("https://ote.domainresellerapi.com/api/v1")
```

`Client` 会保存这些每次调用的状态，因此不要在多个 goroutine 之间共享同一个 `Client`；请为每个 goroutine
创建一个。

## 错误处理

- 网关错误返回 `Code: "API_<HTTP status>"`。当网关发送了自己的错误代码
  （例如 `Dna.DomainService:Domain:10007`）时，该代码会作为 `GatewayCode` 添加，以便调用方区分
  "域名未找到"与"不支持该 TLD"。消息来自网关；空的错误响应体会变为
  `HTTP 429 - Too Many Requests`（以此类推），验证错误会在 `Details` 中逐字段列出
  （`Admin PostalCode - 'PostalCode' is required`），而 `Exception of type '...' was thrown` 之类的后端
  模板化信息会被替换为网关的详细信息，或替换为一条注明错误代码的固定消息 — 绝不会替换为原始响应体。
- 301/302 响应表示凭据被拒绝：`Message: "Invalid API credentials"`。
- 传输故障（超时、DNS、TLS、连接重置）都会作为错误返回，绝不会是空的成功结果。超时返回
  `Code: "CURL_28"` 和 `Could not reach the DomainNameAPI service: the connection
  timed out. (...)`；其他故障使用 curl 风格的编号（`CURL_6` DNS、`CURL_7` 连接、
  `CURL_35` TLS、`CURL_56` 重置、`CURL_60` 证书）以及 `Curl error during request: ...`。
  `API_<n>` 专用于 HTTP 状态码。
- 注册、转移和续费返回 HTTP 200 且 `"success": false` 时视为被拒绝：
  返回 `REGISTER_DOMAIN`、`TRANSFER_DOMAIN` 或 `RENEW_DOMAIN`，并附带网关的消息。
- 空的或 `null` 的成功响应体返回 `Code: "RESPONSE"`。
- 超时设置：每次调用 30 秒，建立连接 10 秒。

## 测试模式

将 `true`（或 `"true"`、`"1"`、`"on"`、`"yes"`）作为构造函数的第三个参数传入，即可使用 OTE
测试网关 `https://ote.domainresellerapi.com/api/v1`。测试模式下的遥测数据会被标记为
`environment: ote`。

## 遥测（错误报告）

当 API 调用失败时，本库可以向 DomainNameAPI 自己的错误跟踪系统（Sentry）发送一份错误报告。
成功调用中会随机抽取一小部分（4%，另外所有耗时超过 1 秒的调用）作为性能记录发送。错误报告**默认开启**。

> ℹ️ `DomainNameAPI-Go` 的 Sentry 项目尚未创建：`DefaultSentryDSN` 为空，因此在某个发布版本设置它之前
> 不会发送任何内容。`dna.WithSentryDSN(...)` 可覆盖 DSN。

**发送的内容：** 错误代码和消息；端点路径和 HTTP 方法；调用耗时；库名称
（`DomainNameAPI-Go`）、库版本和 Go 版本；经销商 ID（Reseller ID）和域名；服务器的公网 IP 地址
（通过 `api.ipify.org` 获取）；操作系统用户名；包括联系人信息在内的 API 请求和响应。仅凭据会被遮蔽：
EPP/auth 代码、API 密钥、密码和令牌（token）会被替换为 `***`。

**绝不会发送的内容：** 可读形式的 EPP/auth 代码、API 密钥、密码或令牌（token）；堆栈跟踪（stack trace）
中的函数调用参数。

**用途和保留期限：** 数据仅用于使用情况统计和问题的早期诊断，在 DomainNameAPI 自己的 Sentry 服务器上
保留 30 天。

**关闭，或仅发送最少信息：**

```go
client := dna.New(resellerID, apiKey, false)
client.SetErrorReportingEnabled(false) // 不发送任何内容

client.SetTelemetryMinimal(true) // 仅发送错误代码/消息、端点、耗时、版本和经销商 ID
```

遥测在后台运行，绝不会长时间阻塞调用，也绝不会在您的代码中引发 panic。

> 🔒 **安全：** TLS 证书验证处于开启状态。在 CA 证书包过时的主机上，可以通过
> `client.SetSSLVerification(false)` 将其关闭 — 不推荐；更新 CA 证书包才是更好的解决办法。

## 高级选项

```go
client := dna.New(resellerID, apiKey, false,
	dna.WithHTTPClient(myHTTPClient),      // 任何实现 Do(*http.Request) (*http.Response, error) 的类型
	dna.WithSentryDSN("https://key@host/1"),
	dna.WithTelemetryClient(myHTTPClient),
)
```

使用自定义 HTTP 客户端时，请在其上配置超时、代理和 TLS；它不得跟随重定向
（请在 `CheckRedirect` 中返回 `http.ErrUseLastResponse`）。

## 测试

主机上无需安装 Go，有 Docker 即可：

```bash
docker run --rm -v "$PWD":/src -w /src golang:1.23 go test ./...
docker run --rm -v "$PWD":/src -w /src golang:1.23 go vet ./...
```

测试套件包含针对 HTTP 层、错误映射、遥测和实用工具的单元测试，以及一个一致性测试，该测试会回放
`testdata/conformance` 中全部 134 个 php-dna 测试用例（线上请求和返回值必须与 php-dna 完全一致）。
没有任何测试会访问网络。

## 返回和错误代码及说明

| 代码 | 描述 | 详细 |
|------|------|------|
| 1000 | 命令成功完成 | 操作成功。 |
| 1001 | 命令成功完成；操作待处理 | 操作成功。但是，操作当前已排队等待完成。 |
| 2003 | 缺少必需参数 | 参数缺失错误。例如：联系信息中没有电话输入。 |
| 2105 | 对象不符合续费条件 | 域名状态不适合续费，已锁定更新操作。状态不应为"clientupdateprohibited"。可能由于其他状态条件。 |
| 2200 | 身份验证错误 | 授权错误，安全代码不正确或域名在另一个注册商处注册。 |
| 2302 | 对象已存在 | 域名或名称服务器信息已存在于数据库中。无法注册。 |
| 2303 | 对象不存在 | 域名或名称服务器信息在数据库中不存在。需要新注册。 |
| 2304 | 对象状态禁止操作 | 域名状态不适合更新，已锁定更新操作。状态不应为"clientupdateprohibited"。可能由于其他状态条件。 |

库代码：`RESPONSE`（空响应）、`CREDENTIALS`（账户端点返回的响应格式无效）、`AVAILABILITY`（未提供域名）、
`BALANCE`（账户未报告该货币）、`TRANSFER_NOT_AVAILABLE`、`REGISTER_DOMAIN` / `TRANSFER_DOMAIN` /
`RENEW_DOMAIN`（网关拒绝了该操作）、`DOMAIN_RENEW`（续费时缺少到期日期）、`ADD_ZONE`、`EDIT_ZONE`、`SET_ZONE`、
`DELETE_ZONE`（不支持的记录类型或未找到该值）、`SET_FORWARD`（不支持的转发类型或地址为空）、
`CURL_<n>`（传输故障）、`API_<status>`（HTTP 错误，当网关发送了自己的代码时附带 `GatewayCode`）。

## 许可证

Apache License 2.0 — 参见 [LICENSE](LICENSE)。
