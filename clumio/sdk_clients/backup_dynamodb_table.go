// Copyright 2026. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK BackupAwsDynamodbTablesV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	sdkBackupDynamoDBTables "github.com/clumio-code/clumio-go-sdk/controllers/backup_aws_dynamodb_tables"
)

type BackupDynamoDBTableClient interface {
	sdkBackupDynamoDBTables.BackupAwsDynamodbTablesV1Client
}

func NewBackupDynamoDBTableClient(config config.Config) BackupDynamoDBTableClient {
	return sdkBackupDynamoDBTables.NewBackupAwsDynamodbTablesV1(config)
}
