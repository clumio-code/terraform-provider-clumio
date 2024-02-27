// Copyright 2024. Clumio, Inc.

// This file hold various utility functions used by the clumio_aws_connection Terraform resource.

package clumio_aws_manual_connection_resources

import (
	"encoding/json"

	"github.com/clumio-code/clumio-go-sdk/models"
)

// stringifyResources accepts the resources struct from API response and converts it into a
// stringified version
func stringifyResources(resources *models.CategorisedResources) *string {
	bytes, err := json.Marshal(resources)
	if err != nil {
			return nil
	}
	stringifiedResources := string(bytes)
	return &stringifiedResources
}
