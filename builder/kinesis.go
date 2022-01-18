package builder

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

type KinesisClient struct {
	kinesisSvc *kinesis.Client
}

func NewKinesisClient(ctx context.Context) (*KinesisClient, error) {
	var err error
	var cfg aws.Config

	if cfg, err = config.LoadDefaultConfig(ctx); err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %w", err)
	}

	kinesisSvc := kinesis.NewFromConfig(cfg)

	return &KinesisClient{
		kinesisSvc: kinesisSvc,
	}, nil
}

func (c *KinesisClient) GetShardCount(ctx context.Context, streamName string) (int, error) {
	var err error
	var out *kinesis.DescribeStreamOutput

	input := &kinesis.DescribeStreamInput{
		StreamName: aws.String(streamName),
	}

	if out, err = c.kinesisSvc.DescribeStream(ctx, input); err != nil {
		return 0, fmt.Errorf("can not describe kinesis stream %s: %w", streamName, err)
	}

	return len(out.StreamDescription.Shards), nil
}
