// Selective Go bindings for Graphviz

package graphviz

/*
#cgo CFLAGS: -DDOT_ONLY=1
#cgo pkg-config: libgvc
#include "wrapper.h"
*/
import "C"

import (
	"fmt"
	"unsafe"
)

type graph struct {
	graph        *C.Graph
	nodes        map[string]unsafe.Pointer
	layoutCalled bool
}
type Subgraph struct {
	graph unsafe.Pointer
}

func MakeGraph() graph {
	return graph{
		graph:        C.makeGraph(),
		nodes:        map[string]unsafe.Pointer{},
		layoutCalled: false,
	}
}

// destroy
func (g *graph) Close() {
	C.freeGraph(g.graph)
	g.graph = nil
}

// Node adds a named node. Name should be unique.
func (g *graph) Node(id string) {
	if g.layoutCalled {
		panic("Can't add nodes after calling layout()")
	}
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))

	g.nodes[id] = C.node(g.graph.graph, cid)
}

func (subg *Subgraph) Node(id string) {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))
	C.node(subg.graph, cid)
}

// Edge adds a directed edge.
// The endpoints have to have been added with Node() before.
func (g *graph) Edge(fromID, toID string) {
	if g.layoutCalled {
		panic("Can't add nodes after calling layout()")
	}
	from := g.nodes[fromID]
	to := g.nodes[toID]
	if from == nil || to == nil {
		panic("Unknown node id")
	}
	C.edge(g.graph.graph, from, to)
}

// Subgraph creates a subgraph
func (g *graph) Subgraph(name string) Subgraph {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return Subgraph{
		graph: C.subgraph(g.graph.graph, cname),
	}
}

// Rank sets the 'rank' attribute. e.g. 'source'
func (g *graph) Rank(rank string) {
	if g.layoutCalled {
		panic("Can't call Rank after calling Layout()")
	}
	cname := C.CString("rank")
	defer C.free(unsafe.Pointer(cname))
	crank := C.CString(rank)
	defer C.free(unsafe.Pointer(crank))
	C.set(g.graph.graph, cname, crank)
}

func (subg *Subgraph) Rank(rank string) {
	cname := C.CString("rank")
	defer C.free(unsafe.Pointer(cname))
	crank := C.CString(rank)
	defer C.free(unsafe.Pointer(crank))
	C.set(subg.graph, cname, crank)
}

// Layout does all the calculations. Pos() will be ready after this.
func (g *graph) Layout() {
	C.layout(g.graph)
	g.layoutCalled = true
}

// Pos only works after Layout()
func (g *graph) Pos(id string) (float32, float32, error) {
	if !g.layoutCalled {
		panic("Can't use pos() before calling layout()")
	}
	node, ok := g.nodes[id]
	if !ok {
		return 0, 0, fmt.Errorf("no such node")
	}
	// pos := C.GoString(C.pos(node))
	pos := C.pos(node)
	// fmt.Printf("Pos for %s: '%.0f %.0f'\n", id, pos.x, pos.y)
	return float32(pos.x), float32(pos.y), nil
}
