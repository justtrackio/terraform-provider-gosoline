package builder

func NewPanelError(appId AppId, gridPos PanelGridPos) Panel {
	return Panel{
		Datasource: "CloudWatch",
		FieldConfig: PanelFieldConfig{
			Defaults: PanelFieldConfigDefaults{
				Custom: PanelFieldConfigDefaultsCustom{
					AxisPlacement: "right",
					LineWidth:     2,
				},
			},
			Overrides: []PanelFieldConfigOverwrite{
				NewColorPropertyOverwrite("Errors", "dark-red"),
			},
		},
		GridPos: gridPos,
		Targets: []PanelTarget{
			{
				Alias:      "Errors",
				Dimensions: map[string]string{},
				MatchExact: false,
				MetricName: "error",
				Namespace:  appId.CloudWatchNamespace(),
				Region:     "default",
				Statistics: []string{
					"Sum",
				},
			},
		},
		Title: "Errors",
		Type:  "timeseries",
	}
}

func NewPanelWarn(appId AppId, gridPos PanelGridPos) Panel {
	return Panel{
		Datasource: "CloudWatch",
		FieldConfig: PanelFieldConfig{
			Defaults: PanelFieldConfigDefaults{
				Custom: PanelFieldConfigDefaultsCustom{
					AxisPlacement: "right",
					LineWidth:     2,
				},
			},
			Overrides: []PanelFieldConfigOverwrite{
				NewColorPropertyOverwrite("Warnings", "dark-yellow"),
			},
		},
		GridPos: gridPos,
		Targets: []PanelTarget{
			{
				Alias:      "Warnings",
				Dimensions: map[string]string{},
				MetricName: "warn",
				Namespace:  appId.CloudWatchNamespace(),
				Region:     "default",
				Statistics: []string{
					"Sum",
				},
			},
		},
		Title: "Warnings",
		Type:  "timeseries",
	}
}
