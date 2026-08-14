// Copyright 2026. Clumio, Inc.

// This file contains the unit tests for the functions in data_source.go

//go:build unit

package clumio_dynamodb_backups

import (
	"context"
	"testing"

	"github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework/common"
	sdkclients "github.com/clumio-code/terraform-provider-clumio/clumio/sdk_clients"

	apiutils "github.com/clumio-code/clumio-go-sdk/api_utils"
	sdkconfig "github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/models"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Unit test for the following cases:
//   - Read DynamoDB table backups success scenario.
//   - SDK API for read DynamoDB table backups returns an error.
//   - SDK API for read DynamoDB table backups returns a nil response.
//   - SDK API for read DynamoDB table backups returns an empty items in response.
func TestDatasourceReadDynamoDBBackups(t *testing.T) {

	ctx := context.Background()
	backupClient := sdkclients.NewMockBackupDynamoDBTableClient(t)
	resourceName := "test_dynamodb_backups"
	tableId := "test-table-id"
	testError := "Test Error"
	latestBackupId := "test-backup-id-latest"
	olderBackupId := "test-backup-id-older"
	latestStartTimestamp := "2026-07-20T00:00:00Z"
	olderStartTimestamp := "2026-07-10T00:00:00Z"
	backupType := backupTypeClumioBackup

	rds := clumioDynamoDBBackupsDataSource{
		name: resourceName,
		client: &common.ApiClient{
			ClumioConfig: sdkconfig.Config{},
		},
		backupClient: backupClient,
	}

	rdsm := &clumioDynamoDBBackupsDataSourceModel{
		TableID:         basetypes.NewStringValue(tableId),
		ClumioType:      basetypes.NewStringValue(backupType),
		BeforeTimestamp: basetypes.NewStringValue("2026-07-25T00:00:00Z"),
		AfterTimestamp:  basetypes.NewStringValue("2026-07-01T00:00:00Z"),
	}

	apiError := &apiutils.APIError{
		ResponseCode: 500,
		Reason:       "test",
		Response:     []byte(testError),
	}

	// Tests the success scenario for DynamoDB table backups read. It should not return
	// Diagnostics and the backups in the model should preserve the order of the API response.
	t.Run("Basic success scenario for read DynamoDB backups", func(t *testing.T) {

		count := int64(2)
		readResponse := &models.ListDynamoDBTableBackupsResponse{
			Embedded: &models.DynamoDBTableBackupListEmbedded{
				Items: []*models.DynamoDBTableBackupWithETag{
					{
						Id:             &latestBackupId,
						StartTimestamp: &latestStartTimestamp,
						ClumioType:     &backupType,
					},
					{
						Id:             &olderBackupId,
						StartTimestamp: &olderStartTimestamp,
						ClumioType:     &backupType,
					},
				},
			},
			CurrentCount: &count,
		}

		// The exact filter and sort expected to be sent to the API for the model above.
		expectedFilter := `{"table_id": {"$eq":"test-table-id"},` +
			`"start_timestamp": {"$lte":"2026-07-25T00:00:00Z","$gt":"2026-07-01T00:00:00Z"},` +
			`"type": {"$all":["clumio_backup"]}}`
		expectedSort := "-start_timestamp"

		// Setup expectations.
		backupClient.EXPECT().ListBackupAwsDynamodbTables(mock.Anything, mock.Anything,
			mock.MatchedBy(func(sort *string) bool {
				return sort != nil && *sort == expectedSort
			}),
			mock.MatchedBy(func(filter *string) bool {
				return filter != nil && *filter == expectedFilter
			})).Times(1).Return(readResponse, nil)

		diags := rds.readDynamoDBBackups(ctx, rdsm)
		assert.Nil(t, diags)

		backups := rdsm.Backups.Elements()
		assert.Len(t, backups, 2)
		firstBackup := backups[0].(types.Object).Attributes()
		assert.Equal(t, latestBackupId, firstBackup[schemaId].(types.String).ValueString())
	})

	// Tests that Diagnostics is returned in case the list DynamoDB table backups API call returns
	// an error.
	t.Run("list DynamoDB backups returns an error", func(t *testing.T) {

		// Setup expectations.
		backupClient.EXPECT().ListBackupAwsDynamodbTables(mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(nil, apiError)

		diags := rds.readDynamoDBBackups(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list DynamoDB table backups API call returns
	// a nil response.
	t.Run("list DynamoDB backups returns a nil response", func(t *testing.T) {

		// Setup expectations.
		backupClient.EXPECT().ListBackupAwsDynamodbTables(mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(nil, nil)

		diags := rds.readDynamoDBBackups(ctx, rdsm)
		assert.NotNil(t, diags)
	})

	// Tests that Diagnostics is returned in case the list DynamoDB table backups API call returns
	// an empty items in response.
	t.Run("list DynamoDB backups returns an empty items in response", func(t *testing.T) {

		readResponse := &models.ListDynamoDBTableBackupsResponse{
			Embedded: &models.DynamoDBTableBackupListEmbedded{
				Items: []*models.DynamoDBTableBackupWithETag{},
			},
		}

		// Setup expectations.
		backupClient.EXPECT().ListBackupAwsDynamodbTables(mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(readResponse, nil)

		diags := rds.readDynamoDBBackups(ctx, rdsm)
		assert.NotNil(t, diags)
	})
}

// Unit test for the datasource Metadata, Configure and Read functions.
func TestDatasourceMetadataConfigureRead(t *testing.T) {

	ctx := context.Background()
	backupId := "test-backup-id"
	startTimestamp := "2026-07-20T00:00:00Z"

	// Tests that the datasource type name is set as part of Metadata().
	t.Run("Metadata sets the datasource type name", func(t *testing.T) {

		ds := &clumioDynamoDBBackupsDataSource{}
		resp := &datasource.MetadataResponse{}
		ds.Metadata(ctx, datasource.MetadataRequest{ProviderTypeName: "clumio"}, resp)
		assert.Equal(t, "clumio_dynamodb_backups", resp.TypeName)
	})

	// Tests that Configure() returns early when no provider data is given and sets up the SDK
	// client when it is.
	t.Run("Configure sets up the SDK client", func(t *testing.T) {

		ds := &clumioDynamoDBBackupsDataSource{}
		ds.Configure(ctx, datasource.ConfigureRequest{}, &datasource.ConfigureResponse{})
		assert.Nil(t, ds.client)

		ds.Configure(ctx, datasource.ConfigureRequest{
			ProviderData: &common.ApiClient{ClumioConfig: sdkconfig.Config{}},
		}, &datasource.ConfigureResponse{})
		assert.NotNil(t, ds.client)
		assert.NotNil(t, ds.backupClient)
	})

	// The Terraform configuration used by the Read() tests below.
	backupObjType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			schemaId:                  tftypes.String,
			schemaStartTimestamp:      tftypes.String,
			schemaExpirationTimestamp: tftypes.String,
			schemaType:                tftypes.String,
		},
	}
	configType := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			schemaTableId:         tftypes.String,
			schemaType:            tftypes.String,
			schemaBeforeTimestamp: tftypes.String,
			schemaAfterTimestamp:  tftypes.String,
			schemaBackups:         tftypes.List{ElementType: backupObjType},
		},
	}
	configVals := map[string]tftypes.Value{
		schemaTableId:         tftypes.NewValue(tftypes.String, "test-table-id"),
		schemaType:            tftypes.NewValue(tftypes.String, nil),
		schemaBeforeTimestamp: tftypes.NewValue(tftypes.String, nil),
		schemaAfterTimestamp:  tftypes.NewValue(tftypes.String, nil),
		schemaBackups:         tftypes.NewValue(tftypes.List{ElementType: backupObjType}, nil),
	}

	schemaResp := &datasource.SchemaResponse{}
	(&clumioDynamoDBBackupsDataSource{}).Schema(ctx, datasource.SchemaRequest{}, schemaResp)
	readRequest := datasource.ReadRequest{
		Config: tfsdk.Config{
			Raw:    tftypes.NewValue(configType, configVals),
			Schema: schemaResp.Schema,
		},
	}

	// Tests the success scenario for the datasource Read(). It should not return Diagnostics and
	// should set the backups in the Terraform state.
	t.Run("Read sets the state from the API response", func(t *testing.T) {

		backupClient := sdkclients.NewMockBackupDynamoDBTableClient(t)
		ds := &clumioDynamoDBBackupsDataSource{
			name:         "clumio_dynamodb_backups",
			backupClient: backupClient,
		}

		readResponse := &models.ListDynamoDBTableBackupsResponse{
			Embedded: &models.DynamoDBTableBackupListEmbedded{
				Items: []*models.DynamoDBTableBackupWithETag{
					{
						Id:             &backupId,
						StartTimestamp: &startTimestamp,
					},
				},
			},
		}

		// Setup expectations.
		backupClient.EXPECT().ListBackupAwsDynamodbTables(mock.Anything, mock.Anything,
			mock.Anything, mock.Anything).Times(1).Return(readResponse, nil)

		resp := &datasource.ReadResponse{
			State: tfsdk.State{
				Raw:    tftypes.NewValue(configType, nil),
				Schema: schemaResp.Schema,
			},
		}
		ds.Read(ctx, readRequest, resp)
		assert.False(t, resp.Diagnostics.HasError())
		assert.False(t, resp.State.Raw.IsNull())
	})
}
