package builder

func NewPanelError(settings PanelSettings) Panel {
	return Panel{
		Datasource: settings.resourceNames.GrafanaCloudWatchDatasourceName,
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
		GridPos: settings.gridPos,
		Targets: []interface{}{
			PanelTargetCloudWatch{
				Alias:      "Errors",
				Dimensions: map[string]string{},
				MatchExact: false,
				MetricName: "error",
				Namespace:  settings.resourceNames.CloudwatchNamespace,
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

func NewPanelWarn(settings PanelSettings) Panel {
	return Panel{
		Datasource: settings.resourceNames.GrafanaCloudWatchDatasourceName,
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
		GridPos: settings.gridPos,
		Targets: []interface{}{
			PanelTargetCloudWatch{
				Alias:      "Warnings",
				Dimensions: map[string]string{},
				MetricName: "warn",
				Namespace:  settings.resourceNames.CloudwatchNamespace,
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
