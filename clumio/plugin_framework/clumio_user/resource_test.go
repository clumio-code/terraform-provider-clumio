// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_user

import (
	"context"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	resourceName          = "test_user"
	id                    = "1"
	idInt                 = int64(1)
	ou                    = "test-ou"
	ouUpdated             = "test-ou-updated"
	testError             = "Test Error"
	email                 = "test-email"
	fullName              = "test-full-name"
	fullNameUpdated       = "test-updated-full-name"
	roleId                = "test-role-id"
	roleId2               = "test-role-id-2"
	inviter               = "test-inviter"
	isConfirmed           = false
	isEnabled             = false
	lastActivityTimestamp = "test-timestamp"
	ouCount               = int64(1)
	invalidId             = "invalid"
	invalidIdInt          = int64(0)
)

// Unit test for the following cases:
//   - Create user success scenario.
//   - SDK API for create user returns an error.
//   - SDK API for create user returns an empty response.
func TestCreateUser(t *testing.T) {

	ctx := context.Background()
	mockUser := sdkclients.NewMockUserClient(t)
	ur := clumioUserResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkUsers: mockUser,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	ouIdsList := []*string{&ou}
	ouIds, conversionDiags := types.SetValueFrom(
		ctx, types.StringType, ouIdsList)
	assert.Nil(t, conversionDiags)

	accessControlModel := []*roleForOrganizationalUnitModel{
		{
			RoleId:                basetypes.NewStringValue(roleId),
			OrganizationalUnitIds: ouIds,
		},
	}
	accessControlList, conversionDiags := types.SetValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaRoleId: types.StringType,
			schemaOrganizationalUnitIds: types.SetType{
				ElemType: types.StringType,
			},
		},
	}, accessControlModel)

	urm := &clumioUserResourceModel{
		Email:                      basetypes.NewStringValue(email),
		FullName:                   basetypes.NewStringValue(fullName),
		AccessControlConfiguration: accessControlList,
	}

	// Tests the success scenario for user create. It should not return Diagnostics.
	t.Run("Basic success scenario for create user", func(t *testing.T) {

		createUserResponse := &models.CreateUserResponse{
			AccessControlConfiguration: []*models.RoleForOrganizationalUnits{
				{
					RoleId:                &roleId,
					OrganizationalUnitIds: ouIdsList,
				},
			},
			Id:                      &id,
			Inviter:                 &inviter,
			IsConfirmed:             &isConfirmed,
			IsEnabled:               &isEnabled,
			LastActivityTimestamp:   &lastActivityTimestamp,
			OrganizationalUnitCount: &ouCount,
		}

		//Setup expectations.
		mockUser.EXPECT().CreateUser(mock.Anything).Times(1).Return(createUserResponse, nil)

		diags := ur.createUser(ctx, urm)
		assert.Nil(t, diags)
		assert.Equal(t, id, urm.Id.ValueString())
		assert.Equal(t, inviter, urm.Inviter.ValueString())
		assert.Equal(t, lastActivityTimestamp, urm.LastActivityTimestamp.ValueString())
		assert.Equal(t, isEnabled, urm.IsEnabled.ValueBool())
		assert.Equal(t, isConfirmed, urm.IsConfirmed.ValueBool())
		assert.Equal(t, ouCount, urm.OrganizationalUnitCount.ValueInt64())
	})

	// Tests that Diagnostics is returned in case the create user API call returns an error.
	t.Run("create user returns an error", func(t *testing.T) {

		// Setup expectations.
		mockUser.EXPECT().CreateUser(mock.Anything).Times(1).Return(nil, apiError)

		diags := ur.createUser(ctx, urm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the create user API call returns an empty
	// response.
	t.Run("create user returns an empty response", func(t *testing.T) {

		// Setup expectations.
		mockUser.EXPECT().CreateUser(mock.Anything).Times(1).Return(nil, nil)

		diags := ur.createUser(ctx, urm)
		assert.NotNil(t, diags)
	})

}

// Unit test for the following cases:
//   - Read user success scenario.
//   - SDK API for read user returns not found error.
//   - SDK API for read user returns an error.
//   - SDK API for read user returns an empty response.
//   - Using invalid user id retunrs an error.
func TestReadUser(t *testing.T) {

	ctx := context.Background()
	mockUser := sdkclients.NewMockUserClient(t)
	ur := clumioUserResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkUsers: mockUser,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	ouIdsList := []*string{&ou}
	ouIds, conversionDiags := types.SetValueFrom(
		ctx, types.StringType, ouIdsList)
	assert.Nil(t, conversionDiags)

	accessControlModel := []*roleForOrganizationalUnitModel{
		{
			RoleId:                basetypes.NewStringValue(roleId),
			OrganizationalUnitIds: ouIds,
		},
	}
	accessControlList, conversionDiags := types.SetValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaRoleId: types.StringType,
			schemaOrganizationalUnitIds: types.SetType{
				ElemType: types.StringType,
			},
		},
	}, accessControlModel)
	assert.Nil(t, conversionDiags)

	urm := &clumioUserResourceModel{
		Id:                         basetypes.NewStringValue(id),
		Email:                      basetypes.NewStringValue(email),
		FullName:                   basetypes.NewStringValue(fullName),
		AccessControlConfiguration: accessControlList,
	}

	// Tests the success scenario for user read. It should not return Diagnostics.
	t.Run("Basic success scenario for read user", func(t *testing.T) {

		readUserResponse := &models.ReadUserResponse{
			AccessControlConfiguration: []*models.RoleForOrganizationalUnits{
				{
					RoleId:                &roleId,
					OrganizationalUnitIds: ouIdsList,
				},
			},
			Email:                   &email,
			FullName:                &fullName,
			Id:                      &id,
			Inviter:                 &inviter,
			IsConfirmed:             &isConfirmed,
			IsEnabled:               &isEnabled,
			LastActivityTimestamp:   &lastActivityTimestamp,
			OrganizationalUnitCount: &ouCount,
		}

		//Setup expectations.
		mockUser.EXPECT().ReadUser(mock.Anything).Times(1).Return(readUserResponse, nil)

		remove, diags := ur.readUser(ctx, urm)
		assert.Nil(t, diags)
		assert.False(t, remove)
		assert.Equal(t, id, urm.Id.ValueString())
		assert.Equal(t, email, urm.Email.ValueString())
		assert.Equal(t, fullName, urm.FullName.ValueString())
		assert.Equal(t, inviter, urm.Inviter.ValueString())
		assert.Equal(t, lastActivityTimestamp, urm.LastActivityTimestamp.ValueString())
		assert.Equal(t, isEnabled, urm.IsEnabled.ValueBool())
		assert.Equal(t, isConfirmed, urm.IsConfirmed.ValueBool())
		assert.Equal(t, ouCount, urm.OrganizationalUnitCount.ValueInt64())
	})

	// Tests that Diagnostics is returned in case the read user API call returns HTTP 404 error.
	t.Run("read user returns http 404 error", func(t *testing.T) {

		// Setup Expectations
		mockUser.EXPECT().ReadUser(mock.Anything).Times(1).Return(nil, apiNotFoundError)

		remove, diags := ur.readUser(ctx, urm)
		assert.Nil(t, diags)
		assert.True(t, remove)
	})

	// Tests that Diagnostics is returned in case the read user API call returns an error.
	t.Run("read user returns an error", func(t *testing.T) {

		// Setup expectations.
		mockUser.EXPECT().ReadUser(mock.Anything).Times(1).Return(nil, apiError)

		remove, diags := ur.readUser(ctx, urm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that Diagnostics is returned in case the read user API call returns an empty
	// response.
	t.Run("read user returns an empty response", func(t *testing.T) {

		// Setup expectations.
		mockUser.EXPECT().ReadUser(mock.Anything).Times(1).Return(nil, nil)

		remove, diags := ur.readUser(ctx, urm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that Diagnostics is returned in case using invalid user id.
	t.Run("using invalid user id returns an error.", func(t *testing.T) {

		invalidUrm := &clumioUserResourceModel{
			Id: basetypes.NewStringValue(invalidId),
		}

		// Setup expectations.
		mockUser.EXPECT().ReadUser(invalidIdInt).Times(1).Return(nil, apiError)

		remove, diags := ur.readUser(ctx, invalidUrm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})
}

// Unit test for the following cases:
//   - Update user success scenario.
//   - SDK API for update user returns an error.
//   - SDK API for update user returns an empty response.
//   - Using invalid user id returns an error.
func TestUpdateUser(t *testing.T) {

	ctx := context.Background()
	mockUser := sdkclients.NewMockUserClient(t)
	ur := clumioUserResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkUsers: mockUser,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	ouIdsList := []*string{&ou}
	ouIds, conversionDiags := types.SetValueFrom(
		ctx, types.StringType, ouIdsList)
	assert.Nil(t, conversionDiags)

	accessControlModel := []*roleForOrganizationalUnitModel{
		{
			RoleId:                basetypes.NewStringValue(roleId),
			OrganizationalUnitIds: ouIds,
		},
	}
	accessControlList, conversionDiags := types.SetValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaRoleId: types.StringType,
			schemaOrganizationalUnitIds: types.SetType{
				ElemType: types.StringType,
			},
		},
	}, accessControlModel)
	assert.Nil(t, conversionDiags)

	plan := &clumioUserResourceModel{
		Id:                         basetypes.NewStringValue(id),
		Email:                      basetypes.NewStringValue(email),
		FullName:                   basetypes.NewStringValue(fullNameUpdated),
		AccessControlConfiguration: accessControlList,
	}

	ouIdsList = []*string{&ouUpdated}
	ouIds, conversionDiags = types.SetValueFrom(
		ctx, types.StringType, ouIdsList)
	assert.Nil(t, conversionDiags)
	accessControlModel = []*roleForOrganizationalUnitModel{
		{
			RoleId:                basetypes.NewStringValue(roleId),
			OrganizationalUnitIds: ouIds,
		},
	}
	accessControlList2, conversionDiags := types.SetValueFrom(ctx, types.ObjectType{
		AttrTypes: map[string]attr.Type{
			schemaRoleId: types.StringType,
			schemaOrganizationalUnitIds: types.SetType{
				ElemType: types.StringType,
			},
		},
	}, accessControlModel)
	assert.Nil(t, conversionDiags)

	state := &clumioUserResourceModel{
		Id:                         basetypes.NewStringValue(id),
		Email:                      basetypes.NewStringValue(email),
		FullName:                   basetypes.NewStringValue(fullName),
		AccessControlConfiguration: accessControlList2,
	}

	// Tests the success scenario for user update. It should not return Diagnostics.
	t.Run("Basic success scenario for update user", func(t *testing.T) {

		updateUserResponse := &models.UpdateUserResponse{
			AccessControlConfiguration: []*models.RoleForOrganizationalUnits{
				{
					RoleId:                &roleId,
					OrganizationalUnitIds: ouIdsList,
				},
			},
			Id:                      &id,
			Inviter:                 &inviter,
			IsConfirmed:             &isConfirmed,
			IsEnabled:               &isEnabled,
			LastActivityTimestamp:   &lastActivityTimestamp,
			OrganizationalUnitCount: &ouCount,
		}

		//Setup expectations.
		mockUser.EXPECT().UpdateUser(idInt, mock.Anything).Times(1).Return(updateUserResponse, nil)

		diags := ur.updateUser(ctx, plan, state)
		assert.Nil(t, diags)
		assert.Equal(t, fullNameUpdated, plan.FullName.ValueString())
	})

	// Tests that Diagnostics is returned in case the update user API call returns an error.
	t.Run("update user returns an error", func(t *testing.T) {

		// Setup expectations.
		mockUser.EXPECT().UpdateUser(idInt, mock.Anything).Times(1).Return(nil, apiError)

		diags := ur.updateUser(ctx, plan, state)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the update user API call returns an empty
	// response.
	t.Run("update user returns an empty response", func(t *testing.T) {

		// Setup expectations.
		mockUser.EXPECT().UpdateUser(idInt, mock.Anything).Times(1).Return(nil, nil)

		diags := ur.updateUser(ctx, plan, state)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned when using invalid user id.
	t.Run("Using invalid user id returns an error", func(t *testing.T) {

		invalidPlan := &clumioUserResourceModel{
			Id:                         basetypes.NewStringValue(invalidId),
			Email:                      basetypes.NewStringValue(email),
			FullName:                   basetypes.NewStringValue(fullNameUpdated),
			AccessControlConfiguration: accessControlList,
		}

		// Setup Expectations
		mockUser.EXPECT().UpdateUser(invalidIdInt, mock.Anything).Times(1).Return(nil, nil)

		diags := ur.updateUser(ctx, invalidPlan, state)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Delete user success scenario.
//   - SDK API for delete user returns an not found error.
//   - SDK API for delete user returns an error.
//   - Polling of delete user task returns an error.
//   - Using invalid user id returns an error.
func TestDeleteUser(t *testing.T) {

	ctx := context.Background()
	mockUser := sdkclients.NewMockUserClient(t)
	ur := clumioUserResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkUsers: mockUser,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	urm := &clumioUserResourceModel{
		Id: basetypes.NewStringValue(id),
	}

	// Tests the success scenario for user deletion. It should not return diag.Diagnostics.
	t.Run("Success scenario for user deletion", func(t *testing.T) {

		// Setup Expectations
		mockUser.EXPECT().DeleteUser(idInt).Times(1).Return(nil, nil)

		diags := ur.deleteUser(ctx, urm)
		assert.Nil(t, diags)
	})

	// Tests that no error is returned if the user does not exist.
	t.Run("User not found should not return error", func(t *testing.T) {

		// Setup Expectations
		mockUser.EXPECT().DeleteUser(idInt).Times(1).Return(nil, apiNotFoundError)

		diags := ur.deleteUser(ctx, urm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned when delete user API call returns error.
	t.Run("delete user returns an error", func(t *testing.T) {

		// Setup Expectations
		mockUser.EXPECT().DeleteUser(idInt).Times(1).Return(nil, apiError)

		diags := ur.deleteUser(ctx, urm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned when using invalid user id.
	t.Run("Using invalid user id returns an error", func(t *testing.T) {

		invalidUrm := &clumioUserResourceModel{
			Id: basetypes.NewStringValue(invalidId),
		}

		// Setup Expectations
		mockUser.EXPECT().DeleteUser(invalidIdInt).Times(1).Return(nil, nil)

		diags := ur.deleteUser(ctx, invalidUrm)
		assert.NotNil(t, diags)
	})
}
