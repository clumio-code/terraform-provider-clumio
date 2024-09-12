// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_policy_rule

import (
	"context"
	"testing"
	"time"

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
	name         = "test-policy-rule"
	resourceName = "test_policy_rule"
	id           = "test-rule-id"
	ou           = "test-ou"
	testError    = "Test Error"
	policyId     = "test-policy-id"
	condition    = "test-condition"
	beforeRuleId = "test-before-rule-id"
	taskid       = "test-task-id"
)

const (
	// The following constants are used as a test name in different tests.
	readTaskError = "read task returns an error"
)

// Unit test for the following cases:
//   - Create policy rule success scenario.
//   - SDK API for create policy rule returns an error.
//   - SDK API for create policy rule returns an empty response.
//   - Polling of create policy rule task returns an error.
func TestCreatePolicyRule(t *testing.T) {

	mockPolicyRule := sdkclients.NewMockPolicyRuleClient(t)
	mockTasks := sdkclients.NewMockTaskClient(t)
	ctx := context.Background()
	pr := policyRuleResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkPolicyRules: mockPolicyRule,
		sdkTasks:       mockTasks,
		pollTimeout:    5 * time.Second,
		pollInterval:   1,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	prm := &policyRuleResourceModel{
		Name:         basetypes.NewStringValue(name),
		Condition:    basetypes.NewStringValue(condition),
		BeforeRuleID: basetypes.NewStringValue(beforeRuleId),
		PolicyID:     basetypes.NewStringValue(policyId),
	}

	// Create the response of the SDK CreatePolicyRule() API.
	createResp := &models.CreateRuleResponse{
		TaskId: &taskid,
		Rule: &models.Rule{
			Id:                   &id,
			OrganizationalUnitId: &ou,
			Condition:            &condition,
		},
	}

	// Tests the success scenario for policy rule create. It should not return Diagnostics.
	t.Run("Basic success scenario for create policy rule", func(t *testing.T) {

		taskStatus := common.TaskSuccess
		readTaskResponse := &models.ReadTaskResponse{
			Status: &taskStatus,
		}

		// Setup expectations.
		mockPolicyRule.EXPECT().CreatePolicyRule(mock.Anything).Times(1).Return(
			createResp, nil)
		mockTasks.EXPECT().ReadTask(taskid).Times(1).Return(readTaskResponse, nil)

		diags := pr.createPolicyRule(ctx, prm)
		assert.Nil(t, diags)
		assert.Equal(t, id, prm.ID.ValueString())
	})

	// Tests that Diagnostics is returned in case the create policy rule API call returns an error.
	t.Run("create policy rule returns an error", func(t *testing.T) {

		// Setup expectations.
		mockPolicyRule.EXPECT().CreatePolicyRule(mock.Anything).Times(1).Return(
			nil, apiError)

		diags := pr.createPolicyRule(ctx, prm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the create policy rule API call returns an empty
	// response.
	t.Run("create policy rule returns an empty response", func(t *testing.T) {

		// Setup expectations.
		mockPolicyRule.EXPECT().CreatePolicyRule(mock.Anything).Times(1).Return(
			nil, nil)

		diags := pr.createPolicyRule(ctx, prm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read task API call returns an error.
	t.Run(readTaskError, func(t *testing.T) {

		// Setup expectations.
		mockPolicyRule.EXPECT().CreatePolicyRule(mock.Anything).Times(1).Return(
			createResp, nil)
		mockTasks.EXPECT().ReadTask(taskid).Times(1).Return(nil, apiError)

		diags := pr.createPolicyRule(ctx, prm)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Read policy rule success scenario.
//   - SDK API for read policy rule returns not found error.
//   - SDK API for read policy rule returns an error.
//   - SDK API for read policy rule returns an empty response.
func TestReadPolicyRule(t *testing.T) {

	mockPolicyRule := sdkclients.NewMockPolicyRuleClient(t)
	mockTasks := sdkclients.NewMockTaskClient(t)
	ctx := context.Background()
	pr := policyRuleResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkPolicyRules: mockPolicyRule,
		sdkTasks:       mockTasks,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	prm := &policyRuleResourceModel{
		ID:           basetypes.NewStringValue(id),
		Name:         basetypes.NewStringValue(name),
		Condition:    basetypes.NewStringValue(condition),
		BeforeRuleID: basetypes.NewStringValue(beforeRuleId),
		PolicyID:     basetypes.NewStringValue(policyId),
	}

	// Create the response of the SDK ReadPolicyRule() API.
	readResp := &models.ReadRuleResponse{
		Id:                   &id,
		OrganizationalUnitId: &ou,
		Condition:            &condition,
		Name:                 &name,
		Priority: &models.RulePriority{
			BeforeRuleId: &beforeRuleId,
		},
		Action: &models.RuleAction{
			AssignPolicy: &models.AssignPolicyAction{
				PolicyId: &policyId,
			},
		},
	}

	// Tests the success scenario for policy rule read. It should not return Diagnostics.
	t.Run("Basic success scenario for read policy rule", func(t *testing.T) {

		// Setup expectations
		mockPolicyRule.EXPECT().ReadPolicyRule(id).Times(1).Return(readResp, nil)

		remove, diags := pr.readPolicyRule(ctx, prm)
		assert.Nil(t, diags)
		assert.False(t, remove)
		assert.Equal(t, id, prm.ID.ValueString())
		assert.Equal(t, name, prm.Name.ValueString())
		assert.Equal(t, condition, prm.Condition.ValueString())
		assert.Equal(t, beforeRuleId, prm.BeforeRuleID.ValueString())
	})

	// Tests that Diagnostics is returned in case the read policy rule API call returns HTTP
	// 404 error.
	t.Run("ReadPolicyRule returns http 404 error", func(t *testing.T) {

		// Setup Expectations
		mockPolicyRule.EXPECT().ReadPolicyRule(id).Times(1).Return(nil, apiNotFoundError)

		remove, diags := pr.readPolicyRule(context.Background(), prm)
		assert.Nil(t, diags)
		assert.True(t, remove)
	})

	// Tests that Diagnostics is returned in case the read policy rule API call returns an
	// error.
	t.Run("ReadPolicyRule returns an error", func(t *testing.T) {

		// Setup Expectations
		mockPolicyRule.EXPECT().ReadPolicyRule(id).Times(1).Return(nil, apiError)

		remove, diags := pr.readPolicyRule(context.Background(), prm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})

	// Tests that Diagnostics is returned in case the read policy rule API call returns an
	// empty response.
	t.Run("ReadPolicyRule returns an empty response", func(t *testing.T) {

		// Setup Expectations
		mockPolicyRule.EXPECT().ReadPolicyRule(id).Times(1).Return(nil, nil)

		remove, diags := pr.readPolicyRule(context.Background(), prm)
		assert.NotNil(t, diags)
		assert.False(t, remove)
	})
}

// Unit test for the following cases:
//   - Update policy rule success scenario.
//   - SDK API for update policy rule returns an error.
//   - SDK API for update policy rule returns an empty response.
//   - Polling of update policy rule task returns an error.
func TestUpdatePolicyRule(t *testing.T) {

	mockPolicyRule := sdkclients.NewMockPolicyRuleClient(t)
	mockTasks := sdkclients.NewMockTaskClient(t)
	ctx := context.Background()
	pr := policyRuleResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkPolicyRules: mockPolicyRule,
		sdkTasks:       mockTasks,
		pollTimeout:    5 * time.Second,
		pollInterval:   1,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	prm := &policyRuleResourceModel{
		ID:           basetypes.NewStringValue(id),
		Name:         basetypes.NewStringValue(name),
		Condition:    basetypes.NewStringValue(condition),
		BeforeRuleID: basetypes.NewStringValue(beforeRuleId),
		PolicyID:     basetypes.NewStringValue(policyId),
	}

	// Create the response of the SDK UpdatePolicyRule() API.
	updateResp := &models.UpdateRuleResponse{
		TaskId: &taskid,
		Rule: &models.Rule{
			Id:                   &id,
			OrganizationalUnitId: &ou,
			Condition:            &condition,
		},
	}

	// Tests the success scenario for policy rule update. It should not return Diagnostics.
	t.Run("Basic success scenario for update policy rule", func(t *testing.T) {

		taskStatus := common.TaskSuccess
		readTaskResponse := &models.ReadTaskResponse{
			Status: &taskStatus,
		}

		// Setup expectations.
		mockPolicyRule.EXPECT().UpdatePolicyRule(id, mock.Anything).Times(1).Return(
			updateResp, nil)
		mockTasks.EXPECT().ReadTask(taskid).Times(1).Return(readTaskResponse, nil)

		diags := pr.updatePolicyRule(ctx, prm)
		assert.Nil(t, diags)
		assert.Equal(t, id, prm.ID.ValueString())
	})

	// Tests that Diagnostics is returned in case the update policy rule API call returns an error.
	t.Run("update policy rule returns an error", func(t *testing.T) {

		// Setup expectations.
		mockPolicyRule.EXPECT().UpdatePolicyRule(id, mock.Anything).Times(1).Return(
			nil, apiError)

		diags := pr.updatePolicyRule(ctx, prm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the update policy rule API call returns an empty
	// response.
	t.Run("update policy rule returns an empty response", func(t *testing.T) {

		// Setup expectations.
		mockPolicyRule.EXPECT().UpdatePolicyRule(id, mock.Anything).Times(1).Return(
			nil, nil)

		diags := pr.updatePolicyRule(ctx, prm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read task API call returns an error.
	t.Run(readTaskError, func(t *testing.T) {

		// Setup expectations.
		mockPolicyRule.EXPECT().UpdatePolicyRule(id, mock.Anything).Times(1).Return(
			updateResp, nil)
		mockTasks.EXPECT().ReadTask(taskid).Times(1).Return(nil, apiError)

		diags := pr.updatePolicyRule(ctx, prm)
		assert.NotNil(t, diags)
	})
}

// Unit test for the following cases:
//   - Delete policy rule success scenario.
//   - SDK API for delete policy rule returns an not found error.
//   - SDK API for delete policy rule returns an error.
//   - Polling of delete policy rule task returns an error.
func TestDeletePolicyRule(t *testing.T) {

	mockPolicyRule := sdkclients.NewMockPolicyRuleClient(t)
	mockTasks := sdkclients.NewMockTaskClient(t)
	ctx := context.Background()
	pr := policyRuleResource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkPolicyRules: mockPolicyRule,
		sdkTasks:       mockTasks,
		pollTimeout:    5 * time.Second,
		pollInterval:   1,
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	apiNotFoundError := &apiutils.APIError{
		ResponseCode: 404,
	}

	prm := &policyRuleResourceModel{
		ID:           basetypes.NewStringValue(id),
		Name:         basetypes.NewStringValue(name),
		Condition:    basetypes.NewStringValue(condition),
		BeforeRuleID: basetypes.NewStringValue(beforeRuleId),
		PolicyID:     basetypes.NewStringValue(policyId),
	}

	// Create the response of the SDK DeletePolicyRule() API.
	deleteResp := &models.DeleteRuleResponse{
		TaskId: &taskid,
	}

	// Tests the success scenario for policy rule deletion. It should not return diag.Diagnostics.
	t.Run("Success scenario for policy rule deletion", func(t *testing.T) {

		taskStatus := common.TaskSuccess
		readTaskResponse := &models.ReadTaskResponse{
			Status: &taskStatus,
		}

		// Setup Expectations
		mockPolicyRule.EXPECT().DeletePolicyRule(id).Times(1).Return(deleteResp, nil)
		mockTasks.EXPECT().ReadTask(taskid).Times(1).Return(readTaskResponse, nil)

		diags := pr.deletePolicyRule(ctx, prm)
		assert.Nil(t, diags)
	})

	// Tests that no error is returned if the policy rule does not exist.
	t.Run("Policy rule not found should not return error", func(t *testing.T) {
		// Setup Expectations
		mockPolicyRule.EXPECT().DeletePolicyRule(id).Times(1).Return(
			nil, apiNotFoundError)

		diags := pr.deletePolicyRule(ctx, prm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned when delete policy rule API call returns error.
	t.Run("delete policy rule returns an error", func(t *testing.T) {
		// Setup Expectations
		mockPolicyRule.EXPECT().DeletePolicyRule(id).Times(1).Return(nil, apiError)

		diags := pr.deletePolicyRule(ctx, prm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the read task API call returns an error.
	t.Run(readTaskError, func(t *testing.T) {

		// Setup expectations.
		mockPolicyRule.EXPECT().DeletePolicyRule(id).Times(1).Return(deleteResp, nil)
		mockTasks.EXPECT().ReadTask(taskid).Times(1).Return(nil, apiError)

		diags := pr.deletePolicyRule(ctx, prm)
		assert.NotNil(t, diags)
	})
}
