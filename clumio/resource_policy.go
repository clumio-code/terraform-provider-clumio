// Copyright 2021. Clumio, Inc.

// clumio_policy resource definition and CRUD implementation.
package clumio

import (
	"context"

	policyDefinitions "github.com/clumio-code/clumio-go-sdk/controllers/policy_definitions"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

const (
	keyActivationStatus = "activation_status"
	keyName = "name"
	keyOperations = "operations"
	keyOrgUnitId = "organizational_unit_id"
	keyActionSetting = "action_setting"
	keyOperationType = "type"
	keyBackupWindow = "backup_window"
	keySlas = "slas"
	keyStartTime = "start_time"
	keyEndTime = "end_time"
	keyRetentionDuration = "retention_duration"
	keyRpoFrequency = "rpo_frequency"
	keyUnit = "unit"
	keyValue = "value"
)

// clumioPolicy returns the resource for Clumio Policy Definition.
func clumioPolicy() *schema.Resource {
	resUnitValue := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"unit": {
				Type: schema.TypeString,
				Required: true,
				Description: "The measurement unit of the SLA parameter. Values include" +
					" hours, days, months, and years.",
			},
			"value": {
				Type: schema.TypeInt,
				Required: true,
				Description: "The measurement value of the SLA parameter.",
			},
		},
	}
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Policy Resource used to schedule backups on Clumio supported" +
			" data sources.",

		CreateContext: clumioPolicyCreate,
		ReadContext:   clumioPolicyRead,
		UpdateContext: clumioPolicyUpdate,
		DeleteContext: clumioPolicyDelete,

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "Policy Id.",
				Type: schema.TypeString,
				Computed: true,
			},
			"lock_status": {
				Description: "Policy Lock Status.",
				Type: schema.TypeString,
				Computed: true,
			},
			"name": {
				Description: "The unique name of the policy.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"activation_status": {
				Type:        schema.TypeString,
				Description: "The status of the policy. Valid values are:" +
					"activated: Backups will take place regularly according to the policy SLA." +
					"deactivated: Backups will not begin until the policy is reactivated." +
					" The assets associated with the policy will have their compliance status set to deactivated.",
				Optional:    true,
			},
			"organizational_unit_id": {
				Type: schema.TypeString,
				Description: "The Clumio-assigned ID of the organizational unit associated with the policy.",
				Optional: true,
			},
			"operations": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"action_setting": {
							Type: schema.TypeString,
							Description: "Determines whether the policy should take action" +
								" now or during the specified backup window. Valid values:" +
								"immediate: to start backup process immediately" +
								"window: to start backup in the specified window",
							Required: true,
						},
						"type": {
							Type: schema.TypeString,
							Description: "The operation to be performed for this SLA set." +
								"Each SLA set corresponds to one and only one operation.",
							Required: true,
						},
						"backup_window":{
							Type:        schema.TypeList,
							Optional:    true,
							MaxItems:    1,
							Description: "The start and end times for the customized" +
								" backup window.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"start_time": {
										Type: schema.TypeString,
										Description: "The time when the backup window opens." +
											" Specify the start time in the format hh:mm," +
											" where hh represents the hour of the day and" +
											" mm represents the minute of the day based on" +
											" the 24 hour clock.",
											Required: true,
									},
									"end_time": {
										Type: schema.TypeString,
										Description: "The time when the backup window closes." +
											" Specify the end time in the format hh:mm," +
											" where hh represents the hour of the day and" +
											" mm represents the minute of the day based on" +
											" the 24 hour clock.",
										Required: true,
									},
								},
							},
						},
						"slas": {
							Type: schema.TypeList,
							Required: true,
							Description: "The service level agreement (SLA) for the policy." +
								" A policy can include one or more SLAs. For example, " +
								"a policy can retain daily backups for a month each, " +
								"and monthly backups for a year each.",
							Elem: &schema.Resource{
								Schema: map[string]*schema.Schema{
									"retention_duration": {
										Type: schema.TypeList,
										Required: true,
										MaxItems: 1,
										Description: "The retention time for this SLA. " +
											"For example, to retain the backup for 1 month," +
											" set unit=months and value=1.",
										Elem: resUnitValue,
									},
									"rpo_frequency": {
										Type: schema.TypeList,
										MaxItems: 1,
										Required: true,
										Description: "The minimum frequency between " +
											"backups for this SLA. Also known as the " +
											"recovery point objective (RPO) interval. For" +
											" example, to configure the minimum frequency" +
											" between backups to be every 2 days, set " +
											"unit=days and value=2. To configure the SLA " +
											"for on-demand backups, set unit=on_demand " +
											"and leave the value field empty.",
										Elem: resUnitValue,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

// clumioPolicyCreate handles the Create action for the Clumio Callback Resource.
func clumioPolicyCreate(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	pd := policyDefinitions.NewPolicyDefinitionsV1(client.clumioConfig)
	activationStatus := getStringValue(d, keyActivationStatus)
	name := getStringValue(d, keyName)
	operationsVal, ok := d.GetOk(keyOperations)
	if !ok{
		return diag.Errorf("Operations is a required attribute")
	}
	policyOperations, _ := mapSchemaOperationsToClumioOperations(operationsVal)
	orgUnitId := getStringValue(d, keyOrgUnitId)
	pdRequest := &models.CreatePolicyDefinitionV1Request{
		ActivationStatus:     &activationStatus,
		Name:                 &name,
		Operations:           policyOperations,
		OrganizationalUnitId: &orgUnitId,
	}
	res, apiErr := pd.CreatePolicyDefinition(pdRequest)
	if apiErr != nil{
		return diag.Errorf("Error creating policy definition %v. Error: %v",
			d.Id(), apiErr.Response)
	}
	d.SetId(*res.Id)
	return nil
}

// clumioPolicyRead handles the Read action for the Clumio Policy Resource.
func clumioPolicyRead(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	pd := policyDefinitions.NewPolicyDefinitionsV1(client.clumioConfig)
	res, apiErr := pd.ReadPolicyDefinition(d.Id(), nil)
	if apiErr != nil{
		return diag.Errorf("Error retrieving policy definition %v. Error: %v",
			d.Id(), apiErr.Response)
	}
	err := d.Set("lock_status", res.LockStatus)
	if err != nil{
		return diag.Errorf("Error setting lock status")
	}
	return nil
}

// clumioPolicyUpdate handles the Update action for the Clumio Policy Resource.
func clumioPolicyUpdate(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	pd := policyDefinitions.NewPolicyDefinitionsV1(client.clumioConfig)
	activationStatus := getStringValue(d, keyActivationStatus)
	name := getStringValue(d, keyName)
	operationsVal, ok := d.GetOk(keyOperations)
	if !ok{
		return diag.Errorf("Operations is a required attribute")
	}
	policyOperations, _ := mapSchemaOperationsToClumioOperations(operationsVal)
	orgUnitId := getStringValue(d, keyOrgUnitId)
	pdRequest := &models.UpdatePolicyDefinitionV1Request{
		ActivationStatus:     &activationStatus,
		Name:                 &name,
		Operations:           policyOperations,
		OrganizationalUnitId: &orgUnitId,
	}
	res, apiErr := pd.UpdatePolicyDefinition(d.Id(), nil, pdRequest)
	if apiErr != nil{
		diag.Errorf("Error updating Policy Definition %v. Error: %v",
			d.Id(), apiErr.Response)
	}
	d.SetId(*res.Id)
	return nil
}

// clumioPolicyDelete handles the Delete action for the Clumio Policy Resource.
func clumioPolicyDelete(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	pd := policyDefinitions.NewPolicyDefinitionsV1(client.clumioConfig)
	_, apiErr := pd.DeletePolicyDefinition(d.Id())
	if apiErr != nil{
		return diag.Errorf("Error deleting policy definition %v. Error: %v",
			d.Id(), apiErr.Response)
	}
	return nil
}

// mapSchemaOperationsToClumioOperations maps the schema operations to the Clumio API
// request operations.
func mapSchemaOperationsToClumioOperations(
	operations interface{}) ([]*models.PolicyOperation, error){
	operationsSlice := operations.([]interface{})
	policyOperations := make([]*models.PolicyOperation, 0)
	for _, operation := range operationsSlice{
		operationAttrMap := operation.(map[string]interface{})
		actionSetting := operationAttrMap[keyActionSetting].(string)
		operationType := operationAttrMap[keyOperationType].(string)
		backupWindowIface, ok := operationAttrMap[keyBackupWindow]
		var backupWindow *models.BackupWindow
		if ok {
			schemaBackupWindow := backupWindowIface.([]interface{})[0].(map[string]interface{})
			schemaBackupWindowStartTime := schemaBackupWindow[keyStartTime].(string)
			schemaBackupWindowEndTime := schemaBackupWindow[keyEndTime].(string)
			backupWindow = &models.BackupWindow{
				EndTime:   &schemaBackupWindowEndTime,
				StartTime: &schemaBackupWindowStartTime,
			}
		}
		backupSLAs := make([]*models.BackupSLA, 0)
		slasIface, ok := operationAttrMap[keySlas]
		if ok {
			schemaSlas := slasIface.([]interface{})
			for _, slaIface := range schemaSlas{
				schemaSla := slaIface.(map[string]interface{})
				retDurationIface, ok := schemaSla[keyRetentionDuration]
				var retentionDuration *models.RetentionBackupSLAParam
				if ok{
					schemaRetDuration := retDurationIface.([]interface{})[0].(map[string]interface{})
					unit := schemaRetDuration[keyUnit].(string)
					value := int64(schemaRetDuration[keyValue].(int))
					retentionDuration = &models.RetentionBackupSLAParam{
						Unit:  &unit,
						Value: &value,
					}
				}
				var rpoFrequency *models.RPOBackupSLAParam
				rpoFrequencyIface, ok := schemaSla[keyRpoFrequency]
				if ok{
					schemaRetDuration := rpoFrequencyIface.([]interface{})[0].(map[string]interface{})
					unit := schemaRetDuration[keyUnit].(string)
					value := int64(schemaRetDuration[keyValue].(int))
					rpoFrequency = &models.RPOBackupSLAParam{
						Unit:  &unit,
						Value: &value,
					}
				}
				backupSLA := &models.BackupSLA{
					RetentionDuration: retentionDuration,
					RpoFrequency: rpoFrequency,
				}
				backupSLAs = append(backupSLAs, backupSLA)
			}
		}
		policyOperation := &models.PolicyOperation{
			ActionSetting: &actionSetting,
			BackupWindow:  backupWindow,
			Slas:          backupSLAs,
			ClumioType:    &operationType,
		}
		policyOperations = append(policyOperations, policyOperation)
	}
	return policyOperations, nil
}

