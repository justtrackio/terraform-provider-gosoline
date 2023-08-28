package builder

func NewPanelLogs(settings PanelSettings) Panel {
	settings.gridPos.W = 24
	settings.gridPos.H = 16

	return Panel{
		Datasource: settings.resourceNames.GrafanaElasticsearchDatasourceName,
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
		GridPos: settings.gridPos,
		Targets: []interface{}{
			PanelTargetElasticsearch{
				RefId: "A",
				Query: "level:[3 TO *]",
				Metrics: []PanelTargetElasticsearchMetric{
					{
						Id:   "1",
						Type: "logs",
						Settings: PanelTargetElasticsearchMetricSettings{
							Limit: "100",
						},
					},
				},
				TimeField: "@timestamp",
			},
		},
		Options: PanelOptionsElasticsearch{
			ShowTime:         true,
			EnableLogDetails: true,
			DedupStrategy:    "none",
			SortOrder:        "Descending",
		},
		Title: "Error & Warning Logs",
		Type:  "logs",
	}
}
