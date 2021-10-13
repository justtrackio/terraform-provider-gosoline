package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/justtrackio/terraform-provider-gosoline/provider"
)

func main() {
	tfsdk.Serve(context.Background(), provider.NewProvider, tfsdk.ServeOpts{
		Name: "gosoline",
	})
}
