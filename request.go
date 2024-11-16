package pf

import (
	"encoding/json"
	"net/http"
)

type Request[T any] struct {
	*http.Request
	Body T
}

func ParseRequest[T any](r *http.Request) (req *Request[T], err error) {
	req = &Request[T]{Request: r}
	switch any(req.Body).(type) {
	case struct{}:
		return req, nil
	default:
		err = json.NewDecoder(r.Body).Decode(req.Body)
		if err != nil {
			return nil, err
		}
		return req, nil
	}
}
