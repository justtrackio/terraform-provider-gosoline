package builder

import "fmt"

func NewPanelKinesisMillisecondsBehind(stream string) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			Datasource: "CloudWatch",
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
			GridPos: gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "",
					Dimensions: map[string]string{
						"StreamName": stream,
					},
					Expression: "",
					Id:         "",
					MatchExact: false,
					MetricName: "MillisecondsBehind",
					Namespace:  appId.CloudWatchNamespace(),
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

func NewPanelKinesisKinsumerMessageCounts(stream string) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			Datasource: "CloudWatch",
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
			GridPos: gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "ReadRecords",
					Dimensions: map[string]string{
						"StreamName": stream,
					},
					MatchExact: false,
					MetricName: "ReadRecords",
					Namespace:  appId.CloudWatchNamespace(),
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
						"StreamName": stream,
					},
					MatchExact: false,
					MetricName: "FailedRecords",
					Namespace:  appId.CloudWatchNamespace(),
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

func NewPanelKinesisReadOperations(stream string, shardCount int) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			Datasource: "CloudWatch",
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
			GridPos: gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Dimensions: map[string]string{
						"StreamName": stream,
					},
					Id:         "m0",
					Hide:       true,
					MatchExact: true,
					MetricName: "ReadRecords",
					Namespace:  appId.CloudWatchNamespace(),
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "ReadCount",
					Dimensions: map[string]string{
						"StreamName": stream,
					},
					Id:         "m1",
					MatchExact: true,
					MetricName: "ReadCount",
					Namespace:  appId.CloudWatchNamespace(),
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias:      "ReadCount Limit",
					Dimensions: map[string]string{},
					Expression: fmt.Sprintf("%d * 5 * PERIOD(m1) * IF(m1, 1, 1)", shardCount),
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

func NewPanelKinesisSuccessRate(stream string) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			Datasource: "CloudWatch",
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
			GridPos: gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Dimensions: map[string]string{
						"StreamName": stream,
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
						"StreamName": stream,
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
						"StreamName": stream,
					},
					Id:         "m3",
					MatchExact: true,
					Hide:       true,
					MetricName: "PutRecords.SuccessfulRecords",
					Namespace:  "AWS/Kinesis",
					RefId:      "D",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
				PanelTargetCloudWatch{
					Dimensions: map[string]string{
						"StreamName": stream,
					},
					Id:         "m4",
					MatchExact: true,
					Hide:       true,
					MetricName: "PutRecords.TotalRecords",
					Namespace:  "AWS/Kinesis",
					RefId:      "E",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
				PanelTargetCloudWatch{
					Alias:      "Put records successful records",
					Dimensions: map[string]string{},
					Expression: "(m3 / m4) * 100",
					RefId:      "F",
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

func NewPanelKinesisStreamGetRecordsBytes(stream string, shardCount int) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			Datasource: "CloudWatch",
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
			GridPos: gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "GetRecordsBytes",
					Dimensions: map[string]string{
						"StreamName": stream,
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
					Expression: fmt.Sprintf("%d * 2097152 * PERIOD(m0) * IF(m0, 1, 1)", shardCount),
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

func NewPanelKinesisStreamIncomingDataBytes(stream string, shardCount int) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			Datasource: "CloudWatch",
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
			GridPos: gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "IncomingBytes",
					Dimensions: map[string]string{
						"StreamName": stream,
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
					Expression: fmt.Sprintf("%d * 1048576 * PERIOD(m0) * IF(m0, 1, 1)", shardCount),
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

func NewPanelKinesisStreamIncomingDataCount(stream string, shardCount int) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			Datasource: "CloudWatch",
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
			GridPos: gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "IncomingRecords",
					Dimensions: map[string]string{
						"StreamName": stream,
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
					Expression: fmt.Sprintf("%d * 1000 * PERIOD(m0) * IF(m0, 1, 1)", shardCount),
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

func NewPanelKinesisRecordWriterPutRecordsCount(stream string) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			Datasource: "CloudWatch",
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
			GridPos: gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "PutRecords",
					Dimensions: map[string]string{
						"StreamName": stream,
					},
					MatchExact: false,
					MetricName: "PutRecords",
					Namespace:  appId.CloudWatchNamespace(),
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
						"StreamName": stream,
					},
					MatchExact: false,
					MetricName: "PutRecordsFailure",
					Namespace:  appId.CloudWatchNamespace(),
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

func NewPanelKinesisRecordWriterPutRecordsBatchSize(stream string, shardCount int) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			Datasource: "CloudWatch",
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
			GridPos: gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "Batch Size",
					Dimensions: map[string]string{
						"StreamName": stream,
					},
					MatchExact: false,
					MetricName: "PutRecordsBatchSize",
					Namespace:  appId.CloudWatchNamespace(),
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
						"StreamName": stream,
					},
					Hide:       true,
					Id:         "m0",
					MatchExact: true,
					MetricName: "PutRecords",
					Namespace:  appId.CloudWatchNamespace(),
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
					Expression: fmt.Sprintf("m0 / %d /PERIOD(m0) * IF(m0, 1, 1)", shardCount),
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

func NewPanelKinesisStreamRecordSize(stream string) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			Datasource: "CloudWatch",
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
			GridPos: gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "IncomingBytes",
					Dimensions: map[string]string{
						"StreamName": stream,
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
						"StreamName": stream,
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
