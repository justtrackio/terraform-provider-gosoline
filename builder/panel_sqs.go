package builder

func NewPanelSqsMessagesVisible(queue MetadataCloudAwsSqsQueue) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GrafanaCloudWatchDatasourceName,
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						SpanNulls:     true,
						LineWidth:     2,
						AxisPlacement: "right",
					},
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverride{},
			},
			GridPos: settings.gridPos,
			Targets: []any{
				PanelTargetCloudWatch{
					Alias: "",
					Dimensions: map[string]string{
						"QueueName": queue.QueueNameFull,
					},
					Expression: "",
					Id:         "",
					MatchExact: false,
					MetricName: "ApproximateNumberOfMessagesVisible",
					Namespace:  "AWS/SQS",
					Period:     "",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Maximum",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Messages In Queue",
			Type:    "timeseries",
		}
	}
}

func NewPanelSqsTraffic(queue MetadataCloudAwsSqsQueue) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GrafanaCloudWatchDatasourceName,
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						SpanNulls:     true,
						LineWidth:     2,
						AxisPlacement: "right",
					},
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverride{},
			},
			GridPos: settings.gridPos,
			Targets: []any{
				PanelTargetCloudWatch{
					Alias: "",
					Dimensions: map[string]string{
						"QueueName": queue.QueueNameFull,
					},
					Expression: "",
					Id:         "",
					MatchExact: false,
					MetricName: "NumberOfMessagesSent",
					Namespace:  "AWS/SQS",
					Period:     "",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "",
					Dimensions: map[string]string{
						"QueueName": queue.QueueNameFull,
					},
					Expression: "",
					Id:         "",
					MatchExact: false,
					MetricName: "NumberOfMessagesReceived",
					Namespace:  "AWS/SQS",
					Period:     "",
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "",
					Dimensions: map[string]string{
						"QueueName": queue.QueueNameFull,
					},
					Expression: "",
					Id:         "",
					MatchExact: false,
					MetricName: "NumberOfMessagesDeleted",
					Namespace:  "AWS/SQS",
					Period:     "",
					RefId:      "C",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Traffic",
			Type:    "timeseries",
		}
	}
}

func NewPanelSqsMessageSize(queue MetadataCloudAwsSqsQueue) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GrafanaCloudWatchDatasourceName,
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						SpanNulls:     true,
						LineWidth:     2,
						AxisPlacement: "right",
					},
					Min:  "0",
					Unit: "bytes",
				},
				Overrides: []PanelFieldConfigOverride{},
			},
			GridPos: settings.gridPos,
			Targets: []any{
				PanelTargetCloudWatch{
					Alias: "Average",
					Dimensions: map[string]string{
						"QueueName": queue.QueueNameFull,
					},
					MatchExact: false,
					MetricName: "SentMessageSize",
					Namespace:  "AWS/SQS",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
				PanelTargetCloudWatch{
					Alias: "Maximum",
					Dimensions: map[string]string{
						"QueueName": queue.QueueNameFull,
					},
					MatchExact: false,
					MetricName: "SentMessageSize",
					Namespace:  "AWS/SQS",
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Maximum",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Message Size",
			Type:    "timeseries",
		}
	}
}
