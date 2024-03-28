// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK AutoUserProvisioningRulesV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	sdkAUPRules "github.com/clumio-code/clumio-go-sdk/controllers/auto_user_provisioning_rules"
)

type AutoUserProvisioningRuleClient interface {
	sdkAUPRules.AutoUserProvisioningRulesV1Client
}

func NewAutoUserProvisioningRuleClient(config config.Config) AutoUserProvisioningRuleClient {
	return sdkAUPRules.NewAutoUserProvisioningRulesV1(config)
}
