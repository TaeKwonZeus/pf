package pf

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

// 1st key: path
// 2nd key: method
type signatures map[string]map[string]*handlerSignature

func (s signatures) add(path string, method string, signature *handlerSignature) {
	if s[path] == nil {
		s[path] = make(map[string]*handlerSignature)
	}
	s[path][method] = signature
}

// Router is a composable router based on chi.Mux that tracks handler request and response body signatures.
type Router struct {
	mux        chi.Router
	subrouters []*Router
	prefix     string

	signatures signatures
}

func NewRouter() *Router {
	return &Router{
		mux:        chi.NewRouter(),
		signatures: make(signatures),
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.mux.ServeHTTP(w, req)
}

func (r *Router) traverseSignatures() signatures {
	out := make(signatures)
	for k, v := range r.signatures {
		out[r.prefix+k] = v
	}

	for _, subrouter := range r.subrouters {
		sub := subrouter.traverseSignatures()
		for k, v := range sub {
			out[r.prefix+k] = v
		}
	}

	return out
}

func Use(r *Router, middlewares ...func(next http.Handler) http.Handler) {
	r.mux.Use(middlewares...)
}

func Method[Req, Res any](r *Router, method string, path string, handler Handler[Req, Res]) {
	h, signature := handler.wrap()
	r.mux.Method(method, path, h)
	r.signatures.add(path, method, signature)
}

func Get[Res any](r *Router, path string, handler Handler[struct{}, Res]) {
	Method(r, http.MethodGet, path, handler)
}

func Post[Req, Res any](r *Router, path string, handler Handler[Req, Res]) {
	Method(r, http.MethodPost, path, handler)
}

func Put[Req, Res any](r *Router, path string, handler Handler[Req, Res]) {
	Method(r, http.MethodPut, path, handler)
}

func Delete[Req, Res any](r *Router, path string, handler Handler[Req, Res]) {
	Method(r, http.MethodDelete, path, handler)
}

func Patch[Req, Res any](r *Router, path string, handler Handler[Req, Res]) {
	Method(r, http.MethodPatch, path, handler)
}

func Head(r *Router, path string, handler Handler[struct{}, struct{}]) {
	Method(r, http.MethodHead, path, handler)
}

func Options[Res any](r *Router, path string, handler Handler[struct{}, Res]) {
	Method(r, http.MethodOptions, path, handler)
}

func Trace[Res any](r *Router, path string, handler Handler[struct{}, Res]) {
	Method(r, http.MethodTrace, path, handler)
}

func Connect(r *Router, path string, handler Handler[struct{}, struct{}]) {
	Method(r, http.MethodConnect, path, handler)
}

func Route(r *Router, path string, fn func(r *Router)) {
	subrouter := NewRouter()
	fn(subrouter)
	Mount(r, path, subrouter)
}

// Handle routes handler to the path. You should not use this for regular function handlers, as handlers added with
// Handle will NOT show up in Swagger.
func Handle(r *Router, path string, handler http.Handler) {
	r.mux.Handle(path, handler)
}

func Mount(r *Router, path string, subrouter *Router) {
	r.mux.Mount(path, subrouter)
	subrouter.prefix = path
	r.subrouters = append(r.subrouters, subrouter)
}
