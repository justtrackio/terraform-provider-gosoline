package builder

import "fmt"

func NewPanelStreamConsumerProcessedCount(consumer string) PanelFactory {
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
					NewColorPropertyOverwrite("Processed", "super-light-blue"),
					NewColorPropertyOverwrite("Error", "dark-red"),
				},
			},
			GridPos: gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "Processed",
					Dimensions: map[string]string{
						"Consumer": consumer,
					},
					Expression: "",
					Id:         "m0",
					MatchExact: true,
					MetricName: "ProcessedCount",
					Namespace:  appId.CloudWatchNamespace(),
					Period:     "",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "Error",
					Dimensions: map[string]string{
						"Consumer": consumer,
					},
					Expression: "",
					Id:         "m1",
					MatchExact: true,
					MetricName: "Error",
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
			Title:   "Processed Count and Errors",
			Type:    "timeseries",
		}
	}
}

func NewPanelStreamConsumerProcessDuration(consumer string) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			Datasource: "CloudWatch",
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						SpanNulls: true,
					},
					Min:  "0",
					Unit: "ms",
				},
			},
			GridPos: gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "Average",
					Dimensions: map[string]string{
						"Consumer": consumer,
					},
					Expression: "",
					Id:         "m0",
					MatchExact: true,
					MetricName: "Duration",
					Namespace:  appId.CloudWatchNamespace(),
					Period:     "",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Duration per consume operation",
			Type:    "timeseries",
		}
	}
}

func NewPanelStreamConsumerRetryActions(consumer string, retryType string) PanelFactory {
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
			},
			GridPos: gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "Processed",
					Dimensions: map[string]string{
						"Consumer": consumer,
					},
					Expression: "",
					Id:         "m0",
					MatchExact: true,
					MetricName: "RetryGetCount",
					Namespace:  appId.CloudWatchNamespace(),
					Period:     "",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "Error",
					Dimensions: map[string]string{
						"Consumer": consumer,
					},
					Expression: "",
					Id:         "m1",
					MatchExact: true,
					MetricName: "RetryPutCount",
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
			Title:   fmt.Sprintf("Retry Actions with type: %s", retryType),
			Type:    "timeseries",
		}
	}
}
