package builder

func NewPanelStreamProducerDaemonSizes(producer string) PanelFactory {
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
						"ProducerDaemon": producer,
					},
					MatchExact: false,
					MetricName: "BatchSize",
					Namespace:  appId.CloudWatchNamespace(),
					Period:     "",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				}, PanelTargetCloudWatch{
					Alias: "Aggregate Size",
					Dimensions: map[string]string{
						"ProducerDaemon": producer,
					},
					MatchExact: false,
					MetricName: "AggregateSize",
					Namespace:  appId.CloudWatchNamespace(),
					Period:     "",
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Average Batch Size / Aggregation Size",
			Type:    "timeseries",
		}
	}
}

func NewPanelStreamProducerMessageCount(producer string) PanelFactory {
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
					Alias: "Message Count",
					Dimensions: map[string]string{
						"ProducerDaemon": producer,
					},
					MatchExact: false,
					MetricName: "MessageCount",
					Namespace:  appId.CloudWatchNamespace(),
					Period:     "",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Message Count",
			Type:    "timeseries",
		}
	}
}
