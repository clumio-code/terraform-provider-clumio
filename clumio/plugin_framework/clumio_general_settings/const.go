// Copyright 2025. Clumio, Inc.

package clumio_general_settings

const (
	// Constants used by the resource model for the clumio_general_settings Terraform resource.
	// These values should match the schema tfsdk tags on the resource model struct in schema.go.
	schemaAutoLogoutDuration = "auto_logout_duration"
	schemaIPAllowlist        = "ip_allowlist"
	schemaPasswordExpiration = "password_expiration_duration"
)

var (
	defaultAutoLogoutDuration = int64(900)     // 15 minutes
	defaultPasswordExpiration = int64(7776000) // 90 days
	defaultIPAllow            = "0.0.0.0/0"    // Allow all IP addresses
)
