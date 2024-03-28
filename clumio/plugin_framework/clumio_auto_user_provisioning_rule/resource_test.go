// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_auto_user_provisioning_rule

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

var (
	resourceName = "test_auto_user_provisioning_rule"
	id           = "test-id"
	ou           = "test-ou"
	name         = "test-auto-user-provisioning-rule"
	condition    = "test-condition"
	roleId       = "test-role-id"

	testError = "Test Error"
)

// Unit test for the following cases:
//   - Create auto user provisioning rule success scenario.
//   - SDK API for create auto user provisioning rule returns error.
//   - SDK API for create auto user provisioning rule returns an empty response.
func TestCreateAutoUserProvisioningRule(t *testing.T) {

	ctx := context.Background()
	aupClient := sdkclients.NewMockAutoUserProvisioningRuleClient(t)

	ar := autoUserProvisioningRuleResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkAUPRules: aupClient,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	ouIds, diags := basetypes.NewSetValueFrom(ctx, types.StringType, []string{ou})
	assert.Nil(t, diags)
	arm := &autoUserProvisioningRuleResourceModel{
		Name:                  basetypes.NewStringValue(name),
		Condition:             basetypes.NewStringValue(condition),
		RoleID:                basetypes.NewStringValue(roleId),
		OrganizationalUnitIDs: ouIds,
	}

	// Tests the success scenario for clumio auto user provisioning rule create. It should not
	// return Diagnostics.
	t.Run("Basic success scenario for create auto user provisioning rule",
		func(t *testing.T) {

			createResponse := &models.CreateAutoUserProvisioningRuleResponse{
				Condition: &condition,
				Name:      &name,
				Provision: &models.RuleProvision{
					OrganizationalUnitIds: []*string{&ou},
					RoleId:                &roleId,
				},
				RuleId: &id,
			}

			//Setup expectations.
			aupClient.EXPECT().CreateAutoUserProvisioningRule(mock.Anything).Times(1).Return(
				createResponse, nil)

			diags = ar.createAutoUserProvisioningRule(ctx, arm)
			assert.Nil(t, diags)
			assert.Equal(t, id, arm.ID.ValueString())
		})

	// Tests that Diagnostics is returned in case the create auto user provisioning rule API call
	// returns an error.
	t.Run("CreateAutoUserProvisioningRule returns an error", func(t *testing.T) {
		// Setup Expectations
		aupClient.EXPECT().CreateAutoUserProvisioningRule(mock.Anything).Times(1).Return(
			nil, apiError)

		diags := ar.createAutoUserProvisioningRule(ctx, arm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the create auto user provisioning rule API call
	// returns an empty response.
	t.Run("CreateAutoUserProvisioningRule returns an empty response", func(t *testing.T) {
		// Setup Expectations
		aupClient.EXPECT().CreateAutoUserProvisioningRule(mock.Anything).Times(1).Return(
			nil, nil)

		diags := ar.createAutoUserProvisioningRule(ctx, arm)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Read auto user provisioning rule success scenario.
//   - SDK API for create auto user provisioning rule returns not found error.
//   - SDK API for create auto user provisioning rule returns error.
//   - SDK API for create auto user provisioning rule returns an empty response.
func TestReadAutoUserProvisioningRule(t *testing.T) {

	ctx := context.Background()
	aupClient := sdkclients.NewMockAutoUserProvisioningRuleClient(t)

	ar := autoUserProvisioningRuleResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkAUPRules: aupClient,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	ouIds, diags := basetypes.NewSetValueFrom(ctx, types.StringType, []string{ou})
	assert.Nil(t, diags)
	arm := &autoUserProvisioningRuleResourceModel{
		ID:                    basetypes.NewStringValue(id),
		Name:                  basetypes.NewStringValue(name),
		Condition:             basetypes.NewStringValue(condition),
		RoleID:                basetypes.NewStringValue(roleId),
		OrganizationalUnitIDs: ouIds,
	}

	// Tests the success scenario for clumio auto user provisioning rule read. It should not
	// return Diagnostics.
	t.Run("Basic success scenario for read auto user provisioning rule", func(t *testing.T) {

		readResponse := &models.ReadAutoUserProvisioningRuleResponse{
			Condition: &condition,
			Name:      &name,
			Provision: &models.RuleProvision{
				OrganizationalUnitIds: []*string{&ou},
				RoleId:                &roleId,
			},
			RuleId: &roleId,
		}

		//Setup expectations.
		aupClient.EXPECT().ReadAutoUserProvisioningRule(mock.Anything).Times(1).Return(
			readResponse, nil)

		remove, diags := ar.readAutoUserProvisioningRule(ctx, arm)
		assert.Nil(t, diags)
		assert.False(t, remove)
	})

	// Tests that Diagnostics is returned in case the create auto user provisioning rule API call
	// returns not found error.
	t.Run("ReadAutoUserProvisioningRule returns not found error", func(t *testing.T) {
		// Setup Expectations
		aupClient.EXPECT().ReadAutoUserProvisioningRule(mock.Anything).Times(1).Return(
			nil, apiNotFoundError)

		remove, diags := ar.readAutoUserProvisioningRule(ctx, arm)
		assert.Nil(t, diags)
		assert.True(t, remove)
	})

	// Tests that Diagnostics is returned in case the create auto user provisioning rule API call
	// returns an error.
	t.Run("ReadAutoUserProvisioningRule returns an error", func(t *testing.T) {
		// Setup Expectations
		aupClient.EXPECT().ReadAutoUserProvisioningRule(mock.Anything).Times(1).Return(
			nil, apiError)

		remove, diags := ar.readAutoUserProvisioningRule(ctx, arm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that Diagnostics is returned in case the create auto user provisioning rule API call
	// returns an empty response.
	t.Run("CreateAutoUserProvisioningRule returns an empty response", func(t *testing.T) {
		// Setup Expectations
		aupClient.EXPECT().ReadAutoUserProvisioningRule(mock.Anything).Times(1).Return(
			nil, nil)

		remove, diags := ar.readAutoUserProvisioningRule(ctx, arm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})
}

// Unit test for the following cases:
//   - Update auto user provisioning rule success scenario.
//   - SDK API for update auto user provisioning rule returns error.
//   - SDK API for update auto user provisioning rule returns an empty response.
func TestUpdateAutoUserProvisioningRule(t *testing.T) {

	ctx := context.Background()
	aupClient := sdkclients.NewMockAutoUserProvisioningRuleClient(t)

	ar := autoUserProvisioningRuleResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkAUPRules: aupClient,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	ouIds, diags := basetypes.NewSetValueFrom(ctx, types.StringType, []string{ou})
	assert.Nil(t, diags)
	arm := &autoUserProvisioningRuleResourceModel{
		ID:                    basetypes.NewStringValue(id),
		Name:                  basetypes.NewStringValue(name),
		Condition:             basetypes.NewStringValue(condition),
		RoleID:                basetypes.NewStringValue(roleId),
		OrganizationalUnitIDs: ouIds,
	}

	// Tests the success scenario for clumio auto user provisioning rule update. It should not
	// return Diagnostics.
	t.Run("Basic success scenario for update auto user provisioning rule",
		func(t *testing.T) {

			createResponse := &models.UpdateAutoUserProvisioningRuleResponse{
				Condition: &condition,
				Name:      &name,
				Provision: &models.RuleProvision{
					OrganizationalUnitIds: []*string{&ou},
					RoleId:                &roleId,
				},
				RuleId: &id,
			}

			//Setup expectations.
			aupClient.EXPECT().UpdateAutoUserProvisioningRule(id, mock.Anything).Times(1).Return(
				createResponse, nil)

			diags = ar.updateAutoUserProvisioningRule(ctx, arm)
			assert.Nil(t, diags)
			assert.Equal(t, id, arm.ID.ValueString())
		})

	// Tests that Diagnostics is returned in case the update auto user provisioning rule API call
	// returns an error.
	t.Run("UpdateAutoUserProvisioningRule returns an error", func(t *testing.T) {
		// Setup Expectations
		aupClient.EXPECT().UpdateAutoUserProvisioningRule(id, mock.Anything).Times(1).Return(
			nil, apiError)

		diags := ar.updateAutoUserProvisioningRule(ctx, arm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the update auto user provisioning rule API call
	// returns an empty response.
	t.Run("UpdateAutoUserProvisioningRule returns an empty response", func(t *testing.T) {
		// Setup Expectations
		aupClient.EXPECT().UpdateAutoUserProvisioningRule(id, mock.Anything).Times(1).Return(
			nil, nil)

		diags := ar.updateAutoUserProvisioningRule(ctx, arm)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Delete auto user provisioning rule success scenario.
//   - Delete auto user provisioning rule should not return error if auto user provisioning rule is
//     not found.
//   - SDK API for delete auto user provisioning rule returns an error.
func TestDeleteAutoUserProvisioningRule(t *testing.T) {

	ctx := context.Background()
	aupClient := sdkclients.NewMockAutoUserProvisioningRuleClient(t)

	ar := autoUserProvisioningRuleResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkAUPRules: aupClient,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	ouIds, diags := basetypes.NewSetValueFrom(ctx, types.StringType, []string{ou})
	assert.Nil(t, diags)
	arm := &autoUserProvisioningRuleResourceModel{
		ID:                    basetypes.NewStringValue(id),
		Name:                  basetypes.NewStringValue(name),
		Condition:             basetypes.NewStringValue(condition),
		RoleID:                basetypes.NewStringValue(roleId),
		OrganizationalUnitIDs: ouIds,
	}

	// Tests the success scenario for clumio auto user provisioning rule delete. It should not
	// return Diagnostics.
	t.Run("Basic success scenario for delete auto user provisioning rule", func(t *testing.T) {

		// Setup Expectations
		aupClient.EXPECT().DeleteAutoUserProvisioningRule(id).Times(1).Return(nil, nil)

		diags = ar.deleteAutoUserProvisioningRule(ctx, arm)
		assert.Nil(t, diags)
	})

	// Tests that no error is returned if the auto user provisioning rule does not exist.
	t.Run("Auto user provisioning rule not found should not return error",
		func(t *testing.T) {

			// Setup Expectations
			aupClient.EXPECT().DeleteAutoUserProvisioningRule(id).Times(1).Return(
				nil, apiNotFoundError)

			diags = ar.deleteAutoUserProvisioningRule(ctx, arm)
			assert.Nil(t, diags)
		})

	// Tests that Diagnostics is returned in case the delete auto user provisioning rule API call
	// returns an error.
	t.Run("DeleteAutoUserProvisioningRule returns an error", func(t *testing.T) {
		// Setup Expectations
		aupClient.EXPECT().DeleteAutoUserProvisioningRule(id).Times(1).Return(nil, apiError)

		diags := ar.deleteAutoUserProvisioningRule(ctx, arm)
		assert.NotNil(t, diags)
	})
}
