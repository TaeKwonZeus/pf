package pf

import (
	"fmt"
	"net/http"
	"reflect"
)

// Handler represents an HTTP callback. Handler takes in a parsed Request
// (optionally, to do nothing set Req to struct{}), as well as a
// ResponseWriter that may marshal the response. Request and ResponseWriter
// contain the underlying http.Request and ResponseWriter from [net/http].
type Handler[Req, Res any] func(w ResponseWriter[Res], r *Request[Req]) error

type handlerSignature struct {
	reqType reflect.Type
	resType reflect.Type

	props []HandlerProperty
}

func (h Handler[Req, Res]) wrap(props []HandlerProperty) (http.HandlerFunc, *handlerSignature) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		req, err := parseRequest[Req](r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusBadRequest)
			return
		}

		err = h(ResponseWriter[Res]{w}, req)
		if err != nil {
			HandleError(w, err)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	}

	return handler, &handlerSignature{
		reqType: reflect.TypeFor[Req](),
		resType: reflect.TypeFor[Res](),
		props:   props,
	}
}
