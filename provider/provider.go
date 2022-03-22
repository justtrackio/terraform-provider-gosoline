package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/justtrackio/terraform-provider-gosoline/builder"
)

type providerData struct {
	MetadataDomain types.String `tfsdk:"metadata_domain"`
}

type GosolineProvider struct {
	configured     bool
	metadataReader *builder.MetadataReader
}

func NewProvider() tfsdk.Provider {
	return &GosolineProvider{}
}

func (p *GosolineProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"metadata_domain": {
				Type:     types.StringType,
				Required: true,
			},
		},
	}, nil
}

func (p *GosolineProvider) Configure(ctx context.Context, request tfsdk.ConfigureProviderRequest, response *tfsdk.ConfigureProviderResponse) {
	var config providerData
	diags := request.Config.Get(ctx, &config)
	response.Diagnostics.Append(diags...)

	if response.Diagnostics.HasError() {
		return
	}

	p.configured = true
	p.metadataReader = builder.NewMetadataReader(config.MetadataDomain.Value)
}

func (p *GosolineProvider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, diag.Diagnostics) {
	return map[string]tfsdk.ResourceType{}, nil
}

func (p *GosolineProvider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, diag.Diagnostics) {
	return map[string]tfsdk.DataSourceType{
		"gosoline_application_dashboard_definition": &ApplicationDashboardDefinitionDatasourceType{},
		"gosoline_application_metadata_definition":  &ApplicationMetadataDefinitionDatasourceType{},
	}, nil
}
