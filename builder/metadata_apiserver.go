package builder

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MetadataApiServerRoutes []MetadataApiServerRoute

func (k MetadataApiServerRoutes) ToValue() attr.Value {
	list := types.List{
		Elems: make([]attr.Value, len(k)),
		ElemType: types.ObjectType{
			AttrTypes: MetadataApiserverRouteAttrTypes(),
		},
	}

	for i, route := range k {
		list.Elems[i] = route.ToValue()
	}

	return list
}

type MetadataApiServer struct {
	Routes MetadataApiServerRoutes `json:"routes"`
}

func (c MetadataApiServer) ToValue() attr.Value {
	return types.Object{
		Attrs: map[string]attr.Value{
			"routes": c.Routes.ToValue(),
		},
		AttrTypes: MetadataApiserverAttrTypes(),
	}
}

type MetadataApiServerRoute struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

func (c MetadataApiServerRoute) ToValue() attr.Value {
	return types.Object{
		Attrs: map[string]attr.Value{
			"method": types.String{Value: c.Method},
			"path":   types.String{Value: c.Path},
		},
		AttrTypes: MetadataApiserverRouteAttrTypes(),
	}
}

func MetadataApiserverAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"routes": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: MetadataApiserverRouteAttrTypes(),
			},
		},
	}
}

func MetadataApiserverRouteAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"method": types.StringType,
		"path":   types.StringType,
	}
}
