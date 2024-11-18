package pf

import "github.com/go-chi/chi/v5"

func URLParam[T any](r Request[T], key string) string {
	return chi.URLParam(r.Request, key)
}
