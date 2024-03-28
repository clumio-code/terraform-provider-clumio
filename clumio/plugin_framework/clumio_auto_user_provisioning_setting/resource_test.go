// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_auto_user_provisioning_setting

import (
	"context"
	"github.com/clumio-code/clumio-go-sdk/models"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	resourceName = "test_auto_user_provisioning_rule"
	id           = "test-id"
	name         = "test-auto-user-provisioning-rule"
	enabled      = true

	testError = "Test Error"
)

// Unit test for the following cases:
//   - Create auto user provisioning setting success scenario.
//   - SDK API for updating auto user provisioning setting returns error.
func TestCreateAutoUserProvisioningSetting(t *testing.T) {

	ctx := context.Background()
	aupClient := sdkclients.NewMockAutoUserProvisioningSettingClient(t)

	ar := autoUserProvisioningSettingResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkAUPSettings: aupClient,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	arm := &autoUserProvisioningSettingResourceModel{
		IsEnabled: basetypes.NewBoolValue(enabled),
	}

	// Tests the success scenario for clumio auto user provisioning setting create. It should not
	// return Diagnostics.
	t.Run("Basic success scenario for create auto user provisioning setting",
		func(t *testing.T) {

			//Setup expectations.
			aupClient.EXPECT().UpdateAutoUserProvisioningSetting(mock.Anything).Times(1).Return(
				nil, nil)

			diags := ar.createAutoUserProvisioningSetting(ctx, arm)
			assert.Nil(t, diags)
		})

	// Tests that Diagnostics is returned in case the update auto user provisioning setting API call
	// returns an error.
	t.Run("UpdateAutoUserProvisioningSetting returns an error", func(t *testing.T) {
		// Setup Expectations
		aupClient.EXPECT().UpdateAutoUserProvisioningSetting(mock.Anything).Times(1).Return(
			nil, apiError)

		diags := ar.createAutoUserProvisioningSetting(ctx, arm)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Read auto user provisioning setting success scenario.
//   - SDK API for reading auto user provisioning setting returns error.
//   - SDK API for reading auto user provisioning setting returns an empty response.
func TestReadAutoUserProvisioningSetting(t *testing.T) {

	ctx := context.Background()
	aupClient := sdkclients.NewMockAutoUserProvisioningSettingClient(t)

	ar := autoUserProvisioningSettingResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkAUPSettings: aupClient,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	arm := &autoUserProvisioningSettingResourceModel{
		IsEnabled: basetypes.NewBoolValue(enabled),
	}

	// Tests the success scenario for clumio auto user provisioning setting read. It should not
	// return Diagnostics.
	t.Run("Basic success scenario for read auto user provisioning setting",
		func(t *testing.T) {

			readResponse := &models.ReadAutoUserProvisioningSettingResponse{
				IsEnabled: &enabled,
			}

			//Setup expectations.
			aupClient.EXPECT().ReadAutoUserProvisioningSetting().Times(1).Return(
				readResponse, nil)

			diags := ar.readAutoUserProvisioningSetting(ctx, arm)
			assert.Nil(t, diags)
		})

	// Tests that Diagnostics is returned in case the read auto user provisioning setting API call
	// returns an error.
	t.Run("ReadAutoUserProvisioningSetting returns an error", func(t *testing.T) {
		// Setup Expectations
		aupClient.EXPECT().ReadAutoUserProvisioningSetting().Times(1).Return(nil, apiError)

		diags := ar.readAutoUserProvisioningSetting(ctx, arm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read auto user provisioning setting API call
	// returns an empty response.
	t.Run("ReadAutoUserProvisioningSetting returns an empty response", func(t *testing.T) {
		// Setup Expectations
		aupClient.EXPECT().ReadAutoUserProvisioningSetting().Times(1).Return(nil, nil)

		diags := ar.readAutoUserProvisioningSetting(ctx, arm)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Update auto user provisioning setting success scenario.
//   - SDK API for updating auto user provisioning setting returns error.
func TestUpdateAutoUserProvisioningSetting(t *testing.T) {

	ctx := context.Background()
	aupClient := sdkclients.NewMockAutoUserProvisioningSettingClient(t)

	ar := autoUserProvisioningSettingResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkAUPSettings: aupClient,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	arm := &autoUserProvisioningSettingResourceModel{
		IsEnabled: basetypes.NewBoolValue(enabled),
	}

	// Tests the success scenario for clumio auto user provisioning setting update. It should not
	// return Diagnostics.
	t.Run("Basic success scenario for update auto user provisioning setting",
		func(t *testing.T) {

			//Setup expectations.
			aupClient.EXPECT().UpdateAutoUserProvisioningSetting(mock.Anything).Times(1).Return(
				nil, nil)

			diags := ar.updateAutoUserProvisioningSetting(ctx, arm)
			assert.Nil(t, diags)
		})

	// Tests that Diagnostics is returned in case the update auto user provisioning setting API call
	// returns an error.
	t.Run("UpdateAutoUserProvisioningSetting returns an error", func(t *testing.T) {
		// Setup Expectations
		aupClient.EXPECT().UpdateAutoUserProvisioningSetting(mock.Anything).Times(1).Return(
			nil, apiError)

		diags := ar.updateAutoUserProvisioningSetting(ctx, arm)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Delete auto user provisioning setting success scenario.
//   - SDK API for updating auto user provisioning setting returns error.
func TestDeleteAutoUserProvisioningSetting(t *testing.T) {

	ctx := context.Background()
	aupClient := sdkclients.NewMockAutoUserProvisioningSettingClient(t)

	ar := autoUserProvisioningSettingResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkAUPSettings: aupClient,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	arm := &autoUserProvisioningSettingResourceModel{
		IsEnabled: basetypes.NewBoolValue(enabled),
	}

	// Tests the success scenario for clumio auto user provisioning setting delete. It should not
	// return Diagnostics.
	t.Run("Basic success scenario for delete auto user provisioning setting",
		func(t *testing.T) {

			//Setup expectations.
			aupClient.EXPECT().UpdateAutoUserProvisioningSetting(mock.Anything).Times(1).Return(
				nil, nil)

			diags := ar.deleteAutoUserProvisioningSetting(ctx, arm)
			assert.Nil(t, diags)
		})

	// Tests that Diagnostics is returned in case the update auto user provisioning setting API call
	// returns an error.
	t.Run("UpdateAutoUserProvisioningSetting returns an error", func(t *testing.T) {
		// Setup Expectations
		aupClient.EXPECT().UpdateAutoUserProvisioningSetting(mock.Anything).Times(1).Return(
			nil, apiError)

		diags := ar.deleteAutoUserProvisioningSetting(ctx, arm)
		assert.NotNil(t, diags)
	})
}
