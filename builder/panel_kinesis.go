package builder

import "fmt"

func NewPanelKinesisKinsumerMillisecondsBehind(stream MetadataCloudAwsKinesisKinsumer) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GetCwDatasourceNameByClientName(stream.AwsClientName),
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						LineWidth:     2,
						AxisPlacement: "right",
					},
					Min:  "0",
					Unit: "ms",
				},
				Overrides: []PanelFieldConfigOverwrite{},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "",
					Dimensions: map[string]string{
						"StreamName": stream.StreamNameFull,
					},
					Expression: "",
					Id:         "",
					MatchExact: false,
					MetricName: "MillisecondsBehind",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
					Period:     "",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Maximum",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "MillisecondsBehind",
			Type:    "timeseries",
		}
	}
}

func NewPanelKinesisKinsumerMessageCounts(stream MetadataCloudAwsKinesisKinsumer) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GetCwDatasourceNameByClientName(stream.AwsClientName),
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						LineWidth:     2,
						AxisPlacement: "right",
					},
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverwrite{
					NewColorPropertyOverwrite("ReadRecords", "semi-dark-blue"),
					NewColorPropertyOverwrite("FailedRecords", "dark-red"),
				},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "ReadRecords",
					Dimensions: map[string]string{
						"StreamName": stream.StreamNameFull,
					},
					MatchExact: false,
					MetricName: "ReadRecords",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
					Period:     "",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "FailedRecords",
					Dimensions: map[string]string{
						"StreamName": stream.StreamNameFull,
					},
					MatchExact: false,
					MetricName: "FailedRecords",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
					Period:     "",
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Message Counts",
			Type:    "timeseries",
		}
	}
}

func NewPanelKinesisKinsumerReadOperations(stream MetadataCloudAwsKinesisKinsumer) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GetCwDatasourceNameByClientName(stream.AwsClientName),
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						LineWidth:     2,
						AxisPlacement: "right",
					},
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverwrite{
					NewColorPropertyOverwrite("ReadCount Limit", "dark-red"),
				},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Dimensions: map[string]string{
						"StreamName": stream.StreamNameFull,
					},
					Id:         "m0",
					Hide:       true,
					MatchExact: true,
					MetricName: "ReadRecords",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "ReadCount",
					Dimensions: map[string]string{
						"StreamName": stream.StreamNameFull,
					},
					Id:         "m1",
					MatchExact: true,
					MetricName: "ReadCount",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias:      "ReadCount Limit",
					Dimensions: map[string]string{},
					Expression: fmt.Sprintf("%d * 5 * PERIOD(m1) * IF(m1, 1, 1)", stream.OpenShardCount),
					MatchExact: true,
					RefId:      "C",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
				PanelTargetCloudWatch{
					Alias:      "Batch Size",
					Dimensions: map[string]string{},
					Expression: "IF(m0, IF(m1, m0 / m1, 0), 0)",
					RefId:      "D",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Read Operations",
			Type:    "timeseries",
		}
	}
}

func NewPanelKinesisKinsumerProcessDuration(stream MetadataCloudAwsKinesisKinsumer) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GetCwDatasourceNameByClientName(stream.AwsClientName),
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						SpanNulls: true,
					},
					Min:  "0",
					Unit: "ms",
				},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "Maximum",
					Dimensions: map[string]string{
						"StreamName": stream.StreamNameFull,
					},
					Id:         "m0",
					MatchExact: true,
					MetricName: "ProcessDuration",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Maximum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "Average",
					Dimensions: map[string]string{
						"StreamName": stream.StreamNameFull,
					},
					Id:         "m1",
					MatchExact: true,
					MetricName: "ProcessDuration",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Process Duration",
			Type:    "timeseries",
		}
	}
}

func NewPanelKinesisStreamSuccessRate(stream KinesisStreamAware) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GetCwDatasourceNameByClientName(stream.GetClientName()),
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						LineWidth:     2,
						AxisPlacement: "right",
					},
					Min:  "0",
					Unit: "percent",
				},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Dimensions: map[string]string{
						"StreamName": stream.GetStreamNameFull(),
					},
					Expression: "",
					Id:         "m0",
					Hide:       true,
					MatchExact: true,
					MetricName: "GetRecords.Success",
					Namespace:  "AWS/Kinesis",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
				PanelTargetCloudWatch{
					Alias:      "Get records success",
					Dimensions: map[string]string{},
					Expression: "m0 * 100",
					Id:         "m1",
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
				PanelTargetCloudWatch{
					Alias: "Put record success",
					Dimensions: map[string]string{
						"StreamName": stream.GetStreamNameFull(),
					},
					Id:         "m2",
					MatchExact: true,
					MetricName: "PutRecord.Success",
					Namespace:  "AWS/Kinesis",
					RefId:      "C",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
				PanelTargetCloudWatch{
					Dimensions: map[string]string{
						"StreamName": stream.GetStreamNameFull(),
					},
					Id:         "m3",
					MatchExact: true,
					Hide:       true,
					MetricName: "PutRecords.Success",
					Namespace:  "AWS/Kinesis",
					RefId:      "D",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
				PanelTargetCloudWatch{
					Alias:      "Put records success",
					Dimensions: map[string]string{},
					Expression: "m3 * 100",
					Id:         "m4",
					RefId:      "E",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Get / Put Success Rate",
			Type:    "timeseries",
		}
	}
}

func NewPanelKinesisStreamGetRecordsBytes(stream KinesisStreamAware) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GetCwDatasourceNameByClientName(stream.GetClientName()),
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						SpanNulls: true,
					},
					Min:  "0",
					Unit: "decbytes",
				},
				Overrides: []PanelFieldConfigOverwrite{
					NewColorPropertyOverwrite("Limit", "dark-red"),
					NewColorPropertyOverwrite("GetRecordsBytes", "super-light-blue"),
				},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "GetRecordsBytes",
					Dimensions: map[string]string{
						"StreamName": stream.GetStreamNameFull(),
					},
					Expression: "",
					Id:         "m0",
					MatchExact: true,
					MetricName: "GetRecords.Bytes",
					Namespace:  "AWS/Kinesis",
					Period:     "",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias:      "Limit",
					Dimensions: map[string]string{},
					Expression: fmt.Sprintf("%d * 2097152 * PERIOD(m0) * IF(m0, 1, 1)", stream.GetOpenShardCount()),
					MatchExact: true,
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Maximum",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Stream get records - sum (Bytes)",
			Type:    "timeseries",
		}
	}
}

func NewPanelKinesisStreamIncomingDataBytes(stream KinesisStreamAware) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GetCwDatasourceNameByClientName(stream.GetClientName()),
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						SpanNulls: true,
					},
					Min:  "0",
					Unit: "decbytes",
				},
				Overrides: []PanelFieldConfigOverwrite{
					NewColorPropertyOverwrite("Limit", "dark-red"),
					NewColorPropertyOverwrite("IncomingBytes", "super-light-blue"),
				},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "IncomingBytes",
					Dimensions: map[string]string{
						"StreamName": stream.GetStreamNameFull(),
					},
					Id:         "m0",
					MatchExact: true,
					MetricName: "IncomingBytes",
					Namespace:  "AWS/Kinesis",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias:      "Limit",
					Dimensions: map[string]string{},
					Expression: fmt.Sprintf("%d * 1048576 * PERIOD(m0) * IF(m0, 1, 1)", stream.GetOpenShardCount()),
					MatchExact: true,
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Maximum",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Stream incoming data - sum (Bytes)",
			Type:    "timeseries",
		}
	}
}

func NewPanelKinesisStreamIncomingDataCount(stream KinesisStreamAware) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GetCwDatasourceNameByClientName(stream.GetClientName()),
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						SpanNulls: true,
					},
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverwrite{
					NewColorPropertyOverwrite("Limit", "dark-red"),
					NewColorPropertyOverwrite("IncomingRecords", "super-light-blue"),
				},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "IncomingRecords",
					Dimensions: map[string]string{
						"StreamName": stream.GetStreamNameFull(),
					},
					Id:         "m0",
					MatchExact: true,
					MetricName: "IncomingRecords",
					Namespace:  "AWS/Kinesis",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias:      "Limit",
					Dimensions: map[string]string{},
					Expression: fmt.Sprintf("%d * 1000 * PERIOD(m0) * IF(m0, 1, 1)", stream.GetOpenShardCount()),
					MatchExact: true,
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Maximum",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Stream incoming data - sum (Count)",
			Type:    "timeseries",
		}
	}
}

func NewPanelKinesisRecordWriterPutRecordsCount(stream MetadataCloudAwsKinesisRecordWriter) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GetCwDatasourceNameByClientName(stream.AwsClientName),
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						LineWidth:     2,
						AxisPlacement: "right",
					},
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverwrite{
					NewColorPropertyOverwrite("PutRecords", "semi-dark-blue"),
					NewColorPropertyOverwrite("PutRecordsFailure", "dark-red"),
				},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "PutRecords",
					Dimensions: map[string]string{
						"StreamName": stream.StreamName,
					},
					MatchExact: false,
					MetricName: "PutRecords",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
					Period:     "",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "PutRecordsFailure",
					Dimensions: map[string]string{
						"StreamName": stream.StreamName,
					},
					MatchExact: false,
					MetricName: "PutRecordsFailure",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
					Period:     "",
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Put Records Statistics",
			Type:    "timeseries",
		}
	}
}

func NewPanelKinesisRecordWriterPutRecordsBatchSize(stream MetadataCloudAwsKinesisRecordWriter) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GetCwDatasourceNameByClientName(stream.AwsClientName),
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						LineWidth:     2,
						AxisPlacement: "right",
					},
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverwrite{},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "Batch Size",
					Dimensions: map[string]string{
						"StreamName": stream.StreamName,
					},
					MatchExact: false,
					MetricName: "PutRecordsBatchSize",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
					Period:     "",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
				PanelTargetCloudWatch{
					Alias: "PutRecords",
					Dimensions: map[string]string{
						"StreamName": stream.StreamName,
					},
					Hide:       true,
					Id:         "m0",
					MatchExact: true,
					MetricName: "PutRecords",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
					Period:     "",
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias:      "Records Per Shard",
					Dimensions: map[string]string{},
					Expression: fmt.Sprintf("m0 / %d /PERIOD(m0) * IF(m0, 1, 1)", stream.OpenShardCount),
					MatchExact: true,
					RefId:      "C",
					Region:     "default",
					Statistics: []string{
						"Maximum",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Average Batch Size / Records per shards",
			Type:    "timeseries",
		}
	}
}

func NewPanelKinesisStreamRecordSize(stream KinesisStreamAware) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GetCwDatasourceNameByClientName(stream.GetClientName()),
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						SpanNulls: true,
					},
					Min:  "0",
					Unit: "decbytes",
				},
				Overrides: []PanelFieldConfigOverwrite{},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "IncomingBytes",
					Dimensions: map[string]string{
						"StreamName": stream.GetStreamNameFull(),
					},
					Id:         "m0",
					Hide:       true,
					MatchExact: true,
					MetricName: "IncomingBytes",
					Namespace:  "AWS/Kinesis",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "IncomingRecords",
					Dimensions: map[string]string{
						"StreamName": stream.GetStreamNameFull(),
					},
					Id:         "m1",
					Hide:       true,
					MatchExact: true,
					MetricName: "IncomingRecords",
					Namespace:  "AWS/Kinesis",
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias:      "Size",
					Dimensions: map[string]string{},
					Expression: "m0 / m1",
					MatchExact: true,
					RefId:      "C",
					Region:     "default",
					Statistics: []string{
						"Maximum",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Average Record Size (Bytes)",
			Type:    "timeseries",
		}
	}
}
