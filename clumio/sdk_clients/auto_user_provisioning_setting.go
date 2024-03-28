// Copyright 2024. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK AutoUserProvisioningSettingsV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	sdkAUPSettings "github.com/clumio-code/clumio-go-sdk/controllers/auto_user_provisioning_settings"
)

type AutoUserProvisioningSettingClient interface {
	sdkAUPSettings.AutoUserProvisioningSettingsV1Client
}

func NewAutoUserProvisioningSettingClient(config config.Config) AutoUserProvisioningSettingClient {
	return sdkAUPSettings.NewAutoUserProvisioningSettingsV1(config)
}
