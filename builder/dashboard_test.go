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
		Application: "dashboard",
	}
	containers := []string{
		appId.Application,
		"log_router",
	}
	db := builder.NewDashboardBuilder(appId, containers)
	db.AddPanel(builder.NewPanelServiceUtilization)
	db.AddPanel(builder.NewPanelTaskDeployment)
	for _, containerName := range containers {
		db.AddPanel(builder.NewPanelContainerCpuFactory(containerName))
		db.AddPanel(builder.NewPanelContainerMemoryFactory(containerName))
	}
	db.AddPanel(builder.NewPanelTaskLogRouterContainerMemory)
	db.AddPanel(builder.NewPanelError)

	dashboard := db.Build()

	body, _ := json.Marshal(dashboard)
	fmt.Println(string(body))
}
