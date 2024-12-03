package pf

import (
	"testing"
)

type TestRequest struct {
	A struct{ Name string }
	B *TestRequestField
	C int
}

type TestRequestField struct {
	Desc string
}

type TestResponse struct {
	A struct{ Name string }
	B *TestResponseField
	C int
}

type TestResponseField struct {
	Desc string
}

func Ping(w ResponseWriter[TestResponse], r *Request[TestRequest]) error {
	return w.OK(TestResponse{})
}

func TestGenerateSpec(t *testing.T) {
	r := NewRouter()
	Post(r, "/get", Ping)

	bytes, err := generateSpec(
		r.traverseSignatures(),
		&SwaggerInfo{
			Title:   "PABLO",
			Version: "v0.0.0.0.0.0.0.1",
		},
	).MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(bytes))
}
