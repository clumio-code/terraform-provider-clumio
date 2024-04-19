// Copyright 2024. Clumio, Inc.

// This file contains the unit tests for the functions in resource.go

//go:build unit

package clumio_policy_rule

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
//   - Read policy rule success scenario.
//   - SDK API for read policy rule returns an error.
//   - SDK API for read policy rule returns an empty response.
func TestDatasourceReadPolicyRule(t *testing.T) {

	ctx := context.Background()
	policyRule := sdkclients.NewMockPolicyRuleClient(t)
	name := "test-policy-rule"
	resourceName := "test_policy_rule"
	id := "test-policy-rule-id"
	policyId := "test-policy-id"
	condition := "test-condition"
	testError := "Test Error"
	beforeRuleId := "test-before-rule-id"

	rds := clumioPolicyRuleDataSource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		sdkPolicyRules: policyRule,
	}

	rdsm := &clumioPolicyRuleDataSourceModel{
		Name:     basetypes.NewStringValue(name),
		PolicyId: basetypes.NewStringValue(policyId),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for policy rules read. It should not return Diagnostics.
	t.Run("Basic success scenario for read policy rules", func(t *testing.T) {

		readResponse := &models.ListRulesResponse{
			Embedded: &models.RuleListEmbedded{
				Items: []*models.Rule{
					{
						Id:   &id,
						Name: &name,
						Action: &models.RuleAction{
							AssignPolicy: &models.AssignPolicyAction{
								PolicyId: &policyId,
							},
						},
						Condition: &condition,
						Priority: &models.RulePriority{
							BeforeRuleId: &beforeRuleId,
						},
					},
				},
			},
		}

		// Setup expectations.
		policyRule.EXPECT().ListPolicyRules(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).
			Return(readResponse, nil)

		diags := rds.readPolicyRule(ctx, rdsm)
		assert.Nil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list policy rules API call returns an error.
	t.Run("list policy rules returns an error", func(t *testing.T) {

		// Setup expectations.
		policyRule.EXPECT().ListPolicyRules(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(nil, apiError)

		diags := rds.readPolicyRule(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list policy rules API call returns an empty
	// response.
	t.Run("list policy rules returns an empty response", func(t *testing.T) {

		// Setup expectations.
		policyRule.EXPECT().ListPolicyRules(mock.Anything, mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(nil, nil)

		diags := rds.readPolicyRule(ctx, rdsm)
		assert.NotNil(t, diags)
	})
}
