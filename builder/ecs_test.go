package builder

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcsClient(t *testing.T) {
	appId := AppId{
		Project:     "myPrj",
		Environment: "production",
		Family:      "biz",
		Application: "fancyBackend",
	}

	client, err := NewEcsClient(context.Background(), appId)
	assert.NoError(t, err)

	balancers, err := client.GetElbTargetGroups(context.Background())
	assert.NoError(t, err)

	fmt.Println(balancers)
}
