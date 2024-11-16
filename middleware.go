package pf

import "net/http"

type Middleware func(next http.Handler) (http.Handler, error)

type Middlewares []Middleware

func (m Middlewares) Handler(h http.Handler) http.Handler {
	for _, middleware := range m {
		h = middleware.wrap()(h)
	}
	return h
}

func (m Middleware) wrap() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h, err := m(next)
			if err != nil {
				handleError(w, err)
				return
			}

			h.ServeHTTP(w, r)
		})
	}
}
