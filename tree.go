package pf

import "net/http"

type leaf struct {
	key       string
	endpoints map[string]http.HandlerFunc
}
