package builder

import "fmt"

func getTaskDefinitionLabelFilter(appId AppId) string {
	return fmt.Sprintf(`container_label_com_amazonaws_ecs_cluster="%s",container_label_com_amazonaws_ecs_task_definition_family="%s-%s"`, appId.EcsClusterName(), appId.EcsClusterName(), appId.Application)
}

func getMainContainerLabelFilter(appId AppId) string {
	return fmt.Sprintf(`container_label_com_amazonaws_ecs_cluster="%s",container_label_com_amazonaws_ecs_container_name="%s"`, appId.EcsClusterName(), appId.Application)
}

func NewPanelTaskCpu(appId AppId, gridPos PanelGridPos) Panel {
	labelFilter := getTaskDefinitionLabelFilter(appId)

	return Panel{
		Datasource: "prometheus",
		FieldConfig: PanelFieldConfig{
			Defaults: PanelFieldConfigDefaults{
				Min: "0",
			},
			Overrides: []PanelFieldConfigOverwrite{
				NewColorPropertyOverwrite("Reserved", "semi-dark-red"),
				NewColorPropertyOverwrite("Minimum", "light-green"),
				NewColorPropertyOverwrite("Average", "light-orange"),
				NewColorPropertyOverwrite("Maximum", "light-red"),
			},
		},
		GridPos: gridPos,
		Targets: []interface{}{
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`avg(sum(container_spec_cpu_shares{%s}) by (instance))`, labelFilter),
				LegendFormat: "Reserved",
				RefId:        "A",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`min(sum(rate(container_cpu_usage_seconds_total{%s}[$__rate_interval])) by (instance))*1024`, labelFilter),
				LegendFormat: "Minimum",
				RefId:        "B",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`avg(sum(rate(container_cpu_usage_seconds_total{%s}[$__rate_interval])) by (instance))*1024`, labelFilter),
				LegendFormat: "Average",
				RefId:        "C",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`max(sum(rate(container_cpu_usage_seconds_total{%s}[$__rate_interval])) by (instance))*1024`, labelFilter),
				LegendFormat: "Maximum",
				RefId:        "D",
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "CPU Utilization",
		Type:    "timeseries",
	}
}

func NewPanelTaskMemory(appId AppId, gridPos PanelGridPos) Panel {
	labelFilter := getTaskDefinitionLabelFilter(appId)

	return Panel{
		Datasource: "prometheus",
		FieldConfig: PanelFieldConfig{
			Defaults: PanelFieldConfigDefaults{
				Min:  "0",
				Unit: "bytes",
			},
			Overrides: []PanelFieldConfigOverwrite{
				NewColorPropertyOverwrite("Reserved", "semi-dark-red"),
				NewColorPropertyOverwrite("Minimum", "light-green"),
				NewColorPropertyOverwrite("Average", "light-orange"),
				NewColorPropertyOverwrite("Maximum", "light-red"),
			},
		},
		GridPos: gridPos,
		Targets: []interface{}{
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`max(sum(container_spec_memory_reservation_limit_bytes{%s}) by (instance))`, labelFilter),
				LegendFormat: "Reserved",
				RefId:        "A",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`min(sum(container_memory_rss{%s}) by (instance))`, labelFilter),
				LegendFormat: "Minimum",
				RefId:        "B",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`avg(sum(container_memory_rss{%s}) by (instance))`, labelFilter),
				LegendFormat: "Average",
				RefId:        "C",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`max(sum(container_memory_rss{%s}) by (instance))`, labelFilter),
				LegendFormat: "Maximum",
				RefId:        "D",
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "Memory Utilization",
		Type:    "timeseries",
	}
}

func NewPanelServiceUtilization(appId AppId, gridPos PanelGridPos) Panel {
	labelFilter := getTaskDefinitionLabelFilter(appId)

	return Panel{
		Datasource: "prometheus",
		FieldConfig: PanelFieldConfig{
			Defaults: PanelFieldConfigDefaults{
				Custom: PanelFieldConfigDefaultsCustom{
					ThresholdsStyle: ThresholdsStyle{
						Mode: "line",
					},
				},
				Max: "200",
				Min: "0",
				Thresholds: PanelFieldConfigDefaultsThresholds{
					Mode: "absolute",
					Steps: []PanelFieldConfigDefaultsThresholdsStep{
						{
							Color: "super-light-green",
						},
						{
							Color: "semi-dark-red",
							Value: 100,
						},
					},
				},
				Unit: "percent",
			},
			Overrides: []PanelFieldConfigOverwrite{
				NewColorPropertyOverwrite("CPU Average", "light-green"),
			},
		},
		GridPos: gridPos,
		Targets: []interface{}{
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`min(sum(rate(container_cpu_usage_seconds_total{%s}[$__rate_interval])) by (instance)/(sum(container_spec_cpu_shares{%s}) by (instance)/1024))*100`, labelFilter, labelFilter),
				LegendFormat: "CPU Minimum",
				RefId:        "A",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`avg(sum(rate(container_cpu_usage_seconds_total{%s}[$__rate_interval])) by (instance)/(sum(container_spec_cpu_shares{%s}) by (instance)/1024))*100`, labelFilter, labelFilter),
				LegendFormat: "CPU Average",
				RefId:        "B",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`max(sum(rate(container_cpu_usage_seconds_total{%s}[$__rate_interval])) by (instance)/(sum(container_spec_cpu_shares{%s}) by (instance)/1024))*100`, labelFilter, labelFilter),
				LegendFormat: "CPU Maximum",
				RefId:        "C",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`avg(sum(container_memory_rss{%s}) by (instance)/sum(container_spec_memory_reservation_limit_bytes{%s}) by (instance))*100`, labelFilter, labelFilter),
				LegendFormat: "Memory Average",
				RefId:        "D",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`max(sum(container_memory_rss{%s}) by (instance)/sum(container_spec_memory_reservation_limit_bytes{%s}) by (instance))*100`, labelFilter, labelFilter),
				LegendFormat: "Memory Maximum",
				RefId:        "E",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`min(sum(container_memory_rss{%s}) by (instance)/sum(container_spec_memory_reservation_limit_bytes{%s}) by (instance))*100`, labelFilter, labelFilter),
				LegendFormat: "Memory Minimum",
				RefId:        "F",
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "Service Utilization",
		Type:    "timeseries",
	}
}

func NewPanelTaskDeployment(appId AppId, gridPos PanelGridPos) Panel {
	labelFilter := getMainContainerLabelFilter(appId)

	return Panel{
		Datasource: "prometheus",
		FieldConfig: PanelFieldConfig{
			Defaults: PanelFieldConfigDefaults{
				Min: "0",
			},
		},
		GridPos: gridPos,
		Targets: []interface{}{
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`count(container_cpu_load_average_10s{%s})`, labelFilter),
				LegendFormat: "RunningTaskCount",
				RefId:        "A",
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "Running Task Count",
		Type:    "timeseries",
	}
}
