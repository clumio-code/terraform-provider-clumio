package common

import clumioConfig "github.com/clumio-code/clumio-go-sdk/config"

// ApiClient defines the APIs/connections required by the resources.
type ApiClient struct {
	ClumioConfig clumioConfig.Config
}
