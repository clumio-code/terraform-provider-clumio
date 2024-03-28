// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_wallet

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

var (
	accountNativeId = "test-aws-account"
	clumioAccountId = "test-clumio-account-id"
	resourceName    = "test_wallet"
	id              = "test-wallet-id"
	token           = "test-token"
	state           = "test-state"

	testError = "Test Error"
)

// Unit test for the following cases:
//   - Create wallet success scenario.
//   - SDK API for create wallet returns error.
//   - SDK API for create wallet returns nil response.
func TestCreateWallet(t *testing.T) {

	mockWallet := sdkclients.NewMockWalletClient(t)
	ctx := context.Background()
	wr := clumioWalletResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkWallets: mockWallet,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Populate the protection group resource model to be used as input to createProtectionGroup()
	wrm := clumioWalletResourceModel{
		AccountNativeId: basetypes.NewStringValue(accountNativeId),
	}

	// Tests the success scenario for clumio wallet create. It should not return Diagnostics.
	t.Run("Basic success scenario for create wallet", func(t *testing.T) {

		createResponse := &models.CreateWalletResponse{
			AccountNativeId:    &accountNativeId,
			ClumioAwsAccountId: &clumioAccountId,
			Id:                 &id,
			Token:              &token,
			State:              &state,
		}

		// Setup Expectations
		mockWallet.EXPECT().CreateWallet(mock.Anything).Times(1).Return(createResponse, nil)

		diags := wr.createWallet(ctx, &wrm)
		assert.Nil(t, diags)
		assert.Equal(t, *createResponse.Id, wrm.Id.ValueString())
		assert.Equal(t, *createResponse.ClumioAwsAccountId, wrm.ClumioAccountId.ValueString())
		assert.Equal(t, *createResponse.Token, wrm.Token.ValueString())
		assert.Equal(t, *createResponse.State, wrm.State.ValueString())

	})

	// Tests that Diagnostics is returned in case the create wallet API call returns an
	// error.
	t.Run("CreateWallet returns an error", func(t *testing.T) {

		// Setup Expectations
		mockWallet.EXPECT().CreateWallet(mock.Anything).Times(1).Return(nil, apiError)

		diags := wr.createWallet(ctx, &wrm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the create wallet API call returns an
	// empty response.
	t.Run("CreateWallet returns an empty response", func(t *testing.T) {

		// Setup Expectations
		mockWallet.EXPECT().CreateWallet(mock.Anything).Times(1).Return(nil, nil)

		diags := wr.createWallet(ctx, &wrm)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Read wallet success scenario.
//   - SDK API for read wallet returns not found error.
//   - SDK API for read Clumio wallet returns error.
//   - SDK API for create Clumio wallet returns nil response.
func TestReadWallet(t *testing.T) {

	mockWallet := sdkclients.NewMockWalletClient(t)
	ctx := context.Background()
	cr := clumioWalletResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkWallets: mockWallet,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	// Populate the protection group resource model to be used as input to createProtectionGroup()
	crm := clumioWalletResourceModel{
		Id:              basetypes.NewStringValue(id),
		AccountNativeId: basetypes.NewStringValue(accountNativeId),
	}

	// Tests the success scenario for wallet read. It should not return Diagnostics.
	t.Run("success scenario for read wallet", func(t *testing.T) {

		readResponse := &models.ReadWalletResponse{
			AccountNativeId:    &accountNativeId,
			ClumioAwsAccountId: &clumioAccountId,
			Id:                 &id,
			Token:              &token,
			State:              &state,
		}

		// Setup Expectations
		mockWallet.EXPECT().ReadWallet(id).Times(1).Return(readResponse, nil)

		remove, diags := cr.readWallet(ctx, &crm)
		assert.Nil(t, diags)
		assert.False(t, remove)
	})

	// Tests that in case the wallet is not found, it returns true to indicate that the wallet
	// should be removed from the state.
	t.Run("read wallet returns not found error", func(t *testing.T) {

		// Setup Expectations
		mockWallet.EXPECT().ReadWallet(id).Times(1).
			Return(nil, apiNotFoundError)

		remove, diags := cr.readWallet(ctx, &crm)
		assert.True(t, remove)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read wallet API call returns an error.
	t.Run("read wallet returns error", func(t *testing.T) {

		// Setup Expectations
		mockWallet.EXPECT().ReadWallet(id).Times(1).
			Return(nil, apiError)

		remove, diags := cr.readWallet(ctx, &crm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that Diagnostics is returned in case the read wallet API call returns an empty
	// response.
	t.Run("read wallet returns nil response", func(t *testing.T) {

		// Setup Expectations
		mockWallet.EXPECT().ReadWallet(id).Times(1).
			Return(nil, nil)

		remove, diags := cr.readWallet(ctx, &crm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})
}

// Unit test for the following cases:
//   - Delete wallet success scenario.
//   - Delete wallet should not return error if wallet is not found.
//   - SDK API for delete wallet returns an error.
func TestDeleteWallet(t *testing.T) {

	mockWallet := sdkclients.NewMockWalletClient(t)
	ctx := context.Background()
	cr := clumioWalletResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkWallets: mockWallet,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	// Populate the protection group resource model to be used as input to createProtectionGroup()
	crm := &clumioWalletResourceModel{
		Id:              basetypes.NewStringValue(id),
		AccountNativeId: basetypes.NewStringValue(accountNativeId),
	}

	// Tests the success scenario for wallet deletion. It should not return diag.Diagnostics.
	t.Run("Success scenario for wallet deletion", func(t *testing.T) {

		// Setup Expectations
		mockWallet.EXPECT().DeleteWallet(id).Times(1).Return(nil, nil)

		diags := cr.deleteWallet(ctx, crm)
		assert.Nil(t, diags)
	})

	// Tests that no error is returned if the policy does not exist.
	t.Run("wallet not found should not return error", func(t *testing.T) {

		// Setup Expectations
		mockWallet.EXPECT().DeleteWallet(id).Times(1).Return(nil, apiNotFoundError)

		diags := cr.deleteWallet(ctx, crm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned when delete policy API call returns an error.
	t.Run("deleteWallet returns an error", func(t *testing.T) {

		// Setup Expectations
		mockWallet.EXPECT().DeleteWallet(id).Times(1).Return(nil, apiError)

		diags := cr.deleteWallet(ctx, crm)
		assert.NotNil(t, diags)
	})

}
