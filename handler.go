package pf

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"reflect"
)

type Handler[Req, Res any] func(w ResponseWriter[Res], r *Request[Req]) error

type anyHandler struct {
	handler http.Handler

	reqType reflect.Type
	resType reflect.Type
}

func (h Handler[Req, Res]) wrap() *anyHandler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		req, err := ParseRequest[Req](r)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse request body: %v", err), http.StatusBadRequest)
			return
		}

		err = h(ResponseWriter[Res]{w}, req)
		if err != nil {
			handleError(w, err)
		}
	}

	var emptyReq [0]Req
	var emptyRes [0]Res

	return &anyHandler{
		handler: http.HandlerFunc(handler),
		reqType: reflect.TypeOf(emptyReq).Elem(),
		resType: reflect.TypeOf(emptyRes).Elem(),
	}
}

func handleError(w http.ResponseWriter, err error) {
	var httpErr httpError
	if errors.As(err, &httpErr) {
		http.Error(w, httpErr.Error(), int(httpErr))
		return
	}

	if err != nil {
		slog.Error("Error in handler", "err", err.Error())
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}
