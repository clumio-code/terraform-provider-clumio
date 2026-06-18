// Copyright 2026. Clumio, Inc.

// This file contains the unit test for the Schema function in schema.go.

//go:build unit

package clumio_gcs_protection_group

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
)

// TestResourceSchema checks the schema returned for the clumio_gcs_protection_group resource.
func TestResourceSchema(t *testing.T) {

	res := &clumioGCSProtectionGroupResource{}
	resp := &resource.SchemaResponse{}
	res.Schema(context.Background(), resource.SchemaRequest{}, resp)
	assert.NotNil(t, resp.Schema)

	// Ensure that all top-level attributes have a description set
	for _, attr := range resp.Schema.Attributes {
		assert.NotEmpty(t, attr.GetDescription())
	}
}
