// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK ProtectionGroupsS3AssetsV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	sdkProtectionGroupS3Assets "github.com/clumio-code/clumio-go-sdk/controllers/protection_groups_s3_assets"
)

type ProtectionGroupS3AssetsClient interface {
	sdkProtectionGroupS3Assets.ProtectionGroupsS3AssetsV1Client
}

func NewProtectionGroupS3AssetsClient(config config.Config) ProtectionGroupS3AssetsClient {
	return sdkProtectionGroupS3Assets.NewProtectionGroupsS3AssetsV1(config)
}
