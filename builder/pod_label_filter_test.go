package builder_test

import (
	"strings"
	"testing"

	"github.com/justtrackio/terraform-provider-gosoline/builder"
)

func TestKubernetesPodLabelFilterRegex(t *testing.T) {
	tests := []struct {
		name        string
		namespace   string
		podName     string
		description string
	}{
		{
			name:        "basic pod name",
			namespace:   "test-ns",
			podName:     "gateway",
			description: "Should generate correct filter for basic pod name",
		},
		{
			name:        "hyphenated pod name",
			namespace:   "prod",
			podName:     "gateway-api",
			description: "Should handle hyphenated pod names correctly",
		},
		{
			name:        "complex pod name",
			namespace:   "staging",
			podName:     "my-service-worker",
			description: "Should handle complex pod names with multiple hyphens",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test through the panel creation since the function is not exported
			resourceNames := &builder.ResourceNames{
				KubernetesNamespace: tt.namespace,
				KubernetesPod:       tt.podName,
				Containers:          []string{"app"},
			}

			// Create a panel and check that the generated expression uses the new regex pattern
			db := builder.NewDashboardBuilder(resourceNames, "kubernetes")
			db.AddPanel(builder.NewPanelKubernetesHealthyPods)
			dashboard := db.Build("test")

			// Find the panel and check its expression
			found := false
			expectedPattern := `pod=~"^` + tt.podName + `-[0-9a-f]+-[0-9a-z]+$"`
			
			for _, panel := range dashboard.Panels {
				if panel.Title == "Healthy Endpoints" {
					targets := panel.Targets
					if len(targets) > 0 {
						if target, ok := targets[0].(builder.PanelTargetPrometheus); ok {
							if strings.Contains(target.Expression, expectedPattern) {
								found = true
								break
							} else {
								t.Errorf("Expression %q does not contain expected pattern %q", target.Expression, expectedPattern)
							}
						}
					}
				}
			}

			if !found {
				t.Error("Could not find the expected pod label filter pattern in any panel")
			}
		})
	}
}

