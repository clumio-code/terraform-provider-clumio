// Copyright 2024. Clumio, Inc.

// This file contains the common test functions which are used by one or more acceptance tests.

package common

import (
	"context"
	"errors"
	"fmt"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"
	"net/http"
	"os"
	"time"

	clumioConfig "github.com/clumio-code/clumio-go-sdk/config"
	protectionGroups "github.com/clumio-code/clumio-go-sdk/controllers/protection_groups"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

// DeletePolicy deletes the policy using the Clumio API. It takes as argument, either the resource
// name or the actual id of the policy.
func DeletePolicy(idOrResourceName string, isResourceName bool) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		id := idOrResourceName
		if isResourceName {
			// retrieve the resource by name from state
			rs, ok := s.RootModule().Resources[idOrResourceName]
			if !ok {
				return fmt.Errorf("Not found: %s", idOrResourceName)
			}

			if rs.Primary.ID == "" {
				return fmt.Errorf("Widget ID is not set")
			}
			id = rs.Primary.ID
		}
		clumioApiToken := os.Getenv(ClumioApiToken)
		clumioApiBaseUrl := os.Getenv(ClumioApiBaseUrl)
		clumioOrganizationalUnitContext := os.Getenv(ClumioOrganizationalUnitContext)
		config := clumioConfig.Config{
			Token:                     clumioApiToken,
			BaseUrl:                   clumioApiBaseUrl,
			OrganizationalUnitContext: clumioOrganizationalUnitContext,
			CustomHeaders: map[string]string{
				"User-Agent": "Clumio-Terraform-Provider-Acceptance-Test",
			},
		}
		pd := sdkclients.NewPolicyDefinitionClient(config)
		res, apiErr := pd.DeletePolicyDefinition(id)
		if apiErr != nil {
			return apiErr
		}
		if res != nil && res.TaskId != nil {
			taskClient := sdkclients.NewTaskClient(config)
			err := PollTask(
				context.Background(), taskClient, *res.TaskId, 300*time.Second, 5*time.Second)
			if err != nil {
				return err
			}
		} else {
			return errors.New("expected task ID in the response")
		}
		return nil
	}
}

// DeleteProtectionGroup deletes the protection group using the Clumio API. It takes as argument,
// either the resource name or the actual id of the policy.
func DeleteProtectionGroup(idOrResourceName string, isResourceName bool) resource.TestCheckFunc {

	return func(s *terraform.State) error {

		id := idOrResourceName
		if isResourceName {
			// retrieve the resource by name from state
			rs, ok := s.RootModule().Resources[idOrResourceName]
			if !ok {
				return fmt.Errorf("Not found: %s", idOrResourceName)
			}

			if rs.Primary.ID == "" {
				return fmt.Errorf("Widget ID is not set")
			}
			id = rs.Primary.ID
		}

		clumioApiToken := os.Getenv(ClumioApiToken)
		clumioApiBaseUrl := os.Getenv(ClumioApiBaseUrl)
		clumioOrganizationalUnitContext := os.Getenv(ClumioOrganizationalUnitContext)
		config := clumioConfig.Config{
			Token:                     clumioApiToken,
			BaseUrl:                   clumioApiBaseUrl,
			OrganizationalUnitContext: clumioOrganizationalUnitContext,
			CustomHeaders: map[string]string{
				"User-Agent": "Clumio-Terraform-Provider-Acceptance-Test",
			},
		}
		pd := protectionGroups.NewProtectionGroupsV1(config)
		_, apiErr := pd.DeleteProtectionGroup(id)
		if apiErr != nil {
			return apiErr
		}
		time.Sleep(3 * time.Second)
		for {
			pg, apiErr := pd.ReadProtectionGroup(id, &DefaultLookBackDays)
			if apiErr != nil {
				if apiErr.ResponseCode == http.StatusNotFound {
					break
				}
				return apiErr
			}
			if !*pg.IsDeleted {
				time.Sleep(2 * time.Second)
			} else {
				break
			}
		}
		return nil
	}
}
