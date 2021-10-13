package builder

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEcsClient(t *testing.T) {
	appId := AppId{
		Project:     "mcoins",
		Environment: "prod",
		Family:      "marketing",
		Application: "attribution-product-adoption",
	}

	client, err := NewEcsClient(context.Background(), appId)
	assert.NoError(t, err)

	balancers, err := client.GetElbTargetGroups(context.Background())
	assert.NoError(t, err)

	fmt.Println(balancers)
	// arn:aws:elasticloadbalancing:eu-central-1:164105964448:loadbalancer/app/mcoins-pr-marketing-playtime/9855ea596f342e28
	// arn:aws:elasticloadbalancing:eu-central-1:164105964448:targetgroup/mcoins-pr-marketing-playtime/96b5a3def5ca7e89
}
