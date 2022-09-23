#
# Copyright 2021. Clumio, Inc.
#

VERSION=0.4.1
ifndef OS_ARCH
OS_ARCH=darwin_amd64
endif

CLUMIO_PROVIDER_DIR=~/.terraform.d/plugins/clumio.com/providers/clumio/${VERSION}/${OS_ARCH}
SWEEP?=us-east-1,us-east-2,us-west-2
SWEEP_DIR?=./clumio

default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

install:
	go mod vendor
	mkdir -p ${CLUMIO_PROVIDER_DIR}
	go build -o ${CLUMIO_PROVIDER_DIR}/terraform-provider-clumio_v${VERSION}

sweep:
	@echo "WARNING: This will destroy infrastructure. Use only in development accounts."
	go test $(SWEEP_DIR) -v -sweep=$(SWEEP) $(SWEEPARGS) -timeout 60m
