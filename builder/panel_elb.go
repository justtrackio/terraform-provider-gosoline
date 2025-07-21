package builder

func NewPanelElbRequestCount(targetGroupIndex int) PanelFactory {
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
			Targets: []any{
				PanelTargetCloudWatch{
					Alias: "Requests",
					Dimensions: map[string]string{
						"TargetGroup":  settings.resourceNames.TargetGroups[targetGroupIndex].TargetGroup,
						"LoadBalancer": settings.resourceNames.TargetGroups[targetGroupIndex].LoadBalancer,
					},
					MatchExact: true,
					MetricName: "RequestCount",
					Namespace:  "AWS/ApplicationELB",
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

func NewPanelElbResponseTime(targetGroupIndex int) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GrafanaCloudWatchDatasourceName,
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Min:  "0",
					Unit: "s",
				},
				Overrides: []PanelFieldConfigOverride{
					NewColorPropertyOverride("Requests", "semi-dark-blue", ""),
				},
			},
			GridPos: settings.gridPos,
			Targets: []any{
				PanelTargetCloudWatch{
					Alias: "Response Time",
					Dimensions: map[string]string{
						"TargetGroup":  settings.resourceNames.TargetGroups[targetGroupIndex].TargetGroup,
						"LoadBalancer": settings.resourceNames.TargetGroups[targetGroupIndex].LoadBalancer,
					},
					MatchExact: true,
					MetricName: "TargetResponseTime",
					Namespace:  "AWS/ApplicationELB",
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

func NewPanelElbHttpStatus(targetGroupIndex int) PanelFactory {
	return func(settings PanelSettings) Panel {
		targetGroup := settings.resourceNames.TargetGroups[targetGroupIndex].TargetGroup
		loadBalancer := settings.resourceNames.TargetGroups[targetGroupIndex].LoadBalancer

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
			Targets: []any{
				PanelTargetCloudWatch{
					Alias: "HTTP 2XX",
					Dimensions: map[string]string{
						"TargetGroup":  targetGroup,
						"LoadBalancer": loadBalancer,
					},
					MatchExact: true,
					MetricName: "HTTPCode_Target_2XX_Count",
					Namespace:  "AWS/ApplicationELB",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "HTTP 3XX",
					Dimensions: map[string]string{
						"TargetGroup":  targetGroup,
						"LoadBalancer": loadBalancer,
					},
					MatchExact: true,
					MetricName: "HTTPCode_Target_3XX_Count",
					Namespace:  "AWS/ApplicationELB",
					RefId:      "B",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "HTTP 4XX",
					Dimensions: map[string]string{
						"TargetGroup":  targetGroup,
						"LoadBalancer": loadBalancer,
					},
					MatchExact: true,
					MetricName: "HTTPCode_Target_4XX_Count",
					Namespace:  "AWS/ApplicationELB",
					RefId:      "C",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
				PanelTargetCloudWatch{
					Alias: "HTTP 5XX",
					Dimensions: map[string]string{
						"TargetGroup":  targetGroup,
						"LoadBalancer": loadBalancer,
					},
					MatchExact: true,
					MetricName: "HTTPCode_Target_5XX_Count",
					Namespace:  "AWS/ApplicationELB",
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

func NewPanelElbHealthyHosts(targetGroupIndex int) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GrafanaCloudWatchDatasourceName,
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverride{},
			},
			GridPos: settings.gridPos,
			Targets: []any{
				PanelTargetCloudWatch{
					Alias: "Hosts",
					Dimensions: map[string]string{
						"TargetGroup":  settings.resourceNames.TargetGroups[targetGroupIndex].TargetGroup,
						"LoadBalancer": settings.resourceNames.TargetGroups[targetGroupIndex].LoadBalancer,
					},
					MatchExact: true,
					MetricName: "HealthyHostCount",
					Namespace:  "AWS/ApplicationELB",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Average",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Healthy Hosts",
			Type:    "timeseries",
		}
	}
}

func NewPanelElbRequestCountPerTarget(targetGroupIndex int) PanelFactory {
	return func(settings PanelSettings) Panel {
		return Panel{
			Datasource: settings.resourceNames.GrafanaCloudWatchDatasourceName,
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverride{},
			},
			GridPos: settings.gridPos,
			Targets: []any{
				PanelTargetCloudWatch{
					Alias: "Requests",
					Dimensions: map[string]string{
						"TargetGroup":  settings.resourceNames.TargetGroups[targetGroupIndex].TargetGroup,
						"LoadBalancer": settings.resourceNames.TargetGroups[targetGroupIndex].LoadBalancer,
					},
					MatchExact: true,
					MetricName: "RequestCountPerTarget",
					Namespace:  "AWS/ApplicationELB",
					RefId:      "A",
					Region:     "default",
					Statistics: []string{
						"Sum",
					},
				},
			},
			Options: &PanelOptionsCloudWatch{},
			Title:   "Request Counts Per Target",
			Type:    "timeseries",
		}
	}
}
