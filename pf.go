package pf

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
)

type Handler[Req, Res any] func(w ResponseWriter[Res], r *Request[Req]) error

func (h Handler[Req, Res]) wrap() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

type Middleware func(next http.Handler) (http.Handler, error)

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

func Get[Res any](r *Router, path string, handler Handler[struct{}, Res]) {
}

func Post[Req, Res any](r *Router, path string, handler Handler[Req, Res]) {
}

func Put[Req, Res any](r *Router, path string, handler Handler[Req, Res]) {
}

func Delete[Req, Res any](r *Router, path string, handler Handler[Req, Res]) {
}

func Patch[Req, Res any](r *Router, path string, handler Handler[Req, Res]) {
}

func Head(r *Router, path string, handler Handler[struct{}, struct{}]) {
}

func Options[Res any](r *Router, path string, handler Handler[struct{}, Res]) {
}

func Trace[Res any](r *Router, path string, handler Handler[struct{}, Res]) {
}

func Connect(r *Router, path string, handler Handler[struct{}, struct{}]) {
}

func Route(r *Router, path string, fn func(r *Router)) {
}

func Handle(r *Router, path string, handler http.Handler) {
}

func Mount(r *Router, path string, handler http.Handler) {
}
