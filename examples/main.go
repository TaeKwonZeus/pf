package main

import (
	"log"
	"net/http"

	"github.com/TaeKwonZeus/pf"
	"github.com/go-chi/chi/v5/middleware"
)

func PingHandler(w pf.ResponseWriter[string], r *pf.Request[struct{}]) error {
	return w.OK("Pong!")
}

type LogHandlerRequest struct {
	Beer   string `json:"beer"`
	Volume int    `json:"volume"`
}

func LogHandler(w pf.ResponseWriter[struct{}], r *pf.Request[LogHandlerRequest]) error {
	log.Println("Received a record:", r.Body)
	return nil
}

func main() {
	r := pf.NewRouter()

	pf.Use(r, middleware.Logger, middleware.Recoverer)

	pf.Get(r, "/ping", PingHandler)
	pf.Route(r, "/wat", func(r *pf.Router) {
		pf.Post(r, "/uploadbeer", LogHandler, pf.WithSummary("This shi logs"))
	})

	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
