// Copyright 2024. Clumio, Inc.

// This file holds the logic to invoke the Clumio Policy Rule SDK API to perform read operation and
// set the attributes from the response of the API in the data source model.

package clumio_policy_rule

import (
	"context"
	"fmt"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// readPolicyRule invokes the API to read the policyRuleClient and from the response populates the
// attributes of the policy rule.
func (r *clumioPolicyRuleDataSource) readPolicyRule(
	_ context.Context, model *clumioPolicyRuleDataSourceModel) diag.Diagnostics {

	items, diags := r.listPolicyRules()
	if diags.HasError() {
		return diags
	}

	// Convert the Clumio API response for the policy rules into the datasource schema model.
	objtype := types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaId:           types.StringType,
			schemaName:         types.StringType,
			schemaPolicyId:     types.StringType,
			schemaBeforeRuleId: types.StringType,
			schemaCondition:    types.StringType,
		},
	}
	attrVals := make([]attr.Value, 0)
	modelName := model.Name.ValueString()
	modelPolicyId := model.PolicyId.ValueString()

	for _, item := range items {
		if modelName != "" {
			if *item.Name != modelName {
				continue
			}
		}
		if modelPolicyId != "" {
			if item.Action == nil || item.Action.AssignPolicy == nil ||
				*item.Action.AssignPolicy.PolicyId != modelPolicyId {
				continue
			}
		}
		attrTypes := make(map[string]attr.Type)
		attrTypes[schemaId] = types.StringType
		attrTypes[schemaName] = types.StringType
		attrTypes[schemaPolicyId] = types.StringType
		attrTypes[schemaBeforeRuleId] = types.StringType
		attrTypes[schemaCondition] = types.StringType

		attrValues := make(map[string]attr.Value)
		attrValues[schemaId] = basetypes.NewStringPointerValue(item.Id)
		attrValues[schemaName] = basetypes.NewStringPointerValue(item.Name)
		attrValues[schemaPolicyId] = basetypes.NewStringPointerValue(
			item.Action.AssignPolicy.PolicyId)
		if item.Priority != nil {
			attrValues[schemaBeforeRuleId] = basetypes.NewStringPointerValue(item.Priority.BeforeRuleId)
		} else {
			attrValues[schemaBeforeRuleId] = basetypes.NewStringNull()
		}
		attrValues[schemaCondition] = basetypes.NewStringPointerValue(item.Condition)
		obj, conversionDiags := types.ObjectValue(attrTypes, attrValues)
		diags.Append(conversionDiags...)
		if diags.HasError() {
			return diags
		}
		diags.Append(conversionDiags...)
		attrVals = append(attrVals, obj)
	}
	if len(attrVals) > 0 {
		setObj, listdiag := types.SetValue(objtype, attrVals)
		diags.Append(listdiag...)
		model.PolicyRules = setObj
	}

	return diags
}

// listPolicyRules invokes the SDK API and returns all the policy rules.
func (r *clumioPolicyRuleDataSource) listPolicyRules() (
	[]*models.Rule, diag.Diagnostics) {

	var diags diag.Diagnostics
	items := make([]*models.Rule, 0)
	// Call the Clumio API to list the policy rules.
	limit := int64(1000)
	var start *string
	for {
		res, apiErr := r.sdkPolicyRules.ListPolicyRules(&limit, start, nil, nil, nil)
		if apiErr != nil {
			summary := fmt.Sprintf("Unable to read %s", r.name)
			detail := common.ParseMessageFromApiError(apiErr)
			diags.AddError(summary, detail)
			return nil, diags
		}
		if res == nil {
			summary := common.NilErrorMessageSummary
			detail := common.NilErrorMessageDetail
			diags.AddError(summary, detail)
			return nil, diags
		}

		if res.Embedded != nil && res.Embedded.Items != nil && len(res.Embedded.Items) > 0 {
			items = append(items, res.Embedded.Items...)
		}
		if res.Links == nil || res.Links.Next == nil {
			break
		} else {
			start = res.Links.Next.Href
		}
	}

	return items, nil
}
