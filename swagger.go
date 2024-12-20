package pf

import (
	"bytes"
	"log/slog"
	"mime/multipart"
	"net/http"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/go-openapi/spec"
	httpSwagger "github.com/swaggo/http-swagger"
)

// AddSwagger generates the OpenAPI 2.0 spec from the handlers already created
// in r and routes the Swagger spec page to endpoint. Only call AddSwagger
// after routing every handler you wish displayed on the page.
func AddSwagger(r *Router, endpoint string, info *SwaggerInfo) error {
	if info == nil {
		info = new(SwaggerInfo)
	}

	s := generateSpec(r.traverseSignatures(), info)
	slog.Info("swagger: generated spec")

	json, err := s.MarshalJSON()
	if err != nil {
		return err
	}

	buffer := bytes.NewReader(json)

	r.mux.Get(
		path.Join(endpoint, "swagger.json"),
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			http.ServeContent(w, r, "swagger.json", time.Time{}, buffer)
		},
	)

	handler := httpSwagger.Handler(httpSwagger.URL("./swagger.json"))

	r.mux.Get(path.Join(endpoint, "*"), handler)

	redirect := func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, path.Join(endpoint, "index.html"), http.StatusFound)
	}
	r.mux.Get(endpoint, redirect)

	slog.Info("swagger: added handler", "endpoint", endpoint)

	return nil
}

// HandlerProperty represents a modification to the handler's metadata
// (summary, description etc.) for Swagger.
type HandlerProperty func(op *spec.Operation)

// WithSummary sets the handler's summary.
func WithSummary(summary string) HandlerProperty {
	return func(op *spec.Operation) {
		op.Summary = summary
	}
}

// WithSummary sets the handler's description.
func WithDescription(description string) HandlerProperty {
	return func(op *spec.Operation) {
		op.Description = description
	}
}

// WithQuery adds query parameters to the handler's metadata.
func WithQuery(query ...string) HandlerProperty {
	return func(op *spec.Operation) {
		for _, q := range query {
			op.Parameters = append(op.Parameters, spec.Parameter{
				ParamProps: spec.ParamProps{
					In:   "query",
					Name: q,
				},
			})
		}
	}
}

// WithConsumes sets the MIME types the handler expects as the request body.
func WithConsumes(mime ...string) HandlerProperty {
	return func(op *spec.Operation) {
		op.Consumes = mime
	}
}

// WithProduces sets the MIME types the handler produces as a response.
func WithProduces(mime ...string) HandlerProperty {
	return func(op *spec.Operation) {
		op.Produces = mime
	}
}

// SwaggerInfo represents metadata about the API.
type SwaggerInfo struct {
	Title          string
	Description    string
	TermsOfService string
	ContactName    string
	ContactURL     string
	ContactEmail   string
	License        string
	LicenseURL     string
	Version        string
}

func toSpecInfo(i *SwaggerInfo) *spec.Info {
	return &spec.Info{
		InfoProps: spec.InfoProps{
			Title:          i.Title,
			Description:    i.Description,
			TermsOfService: i.TermsOfService,
			Contact: &spec.ContactInfo{
				ContactInfoProps: spec.ContactInfoProps{
					Name:  i.ContactName,
					URL:   i.ContactURL,
					Email: i.ContactEmail,
				},
			},
			License: &spec.License{
				LicenseProps: spec.LicenseProps{
					Name: i.License,
					URL:  i.LicenseURL,
				},
			},
			Version: i.Version,
		},
	}
}

type structMap map[reflect.Type]spec.Schema

func generateSpec(signatures signatures, info *SwaggerInfo) *spec.Swagger {
	var s spec.Swagger
	s.Swagger = "2.0"
	s.Info = toSpecInfo(info)
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

func createPathItem(methods map[string]*handlerSignature, structMap structMap) spec.PathItem {
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
			slog.Warn("swagger: unsupported method", "method", method)
		}
	}

	return item
}

func createOperation(sig *handlerSignature, structMap structMap) *spec.Operation {
	var op spec.Operation

	switch sig.reqType {
	case nil, reflect.TypeFor[struct{}]():
	case reflect.TypeFor[*multipart.Form]():
		op.Consumes = []string{"multipart/form-data"}
	default:
		op.Consumes = []string{"application/json"}
		req := getType(sig.reqType, structMap)
		op.Parameters = append(op.Parameters, *spec.BodyParam("body", &req))
	}

	switch sig.resType {
	case nil, reflect.TypeFor[struct{}]():
	default:
		describeResponse(&op, sig.resType, structMap)
	}

	for _, prop := range sig.props {
		prop(&op)
	}

	return &op
}

func describeResponse(op *spec.Operation, typ reflect.Type, structMap structMap) {
	switch typ {
	case reflect.TypeFor[[]byte]():
		op.Produces = []string{"application/octet-stream"}
	case reflect.TypeFor[string]():
		op.Produces = []string{"text/plain"}
	default:
		op.Produces = []string{"application/json"}
		res := getType(typ, structMap)
		op.Responses = &spec.Responses{
			ResponsesProps: spec.ResponsesProps{
				Default: &spec.Response{
					ResponseProps: spec.ResponseProps{
						Schema: &res,
					},
				},
			},
		}
	}
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
