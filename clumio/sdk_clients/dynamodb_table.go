// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK PolicyDefinitionsV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	sdkDynamoDBTable "github.com/clumio-code/clumio-go-sdk/controllers/aws_dynamodb_tables"
)

type DynamoDBTableClient interface {
	sdkDynamoDBTable.AwsDynamodbTablesV1Client
}

func NewDynamoDBTableClient(config config.Config) DynamoDBTableClient {
	return sdkDynamoDBTable.NewAwsDynamodbTablesV1(config)
}
