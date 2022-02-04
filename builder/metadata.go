package builder

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MetadataAppId struct {
	Project     string `json:"project"`
	Environment string `json:"environment"`
	Family      string `json:"family"`
	Application string `json:"application"`
}

func (a MetadataAppId) ToValue() types.Object {
	return types.Object{
		AttrTypes: MetadataAppIdAttrTypes(),
		Attrs: map[string]attr.Value{
			"project":     types.String{Value: a.Project},
			"environment": types.String{Value: a.Environment},
			"family":      types.String{Value: a.Family},
			"application": types.String{Value: a.Application},
		},
	}
}

func MetadataAppIdAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"project":     types.StringType,
		"environment": types.StringType,
		"family":      types.StringType,
		"application": types.StringType,
	}
}

type MetadataApplication struct {
	ApiServer MetadataApiServer `json:"apiserver"`
	Cloud     MetadataCloud     `json:"cloud"`
	Stream    MetadataStream    `json:"stream"`
}

func (a MetadataApplication) ToValue() types.Object {
	return types.Object{
		AttrTypes: MetadataApplicationAttrTypes(),
		Attrs: map[string]attr.Value{
			"cloud":  a.Cloud.ToValue(),
			"stream": a.Stream.ToValue(),
		},
	}
}

func MetadataApplicationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cloud": types.ObjectType{
			AttrTypes: MetadataCloudAttrTypes(),
		},
		"stream": types.ObjectType{
			AttrTypes: MetadataStreamAttrTypes(),
		},
	}
}

type MetadataApiServer struct {
	Routes []MetadataApiServerRoute `json:"routes"`
}

type MetadataApiServerRoute struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

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

type MetadataCloudAwsDynamodb struct {
	Tables []string `json:"tables"`
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

type MetadataCloudAwsKinesisKinsumer struct {
	ClientId       string        `json:"client_id"`
	Name           string        `json:"name"`
	StreamAppId    MetadataAppId `json:"stream_app_id"`
	StreamName     string        `json:"stream_name"`
	StreamNameFull string        `json:"stream_name_full"`
}

func (k MetadataCloudAwsKinesisKinsumer) ToValue() types.Object {
	return types.Object{
		Attrs: map[string]attr.Value{
			"client_id":        types.String{Value: k.ClientId},
			"name":             types.String{Value: k.Name},
			"stream_app_id":    k.StreamAppId.ToValue(),
			"stream_name":      types.String{Value: k.StreamName},
			"stream_name_full": types.String{Value: k.StreamNameFull},
		},
		AttrTypes: MetadataCloudAwsKinesisKinsumerAttrTypes(),
	}
}

func MetadataCloudAwsKinesisKinsumerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"client_id": types.StringType,
		"name":      types.StringType,
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
	StreamName string `json:"stream_name"`
}

func (k MetadataCloudAwsKinesisRecordWriter) ToValue() types.Object {
	return types.Object{
		Attrs: map[string]attr.Value{
			"stream_name": types.String{Value: k.StreamName},
		},
		AttrTypes: MetadataCloudAwsKinesisRecordWriterAttrTypes(),
	}
}

func MetadataCloudAwsKinesisRecordWriterAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"stream_name": types.StringType,
	}
}

type MetadataCloudAwsSqs struct {
	Queues []string `json:"queues"`
}

type MetadataStream struct {
	Consumers MetadataStreamConsumers `json:"consumers"`
	Producers MetadataStreamProducers `json:"producers"`
}

func (s MetadataStream) ToValue() types.Object {
	return types.Object{
		Attrs: map[string]attr.Value{
			"consumers": s.Consumers.ToValue(),
			"producers": s.Producers.ToValue(),
		},
		AttrTypes: MetadataStreamAttrTypes(),
	}
}

func MetadataStreamAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"consumers": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: MetadataStreamConsumerAttrTypes(),
			},
		},
		"producers": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: MetadataStreamProducerAttrTypes(),
			},
		},
	}
}

type MetadataStreamConsumers []MetadataStreamConsumer

func (c MetadataStreamConsumers) ToValue() attr.Value {
	list := types.List{
		Elems: make([]attr.Value, len(c)),
		ElemType: types.ObjectType{
			AttrTypes: MetadataStreamConsumerAttrTypes(),
		},
	}

	for i, consumer := range c {
		list.Elems[i] = consumer.ToValue()
	}

	return list
}

type MetadataStreamConsumer struct {
	Name         string `json:"name"`
	RetryEnabled bool   `json:"retry_enabled"`
	RetryType    string `json:"retry_type"`
}

func (p MetadataStreamConsumer) ToValue() attr.Value {
	return types.Object{
		Attrs: map[string]attr.Value{
			"name":          types.String{Value: p.Name},
			"retry_enabled": types.Bool{Value: p.RetryEnabled},
			"retry_type":    types.String{Value: p.RetryType},
		},
		AttrTypes: MetadataStreamConsumerAttrTypes(),
	}
}

func MetadataStreamConsumerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":          types.StringType,
		"retry_enabled": types.BoolType,
		"retry_type":    types.StringType,
	}
}

type MetadataStreamProducers []MetadataStreamProducer

func (p MetadataStreamProducers) ToValue() types.List {
	list := types.List{
		Elems: make([]attr.Value, len(p)),
		ElemType: types.ObjectType{
			AttrTypes: MetadataStreamProducerAttrTypes(),
		},
	}

	for i, producer := range p {
		list.Elems[i] = producer.ToValue()
	}

	return list
}

type MetadataStreamProducer struct {
	Name          string `json:"name"`
	DaemonEnabled bool   `json:"daemon_enabled"`
}

func (p MetadataStreamProducer) ToValue() attr.Value {
	return types.Object{
		Attrs: map[string]attr.Value{
			"name":           types.String{Value: p.Name},
			"daemon_enabled": types.Bool{Value: p.DaemonEnabled},
		},
		AttrTypes: MetadataStreamProducerAttrTypes(),
	}
}

func MetadataStreamProducerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name":           types.StringType,
		"daemon_enabled": types.BoolType,
	}
}
