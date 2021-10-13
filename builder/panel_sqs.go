package builder

func NewPanelSqsMessagesVisible(queue string) PanelFactory {
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
			Targets: []PanelTarget{
				{
					Alias: "",
					Dimensions: map[string]string{
						"QueueName": queue,
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
			Title: "Messages In Queue",
			Type:  "timeseries",
		}
	}
}

func NewPanelSqsTraffic(queue string) PanelFactory {
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
			Targets: []PanelTarget{
				{
					Alias: "",
					Dimensions: map[string]string{
						"QueueName": queue,
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
				{
					Alias: "",
					Dimensions: map[string]string{
						"QueueName": queue,
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
				{
					Alias: "",
					Dimensions: map[string]string{
						"QueueName": queue,
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
			Title: "Traffic",
			Type:  "timeseries",
		}
	}
}
