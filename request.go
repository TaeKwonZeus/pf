package pf

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
)

// Request wraps http.Request and provides a Body of type T. Body is parsed depending on the type:
// If T is struct{}, then Body is equal to struct{}{}.
// If Body is []byte, then the request body is read into Body.
// If Body is *multipart.Form, then the form data is fetched using ParseMultipartForm.
// Otherwise, the response body is assumed to be JSON and deserialized into Body.
type Request[T any] struct {
	*http.Request
	Body T
}

func ParseRequest[T any](r *http.Request) (*Request[T], error) {
	var body T
	switch any(body).(type) {
	case struct{}:
	case []byte:
		bytes, err := io.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		body = any(bytes)
	case *multipart.Form:
		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			return nil, err
		}
		body = any(r.Form)
	default:
		err := json.NewDecoder(r.Body).Decode(&body)
		if err != nil {
			return nil, err
		}
	}
	return &Request[T]{r, body}, nil
}
