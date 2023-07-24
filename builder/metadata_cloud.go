package builder

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MetadataCloud struct {
	Aws MetadataCloudAws `json:"aws"`
}

func (c MetadataCloud) ToValue() attr.Value {
	return types.Object{
		Attrs: map[string]attr.Value{
			"aws": c.Aws.ToValue(),
		},
		AttrTypes: MetadataCloudAttrTypes(),
	}
}

func MetadataCloudAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"aws": types.ObjectType{
			AttrTypes: MetadataCloudAwsAttrTypes(),
		},
	}
}

type MetadataCloudAws struct {
	Dynamodb MetadataCloudAwsDynamodb `json:"dynamodb"`
	Kinesis  MetadataCloudAwsKinesis  `json:"kinesis"`
	Sqs      MetadataCloudAwsSqs      `json:"sqs"`
}

func (a MetadataCloudAws) ToValue() attr.Value {
	return types.Object{
		Attrs: map[string]attr.Value{
			"kinesis": a.Kinesis.ToValue(),
		},
		AttrTypes: MetadataCloudAwsAttrTypes(),
	}
}

func MetadataCloudAwsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"kinesis": types.ObjectType{
			AttrTypes: MetadataCloudAwsKinesisAttrTypes(),
		},
	}
}

type MetadataCloudAwsDynamodbTable struct {
	AwsClientName string `json:"aws_client_name"`
	TableName     string `json:"table_name"`
}

type MetadataCloudAwsDynamodb struct {
	Tables []MetadataCloudAwsDynamodbTable `json:"tables"`
}

type MetadataCloudAwsKinesis struct {
	Kinsumers     MetadataCloudAwsKinesisKinsumers     `json:"kinsumers"`
	RecordWriters MetadataCloudAwsKinesisRecordWriters `json:"record_writers"`
}

func (k MetadataCloudAwsKinesis) ToValue() attr.Value {
	return types.Object{
		Attrs: map[string]attr.Value{
			"kinsumers":      k.Kinsumers.ToValue(),
			"record_writers": k.RecordWriters.ToValue(),
		},
		AttrTypes: MetadataCloudAwsKinesisAttrTypes(),
	}
}

func MetadataCloudAwsKinesisAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"kinsumers": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: MetadataCloudAwsKinesisKinsumerAttrTypes(),
			},
		},
		"record_writers": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: MetadataCloudAwsKinesisRecordWriterAttrTypes(),
			},
		},
	}
}

type MetadataCloudAwsKinesisKinsumers []MetadataCloudAwsKinesisKinsumer

func (k MetadataCloudAwsKinesisKinsumers) ToValue() types.List {
	list := types.List{
		Elems: make([]attr.Value, len(k)),
		ElemType: types.ObjectType{
			AttrTypes: MetadataCloudAwsKinesisKinsumerAttrTypes(),
		},
	}

	for i, kinsumer := range k {
		list.Elems[i] = kinsumer.ToValue()
	}

	return list
}

type KinesisStreamAware interface {
	GetClientName() string
	GetStreamNameFull() string
	GetOpenShardCount() int
}

type MetadataCloudAwsKinesisKinsumer struct {
	AwsClientName  string        `json:"aws_client_name"`
	ClientId       string        `json:"client_id"`
	Name           string        `json:"name"`
	OpenShardCount int64         `json:"open_shard_count"`
	StreamAppId    MetadataAppId `json:"stream_app_id"`
	StreamArn      string        `json:"stream_arn"`
	StreamName     string        `json:"stream_name"`
	StreamNameFull string        `json:"stream_name_full"`
}

func (k MetadataCloudAwsKinesisKinsumer) GetClientName() string {
	return k.AwsClientName
}

func (k MetadataCloudAwsKinesisKinsumer) GetStreamNameFull() string {
	return k.StreamNameFull
}

func (k MetadataCloudAwsKinesisKinsumer) GetOpenShardCount() int {
	return int(k.OpenShardCount)
}

func (k MetadataCloudAwsKinesisKinsumer) ToValue() types.Object {
	return types.Object{
		Attrs: map[string]attr.Value{
			"aws_client_name":  types.String{Value: k.AwsClientName},
			"client_id":        types.String{Value: k.ClientId},
			"name":             types.String{Value: k.Name},
			"open_shard_count": types.Int64{Value: k.OpenShardCount},
			"stream_app_id":    k.StreamAppId.ToValue(),
			"stream_arn":       types.String{Value: k.StreamArn},
			"stream_name":      types.String{Value: k.StreamName},
			"stream_name_full": types.String{Value: k.StreamNameFull},
		},
		AttrTypes: MetadataCloudAwsKinesisKinsumerAttrTypes(),
	}
}

func MetadataCloudAwsKinesisKinsumerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"aws_client_name":  types.StringType,
		"client_id":        types.StringType,
		"name":             types.StringType,
		"open_shard_count": types.Int64Type,
		"stream_arn":       types.StringType,
		"stream_app_id": types.ObjectType{
			AttrTypes: MetadataAppIdAttrTypes(),
		},
		"stream_name":      types.StringType,
		"stream_name_full": types.StringType,
	}
}

type MetadataCloudAwsKinesisRecordWriters []MetadataCloudAwsKinesisRecordWriter

func (k MetadataCloudAwsKinesisRecordWriters) ToValue() types.List {
	list := types.List{
		Elems: make([]attr.Value, len(k)),
		ElemType: types.ObjectType{
			AttrTypes: MetadataCloudAwsKinesisRecordWriterAttrTypes(),
		},
	}

	for i, kinsumer := range k {
		list.Elems[i] = kinsumer.ToValue()
	}

	return list
}

type MetadataCloudAwsKinesisRecordWriter struct {
	AwsClientName  string `json:"aws_client_name"`
	OpenShardCount int64  `json:"open_shard_count"`
	StreamArn      string `json:"stream_arn"`
	StreamName     string `json:"stream_name"`
}

func (k MetadataCloudAwsKinesisRecordWriter) GetClientName() string {
	return k.AwsClientName
}

func (k MetadataCloudAwsKinesisRecordWriter) GetStreamNameFull() string {
	return k.StreamName
}

func (k MetadataCloudAwsKinesisRecordWriter) GetOpenShardCount() int {
	return int(k.OpenShardCount)
}

func (k MetadataCloudAwsKinesisRecordWriter) ToValue() types.Object {
	return types.Object{
		Attrs: map[string]attr.Value{
			"aws_client_name":  types.String{Value: k.AwsClientName},
			"open_shard_count": types.Int64{Value: k.OpenShardCount},
			"stream_arn":       types.String{Value: k.StreamArn},
			"stream_name":      types.String{Value: k.StreamName},
		},
		AttrTypes: MetadataCloudAwsKinesisRecordWriterAttrTypes(),
	}
}

func MetadataCloudAwsKinesisRecordWriterAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"aws_client_name":  types.StringType,
		"stream_arn":       types.StringType,
		"stream_name":      types.StringType,
		"open_shard_count": types.Int64Type,
	}
}

type MetadataCloudAwsSqsQueue struct {
	AwsClientName string `json:"aws_client_name"`
	QueueArn      string `json:"queue_arn"`
	QueueName     string `json:"queue_name"`
	QueueNameFull string `json:"queue_name_full"`
	QueueUrl      string `json:"queue_url"`
}

type MetadataCloudAwsSqs struct {
	Queues []MetadataCloudAwsSqsQueue `json:"queues"`
}
