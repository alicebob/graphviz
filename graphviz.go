// Selective Go bindings for Graphviz

package graphviz

/*
#cgo CFLAGS: -DDOT_ONLY=1
#cgo pkg-config: libgvc

#include <stdlib.h>
#include <gvc.h>

void*
makeGraph()
{
	return agopen("no name", Agdirected, NIL(Agdisc_t *));
}

pointf
pos(void* node)
{
    return ND_coord((Agnode_t*) node);
}

*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

var mutex sync.Mutex

func init() {
	mutex = sync.Mutex{}
}

type G struct {
	graph unsafe.Pointer // Used for graphs, nodes, edges, &c.
}
type Graph struct {
	G
	gvc   *C.GVC_t
	nodes map[string]unsafe.Pointer
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

type Pos struct {
	X, Y float32
}

func MakeGraph() Graph {
	mutex.Lock()
	defer mutex.Unlock()

	gvc := C.gvContext() // does an implicit aginit()
	return Graph{
		G:     G{graph: C.makeGraph()},
		gvc:   gvc,
		nodes: map[string]unsafe.Pointer{},
	}
}

// destroy
func (g *Graph) Close() {
	mutex.Lock()
	defer mutex.Unlock()
	C.gvFreeLayout(g.gvc, (*C.graph_t)(g.graph))
	C.agclose((*C.Agraph_t)(g.graph))
	C.gvFreeContext(g.gvc)

	g.graph = nil
	g.gvc = nil
}

// Node adds a named node. Name should be unique.
func (g *Graph) Node(id string) Node {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))

	mutex.Lock()
	defer mutex.Unlock()
	node := unsafe.Pointer(C.agnode((*C.Agraph_t)(g.graph), cid, 1 /* create */))

	g.nodes[id] = node
	return Node{G: G{graph: node}}
}

// Note: you'll probably want to add it to the main graph as well, otherwise you
// won't get a position for it back from Layout().
func (subg *Subgraph) Node(id string) Node {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))

	mutex.Lock()
	defer mutex.Unlock()
	node := unsafe.Pointer(C.agnode((*C.Agraph_t)(subg.graph), cid, 1 /* create */))
	return Node{G: G{graph: node}}
}

// Edge adds a directed edge.
// The endpoints have to have been added with Node() before.
func (g *Graph) Edge(fromID, toID string) Edge {
	from := g.nodes[fromID]
	to := g.nodes[toID]
	if from == nil {
		panic(fmt.Sprintf("unknown node id '%v'", fromID))
	}
	if to == nil {
		panic(fmt.Sprintf("unknown node id '%v'", toID))
	}

	mutex.Lock()
	defer mutex.Unlock()
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

	mutex.Lock()
	defer mutex.Unlock()
	sub := unsafe.Pointer(C.agsubg((*C.Agraph_t)(g.graph), cname, 1 /* create */))

	return Subgraph{
		G: G{graph: sub},
	}
}

// Rank sets the 'rank' attribute. e.g. 'source'
func (g *Graph) Rank(rank string) {
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

	mutex.Lock()
	defer mutex.Unlock()
	C.agsafeset(g.graph, cattr, cvalue, cempty)
}

// Layout does all the calculations. Pos() will be ready after this.
func (g *Graph) Layout() map[string]Pos {
	ctype := C.CString("dot")
	defer C.free(unsafe.Pointer(ctype))

	mutex.Lock()
	defer mutex.Unlock()
	C.gvLayout(g.gvc, (*C.graph_t)(g.graph), ctype)

	positions := map[string]Pos{}
	for id, node := range g.nodes {
		pos := C.pos(node)
		positions[id] = Pos{
			X: float32(pos.x),
			Y: float32(pos.y),
		}
	}
	return positions
}
