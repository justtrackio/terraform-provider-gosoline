package provider

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/justtrackio/terraform-provider-gosoline/builder"
)

type ApplicationDashboardDefinitionData struct {
	Project     types.String `tfsdk:"project"`
	Environment types.String `tfsdk:"environment"`
	Family      types.String `tfsdk:"family"`
	Group       types.String `tfsdk:"group"`
	Application types.String `tfsdk:"application"`
	Containers  types.List   `tfsdk:"containers"`
	Body        types.String `tfsdk:"body"`
}

func (d ApplicationDashboardDefinitionData) AppId() builder.AppId {
	return builder.AppId{
		Project:     d.Project.Value,
		Environment: d.Environment.Value,
		Family:      d.Family.Value,
		Group:       d.Group.Value,
		Application: d.Application.Value,
	}
}

type ApplicationDashboardDefinitionDatasourceType struct{}

func (a *ApplicationDashboardDefinitionDatasourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"project": {
				Type:     types.StringType,
				Required: true,
			},
			"environment": {
				Type:     types.StringType,
				Required: true,
			},
			"family": {
				Type:     types.StringType,
				Required: true,
			},
			"group": {
				Type:     types.StringType,
				Required: true,
			},
			"application": {
				Type:     types.StringType,
				Required: true,
			},
			"containers": {
				Type:     types.ListType{ElemType: types.StringType},
				Required: true,
			},
			"body": {
				Type:     types.StringType,
				Computed: true,
			},
		},
	}, nil
}

func (a *ApplicationDashboardDefinitionDatasourceType) NewDataSource(_ context.Context, provider tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return &ApplicationDashboardDefinitionDataSource{
		metadataReader:       provider.(*GosolineProvider).metadataReader,
		resourceNamePatterns: provider.(*GosolineProvider).resourceNamePatterns,
		orchestrator:         provider.(*GosolineProvider).orchestrator,
	}, nil
}

type ApplicationDashboardDefinitionDataSource struct {
	metadataReader       *builder.MetadataReader
	resourceNamePatterns ResourceNamePatterns
	orchestrator         string
}

func (a *ApplicationDashboardDefinitionDataSource) Read(ctx context.Context, request tfsdk.ReadDataSourceRequest, response *tfsdk.ReadDataSourceResponse) {
	state := &ApplicationDashboardDefinitionData{}

	diags := request.Config.Get(ctx, state)
	response.Diagnostics.Append(diags...)

	var err error
	var metadata *builder.MetadataApplication
	if metadata, err = a.metadataReader.ReadMetadata(state.AppId()); err != nil {
		response.Diagnostics.AddError("can not get metadata", err.Error())
		return
	}

	// Always available
	cloudwatchNamespace := builder.Augment(a.resourceNamePatterns.CloudwatchNamespace, state.AppId())
	grafanaCloudWatchDatasourceName := builder.Augment(a.resourceNamePatterns.GrafanaCloudWatchDatasource, state.AppId())
	grafanaElasticsearchDatasourceName := builder.Augment(a.resourceNamePatterns.GrafanaElasticsearchDatasource, state.AppId())

	// only available when orchestrator is ecs => ECS/EC2-LB+TG related
	var targetGroups []builder.ElbTargetGroup
	var ecsTaskDefinitionName string
	var ecsClusterName string
	var ecsServiceName string

	// only available when orchestrator is kubernetes => Kubernetes/Traefik related
	var kubernetesNamespace string
	var kubernetesPod string
	var traefikServiceName string

	containers := make([]string, 0)
	diags = state.Containers.ElementsAs(ctx, &containers, false)
	response.Diagnostics.Append(diags...)

	switch orchestrator := a.orchestrator; orchestrator {
	case orchestratorEcs:
		ecsClusterName = builder.Augment(a.resourceNamePatterns.EcsCluster, state.AppId())
		ecsServiceName = builder.Augment(a.resourceNamePatterns.EcsService, state.AppId())
		targetGroups, ecsTaskDefinitionName, err = getEc2AndEcsData(ctx, response, ecsClusterName, ecsServiceName)
		if err != nil {
			return
		}
	case orchestratorKubernetes:
		kubernetesNamespace = builder.Augment(a.resourceNamePatterns.KubernetesNamespace, state.AppId())
		kubernetesPod = builder.Augment(a.resourceNamePatterns.KubernetesPod, state.AppId())
		traefikServiceName = builder.Augment(a.resourceNamePatterns.TraefikServiceName, state.AppId())
	}

	resourceNames := builder.ResourceNames{
		CloudwatchNamespace:                cloudwatchNamespace,
		EcsCluster:                         ecsClusterName,
		EcsService:                         ecsServiceName,
		EcsTaskDefinition:                  ecsTaskDefinitionName,
		GrafanaCloudWatchDatasourceName:    grafanaCloudWatchDatasourceName,
		GrafanaElasticsearchDatasourceName: grafanaElasticsearchDatasourceName,
		KubernetesNamespace:                kubernetesNamespace,
		KubernetesPod:                      kubernetesPod,
		TraefikServiceName:                 traefikServiceName,
		TargetGroups:                       targetGroups,
		Containers:                         containers,
	}

	db := builder.NewDashboardBuilder(resourceNames, a.orchestrator)
	db.AddServiceAndTask()
	db.AddPanel(builder.NewPanelRow("Errors & Warnings"))
	db.AddPanel(builder.NewPanelError)
	db.AddPanel(builder.NewPanelWarn)
	db.AddPanel(builder.NewPanelLogs)

	for i := range targetGroups {
		db.AddElbTargetGroup(i)
	}

	db.AddTraefikService()

	for _, route := range metadata.ApiServer.Routes {
		if route.Path == "/health" {
			continue
		}

		db.AddApiServerHandler(route.Method, route.Path)
	}

	for _, consumer := range metadata.Stream.Consumers {
		db.AddStreamConsumer(consumer)
	}

	for _, kinsumer := range metadata.Cloud.Aws.Kinesis.Kinsumers {
		db.AddCloudAwsKinesisKinsumer(kinsumer)
	}

	for _, producer := range metadata.Stream.Producers {
		if !producer.DaemonEnabled {
			continue
		}

		db.AddStreamProducerDaemon(producer)
	}

	for _, writer := range metadata.Cloud.Aws.Kinesis.RecordWriters {
		db.AddCloudAwsKinesisRecordWriter(writer)
	}

	for _, stream := range metadata.Cloud.Aws.Kinesis.Kinsumers {
		db.AddCloudAwsKinesisStream(stream)
	}
	for _, stream := range metadata.Cloud.Aws.Kinesis.RecordWriters {
		db.AddCloudAwsKinesisStream(stream)
	}

	for _, queue := range metadata.Cloud.Aws.Sqs.Queues {
		db.AddCloudAwsSqsQueue(queue)
	}

	for _, table := range metadata.Cloud.Aws.Dynamodb.Tables {
		db.AddDynamoDbTable(table)
	}

	dashboard := db.Build()

	body, err := json.Marshal(dashboard)
	if err != nil {
		response.Diagnostics.AddError("can not create dashboard", err.Error())
	}

	state.Body = types.String{
		Value: string(body),
	}

	diags = response.State.Set(ctx, state)
	response.Diagnostics.Append(diags...)
}

func getEc2AndEcsData(ctx context.Context, response *tfsdk.ReadDataSourceResponse, ecsClusterName string, ecsServiceName string) ([]builder.ElbTargetGroup, string, error) {
	var targetGroups []builder.ElbTargetGroup
	var ecsTaskDefinitionName *string

	ecsClient, err := builder.NewEcsClient(ctx, ecsClusterName, ecsServiceName)
	if err != nil {
		response.Diagnostics.AddError("can not get ecs client", err.Error())
		return nil, "", err
	}

	targetGroups, err = ecsClient.GetElbTargetGroups(ctx)
	if err != nil {
		response.Diagnostics.AddError("can not get target groups", err.Error())
		return nil, "", err
	}

	ecsTaskDefinitionName, err = ecsClient.GetTaskDefinitionName(ctx)
	if err != nil {
		response.Diagnostics.AddError("can not get ecs task definition name", err.Error())
		return nil, "", err
	}

	return targetGroups, *ecsTaskDefinitionName, err
}
