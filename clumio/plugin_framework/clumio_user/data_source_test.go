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
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit test for the following cases:
//   - Read user success scenario.
//   - SDK API for read user returns an error.
//   - SDK API for read user returns an empty response.
func TestDatasourceReadUser(t *testing.T) {

	ctx := context.Background()
	pgClient := sdkclients.NewMockUserClient(t)
	name := "test-use"
	resourceName := "test_user"
	id := "test-use-id"
	roleId := "test-role-id"
	ou := "test-ou"
	testError := "Test Error"

	rds := clumioUserDataSource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		userClient: pgClient,
	}

	rdsm := &clumioUserDataSourceModel{
		Name:   basetypes.NewStringValue(name),
		RoleId: basetypes.NewStringValue(roleId),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for user read. It should not return Diagnostics.
	t.Run("Basic success scenario for read user", func(t *testing.T) {

		count := int64(1)
		readResponse := &models.ListUsersResponse{
			Embedded: &models.UserListEmbedded{
				Items: []*models.UserWithETag{
					{
						Id:       &id,
						FullName: &name,
						AccessControlConfiguration: []*models.RoleForOrganizationalUnits{
							{
								RoleId:                &roleId,
								OrganizationalUnitIds: []*string{&ou},
							},
						},
					},
				},
			},
			CurrentCount: &count,
		}

		// Setup expectations.
		pgClient.EXPECT().ListUsers(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(readResponse, nil)

		diags := rds.readUser(ctx, rdsm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list users API call returns an
	// error.
	t.Run("list users returns an error", func(t *testing.T) {

		// Setup expectations.
		pgClient.EXPECT().ListUsers(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(nil, apiError)

		diags := rds.readUser(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list users API call returns an
	// empty response.
	t.Run("list users returns an empty response", func(t *testing.T) {

		// Setup expectations.
		pgClient.EXPECT().ListUsers(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(nil, nil)

		diags := rds.readUser(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list users API call returns an
	// empty response.
	t.Run("list users returns an empty items in response", func(t *testing.T) {

		readResponse := &models.ListUsersResponse{
			Embedded: &models.UserListEmbedded{
				Items: []*models.UserWithETag{},
			},
		}

		// Setup expectations.
		pgClient.EXPECT().ListUsers(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(readResponse, nil)

		diags := rds.readUser(ctx, rdsm)
		assert.NotNil(t, diags)
	})
}
