package main

import (
	"log"
	"net/http"

	"github.com/TaeKwonZeus/pf"
	"github.com/go-chi/chi/v5/middleware"
)

// Declare response types as structs
type PingResponse struct {
	Message string `json:"message"`
}

// Define handlers with generic parameters
func PingHandler(w pf.ResponseWriter[PingResponse], r *pf.Request[struct{}]) error {
	return w.OK(PingResponse{"Pong!"})
}

// Declare request types as structs
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

	// Use middlewares
	pf.Use(r, middleware.Logger, middleware.Recoverer)

	pf.Get(r, "/ping", PingHandler)

	// Subrouting
	pf.Route(r, "/wat", func(r *pf.Router) {
		// Add metadata for Swagger
		pf.Post(r, "/uploadbeer", LogHandler, pf.WithSummary("This shi logs"))
	})

	// Create a Swagger endpoint
	pf.AddSwagger(r, "/swagger", &pf.SwaggerInfo{
		Title:   "bruh",
		Version: "v0.0.1",
	})

	// Listen with net/http
	log.Println("Listening on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
