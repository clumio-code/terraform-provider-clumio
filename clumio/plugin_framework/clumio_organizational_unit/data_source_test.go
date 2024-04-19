// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_organizational_unit

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
//   - Read organizational unit success scenario.
//   - SDK API for read organizational unit returns an error.
//   - SDK API for read organizational unit returns an empty response.
func TestDatasourceReadOrganizationalUnit(t *testing.T) {

	ctx := context.Background()
	ouClient := sdkclients.NewMockOrganizationalUnitClient(t)
	name := "test-organizational-unit"
	resourceName := "test_organizational_unit"
	id := "test-organizational-unit-id"
	roleId := "test-role-id"
	ou := "test-ou"
	testError := "Test Error"

	rds := clumioOrganizationalUnitDataSource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		organizationalUnitClient: ouClient,
	}

	rdsm := &clumioOrganizationalUnitDataSourceModel{
		Name: basetypes.NewStringValue(name),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for organizational unit read. It should not return Diagnostics.
	t.Run("Basic success scenario for read organizational unit", func(t *testing.T) {

		count := int64(1)
		readResponse := &models.ListOrganizationalUnitsResponse{
			Embedded: &models.OrganizationalUnitListEmbedded{
				Items: []*models.OrganizationalUnitWithETag{
					{
						Id:   &id,
						Name: &ou,
						Users: []*models.UserWithRole{
							{
								AssignedRole: &roleId,
								UserId:       &userId,
							},
						},
					},
				},
			},
			CurrentCount: &count,
		}

		// Setup expectations.
		ouClient.EXPECT().ListOrganizationalUnits(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(readResponse, nil)

		diags := rds.readOrganizationalUnit(ctx, rdsm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list organizational units API call returns an
	// error.
	t.Run("list organizational units returns an error", func(t *testing.T) {

		// Setup expectations.
		ouClient.EXPECT().ListOrganizationalUnits(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(nil, apiError)

		diags := rds.readOrganizationalUnit(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list organizational units API call returns an
	// empty response.
	t.Run("list organizational units returns an empty response", func(t *testing.T) {

		// Setup expectations.
		ouClient.EXPECT().ListOrganizationalUnits(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(nil, nil)

		diags := rds.readOrganizationalUnit(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list organizational units API call returns an
	// empty response.
	t.Run("list organizational units returns an empty items in response", func(t *testing.T) {

		readResponse := &models.ListOrganizationalUnitsResponse{
			Embedded: &models.OrganizationalUnitListEmbedded{
				Items: []*models.OrganizationalUnitWithETag{},
			},
		}

		// Setup expectations.
		ouClient.EXPECT().ListOrganizationalUnits(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(readResponse, nil)

		diags := rds.readOrganizationalUnit(ctx, rdsm)
		assert.NotNil(t, diags)
	})
}
