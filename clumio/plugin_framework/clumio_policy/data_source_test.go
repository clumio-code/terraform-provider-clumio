// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_policy

import (
	"context"
	"testing"

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

// Unit test for the following cases:
//   - Read policy success scenario.
//   - SDK API for read policy returns an error.
//   - SDK API for read policy returns an empty response.
func TestDatasourceReadPolicy(t *testing.T) {

	ctx := context.Background()
	policyClient := sdkclients.NewMockPolicyDefinitionClient(t)
	name := "test-policy"
	resourceName := "test_policy"
	id := "test-policy-id"
	activationStatus := "activated"
	operationType := "test-operation-type"
	timezone := "test-timezone"
	ou := "test-ou"
	testError := "Test Error"

	rds := clumioPolicyDataSource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		policyDefinitionClient: policyClient,
	}

	opTypes, diags := basetypes.NewSetValueFrom(ctx, types.StringType, []string{operationType})
	assert.Nil(t, diags)

	rdsm := &clumioPolicyDataSourceModel{
		Name:             basetypes.NewStringValue(name),
		ActivationStatus: basetypes.NewStringValue(activationStatus),
		OperationTypes:   opTypes,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for policy read. It should not return Diagnostics.
	t.Run("Basic success scenario for read policy", func(t *testing.T) {

		readResponse := &models.ListPoliciesResponse{
			Embedded: &models.PolicyListEmbedded{
				Items: []*models.Policy{
					{
						ActivationStatus: &activationStatus,
						Id:               &id,
						Name:             &name,
						Operations: []*models.PolicyOperation{
							{
								ClumioType: &operationType,
							},
						},
						OrganizationalUnitId: &ou,
						Timezone:             &timezone,
					},
				},
			},
		}

		// Setup expectations.
		policyClient.EXPECT().ListPolicyDefinitions(mock.Anything, mock.Anything).Times(1).
			Return(readResponse, nil)

		diags := rds.readPolicy(ctx, rdsm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list policy definitions API call returns an
	// error.
	t.Run("list policies returns an error", func(t *testing.T) {

		// Setup expectations.
		policyClient.EXPECT().ListPolicyDefinitions(mock.Anything, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := rds.readPolicy(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list policy definitions API call returns an
	// empty response.
	t.Run("list policies returns an empty response", func(t *testing.T) {

		// Setup expectations.
		policyClient.EXPECT().ListPolicyDefinitions(mock.Anything, mock.Anything).Times(1).
			Return(nil, nil)

		diags := rds.readPolicy(ctx, rdsm)
		assert.NotNil(t, diags)
	})
}
