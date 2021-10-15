package builder

import "fmt"

func NewPanelLogs(appId AppId, gridPos PanelGridPos) Panel {
	datasource := fmt.Sprintf("elasticsearch-%s-logs-%s-%s-%s", appId.Environment, appId.Project, appId.Family, appId.Application)
	gridPos.W = 24
	gridPos.H = 16

	return Panel{
		Datasource: datasource,
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
