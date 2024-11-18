package pf

import (
	"github.com/go-openapi/spec"
	"reflect"
)

func GenerateSpec(signatures signatures) *spec.Swagger {
	s := new(spec.Swagger)

	// Store definitions to use later
	for _, methods := range signatures {
		for _, sig := range methods {
			_, exists := s.Definitions[sig.reqType.Name()]
			if exists {
				continue
			}

		}
	}

	for path, methods := range signatures {
		s.Paths.Paths[path] = createPathItem(methods)
	}
}

func createPathItem(methods map[string]*handlerSignature) spec.PathItem {
	var props spec.PathItemProps

	return spec.PathItem{PathItemProps: props}
}

func createStructSchema(typ reflect.Type) spec.Schema {
	var props spec.SchemaProps

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		name := field.Tag.Get("json")
		if name == "" {
			name = field.Name
		}
		if field.Anonymous {
			continue
		}
	}
}
