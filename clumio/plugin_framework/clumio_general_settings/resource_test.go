// Copyright 2025. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_general_settings

import (
	"context"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	resourceName                   = "test_general_settings"
	testError                      = "Test Error"
	testAutoLogoutDuration         = int64(1200)
	testPasswordExpirationDuration = int64(7776000)
	testIp1                        = "192.168.1.1"
	testIp2                        = "192.168.1.2"
)

// Unit test for the following cases:
//   - Read general settings success scenario.
//   - SDK API for read general settings returns an error.
//   - SDK API for read general settings returns an empty response.
func TestReadGeneralSettings(t *testing.T) {

	ctx := context.Background()
	mockGeneralSettings := sdkclients.NewMockGeneralSettingsClient(t)
	rcr := &clumioGeneralSettings{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkGeneralSettings: mockGeneralSettings,
	}

	model := &generalSettingsResourceModel{
		AutoLogoutDuration:         types.Int64Value(testAutoLogoutDuration),
		PasswordExpirationDuration: types.Int64Value(testPasswordExpirationDuration),
		IpAllowlist: []types.String{
			types.StringValue(testIp1),
			types.StringValue(testIp2),
		},
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for read general settings. It should not return Diagnostics.
	t.Run("Basic success scenario for read general settings", func(t *testing.T) {

		resp := &models.ReadGeneralSettingsResponseV2{
			AutoLogoutDuration:         &testAutoLogoutDuration,
			PasswordExpirationDuration: &testPasswordExpirationDuration,
			IpAllowlist:                []*string{&testIp1, &testIp2},
		}

		mockGeneralSettings.EXPECT().ReadGeneralSettings().Times(1).Return(resp, nil)
		diags := rcr.readGeneralSettings(ctx, model)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read general settings API call returns
	// an error.
	t.Run("Read general settings returns an error", func(t *testing.T) {

		mockGeneralSettings.EXPECT().ReadGeneralSettings().Times(1).Return(nil, apiError)
		diags := rcr.readGeneralSettings(ctx, model)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read general settings API call returns an
	// empty response.
	t.Run("Read general settings returns an empty response", func(t *testing.T) {

		mockGeneralSettings.EXPECT().ReadGeneralSettings().Times(1).Return(nil, nil)
		diags := rcr.readGeneralSettings(ctx, model)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Update general settings success scenario.
//   - SDK API for update general settings returns an error.
//   - SDK API for update general settings returns an empty response.
func TestUpdateGeneralSettings(t *testing.T) {

	ctx := context.Background()
	mockGeneralSettings := sdkclients.NewMockGeneralSettingsClient(t)
	rcr := &clumioGeneralSettings{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkGeneralSettings: mockGeneralSettings,
	}

	model := &generalSettingsResourceModel{
		AutoLogoutDuration:         types.Int64Value(testAutoLogoutDuration),
		PasswordExpirationDuration: types.Int64Value(testPasswordExpirationDuration),
		IpAllowlist: []types.String{
			types.StringValue(testIp1),
			types.StringValue(testIp2),
		},
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for update general settings. It should not return Diagnostics.
	t.Run("Basic success scenario for update general settings", func(t *testing.T) {

		resp := &models.PatchGeneralSettingsResponseV2{
			AutoLogoutDuration:         &testAutoLogoutDuration,
			PasswordExpirationDuration: &testPasswordExpirationDuration,
			IpAllowlist:                []*string{&testIp1, &testIp2},
		}

		mockGeneralSettings.EXPECT().UpdateGeneralSettings(mock.Anything).Times(1).Return(resp, nil)
		diags := rcr.updateGeneralSettings(ctx, model)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the update general settings API call returns
	// an error.
	t.Run("Update general settings returns an error", func(t *testing.T) {

		mockGeneralSettings.EXPECT().UpdateGeneralSettings(mock.Anything).Times(1).
			Return(nil, apiError)

		diags := rcr.updateGeneralSettings(ctx, model)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the update general settings API call returns
	// an empty response.
	t.Run("Update general settings returns an empty response", func(t *testing.T) {

		mockGeneralSettings.EXPECT().UpdateGeneralSettings(mock.Anything).Times(1).
			Return(nil, nil)

		diags := rcr.updateGeneralSettings(ctx, model)
		assert.NotNil(t, diags)
	})
}
