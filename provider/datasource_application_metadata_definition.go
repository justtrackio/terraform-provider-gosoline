package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/justtrackio/terraform-provider-gosoline/builder"
)

type ApplicationMetadataDefinitionData struct {
	Project     types.String `tfsdk:"project"`
	Environment types.String `tfsdk:"environment"`
	Family      types.String `tfsdk:"family"`
	Group       types.String `tfsdk:"group"`
	Application types.String `tfsdk:"application"`
	Metadata    types.Object `tfsdk:"metadata"`
}

func (d ApplicationMetadataDefinitionData) AppId() builder.AppId {
	return builder.AppId{
		Project:     d.Project.Value,
		Environment: d.Environment.Value,
		Family:      d.Family.Value,
		Group:       d.Group.Value,
		Application: d.Application.Value,
	}
}

type ApplicationMetadataDefinitionDatasourceType struct{}

func (a *ApplicationMetadataDefinitionDatasourceType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
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
			"metadata": {
				Type: types.ObjectType{
					AttrTypes: builder.MetadataApplicationAttrTypes(),
				},
				Computed: true,
			},
		},
	}, nil
}

func (a *ApplicationMetadataDefinitionDatasourceType) NewDataSource(_ context.Context, provider tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	return &ApplicationMetadataDefinitionDatasource{
		metadataReader: provider.(*GosolineProvider).metadataReader,
	}, nil
}

type ApplicationMetadataDefinitionDatasource struct {
	metadataReader *builder.MetadataReader
}

func (a *ApplicationMetadataDefinitionDatasource) Read(ctx context.Context, request tfsdk.ReadDataSourceRequest, response *tfsdk.ReadDataSourceResponse) {
	state := &ApplicationMetadataDefinitionData{}

	diags := request.Config.Get(ctx, state)
	response.Diagnostics.Append(diags...)

	var err error
	var metadata *builder.MetadataApplication

	if metadata, err = a.metadataReader.ReadMetadata(state.AppId()); err != nil {
		response.Diagnostics.AddError("can not get metadata", err.Error())
		return
	}

	state.Metadata = metadata.ToValue()

	diags = response.State.Set(ctx, state)
	response.Diagnostics.Append(diags...)
}
