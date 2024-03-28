// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_organizational_unit

import (
	"context"
	"testing"
	"time"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	name          = "test-organizational-unit"
	resourceName  = "test_organizational_unit"
	id            = "test-ou-id"
	description   = "test-description"
	parentId      = "test-parent-id"
	descendantId  = "test-descendant-id"
	childrenCount = int64(2)
	assignedRole  = "test-role"
	userId        = "test-user-id"
	dataSource    = "test-data-source"
	userCount     = int64(1)
	taskId        = "test-task-id"

	testError = "Test Error"
)

// Unit test for the following cases:
//   - Create organizational unit success scenario with HTTP 200 API response status.
//   - Create organizational unit success scenario with HTTP 202 API response status.
//   - SDK API for create organizational unit returns error.
//   - SDK API for create organizational unit returns nil response.
func TestCreateOrganizationalUnit(t *testing.T) {

	ctx := context.Background()
	mockOU := sdkclients.NewMockOrganizationalUnitClient(t)
	mockTask := sdkclients.NewMockTaskClient(t)

	or := &clumioOrganizationalUnitResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkOrgUnits: mockOU,
		sdkTasks:    mockTask,
	}

	orm := &clumioOrganizationalUnitResourceModel{
		Name:        basetypes.NewStringValue(name),
		Description: basetypes.NewStringValue(description),
		ParentId:    types.String{},
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for organizational unit create with the SDK API returning HTTP 200
	// status code. It should not return Diagnostics.
	t.Run("Success scenario for create organizational unit with HTTP 200 response",
		func(t *testing.T) {

			createResponse := &models.CreateOrganizationalUnitResponseWrapper{
				StatusCode: 200,
				Http200: &models.CreateOrganizationalUnitNoTaskResponse{
					ChildrenCount:             &childrenCount,
					ConfiguredDatasourceTypes: []*string{&dataSource},
					DescendantIds:             []*string{&descendantId},
					Description:               &description,
					Id:                        &id,
					Name:                      &name,
					ParentId:                  &parentId,
					UserCount:                 &userCount,
					Users: []*models.UserWithRole{
						{
							AssignedRole: &assignedRole,
							UserId:       &userId,
						},
					},
				},
			}

			//Setup expectations.
			mockOU.EXPECT().CreateOrganizationalUnit(mock.Anything, mock.Anything).Times(1).
				Return(createResponse, nil)

			diags := or.createOrganizationalUnit(ctx, orm)
			assert.Nil(t, diags)
			assert.Equal(t, id, orm.Id.ValueString())
			assert.Equal(t, name, orm.Name.ValueString())
			assert.Equal(t, description, orm.Description.ValueString())
			assert.Equal(t, parentId, orm.ParentId.ValueString())
			assert.Equal(t, userCount, orm.UserCount.ValueInt64())
			assert.Equal(t, childrenCount, orm.ChildrenCount.ValueInt64())
			usersWithRole := make([]*userWithRole, 0)
			diags = orm.UsersWithRole.ElementsAs(ctx, &usersWithRole, true)
			assert.Nil(t, diags)
			assert.Equal(t, userId, usersWithRole[0].UserId.ValueString())
			assert.Equal(t, assignedRole, usersWithRole[0].AssignedRole.ValueString())
			dataSources := make([]*string, 0)
			diags = orm.ConfiguredDatasourceTypes.ElementsAs(ctx, &dataSources, false)
			assert.Equal(t, dataSource, *dataSources[0])
		})

	// Tests the success scenario for organizational unit create with the SDK API returning HTTP 202
	// status code. It should not return Diagnostics.
	t.Run("Success scenario for create organizational unit with HTTP 202 response",
		func(t *testing.T) {

			createResponse := &models.CreateOrganizationalUnitResponseWrapper{
				StatusCode: 202,
				Http202: &models.CreateOrganizationalUnitResponse{
					ChildrenCount:             &childrenCount,
					ConfiguredDatasourceTypes: []*string{&dataSource},
					DescendantIds:             []*string{&descendantId},
					Description:               &description,
					Id:                        &id,
					Name:                      &name,
					ParentId:                  &parentId,
					UserCount:                 &userCount,
					Users: []*models.UserWithRole{
						{
							AssignedRole: &assignedRole,
							UserId:       &userId,
						},
					},
				},
			}

			//Setup expectations.
			mockOU.EXPECT().CreateOrganizationalUnit(mock.Anything, mock.Anything).Times(1).
				Return(createResponse, nil)

			diags := or.createOrganizationalUnit(ctx, orm)
			assert.Nil(t, diags)
			assert.Equal(t, id, orm.Id.ValueString())
			assert.Equal(t, name, orm.Name.ValueString())
			assert.Equal(t, description, orm.Description.ValueString())
			assert.Equal(t, parentId, orm.ParentId.ValueString())
			assert.Equal(t, userCount, orm.UserCount.ValueInt64())
			assert.Equal(t, childrenCount, orm.ChildrenCount.ValueInt64())
			usersWithRole := make([]*userWithRole, 0)
			diags = orm.UsersWithRole.ElementsAs(ctx, &usersWithRole, true)
			assert.Nil(t, diags)
			assert.Equal(t, userId, usersWithRole[0].UserId.ValueString())
			assert.Equal(t, assignedRole, usersWithRole[0].AssignedRole.ValueString())
			dataSources := make([]*string, 0)
			diags = orm.ConfiguredDatasourceTypes.ElementsAs(ctx, &dataSources, false)
			assert.Equal(t, dataSource, *dataSources[0])
		})

	// Tests that Diagnostics is returned in case the create organizational unit API call returns
	// error.
	t.Run("CreateOrganizationalUnit returns error", func(t *testing.T) {

		// Setup Expectations
		mockOU.EXPECT().CreateOrganizationalUnit(mock.Anything, mock.Anything).Times(1).Return(
			nil, apiError)

		diags := or.createOrganizationalUnit(ctx, orm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the create organizational unit API call returns an
	// empty response.
	t.Run("CreateOrganizationalUnit returns nil response", func(t *testing.T) {

		// Setup Expectations
		mockOU.EXPECT().CreateOrganizationalUnit(mock.Anything, mock.Anything).Times(1).Return(
			nil, nil)

		diags := or.createOrganizationalUnit(ctx, orm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Read organizational unit success scenario.
//   - SDK API for read organizational unit returns not found error.
//   - SDK API for read organizational unit returns error.
//   - SDK API for read organizational unit returns an empty response.
func TestReadOrganizationalUnit(t *testing.T) {

	ctx := context.Background()
	mockOU := sdkclients.NewMockOrganizationalUnitClient(t)

	or := &clumioOrganizationalUnitResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkOrgUnits: mockOU,
	}

	orm := &clumioOrganizationalUnitResourceModel{
		Id:          basetypes.NewStringValue(id),
		Name:        basetypes.NewStringValue(name),
		Description: basetypes.NewStringValue(description),
		ParentId:    types.String{},
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	// Create the response of the SDK ReadOrganizationalUnit() API.
	readResponse := &models.ReadOrganizationalUnitResponse{
		ChildrenCount:             &childrenCount,
		ConfiguredDatasourceTypes: []*string{&dataSource},
		DescendantIds:             []*string{&descendantId},
		Description:               &description,
		Id:                        &id,
		Name:                      &name,
		ParentId:                  &parentId,
		UserCount:                 &userCount,
		Users: []*models.UserWithRole{
			{
				AssignedRole: &assignedRole,
				UserId:       &userId,
			},
		},
	}

	// Tests the success scenario for organizational unit read. It should not return Diagnostics.
	t.Run("success scenario for read organizational unit", func(t *testing.T) {

		// Setup Expectations
		mockOU.EXPECT().ReadOrganizationalUnit(id, mock.Anything).Times(1).
			Return(readResponse, nil)

		remove, diags := or.readOrganizationalUnit(ctx, orm)
		assert.Nil(t, diags)
		assert.False(t, remove)
	})

	// Tests that in case the organizational unit is not found, it returns true to indicate that the
	// organizational unit should be removed from the state.
	t.Run("read organizational unit returns not found error", func(t *testing.T) {

		// Setup Expectations
		mockOU.EXPECT().ReadOrganizationalUnit(id, mock.Anything).Times(1).
			Return(nil, apiNotFoundError)

		remove, diags := or.readOrganizationalUnit(ctx, orm)
		assert.Nil(t, diags)
		assert.True(t, remove)
	})

	// Tests that Diagnostics is returned in case the read organizational unit API call returns
	// error.
	t.Run("read organizational unit returns error", func(t *testing.T) {

		// Setup Expectations
		mockOU.EXPECT().ReadOrganizationalUnit(id, mock.Anything).Times(1).
			Return(nil, apiError)

		remove, diags := or.readOrganizationalUnit(ctx, orm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that Diagnostics is returned in case the read organizational unit API call returns an
	// empty response.
	t.Run("read organizational unit returns nil response", func(t *testing.T) {

		// Setup Expectations
		mockOU.EXPECT().ReadOrganizationalUnit(id, mock.Anything).Times(1).
			Return(nil, nil)

		remove, diags := or.readOrganizationalUnit(ctx, orm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})
}

// Unit test for the following cases:
//   - Update organizational unit success scenario with HTTP 200 API response status.
//   - Update organizational unit success scenario with HTTP 202 API response status.
//   - SDK API for update organizational unit returns error.
//   - SDK API for update organizational unit returns nil response.
func TestUpdateOrganizationalUnit(t *testing.T) {

	ctx := context.Background()
	mockOU := sdkclients.NewMockOrganizationalUnitClient(t)

	or := &clumioOrganizationalUnitResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkOrgUnits: mockOU,
	}

	orm := &clumioOrganizationalUnitResourceModel{
		Id:          basetypes.NewStringValue(id),
		Name:        basetypes.NewStringValue(name),
		Description: basetypes.NewStringValue(description),
		ParentId:    types.String{},
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for organizational unit update with the SDK API returning HTTP 200
	// status code. It should not return Diagnostics.
	t.Run("Success scenario for update organizational unit with HTTP 200 response",
		func(t *testing.T) {

			updateResponse := &models.PatchOrganizationalUnitResponseWrapper{
				StatusCode: 200,
				Http200: &models.PatchOrganizationalUnitNoTaskResponse{
					ChildrenCount:             &childrenCount,
					ConfiguredDatasourceTypes: []*string{&dataSource},
					DescendantIds:             []*string{&descendantId},
					Description:               &description,
					Id:                        &id,
					Name:                      &name,
					ParentId:                  &parentId,
					UserCount:                 &userCount,
					Users: []*models.UserWithRole{
						{
							AssignedRole: &assignedRole,
							UserId:       &userId,
						},
					},
				},
			}

			//Setup expectations.
			mockOU.EXPECT().PatchOrganizationalUnit(id, mock.Anything, mock.Anything).Times(1).
				Return(updateResponse, nil)

			diags := or.updateOrganizationalUnit(ctx, orm)
			assert.Nil(t, diags)
			assert.Equal(t, id, orm.Id.ValueString())
			assert.Equal(t, name, orm.Name.ValueString())
			assert.Equal(t, description, orm.Description.ValueString())
			assert.Equal(t, parentId, orm.ParentId.ValueString())
			assert.Equal(t, userCount, orm.UserCount.ValueInt64())
			assert.Equal(t, childrenCount, orm.ChildrenCount.ValueInt64())
			usersWithRole := make([]*userWithRole, 0)
			diags = orm.UsersWithRole.ElementsAs(ctx, &usersWithRole, true)
			assert.Nil(t, diags)
			assert.Equal(t, userId, usersWithRole[0].UserId.ValueString())
			assert.Equal(t, assignedRole, usersWithRole[0].AssignedRole.ValueString())
			dataSources := make([]*string, 0)
			diags = orm.ConfiguredDatasourceTypes.ElementsAs(ctx, &dataSources, false)
			assert.Equal(t, dataSource, *dataSources[0])
		})

	// Tests the success scenario for organizational unit update with the SDK API returning HTTP 202
	// status code. It should not return Diagnostics.
	t.Run("Success scenario for update organizational unit with HTTP 202 response",
		func(t *testing.T) {

			updateResponse := &models.PatchOrganizationalUnitResponseWrapper{
				StatusCode: 202,
				Http202: &models.PatchOrganizationalUnitResponse{
					ChildrenCount:             &childrenCount,
					ConfiguredDatasourceTypes: []*string{&dataSource},
					DescendantIds:             []*string{&descendantId},
					Description:               &description,
					Id:                        &id,
					Name:                      &name,
					ParentId:                  &parentId,
					UserCount:                 &userCount,
					Users: []*models.UserWithRole{
						{
							AssignedRole: &assignedRole,
							UserId:       &userId,
						},
					},
				},
			}

			//Setup expectations.
			mockOU.EXPECT().PatchOrganizationalUnit(id, mock.Anything, mock.Anything).Times(1).
				Return(updateResponse, nil)

			diags := or.updateOrganizationalUnit(ctx, orm)
			assert.Nil(t, diags)
			assert.Equal(t, id, orm.Id.ValueString())
			assert.Equal(t, name, orm.Name.ValueString())
			assert.Equal(t, description, orm.Description.ValueString())
			assert.Equal(t, parentId, orm.ParentId.ValueString())
			assert.Equal(t, userCount, orm.UserCount.ValueInt64())
			assert.Equal(t, childrenCount, orm.ChildrenCount.ValueInt64())
			usersWithRole := make([]*userWithRole, 0)
			diags = orm.UsersWithRole.ElementsAs(ctx, &usersWithRole, true)
			assert.Nil(t, diags)
			assert.Equal(t, userId, usersWithRole[0].UserId.ValueString())
			assert.Equal(t, assignedRole, usersWithRole[0].AssignedRole.ValueString())
			dataSources := make([]*string, 0)
			diags = orm.ConfiguredDatasourceTypes.ElementsAs(ctx, &dataSources, false)
			assert.Equal(t, dataSource, *dataSources[0])
		})

	// Tests that Diagnostics is returned in case the patch organizational unit API call returns
	// error.
	t.Run("PatchOrganizationalUnit returns error", func(t *testing.T) {

		// Setup Expectations
		mockOU.EXPECT().PatchOrganizationalUnit(id, mock.Anything, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := or.updateOrganizationalUnit(ctx, orm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the patch organizational unit API call returns an
	// empty response.
	t.Run("PatchOrganizationalUnit returns nil response", func(t *testing.T) {

		// Setup Expectations
		mockOU.EXPECT().PatchOrganizationalUnit(id, mock.Anything, mock.Anything).Times(1).
			Return(nil, nil)

		diags := or.updateOrganizationalUnit(ctx, orm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Read organizational unit success scenario.
//   - SDK API for delete organizational unit returns not found error.
//   - SDK API for delete organizational unit returns error.
func TestDeleteOrganizationalUnit(t *testing.T) {

	ctx := context.Background()
	mockOU := sdkclients.NewMockOrganizationalUnitClient(t)
	mockTask := sdkclients.NewMockTaskClient(t)

	or := &clumioOrganizationalUnitResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkOrgUnits:  mockOU,
		sdkTasks:     mockTask,
		pollTimeout:  5 * time.Second,
		pollInterval: 1,
	}

	orm := &clumioOrganizationalUnitResourceModel{
		Id:          basetypes.NewStringValue(id),
		Name:        basetypes.NewStringValue(name),
		Description: basetypes.NewStringValue(description),
		ParentId:    types.String{},
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	deleteResponse := &models.DeleteOrganizationalUnitResponse{
		TaskId: &taskId,
	}

	taskStatus := common.TaskSuccess
	readTaskResponse := &models.ReadTaskResponse{
		Status: &taskStatus,
	}

	// Tests the success scenario for organizational unit deletion. It should not return diag.Diagnostics.
	t.Run("Success scenario for organizational unit deletion", func(t *testing.T) {

		// Setup Expectations
		mockOU.EXPECT().DeleteOrganizationalUnit(id, mock.Anything).Times(1).Return(
			deleteResponse, nil)
		mockTask.EXPECT().ReadTask(taskId).Times(1).Return(readTaskResponse, nil)

		diags := or.deleteOrganizationalUnit(ctx, orm)
		assert.Nil(t, diags)
	})

	// Tests that no error is returned if the organizational unit does not exist.
	t.Run("Organizational unit not found should not return error", func(t *testing.T) {

		// Setup Expectations
		mockOU.EXPECT().DeleteOrganizationalUnit(id, mock.Anything).Times(1).Return(
			nil, apiNotFoundError)

		diags := or.deleteOrganizationalUnit(ctx, orm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned when delete organizational unit API call returns error.
	t.Run("DeleteOrganizationalUnit returns error", func(t *testing.T) {

		// Setup Expectations
		mockOU.EXPECT().DeleteOrganizationalUnit(id, mock.Anything).Times(1).Return(
			nil, apiError)

		diags := or.deleteOrganizationalUnit(ctx, orm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned when organizational unit deletion returns an empty
	// response.
	t.Run("DeleteOrganizationalUnit returns an empty response", func(t *testing.T) {

		// Setup Expectations
		mockOU.EXPECT().DeleteOrganizationalUnit(id, mock.Anything).Times(1).Return(
			nil, nil)

		diags := or.deleteOrganizationalUnit(ctx, orm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned when polling of the delete organizational unit task fails.
	t.Run("Task poll returns error", func(t *testing.T) {
		// Setup Expectations
		mockOU.EXPECT().DeleteOrganizationalUnit(id, mock.Anything).Times(1).Return(
			deleteResponse, nil)
		mockTask.EXPECT().ReadTask(taskId).Times(1).Return(nil, apiError)

		diags := or.deleteOrganizationalUnit(ctx, orm)
		assert.NotNil(t, diags)
	})
}
