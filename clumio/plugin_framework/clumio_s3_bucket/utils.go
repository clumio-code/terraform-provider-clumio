// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_s3_bucket Terraform resource.

package clumio_s3_bucket

import (
	"context"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// populateS3BucketsInDataSourceModel is used to populate the users schema attribute in the data
// source model from the results in the API response.
func populateS3BucketsInDataSourceModel(ctx context.Context,
	model *clumioS3BucketDataSourceModel,
	items []*models.Bucket) diag.Diagnostics {

	var diags diag.Diagnostics

	objtype := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaId:                            types.StringType,
			schemaName:                          types.StringType,
			schemaRegion:                        types.StringType,
			schemaAccountNativeId:               types.StringType,
			schemaProtectionGroupCount:          types.Int64Type,
			schemaEventBridgeEnabled:            types.BoolType,
			schemaLastBackupTimestamp:           types.StringType,
			schemaLastContinuousBackupTimestamp: types.StringType,
		},
	}
	attrVals := make([]attr.Value, 0)
	for _, item := range items {
		attrTypes := make(map[string]attr.Type)
		attrTypes[schemaId] = types.StringType
		attrTypes[schemaName] = types.StringType
		attrTypes[schemaAccountNativeId] = types.StringType
		attrTypes[schemaRegion] = types.StringType
		attrTypes[schemaProtectionGroupCount] = types.Int64Type
		attrTypes[schemaEventBridgeEnabled] = types.BoolType
		attrTypes[schemaLastBackupTimestamp] = types.StringType
		attrTypes[schemaLastContinuousBackupTimestamp] = types.StringType

		attrValues := make(map[string]attr.Value)
		attrValues[schemaId] = basetypes.NewStringPointerValue(item.Id)
		attrValues[schemaName] = basetypes.NewStringPointerValue(item.Name)
		attrValues[schemaAccountNativeId] = basetypes.NewStringPointerValue(item.AccountNativeId)
		attrValues[schemaRegion] = basetypes.NewStringPointerValue(item.AwsRegion)
		attrValues[schemaProtectionGroupCount] = basetypes.NewInt64PointerValue(
			item.ProtectionGroupCount)
		attrValues[schemaEventBridgeEnabled] = basetypes.NewBoolPointerValue(
			item.EventBridgeEnabled)
		attrValues[schemaLastBackupTimestamp] = basetypes.NewStringPointerValue(
			item.LastBackupTimestamp)
		attrValues[schemaLastContinuousBackupTimestamp] = basetypes.NewStringPointerValue(
			item.LastContinuousBackupTimestamp)

		obj, conversionDiags := types.ObjectValue(attrTypes, attrValues)
		diags.Append(conversionDiags...)
		if diags.HasError() {
			return diags
		}
		attrVals = append(attrVals, obj)
	}
	setObj, listdiag := types.SetValue(objtype, attrVals)
	diags.Append(listdiag...)
	model.S3Buckets = setObj

	return diags
}
