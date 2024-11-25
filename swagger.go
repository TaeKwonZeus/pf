package pf

import (
	"log/slog"
	"mime/multipart"
	"net/http"
	"reflect"
	"strings"

	"github.com/go-openapi/spec"
)

type structMap map[reflect.Type]spec.Schema

func GenerateSpec(signatures signatures, info spec.InfoProps) *spec.Swagger {
	var s spec.Swagger
	s.Swagger = "2.0"
	s.Info = &spec.Info{InfoProps: info}
	s.Paths = &spec.Paths{Paths: make(map[string]spec.PathItem)}
	s.Definitions = make(spec.Definitions)

	// Store definitions to use later
	structMap := make(structMap)

	for path, methods := range signatures {
		s.Paths.Paths[path] = createPathItem(methods, structMap)
	}

	for typ, schema := range structMap {
		s.Definitions[typ.Name()] = schema
	}

	return &s
}

func createPathItem(methods map[method]*handlerSignature, structMap structMap) spec.PathItem {
	var item spec.PathItem

	for method, sig := range methods {
		op := createOperation(sig, structMap)

		switch method {
		case http.MethodGet:
			item.Get = op
		case http.MethodPost:
			item.Post = op
		case http.MethodPut:
			item.Put = op
		case http.MethodDelete:
			item.Delete = op
		case http.MethodPatch:
			item.Patch = op
		case http.MethodOptions:
			item.Options = op
		case http.MethodHead:
			item.Head = op
		default:
			slog.Error("swagger: unsupported method", "method", method)
		}
	}

	return item
}

func createOperation(sig *handlerSignature, structMap structMap) *spec.Operation {
	var operation spec.Operation

	switch sig.reqType {
	case reflect.TypeOf(struct{}{}):
	case reflect.TypeOf(&multipart.Form{}):
		operation.Consumes = []string{"multipart/form-data"}
	default:
		operation.Consumes = []string{"encoding/json"}
		req := getType(sig.reqType, structMap)
		operation.Parameters = append(operation.Parameters, *spec.BodyParam("body", &req))
	}

	// TODO add non-JSON responses
	operation.Produces = []string{"encoding/json"}

	res := getType(sig.resType, structMap)
	operation.Responses = &spec.Responses{
		ResponsesProps: spec.ResponsesProps{
			Default: &spec.Response{
				ResponseProps: spec.ResponseProps{
					Schema: &res,
				},
			},
		},
	}

	return &operation
}

// getType marshals a type into a spec.Schema.
// structMap is the map of named structs which should be populated with any named structs encountered.
func getType(typ reflect.Type, structMap structMap) spec.Schema {
	var schema spec.Schema

	switch typ.Kind() {
	case reflect.Bool:
		schema.Type = []string{"boolean"}

	case reflect.Float32:
		schema.Type = []string{"number"}
		schema.Schema = "float"

	case reflect.Float64:
		schema.Type = []string{"number"}
		schema.Schema = "double"

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64, reflect.Int8,
		reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint8, reflect.Uintptr:
		schema.Type = []string{"integer"}

	case reflect.String:
		schema.Type = []string{"string"}

	case reflect.Struct:
		return getStruct(typ, structMap)

	case reflect.Pointer:
		return getType(typ.Elem(), structMap)

	case reflect.Array, reflect.Slice:
		schema.Type = []string{"array"}
		elem := getType(typ.Elem(), structMap)
		schema.Items = &spec.SchemaOrArray{Schema: &elem}

	case reflect.Map:
		if typ.Key().Kind() == reflect.String {
			schema.Type = []string{"object"}
			elem := getType(typ.Elem(), structMap)
			schema.AdditionalProperties = &spec.SchemaOrBool{Schema: &elem}
		} else {
			slog.Error("swagger: support only for maps with string keys; skipping", "type", typ.Name())
		}

	default:
		slog.Error("swagger: unexpected reflect.Kind", "kind", typ.Kind())
	}

	return schema
}

func getStruct(typ reflect.Type, structMap structMap) spec.Schema {
	// Just in case
	if typ.Kind() != reflect.Struct {
		panic("swagger: non-struct type passed to getStruct")
	}

	if typ.NumField() == 0 {
		return spec.Schema{}
	}

	// return ref if found
	if _, ok := structMap[typ]; ok {
		return ref(typ)
	}

	var schema spec.Schema
	schema.Type = []string{"object"}
	schema.Properties = make(spec.SchemaProperties)

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)
		name, required := fieldName(field)

		schema.Properties[name] = getType(field.Type, structMap)
		if required {
			schema.Required = append(schema.Required, name)
		}
	}

	// If the type is inline return the schema itself
	if typ.Name() == "" {
		return schema
	}

	// Otherwise, return a reference
	structMap[typ] = schema
	return ref(typ)
}

func ref(typ reflect.Type) spec.Schema {
	return spec.Schema{
		SchemaProps: spec.SchemaProps{
			Ref: spec.MustCreateRef("#/definitions/" + typ.Name()),
		},
	}
}

func fieldName(field reflect.StructField) (name string, required bool) {
	name = field.Name
	if field.Anonymous {
		name = field.Type.Name()
	}

	tag := field.Tag.Get("json")
	if tag == "" {
		return
	}

	tagName, omitempty, _ := strings.Cut(tag, ",")
	if tagName != "" {
		name = tagName
	}

	return name, omitempty != "omitempty"
}
