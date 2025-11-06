#
# Copyright 2021. Clumio, Inc.
#

# If the version is being changed here, it should also be changed for the variable
# clumioTfProviderVersionValue in the file clumio/plugin_framework/common/const.go.
VERSION=0.16.1
ifndef OS_ARCH
OS_ARCH=darwin_arm64
endif

CLUMIO_PROVIDER_DIR=~/.terraform.d/plugins/clumio.com/providers/clumio/${VERSION}/${OS_ARCH}
SWEEP?=us-east-1,us-east-2,us-west-2
SWEEP_DIR?=./clumio

REPORTS_DIR=build/reports
TESTSUM_ARGS=--format=pkgname-and-test-fails
TESTCOVER_ARGS=-covermode=set -coverprofile=$(REPORTS_DIR)/coverage.out
TESTSUMCOVER_ARGS=--jsonfile=$(REPORTS_DIR)/test-report.out --junitfile=$(REPORTS_DIR)/test-report.xml
.PHONY: all
default: testacc

# Run acceptance tests
.PHONY: testacc
testacc: testacc_basic testacc_post_process

.PHONY: testacc_basic
testacc_basic:
	TF_ACC=1 gotestsum $(TESTSUM_ARGS) -- -vet=off -v ./... $(TESTARGS) -tags="basic" -timeout 120m

.PHONY: testacc_post_process
testacc_post_process:
	TF_ACC=1 gotestsum $(TESTSUM_ARGS) -- -vet=off -v ./... $(TESTARGS) -tags="post_process" -timeout 120m

.PHONY: testacc_sso
testacc_sso:
	TF_ACC=1 gotestsum $(TESTSUM_ARGS) -- -vet=off -v ./... $(TESTARGS) -tags="sso" -timeout 120m

.PHONY: testacc_bucket
testacc_bucket:
	TF_ACC=1 gotestsum $(TESTSUM_ARGS) -- -vet=off -v ./... $(TESTARGS) -tags="bucket" -timeout 120m

.PHONY: testacc_general_settings
testacc_general_settings:
	TF_ACC=1 gotestsum $(TESTSUM_ARGS) -- -vet=off -v ./... $(TESTARGS) -tags="general_settings" -timeout 120m

.PHONY: testunit
testunit:
	gotestsum $(TESTSUM_ARGS) -- -vet=off ./... -v $(TESTARGS) -tags="unit" -timeout 90s

.PHONY: testcover
testcover:
	rm -rf $(REPORTS_DIR) && mkdir -p $(REPORTS_DIR)
	gotestsum $(TESTSUM_ARGS) $(TESTSUMCOVER_ARGS) -- -vet=off $(TESTCOVER_ARGS) ./... $(TESTARGS) -tags="unit" -timeout 90s
# Note that we feed the test results and coverage report to sonar in the ci pipeline.
# The full coverage report is only available for the main branch, while the PRs
# will only report on the modified files.
# The test results are only reported in main and never in PRs by SonarCloud.

.PHONY: install
install:
	go mod vendor
	mkdir -p ${CLUMIO_PROVIDER_DIR}
	go build -o ${CLUMIO_PROVIDER_DIR}/terraform-provider-clumio_v${VERSION}

.PHONY: sweep
sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(SWEEP_DIR) -v -sweep=$(SWEEP) $(SWEEPARGS) -timeout 60m

# Mockery is the actively maintained tool to generate mocks in Go.
# To add mocks update the .mockery.yaml file.
.PHONY: mockery
mockery:
	mockery
