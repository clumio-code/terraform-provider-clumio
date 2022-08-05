// Copyright 2021. Clumio, Inc.

// clumio_policy_rule definition and CRUD implementation.
package clumio_policy_rule

import (
	"context"

	policyRules "github.com/clumio-code/clumio-go-sdk/controllers/policy_rules"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ClumioPolicyRule returns the resource for Clumio Policy Rule.
func ClumioPolicyRule() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Policy Rule Resource used to determine how a policy " +
			"should be assigned to assets.",

		CreateContext: clumioPolicyRuleCreate,
		ReadContext:   clumioPolicyRuleRead,
		UpdateContext: clumioPolicyRuleUpdate,
		DeleteContext: clumioPolicyRuleDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			schemaId: {
				Description: "Policy Rule Id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			schemaName: {
				Description: "The name of the policy rule.",
				Type:        schema.TypeString,
				Required:    true,
			},
			schemaCondition: {
				Description: "The condition of the policy rule. Possible conditions include: " +
					"1) `entity_type` is required and supports `$eq` and `$in` filters. " +
					"2) `aws_account_native_id` and `aws_region` are optional and both support " +
					"`$eq` and `$in` filters. " +
					"3) `aws_tag` is optional and supports `$eq`, `$in`, `$all`, and `$contains` " +
					"filters.",
				Type:     schema.TypeString,
				Required: true,
			},
			schemaBeforeRuleId: {
				Type: schema.TypeString,
				Description: "The policy rule ID before which this policy rule should be " +
					"inserted. An empty value will set the rule to have lowest priority. " +
					"NOTE: If in the Global Organizational Unit, rules can also be prioritized " +
					"against two virtual rules maintained by the system: `asset-level-rule` and " +
					"`child-ou-rule`. `asset-level-rule` corresponds to the priority of Direct " +
					"Assignments (when a policy is applied directly to an asset) whereas " +
					"`child-ou-rule` corresponds to the priority of rules created by child " +
					"organizational units.",
				Required: true,
			},
			schemaPolicyId: {
				Type:        schema.TypeString,
				Description: "The policy ID of the policy to be applied to the assets.",
				Required:    true,
			},
			schemaOrganizationalUnitId: {
				Type: schema.TypeString,
				Description: "The Clumio-assigned ID of the organizational unit" +
					" to be associated with the policy rule.",
				Optional: true,
				Computed: true,
			},
		},
	}
}

// clumioPolicyRuleCreate handles the Create action for the Clumio Policy Rule Resource.
func clumioPolicyRuleCreate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	clumioConfig := common.GetClumioConfigForAPI(client, d)
	pr := policyRules.NewPolicyRulesV1(clumioConfig)
	condition := common.GetStringValue(d, schemaCondition)
	name := common.GetStringValue(d, schemaName)
	beforeRuleId := common.GetStringValue(d, schemaBeforeRuleId)
	priority := &models.RulePriority{
		BeforeRuleId: &beforeRuleId,
	}
	policyId := common.GetStringValue(d, schemaPolicyId)
	action := &models.RuleAction{
		AssignPolicy: &models.AssignPolicyAction{
			PolicyId: &policyId,
		},
	}
	prRequest := &models.CreatePolicyRuleV1Request{
		Action:    action,
		Condition: &condition,
		Name:      &name,
		Priority:  priority,
	}
	res, apiErr := pr.CreatePolicyRule(prRequest)
	if apiErr != nil {
		return diag.Errorf("Error starting task to create policy rule %v. Error: %v",
			name, string(apiErr.Response))
	}
	err := common.PollTask(ctx, client, *res.TaskId, timeoutInSec, intervalInSec)
	if err != nil {
		return diag.Errorf("Error creating policy rule %v. Error: %v",
			d.Id(), err.Error())
	}

	d.SetId(*res.Rule.Id)
	return nil
}

// clumioPolicyRuleUpdate handles the Update action for the Clumio Policy Rule Resource.
func clumioPolicyRuleUpdate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	clumioConfig := common.GetClumioConfigForAPI(client, d)
	pr := policyRules.NewPolicyRulesV1(clumioConfig)
	condition := common.GetStringValue(d, schemaCondition)
	name := common.GetStringValue(d, schemaName)
	beforeRuleId := common.GetStringValue(d, schemaBeforeRuleId)
	priority := &models.RulePriority{
		BeforeRuleId: &beforeRuleId,
	}
	policyId := common.GetStringValue(d, schemaPolicyId)
	action := &models.RuleAction{
		AssignPolicy: &models.AssignPolicyAction{
			PolicyId: &policyId,
		},
	}
	prRequest := &models.UpdatePolicyRuleV1Request{
		Action:    action,
		Condition: &condition,
		Name:      &name,
		Priority:  priority,
	}
	res, apiErr := pr.UpdatePolicyRule(d.Id(), prRequest)
	if apiErr != nil {
		return diag.Errorf("Error starting task to update policy rule %v. Error: %v",
			name, string(apiErr.Response))
	}
	err := common.PollTask(ctx, client, *res.TaskId, timeoutInSec, intervalInSec)
	if err != nil {
		return diag.Errorf("Error updating policy rule %v. Error: %v",
			d.Id(), err.Error())
	}
	d.SetId(*res.Rule.Id)
	return nil
}

// clumioPolicyRuleRead handles the Read action for the Clumio Policy Rule Resource.
func clumioPolicyRuleRead(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	clumioConfig := common.GetClumioConfigForAPI(client, d)
	pr := policyRules.NewPolicyRulesV1(clumioConfig)

	res, apiErr := pr.ReadPolicyRule(d.Id())
	if apiErr != nil {
		return diag.Errorf("Error retrieving policy rule %v. Error: %v",
			d.Get(schemaName).(string), string(apiErr.Response))
	}
	err := d.Set(schemaName, *res.Name)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaName, err)
	}
	err = d.Set(schemaCondition, *res.Condition)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaCondition, err)
	}
	if res.Priority != nil && res.Priority.BeforeRuleId != nil {
		err = d.Set(schemaBeforeRuleId, *res.Priority.BeforeRuleId)
		if err != nil {
			return diag.Errorf(common.SchemaAttributeSetError, schemaBeforeRuleId, err)
		}
	}
	err = d.Set(schemaPolicyId, *res.Action.AssignPolicy.PolicyId)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaPolicyId, err)
	}
	err = d.Set(schemaOrganizationalUnitId, *res.OrganizationalUnitId)
	if err != nil {
		return diag.Errorf(
			common.SchemaAttributeSetError, schemaOrganizationalUnitId, err)
	}
	return nil
}

// clumioPolicyRuleDelete handles the Delete action for the Clumio Policy Rule Resource.
func clumioPolicyRuleDelete(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	clumioConfig := common.GetClumioConfigForAPI(client, d)
	pr := policyRules.NewPolicyRulesV1(clumioConfig)
	res, apiErr := pr.DeletePolicyRule(d.Id())
	if apiErr != nil {
		return diag.Errorf("Error starting task to delete policy rule %v. Error: %v",
			d.Get(schemaName).(string), string(apiErr.Response))
	}
	err := common.PollTask(ctx, client, *res.TaskId, timeoutInSec, intervalInSec)
	if err != nil {
		return diag.Errorf("Error deleting policy rule %v. Error: %v",
			d.Id(), err.Error())
	}
	return nil
}
