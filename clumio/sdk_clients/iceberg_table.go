// Copyright 2026. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK AwsIcebergTablesV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	sdkIcebergTable "github.com/clumio-code/clumio-go-sdk/controllers/aws_iceberg_tables"
)

type IcebergTableClient interface {
	sdkIcebergTable.AwsIcebergTablesV1Client
}

func NewIcebergTableClient(config config.Config) IcebergTableClient {
	return sdkIcebergTable.NewAwsIcebergTablesV1(config)
}
