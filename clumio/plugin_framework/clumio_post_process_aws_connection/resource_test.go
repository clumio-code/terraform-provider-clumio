// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_post_process_aws_connection

import (
	"context"
	"fmt"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var (
	accountId           = "test-aws-account"
	region              = "test-region"
	token               = "test-token"
	externalId          = "test-external-id"
	eventPubId          = "test-event-pub-id"
	roleArn             = "test-role-arn"
	version             = "1.0"
	intermediateRoleArn = "test-intermediate-role-arn"
	eventType           = "test-event-type"

	id = fmt.Sprintf("%s_%s", accountId, region)

	testError = "Test Error"
)

// Unit test for the following cases:
//   - Post-process AWS connection success scenario.
//   - SDK API for post-process AWS connection returns error.
func TestClumioPostProcessAWSConnectionCommon(t *testing.T) {

	ctx := context.Background()
	mockPostProcessConn := sdkclients.NewMockPostProcessAWSConnectionClient(t)

	pr := postProcessAWSConnectionResource{
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkPostProcessConn: mockPostProcessConn,
	}

	props, diags := basetypes.NewMapValueFrom(ctx, types.StringType, map[string]attr.Value{
		"key": basetypes.NewStringValue("value"),
	})
	assert.Nil(t, diags)

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	prm := postProcessAWSConnectionResourceModel{
		ID:                             basetypes.NewStringValue(id),
		AccountID:                      basetypes.NewStringValue(accountId),
		Token:                          basetypes.NewStringValue(token),
		RoleExternalID:                 basetypes.NewStringValue(externalId),
		Region:                         basetypes.NewStringValue(region),
		ClumioEventPubID:               basetypes.NewStringValue(eventPubId),
		RoleArn:                        basetypes.NewStringValue(roleArn),
		ConfigVersion:                  basetypes.NewStringValue(version),
		DiscoverVersion:                basetypes.NewStringValue(version),
		ProtectConfigVersion:           basetypes.NewStringValue(version),
		ProtectEBSVersion:              basetypes.NewStringValue(version),
		ProtectRDSVersion:              basetypes.NewStringValue(version),
		ProtectS3Version:               basetypes.NewStringValue(version),
		ProtectDynamoDBVersion:         basetypes.NewStringValue(version),
		ProtectEC2MssqlVersion:         basetypes.NewStringValue(version),
		ProtectWarmTierVersion:         basetypes.NewStringValue(version),
		ProtectWarmTierDynamoDBVersion: basetypes.NewStringValue(version),
		Properties:                     props,
		IntermediateRoleArn:            basetypes.NewStringValue(intermediateRoleArn),
	}

	// Tests the success scenario for post-process aws connection common function.
	// It should not return Diagnostics.
	t.Run("Success scenario for post-process common", func(t *testing.T) {

		//Setup expectations.
		mockPostProcessConn.EXPECT().PostProcessAwsConnection(mock.Anything).Times(1).
			Return(nil, nil)

		diags = pr.clumioPostProcessAWSConnectionCommon(ctx, prm, eventType)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the post-process aws connection API call returns
	// an error.
	t.Run("PostProcessAwsConnection returns an error", func(t *testing.T) {

		// Setup Expectations
		mockPostProcessConn.EXPECT().PostProcessAwsConnection(mock.Anything).Times(1).
			Return(nil, apiError)

		diags = pr.clumioPostProcessAWSConnectionCommon(ctx, prm, eventType)
		assert.NotNil(t, diags)
	})

}
