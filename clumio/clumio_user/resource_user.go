// Copyright 2021. Clumio, Inc.

// clumio_user definition and CRUD implementation.
package clumio_user

import (
	"context"
	"strconv"
	"strings"

	"github.com/clumio-code/clumio-go-sdk/controllers/users"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/clumio-code/terraform-provider-clumio/clumio/common"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// ClumioUser returns the resource for Clumio AWS Connection.

func ClumioUser() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio User Resource to create and manage users in Clumio.",

		CreateContext: clumioUserCreate,
		ReadContext:   clumioUserRead,
		UpdateContext: clumioUserUpdate,
		DeleteContext: clumioUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			schemaId: {
				Description: "User Id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			schemaEmail: {
				Description: "The email address of the user to be added to Clumio.",
				Type:        schema.TypeString,
				Required:    true,
			},
			schemaFullName: {
				Description: "The full name of the user to be added to Clumio." +
					" For example, enter the user's first name and last name. The name" +
					" appears in the User Management screen and in the body of the" +
					" email invitation.",
				Type:     schema.TypeString,
				Required: true,
			},
			schemaAssignedRole: {
				Description: "The Clumio-assigned ID of the role to assign to the user.",
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
			},
			schemaOrganizationalUnitIds: {
				Description: "The Clumio-assigned IDs of the organizational units" +
					" to be assigned to the user. The Global Organizational Unit ID is " +
					"\"00000000-0000-0000-0000-000000000000\"",
				Type: schema.TypeSet,
				Set:  common.SchemaSetHashString,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			schemaInviter: {
				Description: "The ID number of the user who sent the email invitation.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			schemaIsConfirmed: {
				Description: "Determines whether the user has activated their Clumio" +
					" account. If true, the user has activated the account.",
				Type:     schema.TypeBool,
				Computed: true,
			},
			schemaIsEnabled: {
				Description: "Determines whether the user is enabled (in Activated or" +
					" Invited status) in Clumio. If true, the user is in Activated or" +
					" Invited status in Clumio. Users in Activated status can log in to" +
					" Clumio. Users in Invited status have been invited to log in to" +
					" Clumio via an email invitation and the invitation is pending" +
					" acceptance from the user. If false, the user has been manually" +
					" suspended and cannot log in to Clumio until another Clumio user" +
					" reactivates the account.",
				Type:     schema.TypeBool,
				Computed: true,
			},
			schemaLastActivityTimestamp: {
				Description: "The timestamp of when when the user was last active in" +
					" the Clumio system. Represented in RFC-3339 format.",
				Type:     schema.TypeString,
				Computed: true,
			},
			schemaOrganizationalUnitCount: {
				Description: "The number of organizational units accessible to the user.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

// clumioUserCreate handles the Create action for the Clumio User Resource.
func clumioUserCreate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	usersAPI := users.NewUsersV1(client.ClumioConfig)
	assignedRole := common.GetStringValue(d, schemaAssignedRole)
	fullname := common.GetStringValue(d, schemaFullName)
	email := common.GetStringValue(d, schemaEmail)
	organizationalUnitIds := common.GetStringSliceFromSet(d, schemaOrganizationalUnitIds)
	if len(organizationalUnitIds) == 0 {
		return diag.Errorf(common.SchemaEmptyParameterError, schemaOrganizationalUnitIds)
	}
	res, apiErr := usersAPI.CreateUser(&models.CreateUserV1Request{
		AssignedRole:          &assignedRole,
		Email:                 &email,
		FullName:              &fullname,
		OrganizationalUnitIds: organizationalUnitIds,
	})
	if apiErr != nil {
		return diag.Errorf(
			"Error creating Clumio User. Error: %v", string(apiErr.Response))
	}
	d.SetId(*res.Id)
	return clumioUserRead(ctx, d, meta)
}

// clumioUserRead handles the Read action for the Clumio User Resource.
func clumioUserRead(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	usersAPI := users.NewUsersV1(client.ClumioConfig)
	userId, perr := strconv.ParseInt(d.Id(), 10, 64)
	if perr != nil {
		return diag.Errorf(
			"Invalid user id : %v", d.Id())
	}
	res, apiErr := usersAPI.ReadUser(userId)
	if apiErr != nil {
		if strings.Contains(apiErr.Error(), "The resource is not found.") {
			d.SetId("")
			return nil
		}
		return diag.Errorf(
			"Error retrieving Clumio User. Error: %v", string(apiErr.Response))

	}
	err := d.Set(schemaInviter, *res.Inviter)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaInviter, err)
	}
	err = d.Set(schemaIsConfirmed, *res.IsConfirmed)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaIsConfirmed, err)
	}
	err = d.Set(schemaIsEnabled, *res.IsEnabled)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaIsEnabled, err)
	}
	err = d.Set(schemaLastActivityTimestamp, *res.LastActivityTimestamp)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaLastActivityTimestamp, err)
	}
	err = d.Set(schemaOrganizationalUnitCount, int(*res.OrganizationalUnitCount))
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaOrganizationalUnitCount, err)
	}
	err = d.Set(schemaEmail, *res.Email)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaEmail, err)
	}
	err = d.Set(schemaFullName, *res.FullName)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaFullName, err)
	}
	if res.AssignedRole != nil {
		err = d.Set(schemaAssignedRole, *res.AssignedRole)
		if err != nil {
			return diag.Errorf(common.SchemaAttributeSetError, schemaAssignedRole, err)
		}
	}
	var orgUnitIds *schema.Set
	if res.AssignedOrganizationalUnitIds != nil {
		orgUnitIds = &schema.Set{F: common.SchemaSetHashString}
		for _, orgUnitId := range res.AssignedOrganizationalUnitIds {
			orgUnitIds.Add(*orgUnitId)
		}
	}
	err = d.Set(schemaOrganizationalUnitIds, orgUnitIds)
	if err != nil {
		return diag.Errorf(common.SchemaAttributeSetError, schemaOrganizationalUnitIds, err)
	}
	return nil
}

// clumioUserUpdate handles the Update action for the Clumio User Resource.
func clumioUserUpdate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange(schemaEmail) {
		return diag.Errorf("email is not allowed to be changed")
	}
	client := meta.(*common.ApiClient)
	usersAPI := users.NewUsersV1(client.ClumioConfig)
	updateRequest := &models.UpdateUserV1Request{}
	if d.HasChange(schemaAssignedRole) {
		assignedRole := common.GetStringValue(d, schemaAssignedRole)
		updateRequest.AssignedRole = &assignedRole
	}
	if d.HasChange(schemaFullName) {
		fullname := common.GetStringValue(d, schemaFullName)
		updateRequest.FullName = &fullname
	}
	if d.HasChange(schemaOrganizationalUnitIds) {
		oldValue, newValue := d.GetChange(schemaOrganizationalUnitIds)
		deleted := common.SliceDifference(oldValue.(*schema.Set).List(), newValue.(*schema.Set).List())
		added := common.SliceDifference(newValue.(*schema.Set).List(), oldValue.(*schema.Set).List())
		deletedStrings := common.GetStringSliceFromInterfaceSlice(deleted)
		addedStrings := common.GetStringSliceFromInterfaceSlice(added)
		updateRequest.OrganizationalUnitAssignmentUpdates =
			&models.EntityGroupAssignmentUpdatesV1{
				Add:    addedStrings,
				Remove: deletedStrings,
			}
	}
	userId, perr := strconv.ParseInt(d.Id(), 10, 64)
	if perr != nil {
		return diag.Errorf(
			"Invalid user id : %v", d.Id())
	}
	_, apiErr := usersAPI.UpdateUser(userId, updateRequest)
	if apiErr != nil {
		return diag.Errorf(
			"Error updating Clumio User. Error: %v", string(apiErr.Response))
	}
	return clumioUserRead(ctx, d, meta)
}

// clumioUserDelete handles the Delete action for the Clumio User Resource.
func clumioUserDelete(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*common.ApiClient)
	usersAPI := users.NewUsersV1(client.ClumioConfig)
	userId, perr := strconv.ParseInt(d.Id(), 10, 64)
	if perr != nil {
		return diag.Errorf(
			"Invalid user id : %v", d.Id())
	}
	_, apiErr := usersAPI.DeleteUser(userId)
	if apiErr != nil {
		return diag.Errorf(
			"Error deleting Clumio User %v. Error: %v",
			d.Id(), string(apiErr.Response))
	}
	return nil
}
