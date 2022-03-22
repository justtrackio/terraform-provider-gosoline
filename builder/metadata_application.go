package builder

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MetadataApplication struct {
	ApiServer MetadataApiServer `json:"apiserver"`
	Cloud     MetadataCloud     `json:"cloud"`
	Stream    MetadataStream    `json:"stream"`
}

func (a MetadataApplication) ToValue() types.Object {
	return types.Object{
		AttrTypes: MetadataApplicationAttrTypes(),
		Attrs: map[string]attr.Value{
			"apiserver": a.ApiServer.ToValue(),
			"cloud":     a.Cloud.ToValue(),
			"stream":    a.Stream.ToValue(),
		},
	}
}

func MetadataApplicationAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"apiserver": types.ObjectType{
			AttrTypes: MetadataApiserverAttrTypes(),
		},
		"cloud": types.ObjectType{
			AttrTypes: MetadataCloudAttrTypes(),
		},
		"stream": types.ObjectType{
			AttrTypes: MetadataStreamAttrTypes(),
		},
	}
}
