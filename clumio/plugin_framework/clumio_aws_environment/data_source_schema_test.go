// Copyright 2026. Clumio, Inc.

// This file contains the unit test for the Schema function in data_source_schema.go.

//go:build unit

package clumio_aws_environment

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/stretchr/testify/assert"
)

// TestDatasourceSchema checks the schema returned for the datasource.
func TestDatasourceSchema(t *testing.T) {

	ds := &clumioAWSEnvironmentDataSource{}
	resp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	assert.NotNil(t, resp.Schema)

	// Ensure that all attributes have a description set.
	for _, attr := range resp.Schema.Attributes {
		assert.NotEmpty(t, attr.GetDescription())
	}
}
