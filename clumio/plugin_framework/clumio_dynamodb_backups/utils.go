// Copyright 2026. Clumio, Inc.

// This file holds various utility functions used by the clumio_dynamodb_backups Terraform
// datasource.

package clumio_dynamodb_backups

import (
	"context"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// populateDynamoDBBackupsInDataSourceModel is used to populate the backups schema attribute in the
// data source model from the results in the API response. The order of the items in the API
// response is preserved.
func populateDynamoDBBackupsInDataSourceModel(_ context.Context,
	model *clumioDynamoDBBackupsDataSourceModel,
	items []*models.DynamoDBTableBackupWithETag) diag.Diagnostics {

	var diags diag.Diagnostics

	attrTypes := map[string]attr.Type{
		schemaId:                  types.StringType,
		schemaStartTimestamp:      types.StringType,
		schemaExpirationTimestamp: types.StringType,
		schemaType:                types.StringType,
	}
	objtype := types.ObjectType{
		AttrTypes: attrTypes,
	}
	attrVals := make([]attr.Value, 0)
	for _, item := range items {
		attrValues := map[string]attr.Value{
			schemaId:                  basetypes.NewStringPointerValue(item.Id),
			schemaStartTimestamp:      basetypes.NewStringPointerValue(item.StartTimestamp),
			schemaExpirationTimestamp: basetypes.NewStringPointerValue(item.ExpirationTimestamp),
			schemaType:                basetypes.NewStringPointerValue(item.ClumioType),
		}

		obj, conversionDiags := types.ObjectValue(attrTypes, attrValues)
		diags.Append(conversionDiags...)
		if diags.HasError() {
			return diags
		}
		attrVals = append(attrVals, obj)
	}
	listObj, listdiag := types.ListValue(objtype, attrVals)
	diags.Append(listdiag...)
	model.Backups = listObj

	return diags
}
