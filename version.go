package dna

// Version is the version of this library.
const Version = "1.0.0"

// ApplicationName identifies this library in telemetry events.
const ApplicationName = "DomainNameAPI-Go"

// UserAgent is sent with every gateway request. The gateway sits behind
// Cloudflare, which rejects default client agents (Go-http-client) with
// error 1010, so the header is always set.
const UserAgent = ApplicationName + "/" + Version
