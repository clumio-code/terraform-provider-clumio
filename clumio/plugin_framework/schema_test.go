// Copyright 2024. Clumio, Inc.

// This file contains the unit test for the Schema function in schema.go.

//go:build unit

package clumio_pf

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/stretchr/testify/assert"
)

// TestSchema checks the schema returned for the provider.
func TestSchema(t *testing.T) {

	res := &clumioProvider{}
	resp := &provider.SchemaResponse{}
	res.Schema(context.Background(), provider.SchemaRequest{}, resp)
	assert.NotNil(t, resp.Schema)

	// Ensure that all attributes have a description set.
	for _, attr := range resp.Schema.Attributes {
		assert.NotEmpty(t, attr.GetMarkdownDescription())
	}
}
