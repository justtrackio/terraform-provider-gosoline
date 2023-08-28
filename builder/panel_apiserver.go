package builder

func NewPanelApiServerRequestCount(path string) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GrafanaCloudWatchDatasourceName,
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverwrite{
					NewColorPropertyOverwrite("Requests", "semi-dark-blue"),
				},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "Requests",
					Dimensions: map[string]string{
						"path": path,
					},
					MatchExact: true,
					MetricName: "ApiRequestCount",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Request Count",
			Type:    "timeseries",
		}
	}
}

func NewPanelApiServerResponseTime(path string) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GrafanaCloudWatchDatasourceName,
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Min:  "0",
					Unit: "ms",
				},
				Overrides: []PanelFieldConfigOverwrite{
					NewColorPropertyOverwrite("Requests", "semi-dark-blue"),
				},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "Response Time",
					Dimensions: map[string]string{
						"path": path,
					},
					MatchExact: true,
					MetricName: "ApiRequestResponseTime",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Response Time",
			Type:    "timeseries",
		}
	}
}

func NewPanelApiServerHttpStatus(path string) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GrafanaCloudWatchDatasourceName,
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverwrite{
					NewColorPropertyOverwrite("HTTP 2XX", "semi-dark-green"),
					NewColorPropertyOverwrite("HTTP 3XX", "semi-dark-yellow"),
					NewColorPropertyOverwrite("HTTP 4XX", "semi-dark-orange"),
					NewColorPropertyOverwrite("HTTP 5XX", "dark-red"),
				},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "HTTP 2XX",
					Dimensions: map[string]string{
						"path": path,
					},
					MatchExact: true,
					MetricName: "ApiStatus2XX",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "HTTP 3XX",
					Dimensions: map[string]string{
						"path": path,
					},
					MatchExact: true,
					MetricName: "ApiStatus3XX",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "HTTP 4XX",
					Dimensions: map[string]string{
						"path": path,
					},
					MatchExact: true,
					MetricName: "ApiStatus4XX",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
					RefId:      "C",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "HTTP 5XX",
					Dimensions: map[string]string{
						"path": path,
					},
					MatchExact: true,
					MetricName: "ApiStatus5XX",
					Namespace:  settings.resourceNames.CloudwatchNamespace,
					RefId:      "D",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "HTTP Status Overview",
			Type:    "timeseries",
		}
	}
}
