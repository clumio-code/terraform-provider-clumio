// Copyright 2026. Clumio, Inc.

// This file contains the unit tests for the functions in action.go

//go:build unit

package clumio_restore_dynamodb_table

import (
	"context"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	securevaultBackupAttrTypes = map[string]attr.Type{
		schemaBackupId: types.StringType,
	}
	continuousBackupAttrTypes = map[string]attr.Type{
		schemaTableId:                 types.StringType,
		schemaTimestamp:               types.StringType,
		schemaUseLatestRestorableTime: types.BoolType,
		schemaClumioType:              types.StringType,
	}
	sourceAttrTypes = map[string]attr.Type{
		schemaSecurevaultBackup: types.ObjectType{AttrTypes: securevaultBackupAttrTypes},
		schemaContinuousBackup:  types.ObjectType{AttrTypes: continuousBackupAttrTypes},
	}
	targetAttrTypes = map[string]attr.Type{
		schemaEnvironmentId: types.StringType,
		schemaTableName:     types.StringType,
	}
)

// buildActionModel builds an action model with the given securevault_backup and continuous_backup
// source attribute values and a valid target.
func buildActionModel(
	securevaultBackup attr.Value, continuousBackup attr.Value) *restoreDynamoDBTableActionModel {

	return &restoreDynamoDBTableActionModel{
		Source: types.ObjectValueMust(sourceAttrTypes, map[string]attr.Value{
			schemaSecurevaultBackup: securevaultBackup,
			schemaContinuousBackup:  continuousBackup,
		}),
		Target: types.ObjectValueMust(targetAttrTypes, map[string]attr.Value{
			schemaEnvironmentId: types.StringValue("test-env-id"),
			schemaTableName:     types.StringValue("test-table-name"),
		}),
	}
}

// buildSecurevaultBackup builds a securevault_backup source attribute value with the given backup
// ID.
func buildSecurevaultBackup(backupId string) attr.Value {
	return types.ObjectValueMust(securevaultBackupAttrTypes, map[string]attr.Value{
		schemaBackupId: types.StringValue(backupId),
	})
}

// buildContinuousBackup builds a continuous_backup source attribute value with the given timestamp
// and use_latest_restorable_time values.
func buildContinuousBackup(timestamp types.String, useLatest types.Bool) attr.Value {
	return types.ObjectValueMust(continuousBackupAttrTypes, map[string]attr.Value{
		schemaTableId:                 types.StringValue("test-table-id"),
		schemaTimestamp:               timestamp,
		schemaUseLatestRestorableTime: useLatest,
		schemaClumioType:              types.StringNull(),
	})
}

// Unit test for the following cases:
//   - Restore DynamoDB table from a SecureVault backup success scenario.
//   - Restore DynamoDB table from continuous backup success scenario.
//   - SDK API for restore DynamoDB table returns an error.
//   - SDK API for restore DynamoDB table returns a nil response.
func TestRestoreDynamoDBTable(t *testing.T) {

	ctx := context.Background()
	restoreClient := sdkclients.NewMockRestoredDynamoDBTableClient(t)
	resourceName := "test_restore_dynamodb_table"
	backupId := "test-backup-id"
	taskId := "test-task-id"
	testError := "Test Error"

	ra := clumioRestoreDynamoDBTableAction{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		restoreClient: restoreClient,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for a restore from a SecureVault backup. It should not return
	// Diagnostics and should return the task ID from the API response.
	t.Run("Basic success scenario for restore from SecureVault backup", func(t *testing.T) {

		model := buildActionModel(
			buildSecurevaultBackup(backupId),
			types.ObjectNull(continuousBackupAttrTypes))

		restoreResponse := &models.RestoreDynamoDBTableResponse{
			TaskId: &taskId,
		}

		// Setup expectations, asserting the exact restore request built from the model.
		restoreClient.EXPECT().RestoreAwsDynamodbTable(mock.Anything,
			mock.MatchedBy(func(body models.RestoreAwsDynamodbTableV1Request) bool {
				return body.Source != nil && body.Source.ContinuousBackup == nil &&
					body.Source.SecurevaultBackup != nil &&
					body.Source.SecurevaultBackup.BackupId != nil &&
					*body.Source.SecurevaultBackup.BackupId == backupId &&
					body.Target != nil && body.Target.EnvironmentId != nil &&
					*body.Target.EnvironmentId == "test-env-id" &&
					body.Target.TableName != nil &&
					*body.Target.TableName == "test-table-name"
			})).Times(1).Return(restoreResponse, nil)

		resTaskId, diags := ra.restoreDynamoDBTable(ctx, model)
		assert.Nil(t, diags)
		assert.Equal(t, taskId, resTaskId)
	})

	// Tests the success scenario for a restore from continuous backup with the latest restorable
	// time. It should not return Diagnostics.
	t.Run("Basic success scenario for restore from continuous backup", func(t *testing.T) {

		model := buildActionModel(
			types.ObjectNull(securevaultBackupAttrTypes),
			buildContinuousBackup(types.StringNull(), types.BoolValue(true)))

		restoreResponse := &models.RestoreDynamoDBTableResponse{
			TaskId: &taskId,
		}

		// Setup expectations, asserting the exact restore request built from the model.
		restoreClient.EXPECT().RestoreAwsDynamodbTable(mock.Anything,
			mock.MatchedBy(func(body models.RestoreAwsDynamodbTableV1Request) bool {
				return body.Source != nil && body.Source.SecurevaultBackup == nil &&
					body.Source.ContinuousBackup != nil &&
					body.Source.ContinuousBackup.TableId != nil &&
					*body.Source.ContinuousBackup.TableId == "test-table-id" &&
					body.Source.ContinuousBackup.Timestamp == nil &&
					body.Source.ContinuousBackup.UseLatestRestorableTime != nil &&
					*body.Source.ContinuousBackup.UseLatestRestorableTime
			})).Times(1).Return(restoreResponse, nil)

		resTaskId, diags := ra.restoreDynamoDBTable(ctx, model)
		assert.Nil(t, diags)
		assert.Equal(t, taskId, resTaskId)
	})

	// Tests that Diagnostics is returned in case the restore DynamoDB table API call returns an
	// error.
	t.Run("restore DynamoDB table returns an error", func(t *testing.T) {

		model := buildActionModel(
			buildSecurevaultBackup(backupId),
			types.ObjectNull(continuousBackupAttrTypes))

		// Setup expectations.
		restoreClient.EXPECT().RestoreAwsDynamodbTable(mock.Anything, mock.Anything).Times(1).
			Return(nil, apiError)

		_, diags := ra.restoreDynamoDBTable(ctx, model)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the restore DynamoDB table API call returns a
	// nil response.
	t.Run("restore DynamoDB table returns a nil response", func(t *testing.T) {

		model := buildActionModel(
			buildSecurevaultBackup(backupId),
			types.ObjectNull(continuousBackupAttrTypes))

		// Setup expectations.
		restoreClient.EXPECT().RestoreAwsDynamodbTable(mock.Anything, mock.Anything).Times(1).
			Return(nil, nil)

		_, diags := ra.restoreDynamoDBTable(ctx, model)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the restore response does not contain a task ID.
	t.Run("restore DynamoDB table returns a response without a task ID", func(t *testing.T) {

		model := buildActionModel(
			buildSecurevaultBackup(backupId),
			types.ObjectNull(continuousBackupAttrTypes))

		// Setup expectations.
		restoreClient.EXPECT().RestoreAwsDynamodbTable(mock.Anything, mock.Anything).Times(1).
			Return(&models.RestoreDynamoDBTableResponse{}, nil)

		_, diags := ra.restoreDynamoDBTable(ctx, model)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Exactly one restore source specified is valid.
//   - Both restore sources specified is invalid.
//   - No restore source specified is invalid.
//   - Continuous backup with both timestamp and use_latest_restorable_time is invalid.
//   - Continuous backup with neither timestamp nor use_latest_restorable_time is invalid.
//   - Unknown attribute values are skipped during validation.
func TestValidateRestoreSource(t *testing.T) {

	ctx := context.Background()

	// Tests that no Diagnostics is returned in case only the securevault_backup restore source is
	// specified.
	t.Run("securevault_backup source is valid", func(t *testing.T) {

		model := buildActionModel(
			buildSecurevaultBackup("test-backup-id"),
			types.ObjectNull(continuousBackupAttrTypes))

		diags := validateRestoreSource(ctx, model)
		assert.Nil(t, diags)
	})

	// Tests that no Diagnostics is returned in case only the continuous_backup restore source is
	// specified with a timestamp.
	t.Run("continuous_backup source with timestamp is valid", func(t *testing.T) {

		model := buildActionModel(
			types.ObjectNull(securevaultBackupAttrTypes),
			buildContinuousBackup(types.StringValue("2026-07-20T00:00:00Z"), types.BoolNull()))

		diags := validateRestoreSource(ctx, model)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case both restore sources are specified.
	t.Run("both restore sources specified is invalid", func(t *testing.T) {

		model := buildActionModel(
			buildSecurevaultBackup("test-backup-id"),
			buildContinuousBackup(types.StringNull(), types.BoolValue(true)))

		diags := validateRestoreSource(ctx, model)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case no restore source is specified.
	t.Run("no restore source specified is invalid", func(t *testing.T) {

		model := buildActionModel(
			types.ObjectNull(securevaultBackupAttrTypes),
			types.ObjectNull(continuousBackupAttrTypes))

		diags := validateRestoreSource(ctx, model)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case both the timestamp and the latest restorable
	// time are specified for a continuous backup restore source.
	t.Run("continuous_backup with timestamp and latest time is invalid", func(t *testing.T) {

		model := buildActionModel(
			types.ObjectNull(securevaultBackupAttrTypes),
			buildContinuousBackup(
				types.StringValue("2026-07-20T00:00:00Z"), types.BoolValue(true)))

		diags := validateRestoreSource(ctx, model)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case neither the timestamp nor the latest restorable
	// time is specified for a continuous backup restore source.
	t.Run("continuous_backup without point in time is invalid", func(t *testing.T) {

		model := buildActionModel(
			types.ObjectNull(securevaultBackupAttrTypes),
			buildContinuousBackup(types.StringNull(), types.BoolNull()))

		diags := validateRestoreSource(ctx, model)
		assert.NotNil(t, diags)
	})

	// Tests that no Diagnostics is returned in case attribute values are not yet known, as they
	// are validated again during the invoke phase once known.
	t.Run("unknown attribute values are skipped", func(t *testing.T) {

		model := buildActionModel(
			types.ObjectUnknown(securevaultBackupAttrTypes),
			types.ObjectNull(continuousBackupAttrTypes))

		diags := validateRestoreSource(ctx, model)
		assert.Nil(t, diags)
	})
}

// buildActionConfig builds the Terraform configuration for the action with the given
// securevault_backup and continuous_backup raw attribute values.
func buildActionConfig(
	ctx context.Context, securevaultBackup any, continuousBackup any) tfsdk.Config {

	securevaultBackupType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			schemaBackupId: tftypes.String,
		},
	}
	continuousBackupType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			schemaTableId:                 tftypes.String,
			schemaTimestamp:               tftypes.String,
			schemaUseLatestRestorableTime: tftypes.Bool,
			schemaClumioType:              tftypes.String,
		},
	}
	sourceType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			schemaSecurevaultBackup: securevaultBackupType,
			schemaContinuousBackup:  continuousBackupType,
		},
	}
	targetType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			schemaEnvironmentId: tftypes.String,
			schemaTableName:     tftypes.String,
		},
	}
	configType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			schemaSource: sourceType,
			schemaTarget: targetType,
		},
	}

	configVals := map[string]tftypes.Value{
		schemaSource: tftypes.NewValue(sourceType, map[string]tftypes.Value{
			schemaSecurevaultBackup: tftypes.NewValue(securevaultBackupType, securevaultBackup),
			schemaContinuousBackup:  tftypes.NewValue(continuousBackupType, continuousBackup),
		}),
		schemaTarget: tftypes.NewValue(targetType, map[string]tftypes.Value{
			schemaEnvironmentId: tftypes.NewValue(tftypes.String, "test-env-id"),
			schemaTableName:     tftypes.NewValue(tftypes.String, "test-table-name"),
		}),
	}

	schemaResp := &action.SchemaResponse{}
	(&clumioRestoreDynamoDBTableAction{}).Schema(ctx, action.SchemaRequest{}, schemaResp)
	return tfsdk.Config{
		Raw:    tftypes.NewValue(configType, configVals),
		Schema: schemaResp.Schema,
	}
}

// Unit test for the action Metadata, Configure, ModifyPlan and Invoke functions.
func TestActionMetadataConfigureModifyPlanInvoke(t *testing.T) {

	ctx := context.Background()
	taskId := "test-task-id"
	securevaultBackupVals := map[string]tftypes.Value{
		schemaBackupId: tftypes.NewValue(tftypes.String, "test-backup-id"),
	}

	// Tests that the action type name is set as part of Metadata().
	t.Run("Metadata sets the action type name", func(t *testing.T) {

		ra := &clumioRestoreDynamoDBTableAction{}
		resp := &action.MetadataResponse{}
		ra.Metadata(ctx, action.MetadataRequest{ProviderTypeName: "clumio"}, resp)
		assert.Equal(t, "clumio_restore_dynamodb_table", resp.TypeName)
	})

	// Tests that Configure() returns early when no provider data is given and sets up the SDK
	// client when it is.
	t.Run("Configure sets up the SDK client", func(t *testing.T) {

		ra := &clumioRestoreDynamoDBTableAction{}
		ra.Configure(ctx, action.ConfigureRequest{}, &action.ConfigureResponse{})
		assert.Nil(t, ra.client)

		ra.Configure(ctx, action.ConfigureRequest{
			ProviderData: &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
		}, &action.ConfigureResponse{})
		assert.NotNil(t, ra.client)
		assert.NotNil(t, ra.restoreClient)
	})

	// Tests that ModifyPlan() returns Diagnostics in case no restore source is specified.
	t.Run("ModifyPlan returns Diagnostics for an invalid restore source", func(t *testing.T) {

		ra := &clumioRestoreDynamoDBTableAction{}
		resp := &action.ModifyPlanResponse{}
		ra.ModifyPlan(ctx, action.ModifyPlanRequest{
			Config: buildActionConfig(ctx, nil, nil),
		}, resp)
		assert.True(t, resp.Diagnostics.HasError())
	})

	// Tests the success scenario for the action Invoke(). It should not return Diagnostics and
	// should send a progress event containing the task ID of the restore.
	t.Run("Invoke reports the task ID as a progress event", func(t *testing.T) {

		restoreClient := sdkclients.NewMockRestoredDynamoDBTableClient(t)
		ra := &clumioRestoreDynamoDBTableAction{
			name:          "clumio_restore_dynamodb_table",
			restoreClient: restoreClient,
		}

		// Setup expectations.
		restoreClient.EXPECT().RestoreAwsDynamodbTable(mock.Anything, mock.Anything).Times(1).
			Return(&models.RestoreDynamoDBTableResponse{TaskId: &taskId}, nil)

		progressMessage := ""
		resp := &action.InvokeResponse{
			SendProgress: func(event action.InvokeProgressEvent) {
				progressMessage = event.Message
			},
		}
		ra.Invoke(ctx, action.InvokeRequest{
			Config: buildActionConfig(ctx, securevaultBackupVals, nil),
		}, resp)
		assert.False(t, resp.Diagnostics.HasError())
		// The message is deliberately terse so the task ID is easy to spot in the Terraform
		// output; assert the exact text to catch it growing back.
		assert.Equal(t, "Restore task ID: "+taskId, progressMessage)
	})
}
