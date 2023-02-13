package builder

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcsClient(t *testing.T) {
	t.SkipNow()
	clusterName := "cluster"
	serviceName := "service"

	client, err := NewEcsClient(context.Background(), clusterName, serviceName)
	assert.NoError(t, err)

	balancers, err := client.GetElbTargetGroups(context.Background())
	assert.NoError(t, err)

	fmt.Println(balancers)
}
