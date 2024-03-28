// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_post_process_kms

import (
	"context"
	"fmt"
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
	accountId  = "test-aws-account"
	region     = "test-region"
	token      = "test-token"
	roleId     = "test-role-id"
	externalId = "test-external-id"
	roleArn    = "test-role-arn"
	version    = int64(1)
	eventType  = "test-event-type"
	cmkKeyId   = "test-cmk-key"

	id = fmt.Sprintf("%s_%s", accountId, region)

	testError = "Test Error"
)

// Unit test for the following cases:
//   - Post-process AWS connection success scenario.
//   - SDK API for post-process AWS connection returns error.
func TestClumioPostProcessAWSConnectionCommon(t *testing.T) {

	ctx := context.Background()
	mockPostProcessKms := sdkclients.NewMockPostProcessKMSClient(t)

	pr := clumioPostProcessKmsResource{
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkPostProcessKMS: mockPostProcessKms,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	prm := clumioPostProcessKmsResourceModel{
		Id:                    basetypes.NewStringValue(id),
		AccountId:             basetypes.NewStringValue(accountId),
		Token:                 basetypes.NewStringValue(token),
		RoleId:                basetypes.NewStringValue(roleId),
		Region:                basetypes.NewStringValue(region),
		RoleExternalId:        basetypes.NewStringValue(externalId),
		CreatedMultiRegionCMK: basetypes.NewBoolValue(false),
		MultiRegionCMKKeyId:   basetypes.NewStringValue(cmkKeyId),
		RoleArn:               basetypes.NewStringValue(roleArn),
		TemplateVersion:       basetypes.NewInt64Value(version),
	}

	// Tests the success scenario for post-process kms common function. It should not return
	// Diagnostics.
	t.Run("Success scenario for post-process common", func(t *testing.T) {

		//Setup expectations.
		mockPostProcessKms.EXPECT().PostProcessKms(mock.Anything).Times(1).
			Return(nil, nil)

		diags := pr.clumioPostProcessKmsCommon(ctx, prm, eventType)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the post-process kms API call returns an error.
	t.Run("PostProcessKms returns an error", func(t *testing.T) {

		// Setup Expectations
		mockPostProcessKms.EXPECT().PostProcessKms(mock.Anything).Times(1).
			Return(nil, apiError)

		diags := pr.clumioPostProcessKmsCommon(ctx, prm, eventType)
		assert.NotNil(t, diags)
	})

}
