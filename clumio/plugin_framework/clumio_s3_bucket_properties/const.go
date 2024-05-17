// Copyright 2023. Clumio, Inc.

package clumio_s3_bucket_properties

const (
	schemaId                              = "id"
	schemaBucketId                        = "bucket_id"
	schemaEventBridgeEnabled              = "event_bridge_enabled"
	schemaEventBridgeNotificationDisabled = "event_bridge_notification_disabled"

	setErrorFmt          = "Unable to Set S3 Bucket Properties: %s"
	readS3BucketErrorFmt = "Unable to read properties of S3 bucket with id: %s"
)
