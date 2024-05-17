// Copyright 2024. Clumio, Inc.

// This file contains the unit test for the Schema function in schema.go.

//go:build unit

package clumio_s3_bucket_properties

import (
	"context"
	"testing"
	"time"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"

	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
)

// Unit test for the resource Metadata and Configure functions.
func TestResourceMetadataAndConfigure(t *testing.T) {
	ctx := context.Background()
	res := NewClumioS3BucketPropertiesResource().(*clumioS3BucketPropertiesResource)

	t.Run("Metadata test", func(t *testing.T) {
		res.Metadata(ctx, resource.MetadataRequest{
			ProviderTypeName: "clumio",
		}, &resource.MetadataResponse{})
		assert.Equal(t, "clumio_s3_bucket_properties", res.name)
	})

	t.Run("Configure test with empty provider data", func(t *testing.T) {
		res.Configure(ctx, resource.ConfigureRequest{}, &resource.ConfigureResponse{})
		assert.Nil(t, res.sdkS3BucketClient)
	})

	t.Run("Configure test", func(t *testing.T) {
		apiClient := &common.ApiClient{ClumioConfig: sdkconfig.Config{}}
		res.Configure(
			ctx, resource.ConfigureRequest{ProviderData: apiClient}, &resource.ConfigureResponse{})
		assert.Equal(t, 2*time.Second, res.pollInterval)
		assert.Equal(t, 60*time.Second, res.pollTimeout)
	})
}

func TestCreate(t *testing.T) {
	ctx := context.Background()
	res := NewClumioS3BucketPropertiesResource().(*clumioS3BucketPropertiesResource)
	apiClient := &common.ApiClient{ClumioConfig: sdkconfig.Config{}}
	res.Configure(
		ctx, resource.ConfigureRequest{ProviderData: apiClient}, &resource.ConfigureResponse{})
	schemaResp := &resource.SchemaResponse{}
	res.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	plan := tfsdk.Plan{
		Raw:    tftypes.Value{},
		Schema: schemaResp.Schema,
	}
	createRes := &resource.CreateResponse{}
	res.Create(ctx, resource.CreateRequest{Plan: plan}, createRes)
	assert.True(t, createRes.Diagnostics.HasError())
}

// Test
func TestReadError(t *testing.T) {
	ctx := context.Background()
	res := NewClumioS3BucketPropertiesResource().(*clumioS3BucketPropertiesResource)
	apiClient := &common.ApiClient{ClumioConfig: sdkconfig.Config{}}
	res.Configure(
		ctx, resource.ConfigureRequest{ProviderData: apiClient}, &resource.ConfigureResponse{})
	schemaResp := &resource.SchemaResponse{}
	res.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	state := tfsdk.State{
		Raw:    tftypes.Value{},
		Schema: schemaResp.Schema,
	}
	readRes := &resource.ReadResponse{}
	res.Read(ctx, resource.ReadRequest{State: state}, readRes)
	assert.True(t, readRes.Diagnostics.HasError())
}

func TestUpdateError(t *testing.T) {
	ctx := context.Background()
	res := NewClumioS3BucketPropertiesResource().(*clumioS3BucketPropertiesResource)
	apiClient := &common.ApiClient{ClumioConfig: sdkconfig.Config{}}
	res.Configure(
		ctx, resource.ConfigureRequest{ProviderData: apiClient}, &resource.ConfigureResponse{})
	schemaResp := &resource.SchemaResponse{}
	res.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	plan := tfsdk.Plan{
		Raw:    tftypes.Value{},
		Schema: schemaResp.Schema,
	}
	updateRes := &resource.UpdateResponse{}
	res.Update(ctx, resource.UpdateRequest{Plan: plan}, updateRes)
	assert.True(t, updateRes.Diagnostics.HasError())
}

func TestDeleteError(t *testing.T) {
	ctx := context.Background()
	res := NewClumioS3BucketPropertiesResource().(*clumioS3BucketPropertiesResource)
	apiClient := &common.ApiClient{ClumioConfig: sdkconfig.Config{}}
	res.Configure(
		ctx, resource.ConfigureRequest{ProviderData: apiClient}, &resource.ConfigureResponse{})
	schemaResp := &resource.SchemaResponse{}
	res.Schema(context.Background(), resource.SchemaRequest{}, schemaResp)
	state := tfsdk.State{
		Raw:    tftypes.Value{},
		Schema: schemaResp.Schema,
	}
	deleteRes := &resource.DeleteResponse{}
	res.Delete(ctx, resource.DeleteRequest{State: state}, deleteRes)
	assert.True(t, deleteRes.Diagnostics.HasError())
}
