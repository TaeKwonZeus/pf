package pf

import (
	"github.com/go-chi/chi/v5/middleware"
	"net/http"
	"testing"
)

func Ping(w ResponseWriter[string], r *Request[struct{}]) error {
	return w.OK("Pong!")
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	Use(r, middleware.Logger)
	Get(r, "/ping", Ping)
	Put(r, "/test/c", Ping)
	Route(r, "/test", func(r *Router) {
		Get(r, "/a", Ping)
		Post(r, "/b", Ping)
	})
	t.Log(r.traverseSignatures())
	http.ListenAndServe(":8080", r)
}

func TestChi(t *testing.T) {
}
