package pf

import "testing"

func TestTree(t *testing.T) {
	tree := new(node)

	tree.insert("GET", "/api", &anyHandler{})

	tree.insert("POST", "/api/auth", &anyHandler{})

	tree.insert("PUT", "/api/auth", &anyHandler{})

	t.Log(tree.traverse())
}
