package pf

import (
	"encoding/json"
	"net/http"
)

// ResponseWriter wraps an http.ResponseWriter instance, adding convenience
// methods for marshaling the output.
type ResponseWriter[T any] struct {
	http.ResponseWriter
}

// OK marshals response as JSON and sends an HTTP response with status code 200.
func (w *ResponseWriter[T]) OK(response T) error {
	return w.JSON(http.StatusOK, response)
}

// JSON marshals response as JSON and sends an HTTP response with the status
// code specified by status.
func (w *ResponseWriter[T]) JSON(status int, response T) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w.ResponseWriter).Encode(response)
}
