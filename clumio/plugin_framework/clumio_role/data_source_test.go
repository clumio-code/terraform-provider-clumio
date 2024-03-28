// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_role

import (
	"context"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
)

// Unit test for the following cases:
//   - Read role success scenario.
//   - SDK API for read role returns an error.
//   - SDK API for read role returns an empty response.
func TestReadRole(t *testing.T) {

	ctx := context.Background()
	roleClient := sdkclients.NewMockRoleClient(t)
	name := "test-role"
	resourceName := "test_role"
	id := "test-role-id"
	testError := "Test Error"
	permissionId := "test-permission-id"
	permissionName := "test-permission-name"
	permissionDesc := "test-permission-desc"

	rds := clumioRoleDataSource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		roles: roleClient,
	}

	rdsm := &clumioRoleDataSourceModel{
		Name: basetypes.NewStringValue(name),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for role read. It should not return Diagnostics.
	t.Run("Basic success scenario for read role", func(t *testing.T) {

		readResponse := &models.ListRolesResponse{
			Embedded: &models.RoleListEmbedded{
				Items: []*models.RoleWithETag{
					{
						Name: &name,
						Id:   &id,
						Permissions: []*models.PermissionModel{
							{
								Description: &permissionDesc,
								Id:          &permissionId,
								Name:        &permissionName,
							},
						},
					},
				},
			},
		}
		// Setup expectations.
		roleClient.EXPECT().ListRoles().Times(1).Return(readResponse, nil)

		diags := rds.readRole(ctx, rdsm)
		assert.Nil(t, diags)
		assert.Equal(t, id, rdsm.Id.ValueString())
	})

	// Tests that Diagnostics is returned in case the expected role name is not part of the response
	// of the list roles API call.
	t.Run("Expected role not present in list of roles returns an error", func(t *testing.T) {

		roleName := "some-role"
		readResponse := &models.ListRolesResponse{
			Embedded: &models.RoleListEmbedded{
				Items: []*models.RoleWithETag{
					{
						Name: &roleName,
						Id:   &id,
					},
				},
			},
		}
		// Setup expectations.
		roleClient.EXPECT().ListRoles().Times(1).Return(readResponse, nil)

		diags := rds.readRole(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list roles API call returns an error.
	t.Run("list roles returns an error", func(t *testing.T) {

		// Setup expectations.
		roleClient.EXPECT().ListRoles().Times(1).Return(nil, apiError)

		diags := rds.readRole(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list roles API call returns an empty response.
	t.Run("list roles returns an empty response", func(t *testing.T) {

		// Setup expectations.
		roleClient.EXPECT().ListRoles().Times(1).Return(nil, nil)

		diags := rds.readRole(ctx, rdsm)
		assert.NotNil(t, diags)
	})
}
