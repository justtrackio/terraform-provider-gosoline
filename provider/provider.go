package provider

import (
	"context"
	"fmt"
	"math/big"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/justtrackio/terraform-provider-gosoline/builder"
)

const (
	defaultMetadataHostnameNamePattern               = "{scheme}://{group}-{app}.{family}.{env}.{metadata_domain}:{port}"
	defaultMetadataUseHttps                          = true
	defaultMetadataPort                              = 8070
	defaultEcsClusterNamePattern                     = "{env}"
	defaultEcsServiceNamePattern                     = "{group}-{app}"
	defaultCloudwatchNamespaceNamePattern            = "{project}/{env}/{family}/{group}-{app}"
	defaultGrafanaCloudWatchDatasourceNamePattern    = "cloudwatch-{family}"
	defaultGrafanaElasticsearchDatasourceNamePattern = "elasticsearch-{env}-logs-{project}-{family}-{group}-{app}"
)

type providerData struct {
	Metadata     types.Object `tfsdk:"metadata"`
	NamePatterns types.Object `tfsdk:"name_patterns"`
}

type ResourceNamePatterns struct {
	Hostname                       string
	CloudwatchNamespace            string
	EcsCluster                     string
	EcsService                     string
	GrafanaCloudWatchDatasource    string
	GrafanaElasticsearchDatasource string
}

type MetadataProperties struct {
	Domain   string
	UseHttps bool
	Port     int
}

type GosolineProvider struct {
	metadataReader                *builder.MetadataReader
	resourceNamePatterns          ResourceNamePatterns
	additionalAugmentReplacements map[string]string
}

func NewProvider() tfsdk.Provider {
	return &GosolineProvider{}
}

func (p *GosolineProvider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"metadata": {
				Type: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"domain":    types.StringType,
						"use_https": types.BoolType,
						"port":      types.NumberType,
					},
				},
				Required: true,
				MarkdownDescription: `domain: This is the base domain where your services are available, e.g. example.com. This field is required!
									  use_https: Allows to change from https to http (default: ` + fmt.Sprint(defaultMetadataUseHttps) + `)
									  port: Allows to change the default metadata port (default: ` + fmt.Sprint(defaultMetadataPort) + `)`,
			},
			"name_patterns": {
				Type: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						"hostname":                         types.StringType,
						"cloudwatch_namespace":             types.StringType,
						"ecs_cluster":                      types.StringType,
						"ecs_service":                      types.StringType,
						"grafana_cloudwatch_datasource":    types.StringType,
						"grafana_elasticsearch_datasource": types.StringType,
					},
				},
				Optional: true,
				MarkdownDescription: `hostname: Allows to change the default metadata hostname name pattern (default: ` + defaultMetadataHostnameNamePattern + `)
										  Available placeholders are:
										  * {project}
										  * {env}
										  * {family}
										  * {group}
										  * {app}
										  * {scheme} (http/https, depends on your metadata.use_https provider configuration)
										  * {metadata_domain} (your supplied metadata.domain for the provider configuration)
										  * {port} (depends on your metadata.port provider configuration)
									  cloudwatch_namespace: Allows to change the default cloudwatch namespace name pattern (default: ` + defaultCloudwatchNamespaceNamePattern + `)
										  Available placeholders are:
										  * {project}
										  * {env}
										  * {family}
										  * {group}
										  * {app}
									  ecs_cluster: Allows to change the default ecs cluster name pattern (default: ` + defaultEcsClusterNamePattern + `)
										  Available placeholders are:
										  * {project}
										  * {env}
										  * {family}
										  * {group}
										  * {app}
									  ecs_service: Allows to change the default ecs service name pattern (default: ` + defaultEcsServiceNamePattern + `)
										  Available placeholders are:
										  * {project}
										  * {env}
										  * {family}
										  * {group}
										  * {app}
									  grafana_cloudwatch_datasource: Allows to change the default grafana cloudwatch datasource name pattern (default: ` + defaultGrafanaCloudWatchDatasourceNamePattern + `)
										  Available placeholders are:
										  * {project}
										  * {env}
										  * {family}
										  * {group}
										  * {app}
									  grafana_elasticsearch_datasource: Allows to change the default grafana elasticsearch datasource name pattern (default: ` + defaultGrafanaElasticsearchDatasourceNamePattern + `)
										  Available placeholders are:
										  * {project}
										  * {env}
										  * {family}
										  * {group}
										  * {app}`,
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

	metadataProperties, err := p.getMetadataProperties(ctx, config)
	if err != nil {
		diags.AddError("failed to get metadata properties from attributes", err.Error())
		return
	}

	namepatternProperties, err := p.getNamepatternProperties(ctx, config)
	if err != nil {
		diags.AddError("failed to get metadata properties from attributes", err.Error())
		return
	}

	scheme := "https"
	if !metadataProperties.UseHttps {
		scheme = "http"
	}

	additionalReplacements := map[string]string{
		"metadata_domain": metadataProperties.Domain,
		"scheme":          scheme,
		"port":            fmt.Sprint(metadataProperties.Port),
	}

	p.additionalAugmentReplacements = additionalReplacements
	p.resourceNamePatterns = *namepatternProperties
	p.metadataReader = builder.NewMetadataReader(namepatternProperties.Hostname, additionalReplacements)
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

func (p *GosolineProvider) convertDirectlyToNativeType(ctx context.Context, input attr.Value, dest any) error {
	value, err := input.ToTerraformValue(ctx)
	if err != nil {
		return fmt.Errorf("failed to read terraform value from input: %w", err)
	}

	switch destN := dest.(type) {
	case *string:
		err = value.As(&destN)
	case *bool:
		err = value.As(&destN)
	default:
		return fmt.Errorf("destination type %T not supported", dest)
	}
	if err != nil {
		return fmt.Errorf("failed to convert to native representation: %w", err)
	}

	return nil
}

func (p *GosolineProvider) convertBigFloatToNativeType(ctx context.Context, input attr.Value, dest any) error {
	var bigFloatValue big.Float

	value, err := input.ToTerraformValue(ctx)
	if err != nil {
		return fmt.Errorf("failed to read terraform value from input: %w", err)
	}

	err = value.As(&bigFloatValue)
	if err != nil {
		return fmt.Errorf("failed to convert to native representation: %w", err)
	}

	bigIntValue, _ := bigFloatValue.Int64()
	*dest.(*int) = int(bigIntValue)
	return nil
}

func (p *GosolineProvider) convertToNativeType(ctx context.Context, input attr.Value, dest any) error {
	switch dest.(type) {
	case *string:
		return p.convertDirectlyToNativeType(ctx, input, dest)
	case *bool:
		return p.convertDirectlyToNativeType(ctx, input, dest)
	case *int:
		return p.convertBigFloatToNativeType(ctx, input, dest)
	default:
		return fmt.Errorf("destination type is invalid, supplied %T", dest)
	}
}

func (p *GosolineProvider) getMetadataProperties(ctx context.Context, config providerData) (*MetadataProperties, error) {
	var metadataDomain string

	metadataDomainTfVal := config.Metadata.Attrs["domain"]
	err := p.convertToNativeType(ctx, metadataDomainTfVal, &metadataDomain)
	if err != nil {
		return nil, fmt.Errorf("failed to convert metadata.domain attribute to native type: %w", err)
	}

	var useHttps bool
	err = p.convertToNativeType(ctx, config.Metadata.Attrs["use_https"], &useHttps)
	if err != nil {
		return nil, fmt.Errorf("failed to convert metadata.use_https attribute to native type: %w", err)
	}

	var port int
	err = p.convertToNativeType(ctx, config.Metadata.Attrs["port"], &port)
	if err != nil {
		return nil, fmt.Errorf("failed to convert metadata.port attribute to native type: %w", err)
	}

	props := &MetadataProperties{
		Domain:   metadataDomain,
		UseHttps: useHttps,
		Port:     port,
	}

	return props, nil
}

func (p *GosolineProvider) getNamepatternProperties(ctx context.Context, config providerData) (*ResourceNamePatterns, error) {
	patterns := map[string]string{
		"hostname":                         defaultMetadataHostnameNamePattern,
		"cloudwatch_namespace":             defaultCloudwatchNamespaceNamePattern,
		"ecs_cluster":                      defaultEcsClusterNamePattern,
		"ecs_service":                      defaultEcsServiceNamePattern,
		"grafana_cloudwatch_datasource":    defaultGrafanaCloudWatchDatasourceNamePattern,
		"grafana_elasticsearch_datasource": defaultGrafanaElasticsearchDatasourceNamePattern,
	}

	for key := range patterns {
		if _, ok := config.NamePatterns.Attrs[key]; !ok {
			continue
		}

		var value string
		if err := p.convertToNativeType(ctx, config.NamePatterns.Attrs[key], &value); err != nil {
			return nil, fmt.Errorf("failed to convert name_patterns.%s attribute to native type: %w", key, err)
		}

		patterns[key] = value
	}

	props := &ResourceNamePatterns{
		Hostname:                       patterns["hostname"],
		CloudwatchNamespace:            patterns["cloudwatch_namespace"],
		EcsCluster:                     patterns["ecs_cluster"],
		EcsService:                     patterns["ecs_service"],
		GrafanaCloudWatchDatasource:    patterns["grafana_cloudwatch_datasource"],
		GrafanaElasticsearchDatasource: patterns["grafana_elasticsearch_datasource"],
	}

	return props, nil
}
