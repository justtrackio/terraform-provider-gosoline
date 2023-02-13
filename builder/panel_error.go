package builder

func NewPanelError(resourceNames ResourceNames, gridPos PanelGridPos) Panel {
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
		Targets: []interface{}{
			PanelTargetCloudWatch{
				Alias:      "Errors",
				Dimensions: map[string]string{},
				MatchExact: false,
				MetricName: "error",
				Namespace:  resourceNames.CloudwatchNamespace,
				Region:     "default",
				Statistics: []string{
					"Sum",
				},
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "Errors",
		Type:    "timeseries",
	}
}

func NewPanelWarn(resourceNames ResourceNames, gridPos PanelGridPos) Panel {
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
		Targets: []interface{}{
			PanelTargetCloudWatch{
				Alias:      "Warnings",
				Dimensions: map[string]string{},
				MetricName: "warn",
				Namespace:  resourceNames.CloudwatchNamespace,
				Region:     "default",
				Statistics: []string{
					"Sum",
				},
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "Warnings",
		Type:    "timeseries",
	}
}
