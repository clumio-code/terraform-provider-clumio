// Copyright 2021. Clumio, Inc.

// clumio_user definition and CRUD implementation.
package clumio

import (
	"context"
	"strconv"
	"strings"

	"github.com/clumio-code/clumio-go-sdk/controllers/users"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// clumioUser returns the resource for Clumio AWS Connection.

func clumioUser() *schema.Resource {
	return &schema.Resource{
		// This description is used by the documentation generator and the language server.
		Description: "Clumio AWS Connection Resource used to connect AWS accounts to Clumio.",

		CreateContext: clumioUserCreate,
		ReadContext:   clumioUserRead,
		UpdateContext: clumioUserUpdate,
		DeleteContext: clumioUserDelete,

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"id": {
				Description: "User Id.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"email": {
				Description: "The email address of the user to be added to Clumio.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"full_name": {
				Description: "The full name of the user to be added to Clumio." +
					" For example, enter the user's first name and last name. The name" +
					" appears in the User Management screen and in the body of the" +
					" email invitation.",
				Type:     schema.TypeString,
				Required: true,
			},
			"assigned_role": {
				Description: "The Clumio-assigned ID of the role to assign to the user.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"organizational_unit_ids": {
				Description: "The Clumio-assigned IDs of the organizational units" +
					" to be assigned to the user.",
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"inviter": {
				Description: "The ID number of the user who sent the email invitation.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"is_confirmed": {
				Description: "Determines whether the user has activated their Clumio" +
					" account. If true, the user has activated the account.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"is_enabled": {
				Description: "Determines whether the user is enabled (in Activated or" +
					" Invited status) in Clumio. If true, the user is in Activated or" +
					" Invited status in Clumio. Users in Activated status can log in to" +
					" Clumio. Users in Invited status have been invited to log in to" +
					" Clumio via an email invitation and the invitation is pending" +
					" acceptance from the user. If false, the user has been manually" +
					" suspended and cannot log in to Clumio until another Clumio user" +
					" reactivates the account.",
				Type:        schema.TypeBool,
				Computed:    true,
			},
			"last_activity_timestamp": {
				Description: "The timestamp of when when the user was last active in" +
					" the Clumio system. Represented in RFC-3339 format.",
				Type:        schema.TypeString,
				Computed:    true,
			},
			"organizational_unit_count": {
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
	client := meta.(*apiClient)
	usersAPI := users.NewUsersV1(client.clumioConfig)
	assignedRole := getStringValue(d, "assigned_role")
	fullname := getStringValue(d, "full_name")
	email := getStringValue(d, "email")
	organizationalUnitIds := getStringSlice(d, "organizational_unit_ids")
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
	client := meta.(*apiClient)
	usersAPI := users.NewUsersV1(client.clumioConfig)
	userId, perr := strconv.ParseInt(d.Id(), 10, 64)
	if perr != nil{
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
			"Error creating Clumio AWS Connection. Error: %v", string(apiErr.Response))

	}
	err := d.Set("inviter", *res.Inviter)
	if err != nil {
		return diag.Errorf(
			"Error setting inviter schema attribute. Error: %v", err)
	}
	err = d.Set("is_confirmed", res.IsConfirmed)
	if err != nil {
		return diag.Errorf(
			"Error setting is_confirmed schema attribute. Error: %v", err)
	}
	err = d.Set("is_enabled", res.IsEnabled)
	if err != nil {
		return diag.Errorf("Error setting is_enabled schema attribute."+
			" Error: %v", err)
	}
	err = d.Set("last_activity_timestamp", res.LastActivityTimestamp)
	if err != nil {
		return diag.Errorf("Error setting last_activity_timestamp schema" +
			" attribute. Error: %v", err)
	}
	err = d.Set("organizational_unit_count", int(*res.OrganizationalUnitCount))
	if err != nil {
		return diag.Errorf("Error setting organizational_unit_count schema" +
			" attribute. Error: %v", err)
	}
	return nil
}

// clumioUserUpdate handles the Update action for the Clumio User Resource.
func clumioUserUpdate(
	ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	if d.HasChange("email"){
		return diag.Errorf("email is not allowed to be changed")
	}
	client := meta.(*apiClient)
	usersAPI := users.NewUsersV1(client.clumioConfig)
	updateRequest := &models.UpdateUserV1Request{
	}
	if d.HasChange("assigned_role"){
		assignedRole := getStringValue(d, "assigned_role")
		updateRequest.AssignedRole = &assignedRole
	}
	if d.HasChange("full_name"){
		fullname := getStringValue(d, "full_name")
		updateRequest.FullName = &fullname
	}
	if d.HasChange("organizational_unit_ids"){
		oldValue, newValue := d.GetChange("organizational_unit_ids")
		deleted :=sliceDifference(oldValue.([]interface{}), newValue.([]interface{}))
		added :=sliceDifference(newValue.([]interface{}), oldValue.([]interface{}))
		deletedStrings := getStringSliceFromInterfaceSlice(deleted)
		addedStrings := getStringSliceFromInterfaceSlice(added)
		updateRequest.OrganizationalUnitAssignmentUpdates =
			&models.EntityGroupAssignmetUpdates{
				Add:    addedStrings,
				Remove: deletedStrings,
			}
	}
	userId, perr := strconv.ParseInt(d.Id(), 10, 64)
	if perr != nil{
		return diag.Errorf(
			"Invalid user id : %v", d.Id())
	}
	_, apiErr := usersAPI.UpdateUser(userId, updateRequest)
	if apiErr != nil {
		return diag.Errorf(
			"Error creating Clumio User. Error: %v", string(apiErr.Response))
	}
	return clumioUserRead(ctx, d, meta)
}

// clumioUserDelete handles the Delete action for the Clumio User Resource.
func clumioUserDelete(
	_ context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client := meta.(*apiClient)
	usersAPI := users.NewUsersV1(client.clumioConfig)
	userId, perr := strconv.ParseInt(d.Id(), 10, 64)
	if perr != nil{
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
