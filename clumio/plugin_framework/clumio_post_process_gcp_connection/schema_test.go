// Copyright (c) 2025 Clumio, a Commvault Company All Rights Reserved

//go:build unit

package clumio_post_process_gcp_connection

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/stretchr/testify/assert"
)

// TestSchema checks the schema returned for a given resource.
func TestSchema(t *testing.T) {

	res := &clumioPostProcessGCPConnectionResource{}
	resp := &resource.SchemaResponse{}
	res.Schema(context.Background(), resource.SchemaRequest{}, resp)
	assert.NotNil(t, resp.Schema)

	// Ensure that all attributes have a description set.
	for _, attr := range resp.Schema.Attributes {
		assert.NotEmpty(t, attr.GetDescription())
	}
}
