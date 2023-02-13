package builder

import "fmt"

func getContainerLabelFilter(ecsClusterName, ecsTaskDefinitionName, containerName string) string {
	return fmt.Sprintf(`container_label_com_amazonaws_ecs_cluster="%s", container_label_com_amazonaws_ecs_task_definition_family="%s", container_label_com_amazonaws_ecs_container_name=~"%s"`, ecsClusterName, ecsTaskDefinitionName, containerName)
}

func getTaskDefinitionLabelFilter(ecsClusterName, ecsTaskDefinitionName string) string {
	return fmt.Sprintf(`container_label_com_amazonaws_ecs_cluster="%s", container_label_com_amazonaws_ecs_task_definition_family="%s"`, ecsClusterName, ecsTaskDefinitionName)
}

func NewPanelContainerCpuFactory(containerIndex int) PanelFactory {
	return func(resourceNames ResourceNames, gridPos PanelGridPos) Panel {
		return newPanelContainerCpu(resourceNames, gridPos, containerIndex)
	}
}

func newPanelContainerCpu(resourceNames ResourceNames, gridPos PanelGridPos, containerIndex int) Panel {
	labelFilter := getContainerLabelFilter(resourceNames.EcsCluster, resourceNames.EcsTaskDefinition, resourceNames.Containers[containerIndex])

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
				Expression:   fmt.Sprintf(`max(container_spec_cpu_shares{%s})`, labelFilter),
				LegendFormat: "Reserved",
				RefId:        "reservation",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`min(sum(rate(container_cpu_usage_seconds_total{%s}[$__rate_interval])) by (id))*1024`, labelFilter),
				LegendFormat: "Minimum",
				RefId:        "minimum",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`avg(sum(rate(container_cpu_usage_seconds_total{%s}[$__rate_interval])) by (id))*1024`, labelFilter),
				LegendFormat: "Average",
				RefId:        "average",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`max(sum(rate(container_cpu_usage_seconds_total{%s}[$__rate_interval])) by (id))*1024`, labelFilter),
				LegendFormat: "Maximum",
				RefId:        "maximum",
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   fmt.Sprintf("CPU Utilization (%s)", resourceNames.Containers[containerIndex]),
		Type:    "timeseries",
	}
}

func NewPanelContainerMemoryFactory(containerIndex int) PanelFactory {
	return func(resourceNames ResourceNames, gridPos PanelGridPos) Panel {
		return newPanelContainerMemory(resourceNames, gridPos, containerIndex)
	}
}

func newPanelContainerMemory(resourceNames ResourceNames, gridPos PanelGridPos, containerIndex int) Panel {
	labelFilter := getContainerLabelFilter(resourceNames.EcsCluster, resourceNames.EcsTaskDefinition, resourceNames.Containers[containerIndex])

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
		Title:   fmt.Sprintf("Memory Utilization (%s)", resourceNames.Containers[containerIndex]),
		Type:    "timeseries",
	}
}

func NewPanelServiceUtilization(resourceNames ResourceNames, gridPos PanelGridPos) Panel {
	labelFilter := getTaskDefinitionLabelFilter(resourceNames.EcsCluster, resourceNames.EcsTaskDefinition)

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
				Expression:   fmt.Sprintf(`sum(rate(container_cpu_usage_seconds_total{%s}[$__rate_interval])) by (container_label_com_amazonaws_ecs_container_name)/(sum(container_spec_cpu_shares{%s}) by (container_label_com_amazonaws_ecs_container_name)/1024)*100`, labelFilter, labelFilter),
				LegendFormat: "CPU Average {{container_label_com_amazonaws_ecs_container_name}}",
				RefId:        "cpu_average",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   fmt.Sprintf(`sum(container_memory_rss{%s}) by (container_label_com_amazonaws_ecs_container_name)/sum(container_spec_memory_reservation_limit_bytes{%s}) by (container_label_com_amazonaws_ecs_container_name)*100`, labelFilter, labelFilter),
				LegendFormat: "Memory Average {{container_label_com_amazonaws_ecs_container_name}}",
				RefId:        "memory_average",
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "Service Utilization",
		Type:    "timeseries",
	}
}

func NewPanelTaskDeployment(resourceNames ResourceNames, gridPos PanelGridPos) Panel {
	labelFilter := getContainerLabelFilter(resourceNames.EcsCluster, resourceNames.EcsTaskDefinition, resourceNames.Containers[0])

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
