package builder

import "fmt"

func getContainerLabelFilter(appId AppId, containerName string) string {
	return fmt.Sprintf(`container_label_com_amazonaws_ecs_cluster="%s", container_label_com_amazonaws_ecs_task_definition_family="%s-%s", container_label_com_amazonaws_ecs_container_name="%s"`, appId.EcsClusterName(), appId.EcsClusterName(), appId.Application, containerName)
}

func getTaskDefinitionLabelFilter(appId AppId) string {
	return fmt.Sprintf(`container_label_com_amazonaws_ecs_cluster="%s", container_label_com_amazonaws_ecs_task_definition_family="%s-%s"`, appId.EcsClusterName(), appId.EcsClusterName(), appId.Application)
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

func NewPanelTaskAppContainerMemory(appId AppId, gridPos PanelGridPos) Panel {
	return newPanelTaskContainerMemory(appId, gridPos, appId.Application)
}

func NewPanelTaskLogRouterContainerMemory(appId AppId, gridPos PanelGridPos) Panel {
	return newPanelTaskContainerMemory(appId, gridPos, "log_router")
}

func newPanelTaskContainerMemory(appId AppId, gridPos PanelGridPos, containerName string) Panel {
	labelFilter := getContainerLabelFilter(appId, containerName)

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
				Expression:   fmt.Sprintf(`max(container_spec_memory_reservation_limit_bytes{%s})`, labelFilter),
				LegendFormat: "Reserved",
				RefId:        "A",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`min by (container_label_com_amazonaws_ecs_container_name) (container_memory_rss{%s})`, labelFilter),
				LegendFormat: "Minimum",
				RefId:        "B",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`avg by (container_label_com_amazonaws_ecs_container_name) (container_memory_rss{%s})`, labelFilter),
				LegendFormat: "Average",
				RefId:        "C",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`max by (container_label_com_amazonaws_ecs_container_name) (container_memory_rss{%s})`, labelFilter),
				LegendFormat: "Maximum",
				RefId:        "D",
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   fmt.Sprintf("Memory Utilization (%s)", containerName),
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
	labelFilter := getContainerLabelFilter(appId, appId.Application)

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
