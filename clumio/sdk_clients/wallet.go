// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK UsersV2Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	"github.com/clumio-code/clumio-go-sdk/controllers/wallets"
)

type WalletClient interface {
	wallets.WalletsV1Client
}

func NewWalletClient(config config.Config) WalletClient {
	return wallets.NewWalletsV1(config)
}
