// pf (Pickme Framework) is a REST API library based on chi for use by me and
// my boys on hackathons, olympiads and for personal projects.
package pf

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type method = string

type signatures map[string]map[method]*handlerSignature

func (s signatures) add(path string, method method, signature *handlerSignature) {
	if s[path] == nil {
		s[path] = make(map[string]*handlerSignature)
	}
	s[path][method] = signature
}

// Router is a composable router based on chi.Mux that tracks handler request
// and response body signatures.
type Router struct {
	mux        chi.Router
	subrouters []*Router
	prefix     string

	signatures signatures
}

// NewRouter returns a newly initialized Router.
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

// Use appends one or more middlewares onto the Router stack.
func Use(r *Router, middlewares ...func(next http.Handler) http.Handler) {
	r.mux.Use(middlewares...)
}

// Method adds routes for path that matches the HTTP method specified by method.
// Method also adds metadata, consisting of the request and response type,
// as well as props, for use by Swagger and the like.
func Method[Req, Res any](r *Router, method method, path string, handler Handler[Req, Res], props ...HandlerProperty) {
	h, signature := handler.wrap(props)
	r.mux.Method(method, path, h)
	r.signatures.add(path, method, signature)
}

func Get[Res any](r *Router, path string, handler Handler[struct{}, Res], props ...HandlerProperty) {
	Method(r, http.MethodGet, path, handler, props...)
}

func Post[Req, Res any](r *Router, path string, handler Handler[Req, Res], props ...HandlerProperty) {
	Method(r, http.MethodPost, path, handler, props...)
}

func Put[Req, Res any](r *Router, path string, handler Handler[Req, Res], props ...HandlerProperty) {
	Method(r, http.MethodPut, path, handler, props...)
}

func Delete[Req, Res any](r *Router, path string, handler Handler[Req, Res], props ...HandlerProperty) {
	Method(r, http.MethodDelete, path, handler, props...)
}

func Patch[Req, Res any](r *Router, path string, handler Handler[Req, Res], props ...HandlerProperty) {
	Method(r, http.MethodPatch, path, handler, props...)
}

func Head(r *Router, path string, handler Handler[struct{}, struct{}], props ...HandlerProperty) {
	Method(r, http.MethodHead, path, handler, props...)
}

func Options[Res any](r *Router, path string, handler Handler[struct{}, Res], props ...HandlerProperty) {
	Method(r, http.MethodOptions, path, handler, props...)
}

// Route mounts a sub-router along path.
func Route(r *Router, path string, fn func(r *Router)) {
	subrouter := NewRouter()
	fn(subrouter)
	Mount(r, path, subrouter)
}

// Handle routes handler to the path. You should not use this for regular
// function handlers, as they won't show up in Swagger.
func Handle(r *Router, path string, handler http.Handler) {
	r.mux.Handle(path, handler)
}

// Mount mounts a sub-router along path.
func Mount(r *Router, path string, subrouter *Router) {
	r.mux.Mount(path, subrouter)
	subrouter.prefix = path
	r.subrouters = append(r.subrouters, subrouter)
}
