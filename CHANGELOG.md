# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2026-11-09

### Added

- First release: a faithful Go port of php-dna 4.0.4, the DomainNameAPI REST client.
- `dna.New(resellerID, apiKey, testMode, ...Option)` with trimmed credentials and a test-mode flag that
  also accepts the strings `"true"`, `"1"`, `"on"` and `"yes"`.
- All php-dna methods with the same requests on the wire and the same `map[string]any` / `[]any`
  return values: account (`GetResellerDetails`, `GetCurrentBalance`), search and lists
  (`CheckAvailability`, `GetList`, `GetTldList`), domains (`GetDetails`, `SyncFromRegistry`,
  `GetContacts`, `SaveContacts`, `ModifyNameServer`, theft-protection lock, privacy, `Renew`,
  `RegisterWithContactInfo`), child nameservers, transfers (`Transfer`, `CancelTransfer`,
  `ApproveTransfer`, `RejectTransfer`, `CheckTransfer`), DNS zone management (record and record-set
  methods with zone-call pacing), URL forwarding, utilities (`IsTrTLD`, `ValidateContact`),
  diagnostics and settings.
- php-dna's error mapping: `API_<status>` codes, synthesised messages for empty error bodies,
  flattened validation errors, ABP boilerplate replacement, 301/302 as credential errors and
  transport failures (timeout, DNS, TLS, reset) as errors with curl-style codes.
- `.tr` TRABIS attribute normalization and contact completion/normalization (names, phone numbers,
  country names).
- Sentry telemetry (errors and sampled performance records) with credential masking, a minimal mode
  and an off switch; the DSN is empty until the DomainNameAPI-Go Sentry project exists.
- Ordered `dna.Contacts` input with `AllRoles` and `ContactsFromMap` helpers, and `RegisterOptions`
  for `RegisterWithContactInfo`'s optional arguments.
- Injectable HTTP transport (`WithHTTPClient`), `User-Agent: DomainNameAPI-Go/1.0.0`, 30 s / 10 s
  timeouts and optional TLS verification opt-out.
- Conformance test running all 126 php-dna fixtures, unit tests, runnable examples for every method,
  README in nine languages and GitHub Actions CI.

[1.0.0]: https://github.com/domainreseller/go-dna/releases/tag/v1.0.0
