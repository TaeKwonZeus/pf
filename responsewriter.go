package pf

import (
	"encoding/json"
	"net/http"
)

type ResponseWriter[T any] struct {
	http.ResponseWriter
}

func (w *ResponseWriter[T]) OK(response T) error {
	return w.JSON(http.StatusOK, response)
}

func (w *ResponseWriter[T]) JSON(status int, response T) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w.ResponseWriter).Encode(response)
}
