// Copyright 2021. Clumio, Inc.

package main

import (
	"context"
	"flag"
	"log"

	clumio_pf "github.com/clumio-code/terraform-provider-clumio/clumio/plugin_framework"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

// Run "go generate" to format example terraform files and generate the docs for the registry/website

// If you do not have terraform installed, you can remove the formatting command, but its suggested to
// ensure the documentation is formatted properly.
//go:generate terraform fmt -recursive ./examples/

// Run the docs generation tool, check its repository for more information on how it works and how docs
// can be customized.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	debugFlag := flag.Bool("debug", false, "Start provider in debug mode.")
	flag.Parse()

	err := providerserver.Serve(context.Background(), clumio_pf.New, providerserver.ServeOpts{
		Address: "clumio.com/providers/clumio",
		Debug:   *debugFlag,
	})
	if err != nil {
		log.Fatal(err)
	}
}
