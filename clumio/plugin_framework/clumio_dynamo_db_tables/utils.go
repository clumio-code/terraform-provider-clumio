// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_dynamo_db_tables Terraform resource.

package clumio_dynamo_db_tables

import (
	"context"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// populateDynamoDBTablesInDataSourceModel is used to populate the users schema attribute in the
// data source model from the results in the API response.
func populateDynamoDBTablesInDataSourceModel(ctx context.Context,
	model *clumioDynamoDBTablesDataSourceModel,
	items []*models.DynamoDBTable) diag.Diagnostics {

	var diags diag.Diagnostics

	objtype := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaId:            types.StringType,
			schemaName:          types.StringType,
			schemaTableNativeId: types.StringType,
		},
	}
	attrVals := make([]attr.Value, 0)
	for _, item := range items {
		attrTypes := make(map[string]attr.Type)
		attrTypes[schemaId] = types.StringType
		attrTypes[schemaName] = types.StringType
		attrTypes[schemaTableNativeId] = types.StringType

		attrValues := make(map[string]attr.Value)
		attrValues[schemaId] = basetypes.NewStringPointerValue(item.Id)
		attrValues[schemaName] = basetypes.NewStringPointerValue(item.Name)
		attrValues[schemaTableNativeId] = basetypes.NewStringPointerValue(item.TableNativeId)

		obj, conversionDiags := types.ObjectValue(attrTypes, attrValues)
		diags.Append(conversionDiags...)
		if diags.HasError() {
			return diags
		}
		attrVals = append(attrVals, obj)
	}
	setObj, listdiag := types.SetValue(objtype, attrVals)
	diags.Append(listdiag...)
	model.DynamoDBTables = setObj

	return diags
}
