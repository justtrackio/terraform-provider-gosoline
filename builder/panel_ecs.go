package builder

func NewPanelEcsCpu(appId AppId, gridPos PanelGridPos) Panel {
	return Panel{
		Datasource: "CloudWatch",
		FieldConfig: PanelFieldConfig{
			Defaults: PanelFieldConfigDefaults{
				Min: "0",
			},
			Overrides: []PanelFieldConfigOverwrite{
				NewColorPropertyOverwrite("CpuReserved", "dark-orange"),
				NewColorPropertyOverwrite("CpuUtilized", "light-green"),
			},
		},
		GridPos: gridPos,
		Targets: []interface{}{
			PanelTargetCloudWatch{
				Dimensions: map[string]string{
					"ClusterName": appId.EcsClusterName(),
					"ServiceName": appId.Application,
				},
				Id:         "m1",
				MetricName: "CpuReserved",
				Namespace:  "ECS/ContainerInsights",
				Period:     "60",
				RefId:      "A",
				Region:     "default",
				Statistics: []string{
					"Sum",
				},
			},
			PanelTargetCloudWatch{
				Dimensions: map[string]string{
					"ClusterName": appId.EcsClusterName(),
					"ServiceName": appId.Application,
				},
				MetricName: "CpuUtilized",
				Namespace:  "ECS/ContainerInsights",
				Period:     "60",
				RefId:      "B",
				Region:     "default",
				Statistics: []string{
					"Sum",
				},
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "CPU Utilization",
		Type:    "timeseries",
	}
}

func NewPanelEcsMemory(appId AppId, gridPos PanelGridPos) Panel {
	return Panel{
		Datasource: "CloudWatch",
		FieldConfig: PanelFieldConfig{
			Defaults: PanelFieldConfigDefaults{
				Min:  "0",
				Unit: "decmbytes",
			},
			Overrides: []PanelFieldConfigOverwrite{
				NewColorPropertyOverwrite("MemoryReserved", "dark-orange"),
				NewColorPropertyOverwrite("MemoryUtilized", "light-blue"),
			},
		},
		GridPos: gridPos,
		Targets: []interface{}{
			PanelTargetCloudWatch{
				Dimensions: map[string]string{
					"ClusterName": appId.EcsClusterName(),
					"ServiceName": appId.Application,
				},
				MatchExact: false,
				MetricName: "MemoryReserved",
				Namespace:  "ECS/ContainerInsights",
				Period:     "60",
				RefId:      "A",
				Region:     "default",
				Statistics: []string{
					"Sum",
				},
			},
			PanelTargetCloudWatch{
				Dimensions: map[string]string{
					"ClusterName": appId.EcsClusterName(),
					"ServiceName": appId.Application,
				},
				MatchExact: false,
				MetricName: "MemoryUtilized",
				Namespace:  "ECS/ContainerInsights",
				Period:     "60",
				RefId:      "B",
				Region:     "default",
				Statistics: []string{
					"Sum",
				},
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "Memory Utilization",
		Type:    "timeseries",
	}
}

func NewPanelEcsUtilization(appId AppId, gridPos PanelGridPos) Panel {
	return Panel{
		Datasource: "CloudWatch",
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
			PanelTargetCloudWatch{
				Dimensions: map[string]string{
					"ClusterName": appId.EcsClusterName(),
					"ServiceName": appId.Application,
				},
				Id:         "cpu_reserved",
				Hide:       true,
				MetricName: "CpuReserved",
				Namespace:  "ECS/ContainerInsights",
				RefId:      "A",
				Region:     "default",
				Statistics: []string{
					"Maximum",
				},
			},
			PanelTargetCloudWatch{
				Dimensions: map[string]string{
					"ClusterName": appId.EcsClusterName(),
					"ServiceName": appId.Application,
				},
				Id:         "cpu_average",
				Hide:       true,
				MetricName: "CpuUtilized",
				Namespace:  "ECS/ContainerInsights",
				RefId:      "B",
				Region:     "default",
				Statistics: []string{
					"Average",
				},
			},
			PanelTargetCloudWatch{
				Dimensions: map[string]string{
					"ClusterName": appId.EcsClusterName(),
					"ServiceName": appId.Application,
				},
				Id:         "cpu_maximum",
				Hide:       true,
				MetricName: "CpuUtilized",
				Namespace:  "ECS/ContainerInsights",
				RefId:      "C",
				Region:     "default",
				Statistics: []string{
					"Maximum",
				},
			},
			PanelTargetCloudWatch{
				Alias:      "CPU Average",
				Expression: "100 / MAX(cpu_reserved) * AVG(cpu_average)",
				Dimensions: map[string]string{},
				Id:         "e1",
				MatchExact: true,
				RefId:      "D",
				Region:     "default",
				Statistics: []string{
					"Average",
				},
			},
			PanelTargetCloudWatch{
				Alias:      "CPU Maximum",
				Expression: "100 / MAX(cpu_reserved) * MAX(cpu_maximum)",
				Dimensions: map[string]string{},
				Id:         "e2",
				MatchExact: true,
				RefId:      "E",
				Region:     "default",
				Statistics: []string{
					"Maximum",
				},
			},
			PanelTargetCloudWatch{
				Dimensions: map[string]string{
					"ClusterName": appId.EcsClusterName(),
					"ServiceName": appId.Application,
				},
				Id:         "memory_reserved",
				Hide:       true,
				MetricName: "MemoryReserved",
				Namespace:  "ECS/ContainerInsights",
				RefId:      "F",
				Region:     "default",
				Statistics: []string{
					"Maximum",
				},
			},
			PanelTargetCloudWatch{
				Dimensions: map[string]string{
					"ClusterName": appId.EcsClusterName(),
					"ServiceName": appId.Application,
				},
				Id:         "memory_average",
				Hide:       true,
				MetricName: "MemoryUtilized",
				Namespace:  "ECS/ContainerInsights",
				RefId:      "G",
				Region:     "default",
				Statistics: []string{
					"Average",
				},
			},
			PanelTargetCloudWatch{
				Dimensions: map[string]string{
					"ClusterName": appId.EcsClusterName(),
					"ServiceName": appId.Application,
				},
				Id:         "memory_maximum",
				Hide:       true,
				MetricName: "MemoryUtilized",
				Namespace:  "ECS/ContainerInsights",
				RefId:      "H",
				Region:     "default",
				Statistics: []string{
					"Maximum",
				},
			},
			PanelTargetCloudWatch{
				Alias:      "Memory Average",
				Expression: "100 / MAX(memory_reserved) * AVG(memory_average)",
				Dimensions: map[string]string{},
				Id:         "e3",
				RefId:      "I",
				Region:     "default",
				Statistics: []string{
					"Average",
				},
			},
			PanelTargetCloudWatch{
				Alias:      "Memory Maximum",
				Expression: "100 / AVG(memory_reserved) * MAX(memory_maximum)",
				Dimensions: map[string]string{},
				Id:         "e4",
				RefId:      "J",
				Region:     "default",
				Statistics: []string{
					"Maximum",
				},
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "Service Utilization",
		Type:    "timeseries",
	}
}

func NewPanelEcsDeployment(appId AppId, gridPos PanelGridPos) Panel {
	return Panel{
		Datasource: "CloudWatch",
		FieldConfig: PanelFieldConfig{
			Defaults: PanelFieldConfigDefaults{
				Min: "0",
			},
		},
		GridPos: gridPos,
		Targets: []interface{}{
			PanelTargetCloudWatch{
				Dimensions: map[string]string{
					"ClusterName": appId.EcsClusterName(),
					"ServiceName": appId.Application,
				},
				MetricName: "RunningTaskCount",
				Namespace:  "ECS/ContainerInsights",
				RefId:      "A",
				Region:     "default",
				Statistics: []string{
					"Maximum",
				},
			},
			PanelTargetCloudWatch{
				Dimensions: map[string]string{
					"ClusterName": appId.EcsClusterName(),
					"ServiceName": appId.Application,
				},
				MetricName: "PendingTaskCount",
				Namespace:  "ECS/ContainerInsights",
				RefId:      "B",
				Region:     "default",
				Statistics: []string{
					"Maximum",
				},
			},
			PanelTargetCloudWatch{
				Dimensions: map[string]string{
					"ClusterName": appId.EcsClusterName(),
					"ServiceName": appId.Application,
				},
				MetricName: "DesiredTaskCount",
				Namespace:  "ECS/ContainerInsights",
				RefId:      "C",
				Region:     "default",
				Statistics: []string{
					"Maximum",
				},
			},
			PanelTargetCloudWatch{
				Dimensions: map[string]string{
					"ClusterName": appId.EcsClusterName(),
					"ServiceName": appId.Application,
				},
				MetricName: "DeploymentCount",
				Namespace:  "ECS/ContainerInsights",
				RefId:      "D",
				Region:     "default",
				Statistics: []string{
					"Maximum",
				},
			},
		},
		Options: &PanelOptionsCloudWatch{},
		Title:   "Deployment and Task Count",
		Type:    "timeseries",
	}
}
