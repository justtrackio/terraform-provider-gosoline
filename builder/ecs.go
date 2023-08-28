package builder

import (
	"context"
	"fmt"
	"regexp"
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
	ecsSvc      *ecs.Client
	elbSvc      *elasticloadbalancingv2.Client
	clusterName string
	serviceName string
}

func NewEcsClient(ctx context.Context, clusterName, serviceName string) (*EcsClient, error) {
	var err error
	var cfg aws.Config

	if cfg, err = config.LoadDefaultConfig(ctx); err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %w", err)
	}

	ecsSvc := ecs.NewFromConfig(cfg)
	elbSvc := elasticloadbalancingv2.NewFromConfig(cfg)

	return &EcsClient{
		ecsSvc:      ecsSvc,
		elbSvc:      elbSvc,
		clusterName: clusterName,
		serviceName: serviceName,
	}, nil
}

func (c *EcsClient) GetElbTargetGroups(ctx context.Context) ([]ElbTargetGroup, error) {
	ecsOutput, err := c.ecsSvc.DescribeServices(ctx, &ecs.DescribeServicesInput{
		Cluster:  aws.String(c.clusterName),
		Services: []string{c.serviceName},
	})
	if err != nil {
		return nil, fmt.Errorf("can not describe ecs service %s/%s: %w", c.clusterName, c.serviceName, err)
	}

	if len(ecsOutput.Services) != 1 {
		return nil, fmt.Errorf("there was no ecs service %s/%s found", c.clusterName, c.serviceName)
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
		return nil, fmt.Errorf("can not describe target groups of service %s/%s: %w", c.clusterName, c.serviceName, err)
	}

	targetGroups := make([]ElbTargetGroup, len(elbOutput.TargetGroups))

	for i, targetGroup := range elbOutput.TargetGroups {
		if len(targetGroup.LoadBalancerArns) != 1 {
			return nil, fmt.Errorf("there is more than 1 load balancer for service %s/%s", c.clusterName, c.serviceName)
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

func (c *EcsClient) GetTaskDefinitionName(ctx context.Context) (*string, error) {
	ecsOutput, err := c.ecsSvc.DescribeServices(ctx, &ecs.DescribeServicesInput{
		Cluster:  aws.String(c.clusterName),
		Services: []string{c.serviceName},
	})
	if err != nil {
		return nil, fmt.Errorf("can not describe ecs service %s/%s: %w", c.clusterName, c.serviceName, err)
	}

	if len(ecsOutput.Services) != 1 {
		return nil, fmt.Errorf("there was no ecs service %s/%s found", c.clusterName, c.serviceName)
	}

	taskDefinitionRevisionArn := ecsOutput.Services[0].TaskDefinition
	if taskDefinitionRevisionArn == nil {
		return nil, fmt.Errorf("task definition could not be read from service")
	}

	// extract task definition name from task definition revision arn
	expression := regexp.MustCompile(`task-definition/(.*):\d+`)

	results := expression.FindStringSubmatch(*taskDefinitionRevisionArn)
	if len(results) < 2 {
		return nil, fmt.Errorf("failed to find task definition name in arn %s", *taskDefinitionRevisionArn)
	}

	return &results[1], nil
}

func (c *EcsClient) GetClusterName() string {
	return c.clusterName
}

func (c *EcsClient) GetServiceName() string {
	return c.serviceName
}
