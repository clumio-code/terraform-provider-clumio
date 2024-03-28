// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK PolicyDefinitionsV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	sdkPolicyRules "github.com/clumio-code/clumio-go-sdk/controllers/policy_rules"
)

type PolicyRuleClient interface {
	sdkPolicyRules.PolicyRulesV1Client
}

func NewPolicyRuleClient(config config.Config) PolicyRuleClient {
	return sdkPolicyRules.NewPolicyRulesV1(config)
}
