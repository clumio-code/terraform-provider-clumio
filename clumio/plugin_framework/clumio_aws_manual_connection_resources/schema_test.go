// Copyright 2024. Clumio, Inc.

// This file contains the unit test for the Schema function in schema.go.

//go:build unit

package clumio_aws_manual_connection_resources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stretchr/testify/assert"
)

// TestSchema checks the schema returned for a given resource.
func TestSchema(t *testing.T) {

	res := &clumioAwsManualConnectionResourcesDatasource{}
	resp := &datasource.SchemaResponse{}
	res.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	assert.NotNil(t, resp.Schema)

	// Ensure that all attributes have a description set.
	for _, attr := range resp.Schema.Attributes {
		assert.NotEmpty(t, attr.GetDescription())
	}
}