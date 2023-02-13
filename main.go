package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/justtrackio/terraform-provider-gosoline/provider"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), provider.NewProvider, providerserver.ServeOpts{
		Address: "registry.terraform.io/justtrackio/gosoline",
		Debug:   debug,
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}
