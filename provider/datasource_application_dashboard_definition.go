package provider

import (
	"context"
	"encoding/json"
	"strings"

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

func (a *ApplicationDashboardDefinitionDatasourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
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

func (a *ApplicationDashboardDefinitionDatasourceType) NewDataSource(ctx context.Context, prov tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return &ApplicationDashboardDefinitionDataSource{
		metadataReader: prov.(*GosolineProvider).metadataReader,
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
	var metadata *builder.ApplicationMetadata
	var ecsClient *builder.EcsClient
	var targetGroups []builder.ElbTargetGroup

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
	db.AddPanel(builder.NewPanelRow("Errors & Warnings", false))
	db.AddPanel(builder.NewPanelError)
	db.AddPanel(builder.NewPanelWarn)

	for _, targetGroup := range targetGroups {
		db.AddElbTargetGroup(targetGroup)
	}

	for _, route := range metadata.ApiServer.Routes {
		if route.Path == "/health" {
			continue
		}

		db.AddApiServerHandler(route.Method, route.Path)
	}

	for _, queue := range metadata.Cloud.Aws.Sqs.Queues {
		queue = strings.ReplaceAll(queue, "dev", "prod")
		db.AddSqsQueue(queue)
	}

	for _, table := range metadata.Cloud.Aws.Dynamodb.Tables {
		table = strings.ReplaceAll(table, "dev", "prod")
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
