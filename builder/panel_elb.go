package builder

func NewPanelElbRequestCount(targetGroup ElbTargetGroup) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			Datasource: "CloudWatch",
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverwrite{
					NewColorPropertyOverwrite("Requests", "semi-dark-blue"),
				},
			},
			GridPos: gridPos,
			Targets: []PanelTarget{
				{
					Alias: "Requests",
					Dimensions: map[string]string{
						"TargetGroup":  targetGroup.TargetGroup,
						"LoadBalancer": targetGroup.LoadBalancer,
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
			Title: "Request Count",
			Type:  "timeseries",
		}
	}
}

func NewPanelElbResponseTime(targetGroup ElbTargetGroup) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			Datasource: "CloudWatch",
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Min:  "0",
					Unit: "s",
				},
				Overrides: []PanelFieldConfigOverwrite{
					NewColorPropertyOverwrite("Requests", "semi-dark-blue"),
				},
			},
			GridPos: gridPos,
			Targets: []PanelTarget{
				{
					Alias: "Response Time",
					Dimensions: map[string]string{
						"TargetGroup":  targetGroup.TargetGroup,
						"LoadBalancer": targetGroup.LoadBalancer,
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
			Title: "Response Time",
			Type:  "timeseries",
		}
	}
}

func NewPanelElbHttpStatus(targetGroup ElbTargetGroup) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			Datasource: "CloudWatch",
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
			GridPos: gridPos,
			Targets: []PanelTarget{
				{
					Alias: "HTTP 2XX",
					Dimensions: map[string]string{
						"TargetGroup":  targetGroup.TargetGroup,
						"LoadBalancer": targetGroup.LoadBalancer,
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
				{
					Alias: "HTTP 3XX",
					Dimensions: map[string]string{
						"TargetGroup":  targetGroup.TargetGroup,
						"LoadBalancer": targetGroup.LoadBalancer,
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
				{
					Alias: "HTTP 4XX",
					Dimensions: map[string]string{
						"TargetGroup":  targetGroup.TargetGroup,
						"LoadBalancer": targetGroup.LoadBalancer,
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
				{
					Alias: "HTTP 5XX",
					Dimensions: map[string]string{
						"TargetGroup":  targetGroup.TargetGroup,
						"LoadBalancer": targetGroup.LoadBalancer,
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
			Title: "HTTP Status Overview",
			Type:  "timeseries",
		}
	}
}

func NewPanelElbHealthyHosts(targetGroup ElbTargetGroup) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			Datasource: "CloudWatch",
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverwrite{},
			},
			GridPos: gridPos,
			Targets: []PanelTarget{
				{
					Alias: "Hosts",
					Dimensions: map[string]string{
						"TargetGroup":  targetGroup.TargetGroup,
						"LoadBalancer": targetGroup.LoadBalancer,
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
			Title: "Healthy Hosts",
			Type:  "timeseries",
		}
	}
}

func NewPanelElbRequestCountPerTarget(targetGroup ElbTargetGroup) PanelFactory {
	return func(appId AppId, gridPos PanelGridPos) Panel {
		return Panel{
			Datasource: "CloudWatch",
			FieldConfig: PanelFieldConfig{
				Defaults: PanelFieldConfigDefaults{
					Min: "0",
				},
				Overrides: []PanelFieldConfigOverwrite{},
			},
			GridPos: gridPos,
			Targets: []PanelTarget{
				{
					Alias: "Requests",
					Dimensions: map[string]string{
						"TargetGroup":  targetGroup.TargetGroup,
						"LoadBalancer": targetGroup.LoadBalancer,
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
			Title: "Request Counts Per Target",
			Type:  "timeseries",
		}
	}
}
