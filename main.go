package main

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/justtrackio/terraform-provider-gosoline/provider"
)

func main() {
	err := providerserver.Serve(context.Background(), provider.NewProvider, providerserver.ServeOpts{
		Address: "registry.terraform.io/justtrackio/gosoline",
	})
	if err != nil {
		panic(fmt.Sprintf("failed running justtrack gosoline provider: %s", err.Error()))
	}
}
