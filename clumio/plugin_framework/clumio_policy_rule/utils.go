// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_policy_rule Terraform resource.

package clumio_policy_rule

// clearOUContext resets the OrganizationalUnitContext in the client.
func (r *policyRuleResource) clearOUContext() {

	r.client.ClumioConfig.OrganizationalUnitContext = ""
}
