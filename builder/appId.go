package builder

import "fmt"

type AppId struct {
	Project     string
	Environment string
	Family      string
	Application string
}

func (i AppId) CloudWatchNamespace() string {
	return fmt.Sprintf("%s/%s/%s/%s", i.Project, i.Environment, i.Family, i.Application)
}

func (i AppId) EcsClusterName() string {
	return fmt.Sprintf("%s-%s-%s", i.Project, i.Environment, i.Family)
}
