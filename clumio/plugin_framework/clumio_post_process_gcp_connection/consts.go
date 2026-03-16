// Copyright (c) 2025 Clumio, a Commvault Company All Rights Reserved

package clumio_post_process_gcp_connection

const (
	// Constants used by the resource model for the clumio_post_process_gcp_connection Terraform resource.
	// These values should match the schema tfsdk tags on the resource model struct in schema.go.
	schemaID                  = "id"
	schemaProjectID           = "project_id"
	schemaProjectName         = "project_name"
	schemaProjectNumber       = "project_number"
	schemaToken               = "token"
	schemaServiceAccountEmail = "service_account_email"
	schemaWifPoolId           = "wif_pool_id"
	schemaWifProviderId       = "wif_provider_id"
	schemaConfigVersion       = "config_version"
	schemaProtectGcsVersion   = "protect_gcs_version"
	schemaProperties          = "properties"
)

// RequestType used by GCP post process API
const (
	createRequestType = "CREATE"
	updateRequestType = "UPDATE"
	deleteRequestType = "DELETE"
)
