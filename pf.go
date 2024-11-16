package pf

import (
	"net/http"
)

type Router struct {
	tree *node
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {}

func Use(r *Router, middlewares ...Middleware) {
}

func Method[Req, Res any](r *Router, method string, path string, handler Handler[Req, Res]) {
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
}

func Handle(r *Router, path string, handler http.Handler) {
}

func Mount(r *Router, path string, handler http.Handler) {
}

type Context struct {
	URLParams map[string][]string
}
