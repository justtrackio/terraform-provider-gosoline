package provider

import (
	"context"
	"encoding/json"
	"sort"

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
	Title       types.String `tfsdk:"title"`
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
			"title": {
				Type:     types.StringType,
				Optional: true,
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
	var resourceNames *builder.ResourceNames

	if metadata, err = a.metadataReader.ReadMetadata(state.AppId()); err != nil {
		response.Diagnostics.AddError("can not get metadata", err.Error())

		return
	}

	if resourceNames, err = a.getResourceNames(ctx, state, response); err != nil {
		return
	}

	db := builder.NewDashboardBuilder(resourceNames, a.orchestrator)
	db.AddServiceAndTask()
	db.AddPanel(builder.NewPanelRow("Errors & Warnings"))
	db.AddPanel(builder.NewPanelError)
	db.AddPanel(builder.NewPanelWarn)
	db.AddPanel(builder.NewPanelLogs)

	a.addHttpServers(metadata, resourceNames, db)

	// Sort stream consumers by name for consistent ordering
	streamConsumers := metadata.Stream.Consumers
	sort.Slice(streamConsumers, func(i, j int) bool {
		return streamConsumers[i].Name < streamConsumers[j].Name
	})
	for _, consumer := range streamConsumers {
		db.AddStreamConsumer(consumer)
	}

	// Sort Kinesis kinsumers by name for consistent ordering
	kinsumers := metadata.Cloud.Aws.Kinesis.Kinsumers
	sort.Slice(kinsumers, func(i, j int) bool {
		return kinsumers[i].Name < kinsumers[j].Name
	})
	for _, kinsumer := range kinsumers {
		db.AddCloudAwsKinesisKinsumer(kinsumer)
	}

	// Sort stream producers by name for consistent ordering
	streamProducers := metadata.Stream.Producers
	sort.Slice(streamProducers, func(i, j int) bool {
		return streamProducers[i].Name < streamProducers[j].Name
	})
	for _, producer := range streamProducers {
		if !producer.DaemonEnabled {
			continue
		}

		db.AddStreamProducerDaemon(producer)
	}

	// Sort Kinesis record writers by stream name for consistent ordering
	recordWriters := metadata.Cloud.Aws.Kinesis.RecordWriters
	sort.Slice(recordWriters, func(i, j int) bool {
		return recordWriters[i].StreamName < recordWriters[j].StreamName
	})
	for _, writer := range recordWriters {
		db.AddCloudAwsKinesisRecordWriter(writer)
	}

	// Use already sorted collections for Kinesis streams
	for _, stream := range kinsumers {
		db.AddCloudAwsKinesisStream(stream)
	}
	for _, stream := range recordWriters {
		db.AddCloudAwsKinesisStream(stream)
	}

	// Sort SQS queues by queue name for consistent ordering
	sqsQueues := metadata.Cloud.Aws.Sqs.Queues
	sort.Slice(sqsQueues, func(i, j int) bool {
		return sqsQueues[i].QueueNameFull < sqsQueues[j].QueueNameFull
	})
	for _, queue := range sqsQueues {
		db.AddCloudAwsSqsQueue(queue)
	}

	// Sort DynamoDB tables by table name for consistent ordering
	dynamoTables := metadata.Cloud.Aws.Dynamodb.Tables
	sort.Slice(dynamoTables, func(i, j int) bool {
		return dynamoTables[i].TableName < dynamoTables[j].TableName
	})
	for _, table := range dynamoTables {
		db.AddDynamoDbTable(table)
	}

	dashboard := db.Build(state.Title.Value)

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

func (a *ApplicationDashboardDefinitionDataSource) getResourceNames(ctx context.Context, state *ApplicationDashboardDefinitionData, response *tfsdk.ReadDataSourceResponse) (*builder.ResourceNames, error) {
	var err error

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
	var kubernetesDeployment string
	var traefikServiceName string

	containers := make([]string, 0)
	diags := state.Containers.ElementsAs(ctx, &containers, false)
	response.Diagnostics.Append(diags...)

	switch orchestrator := a.orchestrator; orchestrator {
	case orchestratorEcs:
		ecsClusterName = builder.Augment(a.resourceNamePatterns.EcsCluster, state.AppId())
		ecsServiceName = builder.Augment(a.resourceNamePatterns.EcsService, state.AppId())
		targetGroups, ecsTaskDefinitionName, err = a.getEc2AndEcsData(ctx, response, ecsClusterName, ecsServiceName)
		if err != nil {
			return nil, err
		}
	case orchestratorKubernetes:
		kubernetesNamespace = builder.Augment(a.resourceNamePatterns.KubernetesNamespace, state.AppId())
		kubernetesPod = builder.Augment(a.resourceNamePatterns.KubernetesPod, state.AppId())
		kubernetesDeployment = kubernetesPod // KubernetesPod pattern is actually the deployment name
		traefikServiceName = builder.Augment(a.resourceNamePatterns.TraefikServiceName, state.AppId())
	}

	resourceNames := &builder.ResourceNames{
		CloudwatchNamespace:                cloudwatchNamespace,
		EcsCluster:                         ecsClusterName,
		EcsService:                         ecsServiceName,
		EcsTaskDefinition:                  ecsTaskDefinitionName,
		Environment:                        state.Environment.Value,
		GrafanaCloudWatchDatasourceName:    grafanaCloudWatchDatasourceName,
		GrafanaElasticsearchDatasourceName: grafanaElasticsearchDatasourceName,
		KubernetesDeployment:               kubernetesDeployment,
		KubernetesNamespace:                kubernetesNamespace,
		KubernetesPod:                      kubernetesPod,
		TraefikServiceName:                 traefikServiceName,
		TargetGroups:                       targetGroups,
		Containers:                         containers,
	}

	return resourceNames, nil
}

func (a *ApplicationDashboardDefinitionDataSource) getEc2AndEcsData(ctx context.Context, response *tfsdk.ReadDataSourceResponse, ecsClusterName string, ecsServiceName string) ([]builder.ElbTargetGroup, string, error) {
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

func (a *ApplicationDashboardDefinitionDataSource) addHttpServers(metadata *builder.MetadataApplication, resourceNames *builder.ResourceNames, db *builder.DashboardBuilder) {
	if len(metadata.HttpServers) == 0 {
		return
	}

	for i := range resourceNames.TargetGroups {
		db.AddElbTargetGroup(i)
	}

	db.AddTraefikService()

	// Sort HTTP servers by name for consistent ordering
	httpServers := metadata.HttpServers
	sort.Slice(httpServers, func(i, j int) bool {
		return httpServers[i].Name < httpServers[j].Name
	})

	for _, server := range httpServers {
		// Sort handlers by path and method for consistent ordering
		handlers := server.Handlers
		sort.Slice(handlers, func(i, j int) bool {
			if handlers[i].Path != handlers[j].Path {
				return handlers[i].Path < handlers[j].Path
			}
			return handlers[i].Method < handlers[j].Method

		})

		for _, route := range handlers {
			if route.Path == "/health" {
				continue
			}

			db.AddHttpServerHandler(server.Name, route)
		}
	}
}
