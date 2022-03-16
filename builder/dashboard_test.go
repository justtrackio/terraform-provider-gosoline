package builder_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/justtrackio/terraform-provider-gosoline/builder"
)

func TestDashboardWithError(t *testing.T) {
	db := builder.NewDashboardBuilder(builder.AppId{
		Project:     "gosoline",
		Environment: "test",
		Family:      "monitoring",
		Application: "dashboard",
	})
	db.AddPanel(builder.NewPanelEcsCpu)
	db.AddPanel(builder.NewPanelEcsMemory)
	db.AddPanel(builder.NewPanelError)

	dashboard := db.Build()

	body, _ := json.Marshal(dashboard)
	fmt.Println(string(body))
}
