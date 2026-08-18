// Copyright 2026. Clumio, Inc.

// This file holds the logic to validate the restore source of the action and to invoke the Clumio
// DynamoDB table restore SDK API.

package clumio_restore_dynamodb_table

import (
	"context"
	"fmt"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// validateRestoreSource validates that exactly one restore source is specified and that the
// point-in-time restore options are consistent. It is called from both ModifyPlan and Invoke.
// Attributes whose values are not yet known are skipped, as they can only be validated once known.
func validateRestoreSource(
	ctx context.Context, model *restoreDynamoDBTableActionModel) diag.Diagnostics {

	var diags diag.Diagnostics

	if model.Source.IsNull() || model.Source.IsUnknown() {
		return diags
	}
	var source sourceModel
	diags.Append(model.Source.As(ctx, &source, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}

	// Validate that exactly one of the restore sources is specified.
	if !source.SecurevaultBackup.IsUnknown() && !source.ContinuousBackup.IsUnknown() &&
		source.SecurevaultBackup.IsNull() == source.ContinuousBackup.IsNull() {
		summary := "Invalid restore source"
		detail := fmt.Sprintf("Exactly one of %s or %s must be specified in %s.",
			schemaSecurevaultBackup, schemaContinuousBackup, schemaSource)
		diags.AddAttributeError(path.Root(schemaSource), summary, detail)
		return diags
	}

	// Validate that the point-in-time restore options are consistent. Either a timestamp is
	// given, or the latest restorable time is requested, but not both.
	if source.ContinuousBackup.IsNull() || source.ContinuousBackup.IsUnknown() {
		return diags
	}
	var continuousBackup continuousBackupModel
	diags.Append(source.ContinuousBackup.As(ctx, &continuousBackup, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return diags
	}
	if continuousBackup.Timestamp.IsUnknown() ||
		continuousBackup.UseLatestRestorableTime.IsUnknown() {
		return diags
	}
	hasTimestamp := !continuousBackup.Timestamp.IsNull()
	useLatest := continuousBackup.UseLatestRestorableTime.ValueBool()
	if hasTimestamp == useLatest {
		summary := "Invalid point-in-time restore options"
		detail := fmt.Sprintf("Exactly one of %s or %s (set to true) must be specified in %s.",
			schemaTimestamp, schemaUseLatestRestorableTime, schemaContinuousBackup)
		diags.AddAttributeError(
			path.Root(schemaSource).AtName(schemaContinuousBackup), summary, detail)
	}
	return diags
}

// restoreDynamoDBTable validates the action model, invokes the API to initiate the restore of the
// DynamoDB table and returns the identifier of the Clumio task tracking the restore. Called from
// Invoke after Terraform config values have been decoded into the action model.
func (r *clumioRestoreDynamoDBTableAction) restoreDynamoDBTable(
	ctx context.Context, model *restoreDynamoDBTableActionModel) (string, diag.Diagnostics) {

	// All attribute values are known during the invoke phase so the validation is conclusive.
	diags := validateRestoreSource(ctx, model)
	if diags.HasError() {
		return "", diags
	}

	// Convert the action model into the Clumio API request.
	var source sourceModel
	diags.Append(model.Source.As(ctx, &source, basetypes.ObjectAsOptions{})...)
	var target targetModel
	diags.Append(model.Target.As(ctx, &target, basetypes.ObjectAsOptions{})...)
	if diags.HasError() {
		return "", diags
	}

	restoreSource := &models.DynamoDBTableRestoreSource{}
	if !source.SecurevaultBackup.IsNull() {
		var securevaultBackup securevaultBackupModel
		diags.Append(
			source.SecurevaultBackup.As(ctx, &securevaultBackup, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return "", diags
		}
		restoreSource.SecurevaultBackup = &models.DynamoDBRestoreSourceBackupOptions{
			BackupId: securevaultBackup.BackupID.ValueStringPointer(),
		}
	} else {
		var continuousBackup continuousBackupModel
		diags.Append(
			source.ContinuousBackup.As(ctx, &continuousBackup, basetypes.ObjectAsOptions{})...)
		if diags.HasError() {
			return "", diags
		}
		restoreSource.ContinuousBackup = &models.DynamoDBRestoreSourcePitrOptions{
			TableId:                 continuousBackup.TableID.ValueStringPointer(),
			Timestamp:               continuousBackup.Timestamp.ValueStringPointer(),
			UseLatestRestorableTime: continuousBackup.UseLatestRestorableTime.ValueBoolPointer(),
			ClumioType:              continuousBackup.ClumioType.ValueStringPointer(),
		}
	}
	restoreRequest := models.RestoreAwsDynamodbTableV1Request{
		Source: restoreSource,
		Target: &models.DynamoDBTableRestoreTarget{
			EnvironmentId: target.EnvironmentID.ValueStringPointer(),
			TableName:     target.TableName.ValueStringPointer(),
		},
	}

	// Call the Clumio API to initiate the restore of the DynamoDB table.
	res, apiErr := r.restoreClient.RestoreAwsDynamodbTable(nil, restoreRequest)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to invoke %s", r.name)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return "", diags
	}
	if res == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return "", diags
	}

	// A successful async restore always returns a task ID; treat its absence as an error, matching
	// how the AWS connection resource handles a missing task ID.
	if res.TaskId == nil || *res.TaskId == "" {
		diags.AddError(common.NilErrorMessageSummary, common.NilErrorMessageDetail)
		return "", diags
	}
	return *res.TaskId, diags
}
