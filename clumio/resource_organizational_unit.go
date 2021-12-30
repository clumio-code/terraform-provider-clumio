// Copyright 2021. Clumio, Inc.

// clumio_organizational_unit definition and CRUD implementation.
package clumio

import (
	"context"
	"fmt"
	"strings"

	orgUnits "github.com/clumio-code/clumio-go-sdk/controllers/organizational_units"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// clumioOrganizationalUnit returns the resource for Clumio AWS Connection.

func clumioOrganizationalUnit() *schema.Resource {
	entitySchema := &schema.Resource{
		Schema: map[string]*schema.Schema{
			"id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The Clumio assigned ID of the entity.",
			},
			"type": {
				Type:        schema.TypeInt,
				Required:    true,
				Description: "The entity type.",
			},
		},
	}
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio AWS Connection Resource used to connect AWS accounts to Clumio.",

		CreateContext: clumioOrganizationalUnitCreate,
		ReadContext:   clumioOrganizationalUnitRead,
		UpdateContext: clumioOrganizationalUnitUpdate,
		DeleteContext: clumioOrganizationalUnitDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "OrganizationalUnit Id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"name": {
				Description: "Unique name assigned to the organizational unit.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "A description of the organizational unit.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"parent_id": {
				Description: "The Clumio-assigned ID of the parent organizational unit" +
					" under which the new organizational unit is to be created.",
				Type:     schema.TypeString,
				Optional: true,
			},
			"users": {
				Description: "List of user ids to assign this organizational unit.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"entities": {
				Description: "List of user ids to assign this organizational unit.",
				Type:        schema.TypeList,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"parent_entity": {
							Description: "The parent object of the primary entity" +
								" associated with the organizational unit. For example," +
								" vmware_vcenter is the parent entity of primary entity" +
								" vmware_vm_folder. The parent object is necessary for" +
								" VMware entities and can be omitted for other" +
								" data sources.",
							Type:     schema.TypeList,
							MaxItems: 1,
							Elem:     entitySchema,
							Optional: true,
						},
						"primary_entity": {
							Description: "The primary object associated with the" +
								" organizational unit. Examples of primary entities" +
								" include aws_environment and vmware_vm.",
							Type:     schema.TypeList,
							MaxItems: 1,
							Elem:     entitySchema,
							Optional: true,
						},
					},
				},
				Optional: true,
			},
			"children_count": {
				Description: "Number of immediate children of the organizational unit.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"configured_datasource_types": {
				Description: "Datasource types configured in this organizational unit." +
					" Possible values include aws, microsoft365, vmware, or mssql.",
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"descendant_ids": {
				Description: "List of all recursive descendant organizational units" +
					" of this OU.",
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"user_count": {
				Description: "Number of users to whom this organizational unit or any" +
					" of its descendants have been assigned.",
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

// clumioOrganizationalUnitCreate handles the Create action for the Clumio OrganizationalUnit Resource.
func clumioOrganizationalUnitCreate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	orgUnitsAPI := orgUnits.NewOrganizationalUnitsV1(client.clumioConfig)
	name := getStringValue(d, "name")
	parentId := getStringValue(d, "parent_id")
	request := &models.CreateOrganizationalUnitV1Request{
		Name:     &name,
		ParentId: &parentId,
	}
	description := getStringValue(d, "description")
	if description != "" {
		request.Description = &description
	}
	users := getStringSlice(d, "users")
	request.Users = users
	//TODO: Handle entities after the entityModel is made public in the SDK
	res, apiErr := orgUnitsAPI.CreateOrganizationalUnit(nil, request)
	if apiErr != nil {
		return diag.Errorf(
			"Error creating Clumio OrganizationalUnit. Error: %v", string(apiErr.Response))
	}
	d.SetId(*res.Id)
	return clumioOrganizationalUnitRead(ctx, d, meta)
}

// clumioOrganizationalUnitRead handles the Read action for the Clumio OrganizationalUnit Resource.
func clumioOrganizationalUnitRead(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	orgUnitsAPI := orgUnits.NewOrganizationalUnitsV1(client.clumioConfig)
	res, apiErr := orgUnitsAPI.ReadOrganizationalUnit(d.Id())
	if apiErr != nil {
		if strings.Contains(apiErr.Error(), "The resource is not found.") {
			d.SetId("")
			return nil
		}
		return diag.Errorf(
			"Error creating Clumio AWS Connection. Error: %v", string(apiErr.Response))

	}
	err := d.Set("children_count", int(*res.ChildrenCount))
	if err != nil {
		return diag.Errorf(
			"Error setting children_count schema attribute. Error: %v", err)
	}
	err = d.Set("user_count", int(*res.UserCount))
	if err != nil {
		return diag.Errorf("Error setting user_count schema attribute. Error: %v", err)
	}
	if res.ConfiguredDatasourceTypes != nil {
		configuredDatasourceTypes := make([]string, 0)
		for _, dsType := range res.ConfiguredDatasourceTypes {
			configuredDatasourceTypes = append(configuredDatasourceTypes, *dsType)
		}
		err = d.Set("configured_datasource_types", configuredDatasourceTypes)
		if err != nil {
			return diag.Errorf(
				"Error setting configured_datasource_types schema attribute."+
					" Error: %v", err)
		}
	}
	if res.DescendantIds != nil {
		descendantIds := make([]string, 0)
		for _, dsType := range res.ConfiguredDatasourceTypes {
			descendantIds = append(descendantIds, *dsType)
		}
		err = d.Set("descendant_ids", descendantIds)
		if err != nil {
			return diag.Errorf(
				"Error setting descendant_ids schema attribute. Error: %v", err)
		}
	}
	return nil
}

// clumioOrganizationalUnitUpdate handles the Update action for the Clumio OrganizationalUnit Resource.
func clumioOrganizationalUnitUpdate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("email") {
		return diag.Errorf("email is not allowed to be changed")
	}
	client := meta.(*apiClient)
	orgUnitsAPI := orgUnits.NewOrganizationalUnitsV1(client.clumioConfig)
	updateRequest := &models.PatchOrganizationalUnitV1Request{}
	if d.HasChange("description") {
		description := getStringValue(d, "description")
		updateRequest.Description = &description
	}
	if d.HasChange("name") {
		name := getStringValue(d, "name")
		updateRequest.Name = &name
	}
	if d.HasChange("users") {
		oldValue, newValue := d.GetChange("users")
		deleted := sliceDifference(oldValue.([]interface{}), newValue.([]interface{}))
		added := sliceDifference(newValue.([]interface{}), oldValue.([]interface{}))
		deletedStrings := getStringSliceFromInterfaceSlice(deleted)
		addedStrings := getStringSliceFromInterfaceSlice(added)
		updateRequest.Users =
			&models.UpdateUserAssignments{
				Add:    addedStrings,
				Remove: deletedStrings,
			}
	}
	if d.HasChange("entities") {
		oldValue, newValue := d.GetChange("organizational_unit_ids")
		_ = diffEntities(oldValue.([]interface{}), newValue.([]interface{}))
		_ = diffEntities(newValue.([]interface{}), oldValue.([]interface{}))
		// TODO: Handle entities once the SDK issue with the entityModel is fixed.
	}
	_, apiErr := orgUnitsAPI.PatchOrganizationalUnit(d.Id(), nil, updateRequest)
	if apiErr != nil {
		return diag.Errorf("Error updating Clumio Organizational Unit. Error: %v",
			string(apiErr.Response))
	}
	return clumioOrganizationalUnitRead(ctx, d, meta)
}

// clumioOrganizationalUnitDelete handles the Delete action for the Clumio OrganizationalUnit Resource.
func clumioOrganizationalUnitDelete(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	orgUnitsAPI := orgUnits.NewOrganizationalUnitsV1(client.clumioConfig)
	_, apiErr := orgUnitsAPI.DeleteOrganizationalUnit(d.Id(), nil)
	if apiErr != nil {
		return diag.Errorf("Error deleting Clumio Organizational Unit %v. Error: %v",
			d.Id(), string(apiErr.Response))
	}
	return nil
}

// utility function to generate the diff of the entities schema attribute if there are
// changes.
func diffEntities(slice1 []interface{}, slice2 []interface{}) []interface{} {
	returnSlice := make([]interface{}, 0)
	slice1Map := getEntityMap(slice1)
	slice2Map := getEntityMap(slice2)
	slice1Keys := getMapKeys(slice1Map)
	slice2Keys := getMapKeys(slice2Map)
	diffKeys := sliceDifference(slice1Keys, slice2Keys)
	for _, key := range diffKeys {
		returnSlice = append(returnSlice, slice1Map[key.(string)])
	}
	return returnSlice
}

// Function to generate a entity map from entity slice with the key of the map being
// the combination of the id and type.
func getEntityMap(entitySlice []interface{}) map[string]interface{} {
	entityMap := make(map[string]interface{})
	for _, val1 := range entitySlice {
		valMap := val1.(map[string]interface{})
		mapKey := fmt.Sprintf("%s_%s", valMap["id"].(string), valMap["type"].(string))
		entityMap[mapKey] = val1
	}
	return entityMap
}

// Get the keys for the map as interface slice.
func getMapKeys(entityMap map[string]interface{}) []interface{} {
	interfaceSlice := make([]interface{}, 0)
	for key := range entityMap {
		interfaceSlice = append(interfaceSlice, key)
	}
	return interfaceSlice
}
