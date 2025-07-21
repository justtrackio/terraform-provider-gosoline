package builder_test

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/justtrackio/terraform-provider-gosoline/builder"
)

func TestTraefikResponseTimePanelIntegration(t *testing.T) {
	// Test that the Traefik response time panel generates the correct Prometheus query
	// by creating a dashboard and checking the generated JSON

	resourceNames := &builder.ResourceNames{
		Environment:                        "test",
		GrafanaCloudWatchDatasourceName:    "cw",
		GrafanaElasticsearchDatasourceName: "elastic",
		TraefikServiceName:                 "test-service@kubernetes",
		KubernetesNamespace:                "test-ns",
		KubernetesPod:                      "test-pod",
		Containers:                         []string{"app"},
	}

	db := builder.NewDashboardBuilder(resourceNames, "kubernetes")
	db.AddTraefikService()

	dashboard := db.Build("test dashboard")

	// Marshal to JSON to inspect the generated content
	body, err := json.Marshal(dashboard)
	if err != nil {
		t.Fatal(err)
	}

	jsonStr := string(body)

	// Check that the response time panel exists and has the correct query
	if !strings.Contains(jsonStr, "Response Time") {
		t.Error("Expected to find 'Response Time' panel title")
	}

	// Check that the expression contains the division for average calculation
	expectedQueryParts := []string{
		"traefik_service_request_duration_seconds_sum",
		"traefik_service_requests_total",
		" / ",
	}

	for _, part := range expectedQueryParts {
		if !strings.Contains(jsonStr, part) {
			t.Errorf("Expected to find '%s' in the generated dashboard JSON", part)
		}
	}

	// Verify that the query computes average response time (sum of durations / count of requests)
	// This is the key fix - the expression should divide duration sum by request count
	expectedExpressionPattern := `sum(irate(traefik_service_request_duration_seconds_sum{service=\"test-service@kubernetes\"}[$__rate_interval])) / sum(irate(traefik_service_requests_total{service=\"test-service@kubernetes\"}[$__rate_interval]))`

	if !strings.Contains(jsonStr, expectedExpressionPattern) {
		t.Errorf("Expected to find the corrected expression pattern in JSON:\n%s\nActual JSON:\n%s", expectedExpressionPattern, jsonStr)
	}
}
