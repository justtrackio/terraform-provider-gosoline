package builder

import "fmt"

const (
	datasourcePrometheus = "prometheus"
)

func getTraefikServiceLabelFilter(serviceName string) string {
	return fmt.Sprintf(`service=%q`, serviceName)
}

func getKubernetesPodLabelFilter(namespace, podName string) string {
	return fmt.Sprintf(`namespace=%q, pod=~"%s-.*"`, namespace, podName)
}

func getEcsContainerLabelFilter(ecsClusterName, ecsTaskDefinitionName, containerName string) string {
	return fmt.Sprintf(`container_label_com_amazonaws_ecs_cluster=%q, container_label_com_amazonaws_ecs_task_definition_family=%q, container_label_com_amazonaws_ecs_container_name="%q"`, ecsClusterName, ecsTaskDefinitionName, containerName)
}

func getEcsTaskDefinitionLabelFilter(ecsClusterName, ecsTaskDefinitionName string) string {
	return fmt.Sprintf(`container_label_com_amazonaws_ecs_cluster=%q, container_label_com_amazonaws_ecs_task_definition_family=%q`, ecsClusterName, ecsTaskDefinitionName)
}

func getLabelAndFilters(settings PanelSettings, containerIndex int) (string, string, string) {
	var containerLabelFilter string
	var containerLabel string
	var podLabelFilter string

	switch settings.orchestrator {
	case orchestratorEcs:
		containerLabel = "container_label_com_amazonaws_ecs_container_name"
		containerLabelFilter = getEcsContainerLabelFilter(settings.resourceNames.EcsCluster, settings.resourceNames.EcsTaskDefinition, settings.resourceNames.Containers[containerIndex])
		podLabelFilter = getEcsTaskDefinitionLabelFilter(settings.resourceNames.EcsCluster, settings.resourceNames.EcsTaskDefinition)
	case orchestratorKubernetes:
		containerLabel = "container"
		containerLabelFilter = getKubernetesPodLabelFilter(settings.resourceNames.KubernetesNamespace, settings.resourceNames.KubernetesPod)
		podLabelFilter = getKubernetesPodLabelFilter(settings.resourceNames.KubernetesNamespace, settings.resourceNames.KubernetesPod)
	}

	return containerLabel, containerLabelFilter, podLabelFilter
}

func NewPanelContainerCpuFactory(containerIndex int) PanelFactory {
	return func(settings PanelSettings) Panel {
		return newPanelContainerCpu(settings, containerIndex)
	}
}

func newPanelContainerCpu(settings PanelSettings, containerIndex int) Panel {
	var labelFilter string
	var averageQuery string
	var reservedQuery string
	var minimumQuery string
	var maximumQuery string

	switch settings.orchestrator {
	case orchestratorEcs:
		labelFilter = getEcsContainerLabelFilter(settings.resourceNames.EcsCluster, settings.resourceNames.EcsTaskDefinition, settings.resourceNames.Containers[containerIndex])
		reservedQuery = fmt.Sprintf(`max(container_spec_cpu_shares{%s})`, labelFilter)
		averageQuery = fmt.Sprintf(`avg(sum(rate(container_cpu_usage_seconds_total{%s}[$__rate_interval])) by (id))*1024`, labelFilter)
		maximumQuery = fmt.Sprintf(`max(sum(rate(container_cpu_usage_seconds_total{%s}[$__rate_interval])) by (id))*1024`, labelFilter)
		minimumQuery = fmt.Sprintf(`min(sum(rate(container_cpu_usage_seconds_total{%s}[$__rate_interval])) by (id))*1024`, labelFilter)
	case orchestratorKubernetes:
		labelFilter = getKubernetesPodLabelFilter(settings.resourceNames.KubernetesNamespace, settings.resourceNames.KubernetesPod)
		reservedQuery = fmt.Sprintf(`max(sum(kube_pod_container_resource_requests{resource="cpu",%s} * on(namespace,pod) group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{%s}) by (pod))`, labelFilter, labelFilter)
		averageQuery = fmt.Sprintf(`avg(sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{%s} * on(namespace,pod) group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{%s}) by (pod))`, labelFilter, labelFilter)
		maximumQuery = fmt.Sprintf(`max(sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{%s} * on(namespace,pod) group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{%s}) by (pod))`, labelFilter, labelFilter)
		minimumQuery = fmt.Sprintf(`min(sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{%s} * on(namespace,pod) group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{%s}) by (pod))`, labelFilter, labelFilter)
	}

	return Panel{
		Datasource: datasourcePrometheus,
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
		GridPos: settings.gridPos,
		Targets: []interface{}{
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   reservedQuery,
				LegendFormat: "Reserved",
				RefId:        "reservation",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   minimumQuery,
				LegendFormat: "Minimum",
				RefId:        "minimum",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   averageQuery,
				LegendFormat: "Average",
				RefId:        "average",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   maximumQuery,
				LegendFormat: "Maximum",
				RefId:        "maximum",
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   fmt.Sprintf("CPU Utilization (%s)", settings.resourceNames.Containers[containerIndex]),
		Type:    "timeseries",
	}
}

func NewPanelContainerMemoryFactory(containerIndex int) PanelFactory {
	return func(settings PanelSettings) Panel {
		return newPanelContainerMemory(settings, containerIndex)
	}
}

func newPanelContainerMemory(settings PanelSettings, containerIndex int) Panel {
	containerLabel, containerLabelFilter, podLabelFilter := getLabelAndFilters(settings, containerIndex)
	var averageQuery string
	var reservedQuery string
	var minimumQuery string
	var maximumQuery string

	switch settings.orchestrator {
	case orchestratorEcs:
		reservedQuery = fmt.Sprintf(`max(container_spec_memory_reservation_limit_bytes{%s})`, containerLabelFilter)
		averageQuery = fmt.Sprintf(`avg by (%s) (container_memory_rss{%s})`, containerLabel, containerLabelFilter)
		maximumQuery = fmt.Sprintf(`max by (%s) (container_memory_rss{%s})`, containerLabel, containerLabelFilter)
		minimumQuery = fmt.Sprintf(`min by (%s) (container_memory_rss{%s})`, containerLabel, containerLabelFilter)
	case orchestratorKubernetes:
		reservedQuery = fmt.Sprintf(`max(sum(kube_pod_container_resource_limits{resource="memory",%s} * on(namespace,pod) group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{%s}) by (pod))`, podLabelFilter, podLabelFilter)
		averageQuery = fmt.Sprintf(`avg(sum(container_memory_working_set_bytes{container!="", image!="", %s} * on(namespace,pod) group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{%s}) by (pod))`, podLabelFilter, podLabelFilter)
		maximumQuery = fmt.Sprintf(`max(sum(container_memory_working_set_bytes{container!="", image!="", %s} * on(namespace,pod) group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{%s}) by (pod))`, podLabelFilter, podLabelFilter)
		minimumQuery = fmt.Sprintf(`min(sum(container_memory_working_set_bytes{container!="", image!="", %s} * on(namespace,pod) group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{%s}) by (pod))`, podLabelFilter, podLabelFilter)
	}

	return Panel{
		Datasource: datasourcePrometheus,
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
		GridPos: settings.gridPos,
		Targets: []interface{}{
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   reservedQuery,
				LegendFormat: "Reserved",
				RefId:        "A",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   minimumQuery,
				LegendFormat: "Minimum",
				RefId:        "B",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   averageQuery,
				LegendFormat: "Average",
				RefId:        "C",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   maximumQuery,
				LegendFormat: "Maximum",
				RefId:        "D",
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   fmt.Sprintf("Memory Utilization (%s)", settings.resourceNames.Containers[containerIndex]),
		Type:    "timeseries",
	}
}

func NewPanelServiceUtilization(settings PanelSettings) Panel {
	containerLabel, _, podLabelFilter := getLabelAndFilters(settings, 0)
	var cpuAverageQuery string
	var memoryAverageQuery string

	switch settings.orchestrator {
	case orchestratorEcs:
		cpuAverageQuery = fmt.Sprintf(`sum(rate(container_cpu_usage_seconds_total{%s}[$__rate_interval])) by (%s)/(sum(container_spec_cpu_shares{%s}) by (%s)/1024)*100`, podLabelFilter, containerLabel, podLabelFilter, containerLabel)
		memoryAverageQuery = fmt.Sprintf(`sum(container_memory_rss{%s}) by (%s)/sum(container_spec_memory_reservation_limit_bytes{%s}) by (%s)*100`, podLabelFilter, containerLabel, podLabelFilter, containerLabel)
	case orchestratorKubernetes:
		cpuAverageQuery = fmt.Sprintf(`avg(sum(node_namespace_pod_container:container_cpu_usage_seconds_total:sum_irate{%s} * on(namespace,pod) group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{%s}) by (pod)/sum(kube_pod_container_resource_requests{resource="cpu", %s} * on(namespace,pod) group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{%s}) by (pod)*100)`, podLabelFilter, podLabelFilter, podLabelFilter, podLabelFilter)
		memoryAverageQuery = fmt.Sprintf(`avg(sum(container_memory_working_set_bytes{%s, container!="", image!=""} * on(namespace,pod) group_left(workload, workload_type) namespace_workload_pod:kube_pod_owner:relabel{%s}) by (pod)/ on(pod) cluster:namespace:pod_memory:active:kube_pod_container_resource_requests{resource="memory",%s})*100`, podLabelFilter, podLabelFilter, podLabelFilter)
	}

	cpuAverageLegendFormat := fmt.Sprintf("CPU Average {{%s}}", containerLabel)
	memoryAverageLegendFormat := fmt.Sprintf("Memory Average {{%s}}", containerLabel)

	return Panel{
		Datasource: datasourcePrometheus,
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
				NewColorPropertyOverwrite("CPU Average ", "light-green"), // the trailing slash seems to be important for grafana to match the override due to omitting the {{foo}} part
			},
		},
		GridPos: settings.gridPos,
		Targets: []interface{}{
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   cpuAverageQuery,
				LegendFormat: cpuAverageLegendFormat,
				RefId:        "cpu_average",
			},
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   memoryAverageQuery,
				LegendFormat: memoryAverageLegendFormat,
				RefId:        "memory_average",
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "Service Utilization",
		Type:    "timeseries",
	}
}

func NewPanelTaskDeployment(settings PanelSettings) Panel {
	var labelFilter string
	var query string

	switch settings.orchestrator {
	case orchestratorEcs:
		labelFilter = getEcsContainerLabelFilter(settings.resourceNames.EcsCluster, settings.resourceNames.EcsTaskDefinition, settings.resourceNames.Containers[0])
		query = fmt.Sprintf(`count(container_cpu_load_average_10s{%s})`, labelFilter)
	case orchestratorKubernetes:
		labelFilter = getKubernetesPodLabelFilter(settings.resourceNames.KubernetesNamespace, settings.resourceNames.KubernetesPod)
		query = fmt.Sprintf("count(kube_pod_info{%s})", labelFilter)
	}

	return Panel{
		Datasource: datasourcePrometheus,
		FieldConfig: PanelFieldConfig{
			Defaults: PanelFieldConfigDefaults{
				Min: "0",
			},
		},
		GridPos: settings.gridPos,
		Targets: []interface{}{
			PanelTargetPrometheus{
				Exemplar:     true,
				Expression:   query,
				LegendFormat: "RunningTaskCount",
				RefId:        "A",
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "Running Task Count",
		Type:    "timeseries",
	}
}
