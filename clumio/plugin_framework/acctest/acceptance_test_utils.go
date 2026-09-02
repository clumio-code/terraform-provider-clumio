// Copyright 2024. Clumio, Inc.

// This file contains the common test functions which are used by one or more acceptance tests.
// This package must only be imported from _test.go files so that it stays out of the provider binary.

package acctest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	clumioConfig "github.com/clumio-code/clumio-go-sdk/config"
	protectionGroups "github.com/clumio-code/clumio-go-sdk/controllers/protection_groups"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// idFromState returns the given id as is, or resolves it from the Terraform state when it is a
// resource name.
func idFromState(s *terraform.State, idOrResourceName string, isResourceName bool) (
	string, error) {

	if !isResourceName {
		return idOrResourceName, nil
	}
	rs, ok := s.RootModule().Resources[idOrResourceName]
	if !ok {
		return "", fmt.Errorf("not found: %s", idOrResourceName)
	}
	if rs.Primary.ID == "" {
		return "", fmt.Errorf("widget ID is not set")
	}
	return rs.Primary.ID, nil
}

// apiConfig builds a Clumio SDK config from the environment.
func apiConfig() clumioConfig.Config {
	return clumioConfig.Config{
		Token:                     os.Getenv(common.ClumioApiToken),
		BaseUrl:                   os.Getenv(common.ClumioApiBaseUrl),
		OrganizationalUnitContext: os.Getenv(common.ClumioOrganizationalUnitContext),
		CustomHeaders: map[string]string{
			"User-Agent": "Clumio-Terraform-Provider-Acceptance-Test",
		},
	}
}

// DeletePolicy deletes the policy using the Clumio API. It takes as argument, either the resource
// name or the actual id of the policy.
func DeletePolicy(idOrResourceName string, isResourceName bool) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		id, err := idFromState(s, idOrResourceName, isResourceName)
		if err != nil {
			return err
		}
		config := apiConfig()
		pd := sdkclients.NewPolicyDefinitionClient(config)
		res, apiErr := pd.DeletePolicyDefinition(id)
		if apiErr != nil {
			return apiErr
		}
		if res == nil || res.TaskId == nil {
			return errors.New("expected task ID in the response")
		}
		return common.PollTask(
			context.Background(), sdkclients.NewTaskClient(config), *res.TaskId,
			300*time.Second, 5*time.Second)
	}
}

// DeleteProtectionGroup deletes the protection group using the Clumio API. It takes as argument,
// either the resource name or the actual id of the policy.
func DeleteProtectionGroup(idOrResourceName string, isResourceName bool) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		id, err := idFromState(s, idOrResourceName, isResourceName)
		if err != nil {
			return err
		}
		pd := protectionGroups.NewProtectionGroupsV1(apiConfig())
		if _, apiErr := pd.DeleteProtectionGroup(id); apiErr != nil {
			return apiErr
		}
		time.Sleep(3 * time.Second)
		for {
			pg, apiErr := pd.ReadProtectionGroup(id, nil)
			if apiErr != nil {
				if apiErr.ResponseCode == http.StatusNotFound {
					break
				}
				return apiErr
			}
			if *pg.IsDeleted {
				break
			}
			time.Sleep(2 * time.Second)
		}
		return nil
	}
}
