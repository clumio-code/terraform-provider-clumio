// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_protection_group

import (
	"context"
	"testing"
	"time"

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
	name             = "test-protection-group"
	resourceName     = "test_protection_group"
	id               = "mock-pg-id"
	ou               = "mock-ou"
	testError        = "Test Error"
	description      = "test-description"
	bucketRule       = "test-bucket-rule"
	prefix           = "test-prefix"
	exclPrefix       = "test-excluded-prefix-1"
	exclPrefix2      = "test-excluded-prefix-2"
	storageClass     = "test-storage-class-1"
	storageClass2    = "test-storage-class-2"
	protectionStatus = "protection-status"
	entityId         = "test-entity-id"
	entityType       = "test-entity-type"
	policyId         = "test-policy-id"
)

// Unit test for the following cases:
//   - Create protection group success scenario.
//   - SDK API for create protection group returns an error.
//   - SDK API for create protection group returns an empty response.
//   - Polling of read protection group returns an error.
//   - Polling of read protection group task returns an empty response.
func TestCreateProtectionGroup(t *testing.T) {

	mockProtectionGroup := sdkclients.NewMockProtectionGroupClient(t)
	ctx := context.Background()
	pr := clumioProtectionGroupResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkProtectionGroups: mockProtectionGroup,
		pollTimeout:         5 * time.Second,
		pollInterval:        1,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Populate the protection group resource model to be used as input to createProtectionGroup()
	pgrm := clumioProtectionGroupResourceModel{
		Name:        basetypes.NewStringValue(name),
		Description: basetypes.NewStringValue(description),
		BucketRule:  basetypes.NewStringValue(bucketRule),
		ObjectFilter: []*objectFilterModel{
			{
				LatestVersionOnly: basetypes.NewBoolValue(true),
				PrefixFilters: []*prefixFilterModel{
					{
						ExcludedSubPrefixes: []types.String{
							basetypes.NewStringValue(exclPrefix),
							basetypes.NewStringValue(exclPrefix2),
						},
						Prefix: basetypes.NewStringValue(prefix),
					},
				},
				StorageClasses: []types.String{
					basetypes.NewStringValue(storageClass),
					basetypes.NewStringValue(storageClass2),
				},
			},
		},
	}

	// Create the response of the SDK CreateProtectionGroupDefinition() API.
	createResponse := &models.CreateProtectionGroupResponse{
		Id: &id,
	}

	// Tests the success scenario for protection group create. It should not return Diagnostics.
	t.Run("Basic success scenario for create protection group", func(t *testing.T) {

		// Create the response of the SDK ReadProtectionGroupDefinition() API.
		readResponse := &models.ReadProtectionGroupResponse{
			BucketRule:  &bucketRule,
			Description: &description,
			Id:          &id,
			Name:        &name,
			ObjectFilter: &models.ObjectFilter{
				PrefixFilters: []*models.PrefixFilter{
					{
						ExcludedSubPrefixes: []*string{
							&exclPrefix, &exclPrefix2,
						},
						Prefix: &prefix,
					},
				},
				StorageClasses: []*string{
					&storageClass, &storageClass2,
				},
			},
			OrganizationalUnitId: &ou,
			ProtectionStatus:     &protectionStatus,
		}

		// Setup Expectations
		mockProtectionGroup.EXPECT().CreateProtectionGroup(mock.Anything).Times(1).
			Return(createResponse, nil)
		mockProtectionGroup.EXPECT().ReadProtectionGroup(mock.Anything, mock.Anything).Times(1).
			Return(readResponse, nil)

		diags := pr.createProtectionGroup(ctx, &pgrm)
		assert.Nil(t, diags)
		assert.Equal(t, pgrm.ID.ValueString(), *createResponse.Id)
		assert.Equal(t, pgrm.Name.ValueString(), *readResponse.Name)
		assert.Equal(t, pgrm.OrganizationalUnitID.ValueString(),
			*readResponse.OrganizationalUnitId)
		assert.Equal(t, pgrm.ProtectionStatus.ValueString(), *readResponse.ProtectionStatus)
	})

	// Tests that Diagnostics is returned in case the create protection group API call returns
	// error.
	t.Run("CreateProtectionGroup returns error", func(t *testing.T) {
		pgrm.OrganizationalUnitID = basetypes.NewStringNull()
		// Setup Expectations
		mockProtectionGroup.EXPECT().CreateProtectionGroup(mock.Anything).Times(1).
			Return(nil, apiError)

		diags := pr.createProtectionGroup(context.Background(), &pgrm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the create protection group API call returns an
	// empty response.
	t.Run("CreateProtectionGroup returns an empty response", func(t *testing.T) {
		pgrm.OrganizationalUnitID = basetypes.NewStringNull()
		// Setup Expectations
		mockProtectionGroup.EXPECT().CreateProtectionGroup(mock.Anything).Times(1).
			Return(nil, nil)

		diags := pr.createProtectionGroup(context.Background(), &pgrm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read protection group API call returns an
	// error.
	t.Run("ReadProtectionGroup after create returns an error", func(t *testing.T) {
		pgrm.OrganizationalUnitID = basetypes.NewStringNull()
		// Setup Expectations
		mockProtectionGroup.EXPECT().CreateProtectionGroup(mock.Anything).Times(1).
			Return(createResponse, nil)
		mockProtectionGroup.EXPECT().ReadProtectionGroup(mock.Anything, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := pr.createProtectionGroup(context.Background(), &pgrm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read protection group API call returns an
	// empty response.
	t.Run("ReadProtectionGroup after create returns an empty response", func(t *testing.T) {
		pgrm.OrganizationalUnitID = basetypes.NewStringNull()
		// Setup Expectations
		mockProtectionGroup.EXPECT().CreateProtectionGroup(mock.Anything).Times(1).
			Return(createResponse, nil)
		mockProtectionGroup.EXPECT().ReadProtectionGroup(mock.Anything, mock.Anything).Times(1).
			Return(nil, nil)

		diags := pr.createProtectionGroup(context.Background(), &pgrm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Read protection group success scenario.
//   - SDK API for read protection group returns not found error.
//   - SDK API for read protection group returns an error.
//   - SDK API for read protection group returns an empty response.
func TestReadProtectionGroup(t *testing.T) {

	mockProtectionGroup := sdkclients.NewMockProtectionGroupClient(t)
	ctx := context.Background()
	pr := clumioProtectionGroupResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkProtectionGroups: mockProtectionGroup,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	// Populate the protection group resource model to be used as input to readProtectionGroup()
	pgrm := clumioProtectionGroupResourceModel{
		ID:          basetypes.NewStringValue(id),
		Name:        basetypes.NewStringValue(name),
		Description: basetypes.NewStringValue(description),
		BucketRule:  basetypes.NewStringValue(bucketRule),
		ObjectFilter: []*objectFilterModel{
			{
				LatestVersionOnly: basetypes.NewBoolValue(true),
				PrefixFilters: []*prefixFilterModel{
					{
						ExcludedSubPrefixes: []types.String{
							basetypes.NewStringValue(exclPrefix),
							basetypes.NewStringValue(exclPrefix2),
						},
						Prefix: basetypes.NewStringValue(prefix),
					},
				},
				StorageClasses: []types.String{
					basetypes.NewStringValue(storageClass),
					basetypes.NewStringValue(storageClass2),
				},
			},
		},
	}
	// Tests the success scenario for protection group read. It should not return Diagnostics.
	t.Run("Basic success scenario for read protection group", func(t *testing.T) {

		// Create the response of the SDK ReadProtectionGroupDefinition() API.
		readResponse := &models.ReadProtectionGroupResponse{
			BucketRule:  &bucketRule,
			Description: &description,
			Id:          &id,
			Name:        &name,
			ObjectFilter: &models.ObjectFilter{
				PrefixFilters: []*models.PrefixFilter{
					{
						ExcludedSubPrefixes: []*string{
							&exclPrefix, &exclPrefix2,
						},
						Prefix: &prefix,
					},
				},
				StorageClasses: []*string{
					&storageClass, &storageClass2,
				},
			},
			OrganizationalUnitId: &ou,
			ProtectionStatus:     &protectionStatus,
		}

		// Setup Expectations
		mockProtectionGroup.EXPECT().ReadProtectionGroup(id, mock.Anything).Times(1).
			Return(readResponse, nil)

		remove, diags := pr.readProtectionGroup(ctx, &pgrm)
		assert.False(t, remove)
		assert.Nil(t, diags)
		assert.Equal(t, pgrm.ID.ValueString(), *readResponse.Id)
		assert.Equal(t, pgrm.Name.ValueString(), *readResponse.Name)
		assert.Equal(t, pgrm.OrganizationalUnitID.ValueString(),
			*readResponse.OrganizationalUnitId)
		assert.Equal(t, pgrm.ProtectionStatus.ValueString(), *readResponse.ProtectionStatus)
	})

	// Tests that Diagnostics is returned in case the read protection group API call returns HTTP
	// 404 error.
	t.Run("ReadProtectionGroup returns http 404 error", func(t *testing.T) {
		pgrm.OrganizationalUnitID = basetypes.NewStringNull()
		// Setup Expectations
		mockProtectionGroup.EXPECT().ReadProtectionGroup(mock.Anything, mock.Anything).Times(1).
			Return(nil, apiNotFoundError)

		remove, diags := pr.readProtectionGroup(context.Background(), &pgrm)
		assert.Nil(t, diags)
		assert.True(t, remove)
	})

	// Tests that Diagnostics is returned in case the read protection group API call returns an
	// error.
	t.Run("ReadProtectionGroup returns an error", func(t *testing.T) {
		pgrm.OrganizationalUnitID = basetypes.NewStringNull()
		// Setup Expectations
		mockProtectionGroup.EXPECT().ReadProtectionGroup(mock.Anything, mock.Anything).Times(1).
			Return(nil, apiError)

		remove, diags := pr.readProtectionGroup(context.Background(), &pgrm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that Diagnostics is returned in case the read protection group API call returns an
	// empty response.
	t.Run("ReadProtectionGroup returns nil response", func(t *testing.T) {
		pgrm.OrganizationalUnitID = basetypes.NewStringNull()
		// Setup Expectations
		mockProtectionGroup.EXPECT().ReadProtectionGroup(mock.Anything, mock.Anything).Times(1).
			Return(nil, nil)

		remove, diags := pr.readProtectionGroup(context.Background(), &pgrm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})
}

// Unit test for the following cases:
//   - Update protection group success scenario.
//   - SDK API for read protection group to get version returns error.
//   - SDK API for update protection group returns error.
//   - SDK API for update protection group returns nil response.
//   - Polling of read protection group returns an error.
//   - Polling of read protection group returns an empty response.
func TestUpdateProtectionGroup(t *testing.T) {
	mockProtectionGroup := sdkclients.NewMockProtectionGroupClient(t)
	ctx := context.Background()
	pr := clumioProtectionGroupResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkProtectionGroups: mockProtectionGroup,
		pollTimeout:         5 * time.Second,
		pollInterval:        1,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Populate the protection group resource model to be used as input to createProtectionGroup()
	pgrm := clumioProtectionGroupResourceModel{
		ID:          basetypes.NewStringValue(id),
		Name:        basetypes.NewStringValue(name),
		Description: basetypes.NewStringValue(description),
		BucketRule:  basetypes.NewStringValue(bucketRule),
		ObjectFilter: []*objectFilterModel{
			{
				LatestVersionOnly: basetypes.NewBoolValue(true),
				PrefixFilters: []*prefixFilterModel{
					{
						ExcludedSubPrefixes: []types.String{
							basetypes.NewStringValue(exclPrefix),
							basetypes.NewStringValue(exclPrefix2),
						},
						Prefix: basetypes.NewStringValue(prefix),
					},
				},
				StorageClasses: []types.String{
					basetypes.NewStringValue(storageClass),
					basetypes.NewStringValue(storageClass2),
				},
			},
		},
	}

	oldVersion := int64(1)
	newVersion := int64(2)
	firstResponse := &models.ReadProtectionGroupResponse{
		Name:        &name,
		Description: &description,
		Version:     &oldVersion,
	}

	// Create the response of the SDK UpdateProtectionGroupDefinition() API.
	updateResponse := &models.UpdateProtectionGroupResponse{Id: &id}

	// Tests the success scenario for protection group update. It should not return Diagnostics.
	t.Run("Basic success scenario for update protection group", func(t *testing.T) {

		latestVersionOnly := true
		// Create the response of the SDK ReadProtectionGroupDefinition() API.
		readResponse := &models.ReadProtectionGroupResponse{
			BucketRule:  &bucketRule,
			Description: &description,
			Id:          &id,
			Name:        &name,
			ObjectFilter: &models.ObjectFilter{
				LatestVersionOnly: &latestVersionOnly,
				PrefixFilters: []*models.PrefixFilter{
					{
						ExcludedSubPrefixes: []*string{
							&exclPrefix, &exclPrefix2,
						},
						Prefix: &prefix,
					},
				},
				StorageClasses: []*string{
					&storageClass, &storageClass2,
				},
			},
			OrganizationalUnitId: &ou,
			ProtectionStatus:     &protectionStatus,
			Version:              &newVersion,
		}

		// Setup Expectations
		mockProtectionGroup.EXPECT().ReadProtectionGroup(mock.Anything, mock.Anything).Times(1).
			Return(firstResponse, nil)
		mockProtectionGroup.EXPECT().UpdateProtectionGroup(id, mock.Anything).Times(1).
			Return(updateResponse, nil)
		mockProtectionGroup.EXPECT().ReadProtectionGroup(mock.Anything, mock.Anything).Times(1).
			Return(readResponse, nil)

		diags := pr.updateProtectionGroup(ctx, &pgrm)
		assert.Nil(t, diags)
		assert.Equal(t, pgrm.ID.ValueString(), *updateResponse.Id)
		assert.Equal(t, pgrm.Name.ValueString(), *readResponse.Name)
		assert.Equal(t, pgrm.OrganizationalUnitID.ValueString(),
			*readResponse.OrganizationalUnitId)
		assert.Equal(t, pgrm.ProtectionStatus.ValueString(), *readResponse.ProtectionStatus)
	})

	// Tests that Diagnostics is returned in case the first read protection group API call returns
	// error.
	t.Run("UpdateProtectionGroup returns error", func(t *testing.T) {
		pgrm.OrganizationalUnitID = basetypes.NewStringNull()
		// Setup Expectations
		mockProtectionGroup.EXPECT().ReadProtectionGroup(mock.Anything, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := pr.updateProtectionGroup(context.Background(), &pgrm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the update protection group API call returns
	// error.
	t.Run("UpdateProtectionGroup returns error", func(t *testing.T) {
		pgrm.OrganizationalUnitID = basetypes.NewStringNull()
		// Setup Expectations
		mockProtectionGroup.EXPECT().ReadProtectionGroup(mock.Anything, mock.Anything).Times(1).
			Return(firstResponse, nil)
		mockProtectionGroup.EXPECT().UpdateProtectionGroup(id, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := pr.updateProtectionGroup(context.Background(), &pgrm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the update protection group API call returns an
	// empty response.
	t.Run("UpdateProtectionGroup returns nil response", func(t *testing.T) {
		pgrm.OrganizationalUnitID = basetypes.NewStringNull()
		// Setup Expectations
		mockProtectionGroup.EXPECT().ReadProtectionGroup(mock.Anything, mock.Anything).Times(1).
			Return(firstResponse, nil)
		mockProtectionGroup.EXPECT().UpdateProtectionGroup(id, mock.Anything).Times(1).
			Return(nil, nil)

		diags := pr.updateProtectionGroup(context.Background(), &pgrm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read protection group API call returns an
	// error.
	t.Run("ReadProtectionGroup after update returns an error", func(t *testing.T) {
		pgrm.OrganizationalUnitID = basetypes.NewStringNull()
		// Setup Expectations
		mockProtectionGroup.EXPECT().ReadProtectionGroup(mock.Anything, mock.Anything).Times(1).
			Return(firstResponse, nil)
		mockProtectionGroup.EXPECT().UpdateProtectionGroup(id, mock.Anything).Times(1).
			Return(updateResponse, nil)
		mockProtectionGroup.EXPECT().ReadProtectionGroup(mock.Anything, mock.Anything).Times(1).
			Return(nil, apiError)

		diags := pr.updateProtectionGroup(context.Background(), &pgrm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read protection group API call returns an
	// empty response.
	t.Run("ReadProtectionGroup after update returns an empty response", func(t *testing.T) {
		pgrm.OrganizationalUnitID = basetypes.NewStringNull()
		// Setup Expectations
		mockProtectionGroup.EXPECT().ReadProtectionGroup(mock.Anything, mock.Anything).Times(1).
			Return(firstResponse, nil)
		mockProtectionGroup.EXPECT().UpdateProtectionGroup(id, mock.Anything).Times(1).
			Return(updateResponse, nil)
		mockProtectionGroup.EXPECT().ReadProtectionGroup(mock.Anything, mock.Anything).Times(1).
			Return(nil, nil)

		diags := pr.updateProtectionGroup(context.Background(), &pgrm)
		t.Log(diags)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Delete protection group success scenario.
//   - Delete protection group should not return an error if protection group is not found.
//   - SDK API for delete protection group returns an error.
func TestDeleteProtectionGroup(t *testing.T) {

	mockProtectionGroup := sdkclients.NewMockProtectionGroupClient(t)
	ctx := context.Background()
	pr := clumioProtectionGroupResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkProtectionGroups: mockProtectionGroup,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	// Populate the protection group resource model to be used as input to readProtectionGroup()
	pgrm := &clumioProtectionGroupResourceModel{
		ID:          basetypes.NewStringValue(id),
		Name:        basetypes.NewStringValue(name),
		Description: basetypes.NewStringValue(description),
		BucketRule:  basetypes.NewStringValue(bucketRule),
		ObjectFilter: []*objectFilterModel{
			{
				LatestVersionOnly: basetypes.NewBoolValue(true),
				PrefixFilters: []*prefixFilterModel{
					{
						ExcludedSubPrefixes: []types.String{
							basetypes.NewStringValue(exclPrefix),
							basetypes.NewStringValue(exclPrefix2),
						},
						Prefix: basetypes.NewStringValue(prefix),
					},
				},
				StorageClasses: []types.String{
					basetypes.NewStringValue(storageClass),
					basetypes.NewStringValue(storageClass2),
				},
			},
		},
	}

	// Tests the success scenario for protection group deletion. It should not return
	// diag.Diagnostics.
	t.Run("Success scenario for protection group deletion", func(t *testing.T) {
		// Setup Expectations
		mockProtectionGroup.EXPECT().DeleteProtectionGroup(id).Times(1).Return(mock.Anything, nil)

		diags := pr.deleteProtectionGroup(ctx, pgrm)
		assert.Nil(t, diags)
	})

	// Tests that no error is returned if the protection group does not exist.
	t.Run("Policy not found should not return error", func(t *testing.T) {
		// Setup Expectations
		mockProtectionGroup.EXPECT().DeleteProtectionGroup(id).Times(1).Return(
			nil, apiNotFoundError)

		diags := pr.deleteProtectionGroup(ctx, pgrm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned when delete protection group API call returns error.
	t.Run("deleteProtectionGroup returns an error", func(t *testing.T) {
		// Setup Expectations
		mockProtectionGroup.EXPECT().DeleteProtectionGroup(id).Times(1).Return(nil, apiError)

		diags := pr.deleteProtectionGroup(ctx, pgrm)
		assert.NotNil(t, diags)
	})

}
