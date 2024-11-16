package pf

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"log/slog"
	"net/http"
	"reflect"
)

type Handler[Req, Res any] func(w ResponseWriter[Res], r *Request[Req]) error

type handlerSignature struct {
	reqType reflect.Type
	resType reflect.Type
}

func (h Handler[Req, Res]) wrap() (http.HandlerFunc, handlerSignature) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		req, err := ParseRequest[Req](r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusBadRequest)
			return
		}

		err = h(ResponseWriter[Res]{w}, req)
		if err != nil {
			handleError(w, err)
		}
	}

	var emptyReq [0]Req
	var emptyRes [0]Res

	return handler, handlerSignature{
		reqType: reflect.TypeOf(emptyReq).Elem(),
		resType: reflect.TypeOf(emptyRes).Elem(),
	}
}

type Middleware func(next http.Handler) (http.Handler, error)

type Middlewares []Middleware

func (m Middlewares) Handler(h http.Handler) http.Handler {
	chiMiddlewares := make(chi.Middlewares, 0, len(m))

	for _, middleware := range m {
		chiMiddlewares = append(chiMiddlewares, handleErrorMiddleware(middleware))
	}

	return chiMiddlewares.Handler(h)
}

func handleError(w http.ResponseWriter, err error) {
	var httpErr httpError
	if errors.As(err, &httpErr) {
		http.Error(w, httpErr.Error(), int(httpErr))
		return
	}

	if err != nil {
		slog.Error("Error in handler", "err", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

type Router struct {
	mux    chi.Router
	routes map[*http.HandlerFunc]handlerSignature
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {}

func handleErrorMiddleware(middleware Middleware) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h, err := middleware(next)
			if err != nil {
				handleError(w, err)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}

func Use(r *Router, middlewares ...Middleware) {
}

func Method[Req, Res any](r *Router, method string, path string, handler Handler[Req, Res]) {
	h, signature := handler.wrap()
	r.mux.Method(method, path, h)
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
