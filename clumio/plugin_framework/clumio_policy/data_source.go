// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the Clumio Policy SDK API to perform read operation and
// set the attributes from the response of the API in the data source model.

package clumio_policy

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"strings"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// readPolicy invokes the API to read the policyDefinitionClient and from the response populates the
// attributes of the policy.
func (r *clumioPolicyDataSource) readPolicy(
	ctx context.Context, model *clumioPolicyDataSourceModel) diag.Diagnostics {

	var diags diag.Diagnostics
	nameFilter := ""
	operationTypesFilter := ""
	activationStatusFilter := ""
	filters := make([]string, 0)

	// Prepare the query filter.
	name := model.Name.ValueString()
	if name != "" {
		nameFilter = fmt.Sprintf(`"name": {"$begins_with":"%s"}`, name)
		filters = append(filters, nameFilter)
	}
	if !model.OperationTypes.IsUnknown() && !model.OperationTypes.IsNull() {
		operationTypes := make([]string, 0)
		conversionDiags := model.OperationTypes.ElementsAs(ctx, &operationTypes, false)
		diags.Append(conversionDiags...)
		operationTypesFilter = fmt.Sprintf(
			`"operations.type": {"$in":["%s"]}`, strings.Join(operationTypes, `","`))
		filters = append(filters, operationTypesFilter)
	}
	activationStatus := model.ActivationStatus.ValueString()
	if activationStatus != "" {
		activationStatusFilter = fmt.Sprintf(`"activation_status": {"$eq":"%s"}`, activationStatus)
		filters = append(filters, activationStatusFilter)
	}
	filter := fmt.Sprintf("{%s}", strings.Join(filters, ","))

	// Call the Clumio API to list the policy definitions.
	res, apiErr := r.policyDefinitionClient.ListPolicyDefinitions(&filter, nil)
	if apiErr != nil {
		summary := fmt.Sprintf("Unable to read %s", r.name)
		detail := common.ParseMessageFromApiError(apiErr)
		diags.AddError(summary, detail)
		return diags
	}
	if res == nil {
		summary := common.NilErrorMessageSummary
		detail := common.NilErrorMessageDetail
		diags.AddError(summary, detail)
		return diags
	}

	// Convert the Clumio API response for the policies into the datasource schema model.
	if res.Embedded != nil && res.Embedded.Items != nil && len(res.Embedded.Items) > 0 {
		objtype := types.ObjectType{
			AttrTypes: map[string]attr.Type{
				schemaId:                   types.StringType,
				schemaName:                 types.StringType,
				schemaActivationStatus:     types.StringType,
				schemaTimezone:             types.StringType,
				schemaOrganizationalUnitId: types.StringType,
				schemaOperationTypes:       types.SetType{ElemType: types.StringType},
			},
		}
		attrVals := make([]attr.Value, 0)
		for _, item := range res.Embedded.Items {
			resOperationTypes := make([]string, 0)
			for _, operation := range item.Operations {
				resOperationTypes = append(resOperationTypes, *operation.ClumioType)
			}
			opTypes, conversionDiags := types.SetValueFrom(ctx, types.StringType, resOperationTypes)
			diags.Append(conversionDiags...)
			if diags.HasError() {
				return diags
			}
			attrTypes := make(map[string]attr.Type)
			attrTypes[schemaId] = types.StringType
			attrTypes[schemaName] = types.StringType
			attrTypes[schemaActivationStatus] = types.StringType
			attrTypes[schemaTimezone] = types.StringType
			attrTypes[schemaOrganizationalUnitId] = types.StringType
			attrTypes[schemaOperationTypes] = types.SetType{ElemType: types.StringType}

			attrValues := make(map[string]attr.Value)
			attrValues[schemaId] = basetypes.NewStringPointerValue(item.Id)
			attrValues[schemaName] = basetypes.NewStringPointerValue(item.Name)
			attrValues[schemaActivationStatus] = basetypes.NewStringPointerValue(
				item.ActivationStatus)
			attrValues[schemaTimezone] = basetypes.NewStringPointerValue(item.Timezone)
			attrValues[schemaOrganizationalUnitId] = basetypes.NewStringPointerValue(
				item.OrganizationalUnitId)
			attrValues[schemaOperationTypes] = opTypes
			obj, conversionDiags := types.ObjectValue(attrTypes, attrValues)
			diags.Append(conversionDiags...)
			if diags.HasError() {
				return diags
			}
			diags.Append(conversionDiags...)
			attrVals = append(attrVals, obj)
		}
		setObj, listdiag := types.SetValue(objtype, attrVals)
		diags.Append(listdiag...)
		model.Policies = setObj
	}
	return diags
}
