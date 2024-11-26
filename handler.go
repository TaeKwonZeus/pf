package pf

import (
	"fmt"
	"net/http"
	"reflect"
)

type Handler[Req, Res any] func(w ResponseWriter[Res], r *Request[Req]) error

type handlerSignature struct {
	reqType reflect.Type
	resType reflect.Type

	props []HandlerProperty
}

func (h Handler[Req, Res]) wrap(props []HandlerProperty) (http.HandlerFunc, *handlerSignature) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		req, err := ParseRequest[Req](r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusBadRequest)
			return
		}

		err = h(ResponseWriter[Res]{w}, req)
		if err != nil {
			HandleError(w, err)
		}
	}

	var emptyReq [0]Req
	var emptyRes [0]Res

	return handler, &handlerSignature{
		reqType: reflect.TypeOf(emptyReq).Elem(),
		resType: reflect.TypeOf(emptyRes).Elem(),
		props:   props,
	}
}
