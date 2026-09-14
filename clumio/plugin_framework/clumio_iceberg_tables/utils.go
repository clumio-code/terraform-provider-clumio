// Copyright 2026. Clumio, Inc.

// This file holds various utility functions used by the clumio_iceberg_tables Terraform datasource.

package clumio_iceberg_tables

import (
	"context"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// populateIcebergTablesInDataSourceModel is used to populate the iceberg_tables schema attribute in
// the data source model from the results in the API response.
func populateIcebergTablesInDataSourceModel(ctx context.Context,
	model *clumioIcebergTablesDataSourceModel,
	items []*models.IcebergTable) diag.Diagnostics {

	var diags diag.Diagnostics

	attrTypes := map[string]attr.Type{
		schemaId:          types.StringType,
		schemaName:        types.StringType,
		schemaCatalog:     types.StringType,
		schemaCatalogType: types.StringType,
		schemaNamespace:   types.StringType,
	}
	objtype := types.ObjectType{AttrTypes: attrTypes}

	attrVals := make([]attr.Value, 0)
	for _, item := range items {
		attrValues := map[string]attr.Value{
			schemaId:          basetypes.NewStringPointerValue(item.Id),
			schemaName:        basetypes.NewStringPointerValue(item.Name),
			schemaCatalog:     basetypes.NewStringPointerValue(item.Catalog),
			schemaCatalogType: basetypes.NewStringPointerValue(item.CatalogType),
			schemaNamespace:   basetypes.NewStringPointerValue(item.Namespace),
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
	model.IcebergTables = listObj

	return diags
}
