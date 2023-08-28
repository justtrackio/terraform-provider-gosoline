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
	"github.com/thoas/go-funk"
)

var availableOrchestrators = []string{orchestratorEcs, orchestratorKubernetes}

const (
	orchestratorEcs                                  = "ecs"
	orchestratorKubernetes                           = "kubernetes"
	defaultMetadataHostnameNamePattern               = "{scheme}://{group}-{app}.{family}.{env}.{metadata_domain}:{port}"
	defaultMetadataUseHttps                          = true
	defaultMetadataPort                              = 8070
	defaultOrchestrator                              = orchestratorEcs
	defaultEcsClusterNamePattern                     = "{env}"
	defaultEcsServiceNamePattern                     = "{group}-{app}"
	defaultCloudwatchNamespaceNamePattern            = "{project}/{env}/{family}/{group}-{app}"
	defaultGrafanaCloudWatchDatasourceNamePattern    = "cloudwatch-{family}"
	defaultGrafanaElasticsearchDatasourceNamePattern = "elasticsearch-{env}-logs-{project}-{family}-{group}-{app}"
	defaultKubernetesNamespaceNamePattern            = "{project}"
	defaultKubernetesPodNamePattern                  = "{group}-{app}"
	defaultTraefikServiceNameNamePattern             = "{project}-{group}-{app}-8080@kubernetes"
	propCloudwatchNamespace                          = "cloudwatch_namespace"
	propEcsCluster                                   = "ecs_cluster"
	propEcsService                                   = "ecs_service"
	propGrafanaCloudwatchDatasource                  = "grafana_cloudwatch_datasource"
	propGrafanaElasticsearchDatasource               = "grafana_elasticsearch_datasource"
	propHostname                                     = "hostname"
	propKubernetesNamespace                          = "kubernetes_namespace"
	propKubernetesPod                                = "kubernetes_pod"
	propTraefikServiceName                           = "traefik_service_name"
)

type providerData struct {
	Metadata     types.Object `tfsdk:"metadata"`
	NamePatterns types.Object `tfsdk:"name_patterns"`
	Orchestrator types.String `tfsdk:"orchestrator"`
}

type ResourceNamePatterns struct {
	Hostname                       string
	CloudwatchNamespace            string
	EcsCluster                     string
	EcsService                     string
	GrafanaCloudWatchDatasource    string
	GrafanaElasticsearchDatasource string
	KubernetesNamespace            string
	KubernetesPod                  string
	TraefikServiceName             string
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
	orchestrator                  string
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
			"orchestrator": {
				Type:                types.StringType,
				Optional:            true,
				MarkdownDescription: `orchestrator: Set this to "ecs" for getting ELB/Target-group/ECS related metrics or "kubernetes" to get traefik related metrics inside the grafana dashboard`,
			},
			"name_patterns": {
				Type: types.ObjectType{
					AttrTypes: map[string]attr.Type{
						propHostname:                       types.StringType,
						propCloudwatchNamespace:            types.StringType,
						propEcsCluster:                     types.StringType,
						propEcsService:                     types.StringType,
						propGrafanaCloudwatchDatasource:    types.StringType,
						propGrafanaElasticsearchDatasource: types.StringType,
						propKubernetesNamespace:            types.StringType,
						propKubernetesPod:                  types.StringType,
						propTraefikServiceName:             types.StringType,
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
										  * {app}
									  kubernetes_namespace: Allows to change the default kubernetes namespace name pattern (default: ` + defaultKubernetesNamespaceNamePattern + `)
										  Available placeholders are:
										  * {project}
										  * {env}
										  * {family}
										  * {group}
										  * {app}
									  kubernetes_pod: Allows to change the default kubernetes pod name pattern (default: ` + defaultKubernetesPodNamePattern + `)
										  Available placeholders are:
										  * {project}
										  * {env}
										  * {family}
										  * {group}
										  * {app}
									  traefik_service_name: Allows to change the default traefik service name name pattern (default: ` + defaultTraefikServiceNameNamePattern + `)
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

	p.orchestrator = defaultOrchestrator
	if !config.Orchestrator.IsNull() {
		p.orchestrator = config.Orchestrator.Value
	}
	if !funk.ContainsString(availableOrchestrators, p.orchestrator) {
		response.Diagnostics.AddError("invalid operator", fmt.Sprintf("'%s' is not a valid orchestrator, choose between %v", p.orchestrator, availableOrchestrators))
		return
	}

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
		propHostname:                       defaultMetadataHostnameNamePattern,
		propCloudwatchNamespace:            defaultCloudwatchNamespaceNamePattern,
		propEcsCluster:                     defaultEcsClusterNamePattern,
		propEcsService:                     defaultEcsServiceNamePattern,
		propGrafanaCloudwatchDatasource:    defaultGrafanaCloudWatchDatasourceNamePattern,
		propGrafanaElasticsearchDatasource: defaultGrafanaElasticsearchDatasourceNamePattern,
		propKubernetesNamespace:            defaultKubernetesNamespaceNamePattern,
		propKubernetesPod:                  defaultKubernetesPodNamePattern,
		propTraefikServiceName:             defaultTraefikServiceNameNamePattern,
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
		Hostname:                       patterns[propHostname],
		CloudwatchNamespace:            patterns[propCloudwatchNamespace],
		EcsCluster:                     patterns[propEcsCluster],
		EcsService:                     patterns[propEcsService],
		GrafanaCloudWatchDatasource:    patterns[propGrafanaCloudwatchDatasource],
		GrafanaElasticsearchDatasource: patterns[propGrafanaElasticsearchDatasource],
		KubernetesNamespace:            patterns[propKubernetesNamespace],
		KubernetesPod:                  patterns[propKubernetesPod],
		TraefikServiceName:             patterns[propTraefikServiceName],
	}

	return props, nil
}
