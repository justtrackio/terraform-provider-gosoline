package builder_test

import (
	"testing"

	"github.com/justtrackio/terraform-provider-gosoline/builder"
	"github.com/stretchr/testify/assert"
)

func TestTaskDeploymentPanelUsesCorrectMetric(t *testing.T) {
	// Test data for a "gateway" deployment
	kubernetesNamespace := "test-namespace"
	gatewayDeployment := "gateway"

	resourceNames := &builder.ResourceNames{
		KubernetesNamespace:  kubernetesNamespace,
		KubernetesDeployment: gatewayDeployment,
		KubernetesPod:        gatewayDeployment,
		Containers:           []string{"app"},
	}

	// Create a dashboard builder and generate the task deployment panel
	db := builder.NewDashboardBuilder(resourceNames, "kubernetes")
	db.AddPanel(builder.NewPanelTaskDeployment)
	dashboard := db.Build("test dashboard")

	// Find the task deployment panel by looking through all panels for the Running Task Count title
	var taskPanel builder.Panel
	var found bool
	for _, panel := range dashboard.Panels {
		if panel.Title == "Running Task Count" {
			taskPanel = panel
			found = true

			break
		}
	}

	assert.True(t, found, "Running Task Count panel should be found in dashboard")

	// Verify the panel has the correct Prometheus target
	assert.Len(t, taskPanel.Targets, 1, "Panel should have exactly one target")

	target, ok := taskPanel.Targets[0].(builder.PanelTargetPrometheus)
	assert.True(t, ok, "Target should be a Prometheus target")

	// Verify the query uses the new deployment metric instead of pod_info
	query := target.Expression

	// Should use kube_deployment_status_replicas_ready with sum()
	assert.Contains(t, query, "sum(kube_deployment_status_replicas_ready",
		"Query should use sum(kube_deployment_status_replicas_ready) metric")

	// Should NOT use the old problematic kube_pod_info
	assert.NotContains(t, query, "kube_pod_info",
		"Query should not use the problematic kube_pod_info metric")

	// Should use exact deployment name matching
	assert.Contains(t, query, `deployment="gateway"`,
		"Query should use exact deployment name matching")

	// Should NOT use wildcard pattern
	assert.NotContains(t, query, "gateway-.*",
		"Query should not use wildcard pattern that could match other deployments")

	// Should still include namespace filtering
	assert.Contains(t, query, `namespace="test-namespace"`,
		"Query should include namespace filtering")

	t.Logf("Generated query: %s", query)
}

func TestTaskDeploymentPanelAvoidsCrossTalk(t *testing.T) {
	// Test that "gateway" deployment doesn't match "gateway-abc" deployment

	// Create queries for both deployments
	gatewayResourceNames := &builder.ResourceNames{
		KubernetesNamespace:  "test-ns",
		KubernetesDeployment: "gateway",
		KubernetesPod:        "gateway",
		Containers:           []string{"app"},
	}

	gatewayAbcResourceNames := &builder.ResourceNames{
		KubernetesNamespace:  "test-ns",
		KubernetesDeployment: "gateway-abc",
		KubernetesPod:        "gateway-abc",
		Containers:           []string{"app"},
	}

	// Generate panels for both
	gatewayDb := builder.NewDashboardBuilder(gatewayResourceNames, "kubernetes")
	gatewayDb.AddPanel(builder.NewPanelTaskDeployment)
	gatewayDashboard := gatewayDb.Build("gateway dashboard")

	gatewayAbcDb := builder.NewDashboardBuilder(gatewayAbcResourceNames, "kubernetes")
	gatewayAbcDb.AddPanel(builder.NewPanelTaskDeployment)
	gatewayAbcDashboard := gatewayAbcDb.Build("gateway-abc dashboard")

	// Extract queries
	var gatewayQuery, gatewayAbcQuery string

	for _, panel := range gatewayDashboard.Panels {
		if panel.Title == "Running Task Count" {
			target := panel.Targets[0].(builder.PanelTargetPrometheus)
			gatewayQuery = target.Expression

			break
		}
	}

	for _, panel := range gatewayAbcDashboard.Panels {
		if panel.Title == "Running Task Count" {
			target := panel.Targets[0].(builder.PanelTargetPrometheus)
			gatewayAbcQuery = target.Expression

			break
		}
	}

	// Verify the queries are different and specific
	assert.NotEqual(t, gatewayQuery, gatewayAbcQuery,
		"Queries for different deployments should be different")

	assert.Contains(t, gatewayQuery, `deployment="gateway"`,
		"Gateway query should target exactly 'gateway' deployment")

	assert.Contains(t, gatewayAbcQuery, `deployment="gateway-abc"`,
		"Gateway-abc query should target exactly 'gateway-abc' deployment")

	t.Logf("Gateway query: %s", gatewayQuery)
	t.Logf("Gateway-abc query: %s", gatewayAbcQuery)
}
