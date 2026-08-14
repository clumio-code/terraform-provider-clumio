// Copyright 2026. Clumio, Inc.

// This file holds the type definitions and Schema action function used by the action model for the
// clumio_restore_dynamodb_table Terraform action.

package clumio_restore_dynamodb_table

import (
	"context"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// restoreDynamoDBTableActionModel is the action model for the clumio_restore_dynamodb_table
// Terraform action. It represents the schema of the action and the data it holds. This schema is
// used by customers to configure the action and by the Clumio provider to read the configuration.
type restoreDynamoDBTableActionModel struct {
	Source types.Object `tfsdk:"source"`
	Target types.Object `tfsdk:"target"`
}

// sourceModel is the model for the "source" attribute of the action. Exactly one of its attributes
// must be specified.
type sourceModel struct {
	SecurevaultBackup types.Object `tfsdk:"securevault_backup"`
	ContinuousBackup  types.Object `tfsdk:"continuous_backup"`
}

// securevaultBackupModel is the model for the "securevault_backup" attribute of the source.
type securevaultBackupModel struct {
	BackupID types.String `tfsdk:"backup_id"`
}

// continuousBackupModel is the model for the "continuous_backup" attribute of the source. Exactly
// one of the "timestamp" or "use_latest_restorable_time" attributes must be specified.
type continuousBackupModel struct {
	TableID                 types.String `tfsdk:"table_id"`
	Timestamp               types.String `tfsdk:"timestamp"`
	UseLatestRestorableTime types.Bool   `tfsdk:"use_latest_restorable_time"`
	ClumioType              types.String `tfsdk:"clumio_type"`
}

// targetModel is the model for the "target" attribute of the action.
type targetModel struct {
	EnvironmentID types.String `tfsdk:"environment_id"`
	TableName     types.String `tfsdk:"table_name"`
}

// Schema defines the structure and constraints of the clumio_restore_dynamodb_table Terraform
// action. The schema defines the source of the restore, either a SecureVault backup or continuous
// backup (point-in-time restore), and the target destination of the restore.
func (r *clumioRestoreDynamoDBTableAction) Schema(
	_ context.Context, _ action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			schemaSource: schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					schemaSecurevaultBackup: schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							schemaBackupId: schema.StringAttribute{
								Description: "Unique identifier of the DynamoDB table backup " +
									"to be restored. Use the clumio_dynamodb_backups data " +
									"source to fetch valid values.",
								Required: true,
							},
						},
						Description: "The parameters for initiating a DynamoDB table restore " +
							"from a SecureVault backup.",
						Optional: true,
					},
					schemaContinuousBackup: schema.SingleNestedAttribute{
						Attributes: map[string]schema.Attribute{
							schemaTableId: schema.StringAttribute{
								Description: "Unique identifier of the DynamoDB table to be " +
									"restored. Use the clumio_dynamodb_tables data source to " +
									"fetch valid values.",
								Required: true,
							},
							schemaTimestamp: schema.StringAttribute{
								Description: "The point in time to be restored in RFC-3339 " +
									"format. Exactly one of timestamp or " +
									"use_latest_restorable_time must be specified.",
								Optional: true,
								Validators: []validator.String{
									common.IsTimestamp(),
								},
							},
							schemaUseLatestRestorableTime: schema.BoolAttribute{
								Description: "If set to true, the table is restored to the " +
									"latest possible time. Exactly one of timestamp or " +
									"use_latest_restorable_time must be specified.",
								Optional: true,
							},
							schemaClumioType: schema.StringAttribute{
								Description: "The type of the continuous backup. Possible " +
									"values include `clumio_pitr` and `aws_pitr`. If not " +
									"specified, the type is assumed to be `aws_pitr`.",
								Optional: true,
								Validators: []validator.String{
									stringvalidator.OneOf(clumioTypeClumioPitr, clumioTypeAwsPitr),
								},
							},
						},
						Description: "The parameters for initiating a DynamoDB table " +
							"point-in-time restore.",
						Optional: true,
					},
				},
				Description: "The source of the restore. Exactly one of securevault_backup or " +
					"continuous_backup must be specified.",
				Required: true,
			},
			schemaTarget: schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					schemaEnvironmentId: schema.StringAttribute{
						Description: "Unique identifier of the AWS environment to be used as " +
							"the restore destination. Use the clumio_aws_environment data " +
							"source to fetch valid values.",
						Required: true,
					},
					schemaTableName: schema.StringAttribute{
						Description: "The name of the restored DynamoDB table.",
						Required:    true,
					},
				},
				Description: "The destination of the restore.",
				Required:    true,
			},
		},
		Description: "clumio_restore_dynamodb_table action is used to initiate a restore of a " +
			"DynamoDB table from either a SecureVault backup or continuous backup " +
			"(point-in-time restore). The restore runs asynchronously: the action reports the " +
			"identifier of the Clumio task tracking the restore and does not wait for the " +
			"restore to complete. Invoking actions requires Terraform 1.14 or later.",
	}
}
