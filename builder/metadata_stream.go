package builder

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

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
