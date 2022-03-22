package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/justtrackio/terraform-provider-gosoline/builder"
)

type ApplicationDashboardDefinitionData struct {
	Project     types.String `tfsdk:"project"`
	Environment types.String `tfsdk:"environment"`
	Family      types.String `tfsdk:"family"`
	Application types.String `tfsdk:"application"`
	Body        types.String `tfsdk:"body"`
}

func (d ApplicationDashboardDefinitionData) AppId() builder.AppId {
	return builder.AppId{
		Project:     d.Project.Value,
		Environment: d.Environment.Value,
		Family:      d.Family.Value,
		Application: d.Application.Value,
	}
}

type ApplicationDashboardDefinitionDatasourceType struct {
}

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
			"application": {
				Type:     types.StringType,
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
		metadataReader: provider.(*GosolineProvider).metadataReader,
	}, nil
}

type ApplicationDashboardDefinitionDataSource struct {
	metadataReader *builder.MetadataReader
}

func (a *ApplicationDashboardDefinitionDataSource) Read(ctx context.Context, request tfsdk.ReadDataSourceRequest, response *tfsdk.ReadDataSourceResponse) {
	state := &ApplicationDashboardDefinitionData{}

	diags := request.Config.Get(ctx, state)
	response.Diagnostics.Append(diags...)

	var err error
	var metadata *builder.MetadataApplication
	var ecsClient *builder.EcsClient
	var targetGroups []builder.ElbTargetGroup
	var kinesisData map[string]int

	if metadata, err = a.metadataReader.ReadMetadata(state.AppId()); err != nil {
		response.Diagnostics.AddError("can not get metadata", err.Error())
		return
	}

	if ecsClient, err = builder.NewEcsClient(ctx, state.AppId()); err != nil {
		response.Diagnostics.AddError("can not get ecs client", err.Error())
		return
	}

	if targetGroups, err = ecsClient.GetElbTargetGroups(ctx); err != nil {
		response.Diagnostics.AddError("can not get target groups", err.Error())
		return
	}

	db := builder.NewDashboardBuilder(state.AppId())
	db.AddEcs()
	db.AddPanel(builder.NewPanelRow("Errors & Warnings"))
	db.AddPanel(builder.NewPanelError)
	db.AddPanel(builder.NewPanelWarn)
	db.AddPanel(builder.NewPanelLogs)

	for _, targetGroup := range targetGroups {
		db.AddElbTargetGroup(targetGroup)
	}

	for _, route := range metadata.ApiServer.Routes {
		if route.Path == "/health" {
			continue
		}

		db.AddApiServerHandler(route.Method, route.Path)
	}

	for _, consumer := range metadata.Stream.Consumers {
		db.AddStreamConsumer(consumer)
	}

	if kinesisData, err = a.gatherKinesisData(ctx, metadata.Cloud.Aws.Kinesis); err != nil {
		response.Diagnostics.AddError(fmt.Sprintf("can not get kinesis data"), err.Error())
		return
	}

	for _, kinsumer := range metadata.Cloud.Aws.Kinesis.Kinsumers {
		db.AddCloudAwsKinesisKinsumer(kinsumer.StreamNameFull, kinesisData[kinsumer.StreamNameFull])
	}

	for _, producer := range metadata.Stream.Producers {
		if !producer.DaemonEnabled {
			continue
		}

		db.AddStreamProducerDaemon(producer.Name)
	}

	for _, writer := range metadata.Cloud.Aws.Kinesis.RecordWriters {
		db.AddCloudAwsKinesisRecordWriter(writer.StreamName, kinesisData[writer.StreamName])
	}

	for streamName, shardCount := range kinesisData {
		db.AddCloudAwsKinesisStream(streamName, shardCount)
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

func (a *ApplicationDashboardDefinitionDataSource) gatherKinesisData(ctx context.Context, kinesis builder.MetadataCloudAwsKinesis) (map[string]int, error) {
	var err error
	var kinesisClient *builder.KinesisClient
	var kinesisShardCount int

	if kinesisClient, err = builder.NewKinesisClient(ctx); err != nil {
		return nil, fmt.Errorf("can not get kinesis client: %w", err)
	}

	streams := make(map[string]int)

	for _, kinsumer := range kinesis.Kinsumers {
		if kinesisShardCount, err = kinesisClient.GetShardCount(ctx, kinsumer.StreamNameFull); err != nil {
			return nil, fmt.Errorf("can not get kinesis shard count for stream %s: %w", kinsumer.StreamNameFull, err)
		}

		streams[kinsumer.StreamNameFull] = kinesisShardCount
	}

	for _, writer := range kinesis.RecordWriters {
		if kinesisShardCount, err = kinesisClient.GetShardCount(ctx, writer.StreamName); err != nil {
			return nil, fmt.Errorf("can not get kinesis shard count for stream %s: %w", writer.StreamName, err)
		}

		streams[writer.StreamName] = kinesisShardCount
	}

	return streams, nil
}
