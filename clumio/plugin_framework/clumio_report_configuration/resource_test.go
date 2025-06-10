// Copyright 2025. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_report_configuration

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
	resourceName                      = "test_report_configuration"
	testError                         = "Test Error"
	testId                            = "test-id"
	testCreatedAt                     = "2025-01-01T00:00:00Z"
	testEmail1                        = "email1"
	testEmail2                        = "email2"
	testStartTime                     = "2025-01-01T00:00:00Z"
	testTimeUnit                      = "days"
	testTimeValue                     = int32(7)
	testShouldIgnoreDeactivatedPolicy = true
	testRegion                        = "us-west-2"
	testTagOpMode                     = "equal"
	testTagKey                        = "test_tag_key"
	testTagValue                      = "test_tag_value"
	testAssetType1                    = "asset_type_1"
	testAssetType2                    = "asset_type_2"
	testDataSource1                   = "data_source_1"
	testDataSource2                   = "data_source_2"
	testOrganizationalUnit1           = "organizational_unit_1"
	testOrganizationalUnit2           = "organizational_unit_2"
	testTimeUnitParam                 = &models.TimeUnitParam{
		Unit:  &testTimeUnit,
		Value: &testTimeValue,
	}
)

const (
	// The following constants are used as a test name in different tests.
	readPolicyError                  = "Read policy definition returns an error"
	readPolicyNotFoundError          = "Read policy definition returns not found error"
	setPolicyAssignmentPollingError  = "Polling for set policy assignment task returns an error"
	readProtectionGroupError         = "Read protection group returns an error"
	readProtectionGroupEmptyResponse = "Read protection group returns an empty response"
	readProtectionGroupNotFoundError = "Read protection group returns not found error"
)

// Unit test for the following cases:
//   - Create report configuration success scenario.
//   - SDK API for create report configuration returns an error.
//   - SDK API for create report configuration returns an empty response.
func TestCreateReportConfiguration(t *testing.T) {

	ctx := context.Background()
	mockReportConfigurations := sdkclients.NewMockReportConfigurationClient(t)
	rcr := &clumioReportConfigurationResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkReportConfigurations: mockReportConfigurations,
	}

	model := &reportConfigurationResourceModel{
		Name:        types.StringValue("test_name"),
		Description: types.StringValue("Test Description"),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for create report configuration. It should not return Diagnostics.
	t.Run("Basic success scenario for create report configuration", func(t *testing.T) {

		resp := &models.CreateComplianceConfigurationResponse{
			Id:      &testId,
			Created: &testCreatedAt,
		}

		mockReportConfigurations.EXPECT().CreateComplianceReportConfiguration(mock.Anything).
			Times(1).Return(resp, nil)

		diags := rcr.createReportConfiguration(ctx, model)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the create report configuration API call returns
	// an error.
	t.Run("Create report configuration returns an error", func(t *testing.T) {

		mockReportConfigurations.EXPECT().CreateComplianceReportConfiguration(mock.Anything).
			Times(1).Return(nil, apiError)

		diags := rcr.createReportConfiguration(ctx, model)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the create report configuration API call returns
	// an empty response.
	t.Run("Create report configuration returns an empty response", func(t *testing.T) {

		mockReportConfigurations.EXPECT().CreateComplianceReportConfiguration(mock.Anything).
			Times(1).Return(nil, nil)

		diags := rcr.createReportConfiguration(ctx, model)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Read report configuration success scenario.
//   - SDK API for read report configuration returns an error.
//   - SDK API for read report configuration returns not found error.
//   - SDK API for read report configuration returns an empty response.
func TestReadReportConfiguration(t *testing.T) {

	ctx := context.Background()
	mockReportConfigurations := sdkclients.NewMockReportConfigurationClient(t)
	rcr := &clumioReportConfigurationResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkReportConfigurations: mockReportConfigurations,
	}

	model := &reportConfigurationResourceModel{
		ID:          types.StringValue(testId),
		Description: types.StringValue("Test Description"),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	// Tests the success scenario for read report configuration. It should not return Diagnostics.
	t.Run("Basic success scenario for read report configuration", func(t *testing.T) {

		resp := &models.ReadComplianceConfigurationResponse{
			Id:          &testId,
			Name:        model.Name.ValueStringPointer(),
			Description: model.Description.ValueStringPointer(),
			Created:     &testCreatedAt,
			Notification: &models.NotificationSetting{
				EmailList: []*string{&testEmail1, &testEmail2},
			},
			Parameter: &models.Parameter{
				Controls: &models.ComplianceControls{
					AssetBackup: &models.AssetBackupControl{
						LookBackPeriod:           testTimeUnitParam,
						MinimumRetentionDuration: testTimeUnitParam,
						WindowSize:               testTimeUnitParam,
					},
					AssetProtection: &models.AssetProtectionControl{
						ShouldIgnoreDeactivatedPolicy: &testShouldIgnoreDeactivatedPolicy,
					},
					Policy: &models.PolicyControl{
						MinimumRetentionDuration: testTimeUnitParam,
						MinimumRpoFrequency:      testTimeUnitParam,
					},
				},
				Filters: &models.ComplianceFilters{
					Asset: &models.AssetFilter{
						Groups: []*models.AssetGroupFilter{
							{
								Id:     &testId,
								Region: &testRegion,
							},
						},
						TagOpMode: &testTagOpMode,
						Tags: []*models.Tag{
							{
								Key:   &testTagKey,
								Value: &testTagValue,
							},
						},
					},
					Common: &models.CommonFilter{
						AssetTypes:  []*string{&testAssetType1, &testAssetType2},
						DataSources: []*string{&testDataSource1, &testDataSource2},
						OrganizationalUnits: []*string{
							&testOrganizationalUnit1, &testOrganizationalUnit2,
						},
					},
				},
			},
			Schedule: &models.ScheduleSetting{
				StartTime: &testStartTime,
			},
		}

		mockReportConfigurations.EXPECT().ReadComplianceReportConfiguration(testId).Times(1).
			Return(resp, nil)

		remove, diags := rcr.readReportConfiguration(ctx, model)
		assert.Nil(t, diags)
		assert.False(t, remove)
	})

	// Tests that Diagnostics is returned in case the read report configuration API call returns
	// an error.
	t.Run("Read report configuration returns an error", func(t *testing.T) {

		mockReportConfigurations.EXPECT().ReadComplianceReportConfiguration(testId).Times(1).
			Return(nil, apiError)

		remove, diags := rcr.readReportConfiguration(ctx, model)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that read report configuration returns true to indicate that the resource should be
	// removed when read report configuration API call returns not found error.
	t.Run("Read report configuration returns not found error", func(t *testing.T) {

		mockReportConfigurations.EXPECT().ReadComplianceReportConfiguration(testId).Times(1).
			Return(nil, apiNotFoundError)

		remove, diags := rcr.readReportConfiguration(ctx, model)
		assert.Nil(t, diags)
		assert.True(t, remove)
	})

	// Tests that Diagnostics is returned in case the read report configuration API call returns an
	// empty response.
	t.Run("Read report configuration returns an empty response", func(t *testing.T) {

		mockReportConfigurations.EXPECT().ReadComplianceReportConfiguration(testId).Times(1).
			Return(nil, nil)

		remove, diags := rcr.readReportConfiguration(ctx, model)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})
}

// Unit test for the following cases:
//   - Update report configuration success scenario.
//   - SDK API for update report configuration returns an error.
//   - SDK API for update report configuration returns an empty response.
func TestUpdateReportConfiguration(t *testing.T) {

	ctx := context.Background()
	mockReportConfigurations := sdkclients.NewMockReportConfigurationClient(t)
	rcr := &clumioReportConfigurationResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkReportConfigurations: mockReportConfigurations,
	}

	model := &reportConfigurationResourceModel{
		ID: types.StringValue(testId),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for update report configuration. It should not return Diagnostics.
	t.Run("Basic success scenario for update report configuration", func(t *testing.T) {

		resp := &models.UpdateComplianceConfigurationResponse{
			Id:      &testId,
			Created: &testCreatedAt,
		}

		mockReportConfigurations.EXPECT().UpdateComplianceReportConfiguration(testId, mock.Anything).
			Times(1).Return(resp, nil)

		diags := rcr.updateReportConfiguration(ctx, model)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the update report configuration API call returns
	// an error.
	t.Run("Update report configuration returns an error", func(t *testing.T) {

		mockReportConfigurations.EXPECT().UpdateComplianceReportConfiguration(testId, mock.Anything).
			Times(1).Return(nil, apiError)

		diags := rcr.updateReportConfiguration(ctx, model)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the update report configuration API call returns
	// an empty response.
	t.Run("Update report configuration returns an empty response", func(t *testing.T) {

		mockReportConfigurations.EXPECT().UpdateComplianceReportConfiguration(testId, mock.Anything).
			Times(1).Return(nil, nil)

		diags := rcr.updateReportConfiguration(ctx, model)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Delete report configuration success scenario.
//   - SDK API for delete report configuration returns an error.
//   - SDK API for delete report configuration returns not found error.
func TestDeleteReportConfiguration(t *testing.T) {

	ctx := context.Background()
	mockReportConfigurations := sdkclients.NewMockReportConfigurationClient(t)
	rcr := &clumioReportConfigurationResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkReportConfigurations: mockReportConfigurations,
	}

	model := &reportConfigurationResourceModel{
		ID: types.StringValue(testId),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	// Tests the success scenario for delete report configuration. It should not return Diagnostics.
	t.Run("Basic success scenario for delete report configuration", func(t *testing.T) {

		mockReportConfigurations.EXPECT().DeleteComplianceReportConfiguration(testId).Times(1).
			Return(nil, nil)

		diags := rcr.deleteReportConfiguration(ctx, model)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the delete report configuration API call returns
	// an error.
	t.Run("Delete report configuration returns an error", func(t *testing.T) {

		mockReportConfigurations.EXPECT().DeleteComplianceReportConfiguration(testId).Times(1).
			Return(nil, apiError)

		diags := rcr.deleteReportConfiguration(ctx, model)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the delete report configuration API call returns
	// not found error.
	t.Run("Delete report configuration returns not found error", func(t *testing.T) {

		mockReportConfigurations.EXPECT().DeleteComplianceReportConfiguration(testId).Times(1).
			Return(nil, apiNotFoundError)

		diags := rcr.deleteReportConfiguration(ctx, model)
		assert.Nil(t, diags)
	})
}
