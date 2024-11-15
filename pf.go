package pf

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
)

type ResponseWriter[T any] struct {
	w http.ResponseWriter
}

func (w *ResponseWriter[T]) JSON(response T) error {
	w.w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w.w).Encode(response)
}

type Handler[Req, Res any] func(w ResponseWriter[Res], r *http.Request, body Req) error

func (h Handler[Req, Res]) wrap() http.HandlerFunc {
	var req [0]Req

	switch any(req).(type) {
	case [0]struct{}:
		return func(w http.ResponseWriter, r *http.Request) {
			var empty Req
			err := h(ResponseWriter[Res]{w}, r, empty)
			if err != nil {
				handleError(w, err)
			}
		}
	default:
		return func(w http.ResponseWriter, r *http.Request) {
			var payload Req
			err := json.NewDecoder(r.Body).Decode(&payload)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = h(ResponseWriter[Res]{w}, r, payload)
			if err != nil {
				handleError(w, err)
			}
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

func NewRouter() *Router { return nil }

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {}

func Use(r *Router, middlewares ...Middleware) {}

func Get[Res any](r *Router, path string, handler Handler[struct{}, Res]) {}

func Post[Req, Res any](r *Router, path string, handler Handler[Req, Res]) {}

func Put[Req, Res any](r *Router, path string, handler Handler[Req, Res]) {}

func Delete[Req, Res any](r *Router, path string, handler Handler[Req, Res]) {}

func Patch[Req, Res any](r *Router, path string, handler Handler[Req, Res]) {}

func Group(r *Router, path string, fn func(r *Router)) {}

func Handle(r *Router, path string, handler http.Handler) {}
