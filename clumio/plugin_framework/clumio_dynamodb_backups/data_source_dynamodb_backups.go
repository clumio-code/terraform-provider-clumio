// Copyright 2026. Clumio, Inc.

// This file holds the datasource implementation for the clumio_dynamodb_backups Terraform
// datasource. This datasource is used to retrieve the backups of a DynamoDB table based on the
// specified attributes.

package clumio_dynamodb_backups

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &clumioDynamoDBBackupsDataSource{}
	_ datasource.DataSourceWithConfigure = &clumioDynamoDBBackupsDataSource{}
)

// clumioDynamoDBBackupsDataSource is the struct backing the clumio_dynamodb_backups Terraform
// datasource. It holds the Clumio API client and any other required state needed to read the
// DynamoDB table backups within Clumio.
type clumioDynamoDBBackupsDataSource struct {
	name         string
	client       *common.ApiClient
	backupClient sdkclients.BackupDynamoDBTableClient
}

// NewClumioDynamoDBBackupsDataSource creates a new instance of clumioDynamoDBBackupsDataSource.
// Its attributes are initialized later by Terraform via Metadata and Configure once the Provider
// is initialized.
func NewClumioDynamoDBBackupsDataSource() datasource.DataSource {
	return &clumioDynamoDBBackupsDataSource{}
}

// Metadata returns the name of the datasource type. This is used by Terraform configurations to
// instantiate the datasource.
func (r *clumioDynamoDBBackupsDataSource) Metadata(
	_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	r.name = req.ProviderTypeName + "_dynamodb_backups"
	resp.TypeName = r.name
}

// Configure sets up the datasource with the Clumio API client and any other required state. It is
// called by Terraform once the Provider is initialized.
func (r *clumioDynamoDBBackupsDataSource) Configure(
	_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	r.client = req.ProviderData.(*common.ApiClient)
	r.backupClient = sdkclients.NewBackupDynamoDBTableClient(r.client.ClumioConfig)
}

// Read retrieves the datasource from the Clumio API and sets the Terraform state.
func (r *clumioDynamoDBBackupsDataSource) Read(
	ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	// Retrieve the schema from the current Terraform configuration.
	var state clumioDynamoDBBackupsDataSourceModel
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = r.readDynamoDBBackups(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the schema into the Terraform state.
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}
