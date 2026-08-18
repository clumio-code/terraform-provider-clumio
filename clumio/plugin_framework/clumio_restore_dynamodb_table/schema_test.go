// Copyright 2026. Clumio, Inc.

// This file contains the unit test for the Schema function in schema.go.

//go:build unit

package clumio_restore_dynamodb_table

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/stretchr/testify/assert"
)

// TestActionSchema checks the schema returned for the action.
func TestActionSchema(t *testing.T) {

	ra := &clumioRestoreDynamoDBTableAction{}
	resp := &action.SchemaResponse{}
	ra.Schema(context.Background(), action.SchemaRequest{}, resp)
	assert.NotNil(t, resp.Schema)

	// Ensure that all attributes have a description set.
	for _, attr := range resp.Schema.Attributes {
		assert.NotEmpty(t, attr.GetDescription())
	}
}
