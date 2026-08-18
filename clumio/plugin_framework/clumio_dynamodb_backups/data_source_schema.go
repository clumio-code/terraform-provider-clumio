// Copyright 2026. Clumio, Inc.

// This file holds the type definition and Schema datasource function used by the datasource model
// for the clumio_dynamodb_backups Terraform datasource.

package clumio_dynamodb_backups

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// clumioDynamoDBBackupsDataSourceModel is the datasource model for the clumio_dynamodb_backups
// Terraform datasource. It represents the schema of the datasource and the data it holds. This
// schema is used by customers to configure the datasource and by the Clumio provider to read and
// write the datasource.
type clumioDynamoDBBackupsDataSourceModel struct {
	TableID         types.String `tfsdk:"table_id"`
	ClumioType      types.String `tfsdk:"type"`
	BeforeTimestamp types.String `tfsdk:"before_timestamp"`
	AfterTimestamp  types.String `tfsdk:"after_timestamp"`
	Backups         types.List   `tfsdk:"backups"`
}

// Schema defines the structure and constraints of the clumio_dynamodb_backups Terraform
// datasource. The schema defines various attributes such as the table_id, type, before_timestamp
// and after_timestamp, which are used to determine the DynamoDB table backups to retrieve, and the
// computed attribute 'backups' which holds the matched backups, sorted newest first.
func (r *clumioDynamoDBBackupsDataSource) Schema(
	_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			schemaTableId: schema.StringAttribute{
				Description: "Unique identifier of the DynamoDB table whose backups are to be " +
					"retrieved.",
				Required: true,
			},
			schemaType: schema.StringAttribute{
				Description: "The type of backups to retrieve. Possible values include " +
					"`clumio_backup` and `aws_snapshot`. If not specified, backups of all types " +
					"are retrieved.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.OneOf(backupTypeClumioBackup, backupTypeAwsSnapshot),
				},
			},
			schemaBeforeTimestamp: schema.StringAttribute{
				Description: "If specified, only backups whose start timestamp is at or before " +
					"the given timestamp are retrieved. Represented in RFC-3339 format.",
				Optional: true,
				Validators: []validator.String{
					common.IsTimestamp(),
				},
			},
			schemaAfterTimestamp: schema.StringAttribute{
				Description: "If specified, only backups whose start timestamp is after the " +
					"given timestamp are retrieved. Represented in RFC-3339 format.",
				Optional: true,
				Validators: []validator.String{
					common.IsTimestamp(),
				},
			},
			schemaBackups: schema.ListNestedAttribute{
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						schemaId: schema.StringAttribute{
							Description: "Unique identifier of the DynamoDB table backup.",
							Computed:    true,
						},
						schemaStartTimestamp: schema.StringAttribute{
							Description: "The timestamp of when this backup started. " +
								"Represented in RFC-3339 format.",
							Computed: true,
						},
						schemaExpirationTimestamp: schema.StringAttribute{
							Description: "The timestamp of when this backup expires. " +
								"Represented in RFC-3339 format.",
							Computed: true,
						},
						schemaType: schema.StringAttribute{
							Description: "The type of the backup. Possible values include " +
								"`clumio_backup` and `aws_snapshot`.",
							Computed: true,
						},
					},
				},
				Computed: true,
				Description: "List of DynamoDB table backups which matched the query criteria, " +
					"sorted by start timestamp with the most recent backup first.",
			},
		},
		Description: "clumio_dynamodb_backups data source is used to retrieve details of the " +
			"backups of a DynamoDB table for use in other resources or actions.",
	}
}
