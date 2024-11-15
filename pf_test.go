package pf

import (
	"log/slog"
	"net/http"
)

func Ping(w ResponseWriter[string], r *http.Request, body struct{}) error {
	return w.JSON("Pong!")
}

func main() {
	r := NewRouter()
	Get(r, "/ping", Ping)

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		slog.Error("Error starting server", "err", err)
	}
}
