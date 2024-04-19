// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_protection_group

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
//   - Read protection group success scenario.
//   - SDK API for read protection group returns an error.
//   - SDK API for read protection group returns an empty response.
func TestDatasourceReadProtectionGroup(t *testing.T) {

	ctx := context.Background()
	pgClient := sdkclients.NewMockProtectionGroupClient(t)
	name := "test-protection-group"
	resourceName := "test_protection_group"
	id := "test-protection-group-id"
	testError := "Test Error"

	rds := clumioProtectionGroupDataSource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		protectionGroupClient: pgClient,
	}

	rdsm := &clumioProtectionGroupDataSourceModel{
		Name: basetypes.NewStringValue(name),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for protection group read. It should not return Diagnostics.
	t.Run("Basic success scenario for read protection group", func(t *testing.T) {

		count := int64(1)
		readResponse := &models.ListProtectionGroupsResponse{
			Embedded: &models.ProtectionGroupListEmbedded{
				Items: []*models.ProtectionGroup{
					{
						Id:   &id,
						Name: &name,
					},
				},
			},
			CurrentCount: &count,
		}

		// Setup expectations.
		pgClient.EXPECT().ListProtectionGroups(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(readResponse, nil)

		diags := rds.readProtectionGroup(ctx, rdsm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list protection groups API call returns an
	// error.
	t.Run("list protection groups returns an error", func(t *testing.T) {

		// Setup expectations.
		pgClient.EXPECT().ListProtectionGroups(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(nil, apiError)

		diags := rds.readProtectionGroup(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list protection groups API call returns an
	// empty response.
	t.Run("list protection groups returns an empty response", func(t *testing.T) {

		// Setup expectations.
		pgClient.EXPECT().ListProtectionGroups(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(nil, nil)

		diags := rds.readProtectionGroup(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list protection groups API call returns an
	// empty response.
	t.Run("list protection groups returns an empty items in response", func(t *testing.T) {

		readResponse := &models.ListProtectionGroupsResponse{
			Embedded: &models.ProtectionGroupListEmbedded{
				Items: []*models.ProtectionGroup{},
			},
		}

		// Setup expectations.
		pgClient.EXPECT().ListProtectionGroups(mock.Anything, mock.Anything, mock.Anything).
			Times(1).Return(readResponse, nil)

		diags := rds.readProtectionGroup(ctx, rdsm)
		assert.NotNil(t, diags)
	})
}
