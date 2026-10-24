// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_aws_manual_connection_resources

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
//   - Read AWS manual connection resources success scenario.
//   - SDK API for read AWS manual connection resources returns an error.
//   - SDK API for read AWS manual connection resources returns an empty response.
func TestReadAWSManualConnectionResources(t *testing.T) {

	ctx := context.Background()
	templatesClient := sdkclients.NewMockAWSTemplatesClient(t)
	resourceName := "test_resources"
	accountId := "test-aws-account"
	region := "test-region"

	testError := "Test Error"

	rds := clumioAwsManualConnectionResourcesDatasource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		awsTemplates: templatesClient,
	}

	rdsm := &clumioAwsManualConnectionResourcesModel{
		AccountId: basetypes.NewStringValue(accountId),
		AwsRegion: basetypes.NewStringValue(region),
		AssetsEnabled: &assetTypesEnabledModel{
			EBS:      basetypes.NewBoolValue(true),
			RDS:      basetypes.NewBoolValue(true),
			DynamoDB: basetypes.NewBoolValue(true),
			S3:       basetypes.NewBoolValue(true),
			EC2MSSQL: basetypes.NewBoolValue(true),
		},
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for role read. It should not return Diagnostics.
	t.Run("Basic success scenario for read AWS manual connection resources",
		func(t *testing.T) {

			readResponse := &models.CreateAWSTemplateV2Response{
				Resources: &models.CategorisedResources{},
			}
			// Setup expectations.
			templatesClient.EXPECT().CreateConnectionTemplate(mock.Anything).Times(1).Return(
				readResponse, nil)

			diags := rds.readAWSManualConnectionResources(ctx, rdsm)
			assert.Nil(t, diags)
		})

	// Tests that Diagnostics is returned in case the create connection template API call returns an
	// error.
	t.Run("CreateConnectionTemplate returns an error", func(t *testing.T) {

		// Setup expectations.
		templatesClient.EXPECT().CreateConnectionTemplate(mock.Anything).Times(1).Return(
			nil, apiError)

		diags := rds.readAWSManualConnectionResources(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the create connection template API call returns an
	// empty response.
	t.Run("CreateConnectionTemplate returns an empty response", func(t *testing.T) {

		// Setup expectations.
		templatesClient.EXPECT().CreateConnectionTemplate(mock.Anything).Times(1).Return(nil, nil)

		diags := rds.readAWSManualConnectionResources(ctx, rdsm)
		assert.NotNil(t, diags)
	})
}
