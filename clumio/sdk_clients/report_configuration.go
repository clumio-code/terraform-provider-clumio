// Copyright 2025. Clumio, Inc.

// Contains the wrapper interface for Clumio GO SDK ReportComplianceV1Client.

package sdkclients

import (
	"github.com/clumio-code/clumio-go-sdk/config"
	sdkReportCompliances "github.com/clumio-code/clumio-go-sdk/controllers/report_compliance"
)

type ReportConfigurationClient interface {
	sdkReportCompliances.ReportComplianceV1Client
}

func NewReportConfigurationClient(config config.Config) ReportConfigurationClient {
	return sdkReportCompliances.NewReportComplianceV1(config)
}
