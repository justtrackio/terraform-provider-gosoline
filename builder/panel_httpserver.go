package builder

func NewPanelHttpServerRequestCount(serverName string, handler MetadataHttpServerHandler) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GrafanaCloudWatchDatasourceName,
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverride{
					NewColorPropertyOverride("Requests", "semi-dark-blue", ""),
				},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "Requests",
					Dimensions: map[string]string{
						"Method":     handler.Method,
						"Path":       handler.Path,
						"ServerName": serverName,
					},
					MatchExact: true,
					MetricName: "HttpRequestCountPerRoute",
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

func NewPanelHttpServerResponseTime(serverName string, handler MetadataHttpServerHandler) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GrafanaCloudWatchDatasourceName,
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Min:  "0",
					Unit: "ms",
				},
				Overrides: []PanelFieldConfigOverride{
					NewColorPropertyOverride("Requests", "semi-dark-blue", ""),
				},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "Response Time",
					Dimensions: map[string]string{
						"Method":     handler.Method,
						"Path":       handler.Path,
						"ServerName": serverName,
					},
					MatchExact: true,
					MetricName: "HttpRequestResponseTimePerRoute",
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

func NewPanelHttpServerHttpStatus(serverName string, handler MetadataHttpServerHandler) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GrafanaCloudWatchDatasourceName,
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverride{
					NewColorPropertyOverride("HTTP 2XX", "semi-dark-green", ""),
					NewColorPropertyOverride("HTTP 3XX", "semi-dark-yellow", ""),
					NewColorPropertyOverride("HTTP 4XX", "semi-dark-orange", ""),
					NewColorPropertyOverride("HTTP 5XX", "dark-red", ""),
				},
			},
			GridPos: settings.gridPos,
			Targets: []interface{}{
				PanelTargetCloudWatch{
					Alias: "HTTP 2XX",
					Dimensions: map[string]string{
						"Method":     handler.Method,
						"Path":       handler.Path,
						"ServerName": serverName,
					},
					MatchExact: true,
					MetricName: "HttpStatus2XXPerRoute",
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
						"Method":     handler.Method,
						"Path":       handler.Path,
						"ServerName": serverName,
					},
					MatchExact: true,
					MetricName: "HttpStatus3XXPerRoute",
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
						"Method":     handler.Method,
						"Path":       handler.Path,
						"ServerName": serverName,
					},
					MatchExact: true,
					MetricName: "HttpStatus4XXPerRoute",
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
						"Method":     handler.Method,
						"Path":       handler.Path,
						"ServerName": serverName,
					},
					MatchExact: true,
					MetricName: "HttpStatus5XXPerRoute",
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
