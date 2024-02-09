package builder

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MetadataApplication struct {
	Cloud       MetadataCloud       `json:"cloud"`
	HttpServers MetadataHttpServers `json:"httpservers"`
	Stream      MetadataStream      `json:"stream"`
}

func (a MetadataApplication) ToValue() types.Object {
	return types.Object{
		AttrTypes: MetadataApplicationAttrTypes(),
		Attrs: map[string]attr.Value{
			"cloud":       a.Cloud.ToValue(),
			"httpservers": a.HttpServers.ToValue(),
			"stream":      a.Stream.ToValue(),
		},
	}
}

func MetadataApplicationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"cloud": types.ObjectType{
			AttrTypes: MetadataCloudAttrTypes(),
		},
		"httpservers": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: MetadataHttpServerAttrTypes(),
			},
		},
		"stream": types.ObjectType{
			AttrTypes: MetadataStreamAttrTypes(),
		},
	}
}
