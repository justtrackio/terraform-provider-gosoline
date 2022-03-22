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
