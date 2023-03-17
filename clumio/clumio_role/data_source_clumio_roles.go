// Copyright 2021. Clumio, Inc.

// clumio_role Data Source implementation.

package clumio_role

import (
	"context"

	"github.com/clumio-code/clumio-go-sdk/controllers/roles"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceClumioRole returns the resource for Clumio Roles.
func DataSourceClumioRole() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio Roles Data Source used to list the Clumio Roles.",
		ReadContext: dataSourceClumioRoleRead,
		Schema: map[string]*schema.Schema{

			schemaId: {
				Type:        schema.TypeString,
				Description: "The Clumio-assigned ID of the role.",
				Computed:    true,
			},
			schemaName: {
				Type:        schema.TypeString,
				Description: "Unique name assigned to the role.",
				Required:    true,
			},
			schemaDescription: {
				Type:        schema.TypeString,
				Description: "A description of the role.",
				Computed:    true,
			},
			schemaUserCount: {
				Type:        schema.TypeInt,
				Description: "Number of users to whom the role has been assigned.",
				Computed:    true,
			},
			schemaPermissions: {
				Type:        schema.TypeList,
				Description: "Permissions contained in the role.",
				Computed:    true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						schemaDescription: {
							Type:        schema.TypeString,
							Description: "Description of the permission.",
							Computed:    true,
						},
						schemaId: {
							Type:        schema.TypeString,
							Description: "The Clumio-assigned ID of the permission.",
							Computed:    true,
						},
						schemaName: {
							Type:        schema.TypeString,
							Description: "Name of the permission.",
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

// dataSourceClumioRoleRead handles the Read action for the Clumio Roles Data Source.
func dataSourceClumioRoleRead(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	rolesApi := roles.NewRolesV1(client.ClumioConfig)
	rolesRes, apiErr := rolesApi.ListRoles()
	if apiErr != nil {
		return diag.Errorf(
			"Error listing Clumio Roles. Error: %v", string(apiErr.Response))
	}
	name, ok := d.GetOk(schemaName)
	if !ok {
		return diag.Errorf("name is a required schema attribute.")
	}
	var expectedRole *models.RoleWithETag
	for _, roleItem := range rolesRes.Embedded.Items {
		if *roleItem.Name == name.(string) {
			expectedRole = roleItem
			break
		}
	}
	if expectedRole == nil {
		return diag.Errorf("Role with name %s is not found.", name)
	}
	d.SetId(*expectedRole.Id)
	err := d.Set(schemaDescription, expectedRole.Description)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaDescription, err)
	}
	err = d.Set(schemaUserCount, expectedRole.UserCount)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaUserCount, err)
	}

	err = d.Set(schemaPermissions, getSchemaPermissions(expectedRole))
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaUserCount, err)
	}
	return nil
}

// getSchemaPermissions returns the value to set for the permissions schema attribute
func getSchemaPermissions(expectedRole *models.RoleWithETag) []interface{} {
	permissionsList := make([]interface{}, 0)
	for _, permission := range expectedRole.Permissions {
		permissionMap := make(map[string]interface{})
		permissionMap[schemaId] = *permission.Id
		permissionMap[schemaName] = *permission.Name
		permissionMap[schemaDescription] = *permission.Description
		permissionsList = append(permissionsList, permissionMap)
	}
	return permissionsList
}
