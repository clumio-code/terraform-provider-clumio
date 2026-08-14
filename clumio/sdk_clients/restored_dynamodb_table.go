// Copyright 2026. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK RestoredAwsDynamodbTablesV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	sdkRestoredDynamoDBTables "github.com/clumio-code/clumio-go-sdk/controllers/restored_aws_dynamodb_tables"
)

type RestoredDynamoDBTableClient interface {
	sdkRestoredDynamoDBTables.RestoredAwsDynamodbTablesV1Client
}

func NewRestoredDynamoDBTableClient(config config.Config) RestoredDynamoDBTableClient {
	return sdkRestoredDynamoDBTables.NewRestoredAwsDynamodbTablesV1(config)
}
