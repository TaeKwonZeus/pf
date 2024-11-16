package pf

import (
	"maps"
	"sort"
	"strings"
)

type endpoints = map[string]*anyHandler

type leafNode struct {
	key       string
	endpoints endpoints
}

type edge struct {
	label byte
	node  *node
}

type edges []edge

func (e edges) Len() int {
	return len(e)
}

func (e edges) Less(i, j int) bool {
	return e[i].label < e[j].label
}

func (e edges) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e edges) Sort() {
	sort.Sort(e)
}

// longestPrefix finds the length of the shared prefix
// of two strings
func longestPrefix(k1, k2 string) int {
	minLen := min(len(k1), len(k2))

	for i := 0; i < min(len(k1), len(k2)); i++ {
		if k1[i] != k2[i] {
			return i
		}
	}

	return minLen
}

type node struct {
	leaf *leafNode

	prefix string

	edges []edge
}

func (n *node) isLeaf() bool {
	return n.leaf != nil
}

func (n *node) addEdge(e edge) {
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool {
		return n.edges[i].label >= e.label
	})

	n.edges = append(n.edges, edge{})
	copy(n.edges[idx+1:], n.edges[idx:])
	n.edges[idx] = e
}

func (n *node) updateEdge(label byte, node *node) {
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool {
		return n.edges[i].label >= label
	})
	if idx < num && n.edges[idx].label == label {
		n.edges[idx].node = node
		return
	}
	panic("replacing missing edge")
}

func (n *node) getEdge(label byte) *node {
	num := len(n.edges)
	idx := sort.Search(num, func(i int) bool {
		return n.edges[i].label >= label
	})
	if idx < num && n.edges[idx].label == label {
		return n.edges[idx].node
	}
	return nil
}

func (n *node) insert(method string, path string, handler *anyHandler) {
	// TODO split the handler types like in chi and add wildcards when leaf already exists
	var parent *node
	traverse := n
	search := path

	for {
		// Handle key exhaustion
		if len(search) == 0 {
			if traverse.isLeaf() {
				traverse.leaf.endpoints[method] = handler
				return
			}

			traverse.leaf = &leafNode{
				key:       path,
				endpoints: endpoints{method: handler},
			}
			return
		}

		// Look for the edge
		parent = traverse
		traverse = traverse.getEdge(search[0])

		// No edge, create one
		if traverse == nil {
			e := edge{
				label: search[0],
				node: &node{
					leaf: &leafNode{
						key:       path,
						endpoints: endpoints{method: handler},
					},
					prefix: search,
				},
			}
			parent.addEdge(e)
			return
		}

		// Determine the longest prefix of the search key on match
		commonPrefix := longestPrefix(search, traverse.prefix)
		if commonPrefix == len(traverse.prefix) {
			search = search[commonPrefix:]
			continue
		}

		child := &node{
			prefix: search[:commonPrefix],
		}
		parent.updateEdge(search[0], child)

		// Restore the existing node
		child.addEdge(edge{
			label: traverse.prefix[commonPrefix],
			node:  traverse,
		})
		traverse.prefix = traverse.prefix[commonPrefix:]

		// Create a new leaf node
		leaf := &leafNode{
			key:       path,
			endpoints: endpoints{method: handler},
		}

		// If the new key is a subset, add to this node
		search = search[commonPrefix:]
		if len(search) == 0 {
			child.leaf = leaf
			return
		}

		// Create a new edge for the node
		child.addEdge(edge{
			label: search[0],
			node: &node{
				leaf:   leaf,
				prefix: search,
			},
		})
		return
	}
}

// get finds a handler with the specified method and path or returns nil
func (n *node) get(method string, path string) *anyHandler {
	traverse := n
	search := path

	for {
		// Check for key exhaustion
		if len(search) == 0 {
			if traverse.isLeaf() {
				return traverse.leaf.endpoints[method]
			}
			break
		}

		// Look for an edge
		traverse = traverse.getEdge(search[0])
		if traverse == nil {
			break
		}

		// Consume the search prefix
		if strings.HasPrefix(search, traverse.prefix) {
			search = search[len(traverse.prefix):]
		} else {
			break
		}
	}

	return nil
}

// traverse gets a map with full paths as keys and handlers as values
func (n *node) traverse() map[string]endpoints {
	// return endpoints if leaf node
	if n.leaf != nil {
		return map[string]endpoints{n.leaf.key: n.leaf.endpoints}
	}

	out := make(map[string]endpoints)

	for _, e := range n.edges {
		nodeMap := e.node.traverse()

		// Add all edges from child nodes
		maps.Insert(out, maps.All(nodeMap))
	}

	return out
}
