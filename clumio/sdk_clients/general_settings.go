// Copyright 2025. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK GeneralSettingsV2Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	sdkGeneralSettings "github.com/clumio-code/clumio-go-sdk/controllers/general_settings"
)

type GeneralSettingsClient interface {
	sdkGeneralSettings.GeneralSettingsV2Client
}

func NewGeneralSettingsClient(config config.Config) GeneralSettingsClient {
	return sdkGeneralSettings.NewGeneralSettingsV2(config)
}
