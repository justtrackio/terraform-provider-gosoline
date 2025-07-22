package builder

type ResourceNames struct {
	CloudwatchNamespace                string
	Containers                         []string
	EcsCluster                         string
	EcsService                         string
	EcsTaskDefinition                  string
	Environment                        string
	GrafanaCloudWatchDatasourceName    string
	GrafanaElasticsearchDatasourceName string
	KubernetesDeployment               string
	KubernetesNamespace                string
	KubernetesPod                      string
	TargetGroups                       []ElbTargetGroup
	TraefikServiceName                 string
}
