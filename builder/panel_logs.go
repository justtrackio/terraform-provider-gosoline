package builder

func NewPanelLogs(resourceNames ResourceNames, gridPos PanelGridPos) Panel {
	gridPos.W = 24
	gridPos.H = 16

	return Panel{
		Datasource: resourceNames.GrafanaElasticsearchDatasourceName,
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
		GridPos: gridPos,
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
