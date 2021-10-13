package builder

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ecs"
	"github.com/aws/aws-sdk-go-v2/service/elasticloadbalancingv2"
)

type ElbTargetGroup struct {
	LoadBalancer string
	TargetGroup  string
}

type EcsClient struct {
	appId  AppId
	ecsSvc *ecs.Client
	elbSvc *elasticloadbalancingv2.Client
}

func NewEcsClient(ctx context.Context, appId AppId) (*EcsClient, error) {
	var err error
	var cfg aws.Config

	if cfg, err = config.LoadDefaultConfig(ctx); err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %w", err)
	}

	ecsSvc := ecs.NewFromConfig(cfg)
	elbSvc := elasticloadbalancingv2.NewFromConfig(cfg)

	return &EcsClient{
		appId:  appId,
		ecsSvc: ecsSvc,
		elbSvc: elbSvc,
	}, nil
}

func (c *EcsClient) GetElbTargetGroups(ctx context.Context) ([]ElbTargetGroup, error) {
	ecsOutput, err := c.ecsSvc.DescribeServices(ctx, &ecs.DescribeServicesInput{
		Cluster:  aws.String(c.appId.EcsClusterName()),
		Services: []string{c.appId.Application},
	})
	if err != nil {
		return nil, fmt.Errorf("can not describe ecs service %s/%s: %w", c.appId.EcsClusterName(), c.appId.Application, err)
	}

	if len(ecsOutput.Services) != 1 {
		return nil, fmt.Errorf("there was no ecs service %s/%s found", c.appId.EcsClusterName(), c.appId.Application)
	}

	loadbalancers := ecsOutput.Services[0].LoadBalancers
	targetGroupArns := make([]string, len(loadbalancers))

	if len(loadbalancers) == 0 {
		return []ElbTargetGroup{}, nil
	}

	for i, loadbalancer := range loadbalancers {
		targetGroupArns[i] = *loadbalancer.TargetGroupArn
	}

	elbOutput, err := c.elbSvc.DescribeTargetGroups(ctx, &elasticloadbalancingv2.DescribeTargetGroupsInput{
		TargetGroupArns: targetGroupArns,
	})
	if err != nil {
		return nil, fmt.Errorf("can not describe target groups of service %s/%s: %w", c.appId.EcsClusterName(), c.appId.Application, err)
	}

	targetGroups := make([]ElbTargetGroup, len(elbOutput.TargetGroups))

	for i, targetGroup := range elbOutput.TargetGroups {
		if len(targetGroup.LoadBalancerArns) != 1 {
			return nil, fmt.Errorf("there is more than 1 load balancer for service %s/%s", c.appId.EcsClusterName(), c.appId.Application)
		}

		k := strings.LastIndex(targetGroup.LoadBalancerArns[0], ":")
		l := strings.LastIndex(*targetGroup.TargetGroupArn, ":")

		targetGroups[i] = ElbTargetGroup{
			LoadBalancer: targetGroup.LoadBalancerArns[0][k+14:],
			TargetGroup:  (*targetGroup.TargetGroupArn)[l+1:],
		}
	}

	return targetGroups, nil
}
