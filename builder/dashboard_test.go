package builder_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/justtrackio/terraform-provider-gosoline/builder"
)

func TestEcsDashboardWithError(t *testing.T) {
	appId := builder.AppId{
		Project:     "gosoline",
		Environment: "test",
		Family:      "monitoring",
		Group:       "grp",
		Application: "dashboard",
	}
	cloudwatchNamespace := "a/b/c/d/e"
	ecsClusterName := "cluster"
	ecsServiceName := "service"
	ecsTaskDefinitionName := "task-def"
	grafanaElasticsearchDatasourceName := "elastic"
	grafanaCloudWatchDatasourceName := "cw"
	containers := []string{
		appId.Application,
		"log_router",
	}
	targetGroups := []builder.ElbTargetGroup{
		{
			LoadBalancer: "foo",
			TargetGroup:  "bar",
		},
	}

	resourceNames := builder.ResourceNames{
		CloudwatchNamespace:                cloudwatchNamespace,
		EcsCluster:                         ecsClusterName,
		EcsService:                         ecsServiceName,
		EcsTaskDefinition:                  ecsTaskDefinitionName,
		GrafanaElasticsearchDatasourceName: grafanaElasticsearchDatasourceName,
		GrafanaCloudWatchDatasourceName:    grafanaCloudWatchDatasourceName,
		TargetGroups:                       targetGroups,
		Containers:                         containers,
	}
	db := builder.NewDashboardBuilder(resourceNames, "ecs")
	db.AddPanel(builder.NewPanelServiceUtilization)
	db.AddPanel(builder.NewPanelTaskDeployment)
	for i := range containers {
		db.AddPanel(builder.NewPanelContainerCpuFactory(i))
		db.AddPanel(builder.NewPanelContainerMemoryFactory(i))
	}
	db.AddPanel(builder.NewPanelWarn)
	db.AddPanel(builder.NewPanelError)

	dashboard := db.Build()

	body, _ := json.Marshal(dashboard)
	fmt.Println(string(body))
}

func TestKubernetesDashboardWithError(t *testing.T) {
	appId := builder.AppId{
		Project:     "gosoline",
		Environment: "test",
		Family:      "monitoring",
		Group:       "grp",
		Application: "dashboard",
	}
	cloudwatchNamespace := "a/b/c/d/e"
	ecsClusterName := ""
	ecsServiceName := ""
	ecsTaskDefinitionName := ""
	grafanaElasticsearchDatasourceName := "elastic"
	grafanaCloudWatchDatasourceName := "cw"
	containers := []string{
		appId.Application,
		"log_router",
	}
	traefikServiceName := "ns-service-8070@kubernetes"
	kubernetesNamespace := "foo"
	kubernetesPod := "pod"
	var targetGroups []builder.ElbTargetGroup

	resourceNames := builder.ResourceNames{
		CloudwatchNamespace:                cloudwatchNamespace,
		EcsCluster:                         ecsClusterName,
		EcsService:                         ecsServiceName,
		EcsTaskDefinition:                  ecsTaskDefinitionName,
		GrafanaCloudWatchDatasourceName:    grafanaCloudWatchDatasourceName,
		GrafanaElasticsearchDatasourceName: grafanaElasticsearchDatasourceName,
		KubernetesNamespace:                kubernetesNamespace,
		KubernetesPod:                      kubernetesPod,
		TargetGroups:                       targetGroups,
		TraefikServiceName:                 traefikServiceName,
		Containers:                         containers,
	}
	db := builder.NewDashboardBuilder(resourceNames, "kubernetes")
	db.AddPanel(builder.NewPanelServiceUtilization)
	db.AddPanel(builder.NewPanelTaskDeployment)
	for i := range containers {
		db.AddPanel(builder.NewPanelContainerCpuFactory(i))
		db.AddPanel(builder.NewPanelContainerMemoryFactory(i))
	}
	db.AddPanel(builder.NewPanelWarn)
	db.AddPanel(builder.NewPanelError)
	db.AddTraefikService()

	dashboard := db.Build()

	body, _ := json.Marshal(dashboard)
	fmt.Println(string(body))
}
