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
