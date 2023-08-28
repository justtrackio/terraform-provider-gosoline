package builder

import "fmt"

const defaultClientName = "default"

type ResourceNames struct {
	CloudwatchNamespace                string
	EcsCluster                         string
	EcsService                         string
	EcsTaskDefinition                  string
	GrafanaCloudWatchDatasourceName    string
	GrafanaElasticsearchDatasourceName string
	KubernetesNamespace                string
	KubernetesPod                      string
	TargetGroups                       []ElbTargetGroup
	TraefikServiceName                 string
	Containers                         []string
}

func (r *ResourceNames) GetCwDatasourceNameByClientName(clientName string) string {
	if clientName == defaultClientName {
		return r.GrafanaCloudWatchDatasourceName
	}

	return fmt.Sprintf("cloudwatch-%s", clientName)
}
