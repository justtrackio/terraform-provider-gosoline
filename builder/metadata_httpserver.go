package builder

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type MetadataHttpServers []MetadataHttpServer

func (s MetadataHttpServers) ToValue() attr.Value {
	list := types.List{
		Elems: make([]attr.Value, len(s)),
		ElemType: types.ObjectType{
			AttrTypes: MetadataHttpServerAttrTypes(),
		},
	}

	for i, server := range s {
		list.Elems[i] = server.ToValue()
	}

	return list
}

type MetadataHttpServer struct {
	Name     string                     `json:"name"`
	Handlers MetadataHttpServerHandlers `json:"handlers"`
}

func (c MetadataHttpServer) ToValue() attr.Value {
	return types.Object{
		Attrs: map[string]attr.Value{
			"name":     types.String{Value: c.Name},
			"handlers": c.Handlers.ToValue(),
		},
		AttrTypes: MetadataHttpServerAttrTypes(),
	}
}

func MetadataHttpServerAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"name": types.StringType,
		"handlers": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: MetadataHttpServerRouteAttrTypes(),
			},
		},
	}
}

type MetadataHttpServerHandlers []MetadataHttpServerHandler

func (k MetadataHttpServerHandlers) ToValue() attr.Value {
	list := types.List{
		Elems: make([]attr.Value, len(k)),
		ElemType: types.ObjectType{
			AttrTypes: MetadataHttpServerRouteAttrTypes(),
		},
	}

	for i, route := range k {
		list.Elems[i] = route.ToValue()
	}

	return list
}

type MetadataHttpServerHandler struct {
	Method string `json:"method"`
	Path   string `json:"path"`
}

func (c MetadataHttpServerHandler) ToValue() attr.Value {
	return types.Object{
		Attrs: map[string]attr.Value{
			"method": types.String{Value: c.Method},
			"path":   types.String{Value: c.Path},
		},
		AttrTypes: MetadataHttpServerRouteAttrTypes(),
	}
}

func MetadataHttpServerRouteAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"method": types.StringType,
		"path":   types.StringType,
	}
}
