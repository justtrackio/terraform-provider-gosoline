package builder

import "fmt"

func NewPanelStreamConsumerProcessedCount(consumer MetadataStreamConsumer) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GrafanaCloudWatchDatasourceName,
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
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "Processed",
					Dimensions: map[string]string{
						"Consumer": consumer.Name,
					},
					Expression: "",
					Id:         "m0",
					MatchExact: true,
					MetricName: "ProcessedCount",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
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
						"Consumer": consumer.Name,
					},
					Expression: "",
					Id:         "m1",
					MatchExact: true,
					MetricName: "Error",
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
			Title:   "Processed Count and Errors",
			Type:    "timeseries",
		}
	}
}

func NewPanelStreamConsumerProcessDuration(consumer MetadataStreamConsumer) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GrafanaCloudWatchDatasourceName,
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
					Alias: "Average",
					Dimensions: map[string]string{
						"Consumer": consumer.Name,
					},
					Expression: "",
					Id:         "m0",
					MatchExact: true,
					MetricName: "Duration",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
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

func NewPanelStreamConsumerRetryActions(consumer MetadataStreamConsumer) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GrafanaCloudWatchDatasourceName,
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Custom: PanelFieldConfigDefaultsCustom{
						SpanNulls: true,
					},
					Min: "0",
				},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "Processed",
					Dimensions: map[string]string{
						"Consumer": consumer.Name,
					},
					Expression: "",
					Id:         "m0",
					MatchExact: true,
					MetricName: "RetryGetCount",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
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
						"Consumer": consumer.Name,
					},
					Expression: "",
					Id:         "m1",
					MatchExact: true,
					MetricName: "RetryPutCount",
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
			Title:   fmt.Sprintf("Retry Actions with type: %s", consumer.RetryType),
			Type:    "timeseries",
		}
	}
}
