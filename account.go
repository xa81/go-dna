package dna

import (
	"sort"
	"strconv"
	"strings"
)

// GetResellerDetails returns the reseller account with its balances.
//
//	{"result": "OK", "id": "<uuid>", "active": true, "name": "...",
//	 "balance": "601.64", "currency": "USD", "symbol": "$",
//	 "balances": [{"balance": "601.64", "currency": "USD", "symbol": "$"},
//	              {"balance": "1250.5", "currency": "TL", "symbol": "TL"}]}
func (c *Client) GetResellerDetails() map[string]any {
	response, err := c.request("GET", "deposit/accounts/me", nil)
	if err != nil {
		return c.errorEnvelope(err)
	}
	if dig(response, "resellerId") == nil {
		return map[string]any{
			"result": ResultError,
			"error":  setError("CREDENTIALS", "Invalid response format", "Response does not contain required fields"),
		}
	}
	usd := phpString(coalesce(dig(response, "usdBalance"), int64(0)))
	try := phpString(coalesce(dig(response, "tryBalance"), int64(0)))
	return map[string]any{
		"result":   ResultOK,
		"id":       dig(response, "resellerId"),
		"active":   true,
		"name":     coalesce(dig(response, "resellerName"), ""),
		"balance":  usd,
		"currency": "USD",
		"symbol":   "$",
		"balances": []any{
			map[string]any{"balance": usd, "currency": "USD", "symbol": "$"},
			map[string]any{"balance": try, "currency": "TL", "symbol": "TL"},
		},
	}
}

type currencyInfo struct {
	id     int64
	name   string
	symbol string
}

var currencyMap = map[string]currencyInfo{
	"USD": {2, "USD", "$"},
	"TRY": {1, "TL", "TL"},
	"EUR": {3, "EUR", "€"},
	"GBP": {4, "GBP", "£"},
}

// GetCurrentBalance returns the balance in one currency (USD, TRY, EUR,
// GBP). An empty currency means "USD", php-dna's default. It uses the
// legacy OperationResult envelope, also on failure:
//
//	{"ErrorCode": 0, "OperationMessage": "Command completed succesfully.",
//	 "OperationResult": "SUCCESS", "Balance": "601.64", "CurrencyId": 2,
//	 "CurrencyInfo": nil, "CurrencyName": "USD", "CurrencySymbol": "$"}
func (c *Client) GetCurrentBalance(currency string) map[string]any {
	if currency == "" {
		currency = "USD"
	}
	name := asciiUpper(currency)
	mapped, ok := currencyMap[name]
	if !ok {
		mapped = currencyInfo{0, name, ""}
	}
	response, err := c.request("GET", "deposit/accounts/me", params{{"currency", name}})
	if err != nil {
		return map[string]any{
			"ErrorCode":        formatErrorCode(err.code),
			"OperationMessage": err.message,
			"OperationResult":  "FAILED",
			"Balance":          "0.00",
			"CurrencyId":       mapped.id,
			"CurrencyInfo":     nil,
			"CurrencyName":     mapped.name,
			"CurrencySymbol":   mapped.symbol,
		}
	}
	balanceKey := asciiLower(currency) + "Balance"
	return map[string]any{
		"ErrorCode":        int64(0),
		"OperationMessage": "Command completed succesfully.",
		"OperationResult":  "SUCCESS",
		"Balance":          numberFormat(phpFloat(coalesce(dig(response, balanceKey), int64(0))), 2),
		"CurrencyId":       mapped.id,
		"CurrencyInfo":     nil,
		"CurrencyName":     mapped.name,
		"CurrencySymbol":   mapped.symbol,
	}
}

// CheckAvailability checks every domain with every extension in one bulk
// search. An empty command means "create". On success it returns a bare
// list ([]any) of
//
//	{"TLD": "com", "DomainName": "hello", "Status": "available"|"notavailable",
//	 "Command": "create", "Period": 1, "IsFee": false, "Price": "11.9900",
//	 "Currency": "USD", "Reason": nil}
//
// and on failure the usual error envelope (map[string]any).
func (c *Client) CheckAvailability(domains, extensions []string, period int, command string) any {
	if command == "" {
		command = "create"
	}
	var queries []any
	for _, d := range domains {
		for _, tld := range extensions {
			queries = append(queries, map[string]any{"domainName": d + "." + strings.TrimLeft(tld, ".")})
		}
	}
	if len(queries) == 0 {
		return map[string]any{
			"result": ResultError,
			"error":  setError("AVAILABILITY", "Domain names not found.", "No domain names provided for availability check"),
		}
	}
	response, err := c.request("POST", "domains/bulk-search", queries)
	if err != nil {
		return c.errorEnvelope(err)
	}

	out := []any{}
	items := coalesce(dig(response, "infos"), response)
	for _, item := range phpValues(items) {
		if !isArray(item) {
			continue
		}
		domainRaw := phpString(coalesce(dig(item, "domainName"), ""))
		var tld string
		if t := dig(item, "tld"); t != nil {
			tld = asciiLower(phpString(t))
		} else if i := strings.LastIndexByte(domainRaw, '.'); i >= 0 {
			tld = asciiLower(domainRaw[i+1:])
		}
		domainName := strings.ReplaceAll(asciiLower(domainRaw), "."+tld, "")
		status := asciiLower(phpString(coalesce(dig(item, "status"), "")))
		available := status == "available" || status == "1" || status == "true"

		st := "notavailable"
		if available {
			st = "available"
		}
		var price any
		if p := dig(item, "price"); p != nil {
			price = numberFormat(phpFloat(p), 4)
		}
		out = append(out, map[string]any{
			"TLD":        tld,
			"DomainName": domainName,
			"Status":     st,
			"Command":    command,
			"Period":     coalesce(dig(item, "period"), int64(period)),
			"IsFee":      coalesce(dig(item, "isPremium"), false),
			"Price":      price,
			"Currency":   dig(item, "currency"),
			"Reason":     dig(item, "reason"),
		})
	}
	return out
}

var domainStatusMap = map[int64]string{
	0: "Active", 1: "WaitingForRegistration", 2: "WaitingForDocument",
	3: "WaitingForIncomingTransfer", 4: "TransferredOut", 7: "PendingDelete",
	8: "Deleted", 9: "ConfirmationEmailSend", 11: "WaitingForOutgoingTransfer",
	12: "PendingHold", 15: "ModificationPending", 18: "InCase", 19: "PendingRestore",
}

// normalizeDomainStatus turns the numeric DomainStatus enum into its label;
// labels pass through.
func normalizeDomainStatus(v any) string {
	if isNumeric(v) {
		if s, ok := domainStatusMap[phpInt(v)]; ok {
			return s
		}
	}
	return phpString(v)
}

// GetList lists the domains in the account. extra holds gateway filters
// (DomainName, MaxResultCount, SkipCount, Sorting, ...) and php-dna's legacy
// paging keys PageNumber (0-based), PageSize, OrderColumn and
// OrderDirection, which are translated. nil lists the first 200 domains.
func (c *Client) GetList(extra map[string]any) map[string]any {
	ep, _ := normalizeInput(extra).(map[string]any)
	if ep == nil {
		ep = map[string]any{}
	}
	if ep["PageSize"] != nil {
		ep["MaxResultCount"] = phpInt(ep["PageSize"])
	}
	if ep["PageNumber"] != nil {
		pageSize := phpInt(coalesce(ep["PageSize"], int64(200)))
		ep["SkipCount"] = phpInt(ep["PageNumber"]) * pageSize
	}
	if ep["OrderColumn"] != nil {
		direction := "ASC"
		if asciiUpper(phpString(coalesce(ep["OrderDirection"], "ASC"))) == "DESC" {
			direction = "DESC"
		}
		ep["Sorting"] = phpString(ep["OrderColumn"]) + " " + direction
	}
	delete(ep, "PageNumber")
	delete(ep, "PageSize")
	delete(ep, "OrderColumn")
	delete(ep, "OrderDirection")

	p := params{{"MaxResultCount", int64(200)}, {"SkipCount", int64(0)}}
	for i := range p {
		if v, ok := ep[p[i].k]; ok {
			p[i].v = v
			delete(ep, p[i].k)
		}
	}
	for _, k := range sortedKeys(ep) {
		p = append(p, kv{k, ep[k]})
	}

	response, err := c.request("GET", "domains", p)
	if err != nil {
		return c.errorEnvelope(err)
	}

	domains := []any{}
	if items := dig(response, "items"); items != nil {
		for _, item := range phpValues(items) {
			var status string
			if sc := dig(item, "statusCode"); sc != nil && sc != "" && !isNumeric(sc) {
				status = phpString(sc)
			} else if st := dig(item, "statusText"); st != nil && st != "" && !isNumeric(st) {
				status = phpString(st)
			} else {
				status = normalizeDomainStatus(coalesce(dig(item, "status"), ""))
			}
			domains = append(domains, map[string]any{
				"ID":                      phpInt(coalesce(dig(item, "id"), int64(0))),
				"Status":                  status,
				"DomainName":              coalesce(dig(item, "domainName"), ""),
				"AuthCode":                coalesce(dig(item, "authCode"), ""),
				"LockStatus":              boolString(dig(item, "lockStatus")),
				"PrivacyProtectionStatus": boolString(dig(item, "privacyProtectionStatus")),
				"IsChildNameServer":       boolString(dig(item, "hosts")),
				"Contacts": map[string]any{
					"Billing":        map[string]any{"ID": ""},
					"Technical":      map[string]any{"ID": ""},
					"Administrative": map[string]any{"ID": ""},
					"Registrant":     map[string]any{"ID": ""},
				},
				"Dates": map[string]any{
					"Start":         formatAPIDate(dig(item, "startDate")),
					"Expiration":    formatAPIDate(dig(item, "expirationDate")),
					"RemainingDays": phpInt(coalesce(dig(item, "remainingDay"), int64(0))),
				},
				"NameServers":      coalesce(dig(item, "nameServers"), []any{}),
				"Additional":       []any{},
				"ChildNameServers": []any{},
			})
		}
	}
	return map[string]any{
		"data":       map[string]any{"Domains": domains},
		"result":     ResultOK,
		"TotalCount": phpInt(coalesce(dig(response, "totalCount"), int64(0))),
	}
}

// boolString is php-dna's `!empty($x) ? 'true' : 'false'`.
func boolString(v any) string {
	if phpEmpty(v) {
		return "false"
	}
	return "true"
}

var tldPriceTypes = []kv{
	{"register", "registration"},
	{"renew", "renew"},
	{"transfer", "transfer"},
	{"restore", "restore"},
	{"refund", "refund"},
	{"backorder", "backorder"},
}

// GetTldList returns the TLD list with the pricing matrix. count is the
// number of TLDs (0 means 20, php-dna's default). Prices are strings with
// four decimals, keyed by period:
//
//	{"result": "OK", "data": [{"id": 1, "status": "Active", "maxchar": 63,
//	  "maxperiod": 10, "minchar": 1, "minperiod": 1, "tld": "com",
//	  "pricing": {"registration": {"1": "10.8100"}, ...},
//	  "currencies": {"registration": "USD", ...}}]}
func (c *Client) GetTldList(count int) map[string]any {
	if count == 0 {
		count = 20
	}
	response, err := c.get("products/tlds", params{{"MaxResultCount", int64(count)}, {"SkipCount", int64(0)}})
	if err != nil {
		return c.errorEnvelope(err)
	}

	tldData := []any{}
	id := int64(1)
	if items := dig(response, "items"); isArray(items) {
		for _, tld := range phpValues(items) {
			pricing := map[string]any{}
			currencies := map[string]any{}
			if first := firstElement(dig(tld, "prices")); isArray(first) {
				for _, pt := range tldPriceTypes {
					apiValue := dig(first, pt.k)
					if apiValue == nil {
						continue
					}
					outType := pt.v.(string)
					addPrice := func(info any) {
						period := phpInt(coalesce(dig(info, "period"), int64(1)))
						if period < 1 {
							period = 1
						}
						price := "0.0000"
						if p := dig(info, "price"); p != nil {
							price = numberFormat(phpFloat(p), 4)
						}
						periods, _ := pricing[outType].(map[string]any)
						if periods == nil {
							periods = map[string]any{}
							pricing[outType] = periods
						}
						periods[strconv.FormatInt(period, 10)] = price
						currencies[outType] = coalesce(dig(info, "currency"), "")
					}
					if list, ok := apiValue.([]any); ok && len(list) > 0 {
						for _, info := range list {
							if isArray(info) {
								addPrice(info)
							}
						}
					} else if m, ok := apiValue.(map[string]any); ok && m["0"] != nil {
						for _, k := range sortedNumericKeys(m) {
							if isArray(m[k]) {
								addPrice(m[k])
							}
						}
					} else if isArray(apiValue) {
						addPrice(apiValue)
					}
				}
			}
			tldData = append(tldData, map[string]any{
				"id":         id,
				"status":     coalesce(dig(tld, "status"), "Active"),
				"maxchar":    coalesce(dig(tld, "constraints", "maxLenght"), int64(63)),
				"maxperiod":  coalesce(dig(tld, "maxRegistrationPeriod"), int64(10)),
				"minchar":    coalesce(dig(tld, "constraints", "minLength"), int64(1)),
				"minperiod":  coalesce(dig(tld, "minRegistrationPeriod"), int64(1)),
				"tld":        dig(tld, "name"),
				"pricing":    pricing,
				"currencies": currencies,
			})
			id++
		}
	}
	return map[string]any{"data": tldData, "result": ResultOK}
}

// firstElement is `$v[0]` for a list (or a map with key "0").
func firstElement(v any) any {
	switch t := v.(type) {
	case []any:
		if len(t) > 0 {
			return t[0]
		}
	case map[string]any:
		return t["0"]
	}
	return nil
}

func sortedNumericKeys(m map[string]any) []string {
	keys := sortedKeys(m)
	sort.SliceStable(keys, func(i, j int) bool { return phpInt(keys[i]) < phpInt(keys[j]) })
	return keys
}
