// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_aws_manual_connection

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
	accountId    = "test-aws-account"
	region       = "test-region"
	resourceName = "test_aws_connection"
	id           = fmt.Sprintf("%s_%s", accountId, region)
	someArn      = "test-arn"

	testError = "Test Error"
)

// Unit test for the following cases:
//   - Create AWS manual connection success scenario.
//   - SDK API for update AWS connection returns error.
func TestCreateAWSManualConnection(t *testing.T) {

	mockAwsConnClient := sdkclients.NewMockAWSConnectionClient(t)
	ctx := context.Background()
	cr := clumioAWSManualConnectionResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkConnections: mockAwsConnClient,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Populate the protection group resource model to be used as input to createProtectionGroup()
	crm := clumioAWSManualConnectionResourceModel{
		AccountId: basetypes.NewStringValue(accountId),
		AwsRegion: basetypes.NewStringValue(region),
		AssetsEnabled: &AssetsEnabledModel{
			EBS:      basetypes.NewBoolValue(true),
			RDS:      basetypes.NewBoolValue(true),
			DynamoDB: basetypes.NewBoolValue(true),
			S3:       basetypes.NewBoolValue(true),
			EC2MSSQL: basetypes.NewBoolValue(true),
		},
		Resources: &ResourcesModel{
			ClumioIAMRoleArn:     basetypes.NewStringValue(someArn),
			ClumioSupportRoleArn: basetypes.NewStringValue(someArn),
			ClumioEventPubArn:    basetypes.NewStringValue(someArn),
			EventRules: &EventRules{
				CloudtrailRuleArn: basetypes.NewStringValue(someArn),
				CloudwatchRuleArn: basetypes.NewStringValue(someArn),
			},
			ServiceRoles: &ServiceRoles{
				Mssql: &MssqlServiceRoles{
					SsmNotificationRoleArn:   basetypes.NewStringValue(someArn),
					Ec2SsmInstanceProfileArn: basetypes.NewStringValue(someArn),
				},
				S3: &S3ServiceRoles{
					ContinuousBackupsRoleArn: basetypes.NewStringValue(someArn),
				},
			},
		},
	}

	// Tests the success scenario for clumio aws manual connection create. It should not return
	// Diagnostics.
	t.Run("Basic success scenario for create aws connection", func(t *testing.T) {

		// Setup Expectations
		mockAwsConnClient.EXPECT().UpdateAwsConnection(id, mock.Anything).Times(1).
			Return(nil, nil)

		diags := cr.createAWSManualConnection(ctx, &crm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the update aws connection API call returns an
	// error.
	t.Run("CreateAWSManualConnection returns an error", func(t *testing.T) {

		// Setup Expectations
		mockAwsConnClient.EXPECT().UpdateAwsConnection(id, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := cr.createAWSManualConnection(ctx, &crm)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Update AWS manual connection success scenario.
//   - SDK API for update Clumio AWS connection returns an error.
//   - AssetEnabled downgrade returns an error.
func TestUpdateAWSManualConnection(t *testing.T) {

	mockAwsConnClient := sdkclients.NewMockAWSConnectionClient(t)
	ctx := context.Background()
	cr := clumioAWSManualConnectionResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkConnections: mockAwsConnClient,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Populate the protection group resource model to be used as input to createProtectionGroup()
	plan := clumioAWSManualConnectionResourceModel{
		ID:        basetypes.NewStringValue(id),
		AccountId: basetypes.NewStringValue(accountId),
		AwsRegion: basetypes.NewStringValue(region),
		AssetsEnabled: &AssetsEnabledModel{
			EBS:      basetypes.NewBoolValue(true),
			RDS:      basetypes.NewBoolValue(true),
			DynamoDB: basetypes.NewBoolValue(true),
			S3:       basetypes.NewBoolValue(true),
			EC2MSSQL: basetypes.NewBoolValue(true),
		},
		Resources: &ResourcesModel{
			ClumioIAMRoleArn:     basetypes.NewStringValue(someArn),
			ClumioSupportRoleArn: basetypes.NewStringValue(someArn),
			ClumioEventPubArn:    basetypes.NewStringValue(someArn),
			EventRules: &EventRules{
				CloudtrailRuleArn: basetypes.NewStringValue(someArn),
				CloudwatchRuleArn: basetypes.NewStringValue(someArn),
			},
			ServiceRoles: &ServiceRoles{
				Mssql: &MssqlServiceRoles{
					SsmNotificationRoleArn:   basetypes.NewStringValue(someArn),
					Ec2SsmInstanceProfileArn: basetypes.NewStringValue(someArn),
				},
				S3: &S3ServiceRoles{
					ContinuousBackupsRoleArn: basetypes.NewStringValue(someArn),
				},
			},
		},
	}

	// Populate the protection group resource model to be used as input to createProtectionGroup()
	state := clumioAWSManualConnectionResourceModel{
		ID: basetypes.NewStringValue(id),
		AssetsEnabled: &AssetsEnabledModel{
			EBS:      basetypes.NewBoolValue(true),
			RDS:      basetypes.NewBoolValue(true),
			DynamoDB: basetypes.NewBoolValue(true),
			S3:       basetypes.NewBoolValue(true),
			EC2MSSQL: basetypes.NewBoolValue(true),
		},
	}

	// Tests the success scenario for AWS connection update. It should not return Diagnostics.
	t.Run("success scenario for update aws manual connection", func(t *testing.T) {

		mockAwsConnClient.EXPECT().UpdateAwsConnection(id, mock.Anything).Times(1).
			Return(nil, nil)

		diags := cr.updateAWSManualConnection(ctx, &plan, &state)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the update AWS connection API call returns an
	// error.
	t.Run("update aws connection returns an error", func(t *testing.T) {

		// Setup Expectations
		mockAwsConnClient.EXPECT().UpdateAwsConnection(id, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := cr.updateAWSManualConnection(ctx, &plan, &state)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case some enabled asset is removed while updating.
	t.Run("Downgrading assets enabled returns an error", func(t *testing.T) {

		// Setup Expectations
		mockAwsConnClient.EXPECT().UpdateAwsConnection(id, mock.Anything).Times(1).
			Return(nil, apiError)

		plan.AssetsEnabled.EBS = basetypes.NewBoolValue(false)
		diags := cr.updateAWSManualConnection(ctx, &plan, &state)
		assert.NotNil(t, diags)
	})
}
