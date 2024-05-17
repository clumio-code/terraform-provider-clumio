// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_post_process_aws_connection

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
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
	invalidVersion      = "1.2.3"
	emptyVersion        = ""
	intermediateRoleArn = "test-intermediate-role-arn"
	eventType           = eventTypeUpdate

	id = fmt.Sprintf("%s_%s", accountId, region)

	testError = "Test Error"
)

// Unit test for the following cases:
//   - Post-process AWS connection success scenario.
//   - Get template configuration returns an error.
//   - SDK API for post-process AWS connection returns an error.
func TestClumioPostProcessAWSConnectionCommon(t *testing.T) {

	ctx := context.Background()
	mockPostProcessConn := sdkclients.NewMockPostProcessAWSConnectionClient(t)
	mockAWSConn := sdkclients.NewMockAWSConnectionClient(t)

	pr := postProcessAWSConnectionResource{
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkPostProcessConn: mockPostProcessConn,
		sdkAWSConnection:   mockAWSConn,
		pollInterval:       1,
		pollTimeout:        5 * time.Second,
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

	// Tests the success scenario for post-process aws connection common function with event type
	// as Create. It should not return Diagnostics.
	t.Run("Success scenario for post-process common - Create", func(t *testing.T) {

		eventType := eventTypeCreate
		connStatusConnected := connected
		//Setup expectations.
		mockAWSConn.EXPECT().ReadAwsConnection(mock.Anything, mock.Anything).Times(1).Return(
			&models.ReadAWSConnectionResponse{ConnectionStatus: &connStatusConnected}, nil)
		mockPostProcessConn.EXPECT().PostProcessAwsConnection(mock.Anything).Times(1).
			Return(nil, nil)

		diags = pr.clumioPostProcessAWSConnectionCommon(ctx, prm, eventType)
		assert.Nil(t, diags)
	})

	// Tests the success scenario for post-process aws connection common function.
	// It should not return Diagnostics.
	t.Run("Success scenario for post-process common with waitForIngestion"+
		" and waitForTargetStatus", func(t *testing.T) {

		ingestionStatus := "completed"
		targetSetupStatus := "completed"
		prm.WaitForIngestion = basetypes.NewBoolValue(true)
		prm.WaitForDataPlaneResources = basetypes.NewBoolValue(true)
		//Setup expectations.
		mockPostProcessConn.EXPECT().PostProcessAwsConnection(mock.Anything).Times(1).
			Return(nil, nil)
		mockAWSConn.EXPECT().ReadAwsConnection(mock.Anything, mock.Anything).Times(1).
			Return(&models.ReadAWSConnectionResponse{
				AccountNativeId:   &accountId,
				AwsRegion:         &region,
				IngestionStatus:   &ingestionStatus,
				TargetSetupStatus: &targetSetupStatus,
			}, nil)

		diags = pr.clumioPostProcessAWSConnectionCommon(ctx, prm, eventType)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case ingestionStatus is true and the value returned for
	// ingestionStatus is failed.
	t.Run("Error scenario for post-process common with waitForIngestion failure",
		func(t *testing.T) {

			ingestionStatus := "failed"
			targetSetupStatus := "completed"
			prm.WaitForIngestion = basetypes.NewBoolValue(true)
			prm.WaitForDataPlaneResources = basetypes.NewBoolValue(true)
			//Setup expectations.
			mockPostProcessConn.EXPECT().PostProcessAwsConnection(mock.Anything).Times(1).
				Return(nil, nil)
			mockAWSConn.EXPECT().ReadAwsConnection(mock.Anything, mock.Anything).Times(1).
				Return(&models.ReadAWSConnectionResponse{
					AccountNativeId:   &accountId,
					AwsRegion:         &region,
					IngestionStatus:   &ingestionStatus,
					TargetSetupStatus: &targetSetupStatus,
				}, nil)

			diags = pr.clumioPostProcessAWSConnectionCommon(ctx, prm, eventType)
			assert.NotNil(t, diags)
		})

	// Tests that Diagnostics is returned in case waitForDataPlaneResources is true and the value
	// returned for targetSetupStatus is failed.
	t.Run("Error scenario for post-process common with waitForDataPlaneResources failure",
		func(t *testing.T) {

			ingestionStatus := "completed"
			targetSetupStatus := "failed"
			prm.WaitForIngestion = basetypes.NewBoolValue(true)
			prm.WaitForDataPlaneResources = basetypes.NewBoolValue(true)
			//Setup expectations.
			mockPostProcessConn.EXPECT().PostProcessAwsConnection(mock.Anything).Times(1).
				Return(nil, nil)
			mockAWSConn.EXPECT().ReadAwsConnection(mock.Anything, mock.Anything).Times(1).
				Return(&models.ReadAWSConnectionResponse{
					AccountNativeId:   &accountId,
					AwsRegion:         &region,
					IngestionStatus:   &ingestionStatus,
					TargetSetupStatus: &targetSetupStatus,
				}, nil)

			diags = pr.clumioPostProcessAWSConnectionCommon(ctx, prm, eventType)
			assert.NotNil(t, diags)
		})

	// Tests that Diagnostics is returned in case getting the template configuration returns an
	// error.
	t.Run("GetTemplateConfiguration returns an error", func(t *testing.T) {

		prmWithInvalidVersion := postProcessAWSConnectionResourceModel{
			ConfigVersion: basetypes.NewStringValue(invalidVersion),
		}

		diags = pr.clumioPostProcessAWSConnectionCommon(ctx, prmWithInvalidVersion, eventType)
		assert.NotNil(t, diags)
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
