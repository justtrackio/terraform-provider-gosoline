package builder_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/justtrackio/terraform-provider-gosoline/builder"
)

func TestDashboardWithError(t *testing.T) {
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
	db := builder.NewDashboardBuilder(resourceNames)
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
