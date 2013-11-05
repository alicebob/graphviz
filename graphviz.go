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

type G struct {
	graph unsafe.Pointer // Used for graphs, nodes, edges, &c.
}
type Graph struct {
	G
	gvc          *[0]byte
	nodes        map[string]unsafe.Pointer
	layoutCalled bool
}
type Subgraph struct {
	G
}
type Node struct {
	G
}
type Edge struct {
	G
}

func MakeGraph() Graph {
	gvc := C.gvContext() // does an implicit aginit()
	return Graph{
		G:            G{graph: C.makeGraph()},
		gvc:          gvc,
		nodes:        map[string]unsafe.Pointer{},
		layoutCalled: false,
	}
}

// destroy
func (g *Graph) Close() {
	C.gvFreeLayout(g.gvc, (*C.graph_t)(g.graph))
	C.agclose((*C.Agraph_t)(g.graph))
	C.gvFreeContext(g.gvc)

	// C.freeGraph(g.graph)
	g.graph = nil
	g.gvc = nil
}

// Node adds a named node. Name should be unique.
func (g *Graph) Node(id string) Node {
	if g.layoutCalled {
		panic("Can't add nodes after calling layout()")
	}
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))

	node := unsafe.Pointer(C.agnode((*C.Agraph_t)(g.graph), cid, 1 /* create */))

	g.nodes[id] = node
	return Node{G: G{graph: node}}
}

func (subg *Subgraph) Node(id string) Node {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))
	node := unsafe.Pointer(C.agnode((*C.Agraph_t)(subg.graph), cid, 1 /* create */))
	return Node{G: G{graph: node}}
}

// Edge adds a directed edge.
// The endpoints have to have been added with Node() before.
func (g *Graph) Edge(fromID, toID string) Edge {
	if g.layoutCalled {
		panic("Can't add nodes after calling layout()")
	}
	from := g.nodes[fromID]
	to := g.nodes[toID]
	if from == nil || to == nil {
		panic("Unknown node id")
	}
	edge := unsafe.Pointer(C.agedge((*C.Agraph_t)(g.graph), (*C.Agnode_t)(from),
		(*C.Agnode_t)(to), nil, 1 /* create */))
	return Edge{
		G: G{graph: edge},
	}
}

// Subgraph creates a subgraph
func (g *Graph) Subgraph(name string) Subgraph {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	sub := unsafe.Pointer(C.agsubg((*C.Agraph_t)(g.graph), cname, 1 /* create */))

	return Subgraph{
		G: G{graph: sub},
	}
}

// Rank sets the 'rank' attribute. e.g. 'source'
func (g *Graph) Rank(rank string) {
	if g.layoutCalled {
		panic("Can't call Rank after calling Layout()")
	}
	g.Set("rank", rank)
}

func (subg *Subgraph) Rank(rank string) {
	subg.Set("rank", rank)
}

func (g *G) Set(attr string, value string) {
	cattr := C.CString(attr)
	defer C.free(unsafe.Pointer(cattr))
	cvalue := C.CString(value)
	defer C.free(unsafe.Pointer(cvalue))
	cempty := C.CString("")
	defer C.free(unsafe.Pointer(cempty))
	// C.set(g.graph, cname, cvalue)
	C.agsafeset(g.graph, cattr, cvalue, cempty)
}

// Layout does all the calculations. Pos() will be ready after this.
func (g *Graph) Layout() {
	// C.layout(g.gvc, unsafe.Pointer(g.graph))
	ctype := C.CString("dot")
	defer C.free(unsafe.Pointer(ctype))
	C.gvLayout(g.gvc, (*C.graph_t)(g.graph), ctype)
	g.layoutCalled = true
}

// Pos only works after Layout()
func (g *Graph) Pos(id string) (float32, float32, error) {
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
