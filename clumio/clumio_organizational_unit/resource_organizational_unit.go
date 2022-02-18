// Copyright 2021. Clumio, Inc.

// clumio_organizational_unit definition and CRUD implementation.

package clumio_organizational_unit

import (
	"context"
	"strings"

	orgUnits "github.com/clumio-code/clumio-go-sdk/controllers/organizational_units"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ClumioOrganizationalUnit returns the resource for Clumio Organizational Unit.
func ClumioOrganizationalUnit() *schema.Resource {
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
			schemaId: {
				Description: "OrganizationalUnit Id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			schemaName: {
				Description: "Unique name assigned to the organizational unit.",
				Type:        schema.TypeString,
				Required:    true,
			},
			schemaDescription: {
				Description: "A description of the organizational unit.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			schemaParentId: {
				Description: "The Clumio-assigned ID of the parent organizational unit" +
					" under which the new organizational unit is to be created.",
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},
			schemaChildrenCount: {
				Description: "Number of immediate children of the organizational unit.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			schemaConfiguredDatasourceTypes: {
				Description: "Datasource types configured in this organizational unit." +
					" Possible values include aws, microsoft365, vmware, or mssql.",
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			schemaDescendantIds: {
				Description: "List of all recursive descendant organizational units" +
					" of this OU.",
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
			schemaUserCount: {
				Description: "Number of users to whom this organizational unit or any" +
					" of its descendants have been assigned.",
				Type:     schema.TypeInt,
				Computed: true,
			},
			schemaUsers: {
				Description: "List of user ids to assign this organizational unit.",
				Type:        schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

// clumioOrganizationalUnitCreate handles the Create action for the Clumio Organizational Unit Resource.
func clumioOrganizationalUnitCreate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	orgUnitsAPI := orgUnits.NewOrganizationalUnitsV1(client.ClumioConfig)
	name := common.GetStringValue(d, schemaName)
	parentId := common.GetStringValue(d, schemaParentId)
	request := &models.CreateOrganizationalUnitV1Request{
		Name:     &name,
		ParentId: &parentId,
	}
	description := common.GetStringValue(d, schemaDescription)
	if description != "" {
		request.Description = &description
	}
	res, apiErr := orgUnitsAPI.CreateOrganizationalUnit(nil, request)
	if apiErr != nil {
		return diag.Errorf(
			"Error creating Clumio OrganizationalUnit. Error: %v", string(apiErr.Response))
	}
	d.SetId(*res.Id)
	return clumioOrganizationalUnitRead(ctx, d, meta)
}

// clumioOrganizationalUnitRead handles the Read action for the Clumio Organizational Unit Resource.
func clumioOrganizationalUnitRead(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	orgUnitsAPI := orgUnits.NewOrganizationalUnitsV1(client.ClumioConfig)
	res, apiErr := orgUnitsAPI.ReadOrganizationalUnit(d.Id())
	if apiErr != nil {
		if strings.Contains(apiErr.Error(), "The resource is not found.") {
			d.SetId("")
			return nil
		}
		return diag.Errorf(
			"Error reading Clumio Organizational Unit. Error: %v", string(apiErr.Response))

	}
	err := d.Set(schemaChildrenCount, int(*res.ChildrenCount))
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaChildrenCount, err)
	}
	err = d.Set(schemaUserCount, int(*res.UserCount))
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaUserCount, err)
	}
	var configuredDatasourceTypes []string
	if res.ConfiguredDatasourceTypes != nil {
		configuredDatasourceTypes = make([]string, 0)
		for _, dsType := range res.ConfiguredDatasourceTypes {
			configuredDatasourceTypes = append(configuredDatasourceTypes, *dsType)
		}
	}
	err = d.Set(schemaConfiguredDatasourceTypes, configuredDatasourceTypes)
	if err != nil {
		return diag.Errorf(
			common.SchemaAttributeSetError, schemaConfiguredDatasourceTypes, err)
	}
	var descendantIds []string
	if res.DescendantIds != nil {
		descendantIds = make([]string, 0)
		for _, dsType := range res.ConfiguredDatasourceTypes {
			descendantIds = append(descendantIds, *dsType)
		}
	}
	err = d.Set(schemaDescendantIds, descendantIds)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaDescendantIds, err)
	}
	var users []string
	if res.Users != nil {
		users = make([]string, 0)
		for _, user := range res.Users {
			users = append(users, *user)
		}
	}
	err = d.Set(schemaUsers, users)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaUsers, err)
	}
	err = d.Set(schemaName, *res.Name)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaName, err)
	}
	if res.Description != nil {
		err = d.Set(schemaDescription, *res.Description)
		if err != nil {
			return diag.Errorf(common.SchemaAttributeSetError, schemaDescription, err)
		}
	}
	if res.ParentId != nil {
		err = d.Set(schemaParentId, *res.ParentId)
		if err != nil {
			return diag.Errorf(common.SchemaAttributeSetError, schemaParentId, err)
		}
	}

	return nil
}

// clumioOrganizationalUnitUpdate handles the Update action for the Clumio Organizational Unit Resource.
func clumioOrganizationalUnitUpdate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	orgUnitsAPI := orgUnits.NewOrganizationalUnitsV1(client.ClumioConfig)
	updateRequest := &models.PatchOrganizationalUnitV1Request{}
	description := common.GetStringValue(d, schemaDescription)
	updateRequest.Description = &description
	name := common.GetStringValue(d, schemaName)
	updateRequest.Name = &name

	_, apiErr := orgUnitsAPI.PatchOrganizationalUnit(d.Id(), nil, updateRequest)
	if apiErr != nil {
		return diag.Errorf("Error updating Clumio Organizational Unit. Error: %v",
			string(apiErr.Response))
	}
	return clumioOrganizationalUnitRead(ctx, d, meta)
}

// clumioOrganizationalUnitDelete handles the Delete action for the Clumio Organizational Unit Resource.
func clumioOrganizationalUnitDelete(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	orgUnitsAPI := orgUnits.NewOrganizationalUnitsV1(client.ClumioConfig)
	_, apiErr := orgUnitsAPI.DeleteOrganizationalUnit(d.Id(), nil)
	if apiErr != nil {
		return diag.Errorf("Error deleting Clumio Organizational Unit %v. Error: %v",
			d.Id(), string(apiErr.Response))
	}
	return nil
}
