package pf

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"testing"
)

func Ping(w ResponseWriter[string], r *Request[struct{}]) error {
	return w.OK("Pong!")
}

func main() {
	r := NewRouter()
	Get(r, "/ping", Ping)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		slog.Error("Error starting server", "err", err)
	}
}

func TestChi(t *testing.T) {
	r := chi.NewRouter()
	r.Get("/api/a", func(w http.ResponseWriter, r *http.Request) {
		t.Log("AAA")
	})
	r.Route("/api/", func(r chi.Router) {
		r.Use(middleware.RequestID)

		r.Get("/b", func(w http.ResponseWriter, r *http.Request) {
			t.Log("BBB")
		})
		r.Get("/c", func(w http.ResponseWriter, r *http.Request) {
			t.Log("CCC")
		})
	})
	http.ListenAndServe(":8080", r)
}
