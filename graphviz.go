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
pos(Agnode_t* node)
{
    return ND_coord(node);
}

*/
import "C"

import (
	"fmt"
	"sync"
	"unsafe"
)

var mutex sync.Mutex
var gvc *C.GVC_t

func init() {
	mutex = sync.Mutex{}
	gvc = C.gvContext() // does an implicit aginit()
}

type G struct {
	graph unsafe.Pointer // Used for graphs, nodes, edges, &c.
}
type Graph struct {
	G
}
type Subgraph struct {
	G
	main *Graph
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

	// gvc := C.gvContext() // does an implicit aginit()
	return Graph{
		G: G{graph: C.makeGraph()},
	}
}

// destroy
func (g *Graph) Close() {
	mutex.Lock()
	defer mutex.Unlock()
	C.agclose((*C.Agraph_t)(g.graph))
	// C.gvFreeContext(g.gvc)

	g.graph = nil
	// g.gvc = nil
}

// Node adds a named node. Will return an old one if it existed.
func (g *Graph) Node(id string) Node {
	mutex.Lock()
	defer mutex.Unlock()

	node := g.node(id)
	return Node{G: G{graph: node}}
}

// Internal node lookup. Needs to be locked!
func (g *Graph) node(id string) unsafe.Pointer {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))

	return unsafe.Pointer(C.agnode((*C.Agraph_t)(g.graph), cid, 1 /* create */))
}

func (subg *Subgraph) Node(id string) Node {
	mutex.Lock()
	defer mutex.Unlock()

	node := subg.main.node(id)

	subnode := unsafe.Pointer(C.agsubnode((*C.Agraph_t)(subg.graph), (*C.Agnode_t)(node), 1 /* create */))
	// node := unsafe.Pointer(C.agnode((*C.Agraph_t)(subg.graph), cid, 1 /* create */))
	return Node{G: G{graph: subnode}}
}

// Edge adds a directed edge.
// The endpoints have to have been added with Node() before.
func (g *Graph) Edge(fromID, toID string) Edge {
	mutex.Lock()
	defer mutex.Unlock()

	from := g.node(fromID)
	to := g.node(toID)
	if from == nil {
		panic(fmt.Sprintf("unknown node id '%v'", fromID))
	}
	if to == nil {
		panic(fmt.Sprintf("unknown node id '%v'", toID))
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

	mutex.Lock()
	defer mutex.Unlock()
	sub := unsafe.Pointer(C.agsubg((*C.Agraph_t)(g.graph), cname, 1 /* create */))

	return Subgraph{
		G:    G{graph: sub},
		main: g,
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

// Layout does all the calculations. Returns the positions of all nodes.
func (g *Graph) Layout() map[string]Pos {
	ctype := C.CString("dot")
	defer C.free(unsafe.Pointer(ctype))

	mutex.Lock()
	defer mutex.Unlock()
	C.gvLayout(gvc, (*C.graph_t)(g.graph), ctype)

	positions := map[string]Pos{}
	node := C.agfstnode((*C.Agraph_t)(g.graph))
	for node != nil {
		pos := C.pos(node)
		name := C.GoString(C.agnameof(unsafe.Pointer(node)))
		positions[name] = Pos{
			X: float32(pos.x),
			Y: float32(pos.y),
		}
		node = C.agnxtnode((*C.Agraph_t)(g.graph), (*C.Agnode_t)(node))
	}
	C.gvFreeLayout(gvc, (*C.graph_t)(g.graph))
	return positions
}
