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

void
closeGraph(void* g)
{
	agclose(g);
}

void*
makeSubgraph(void* g, char* name)
{
	return agsubg((Agraph_t*) g, name, 1);
}

void*
makeEdge(void* g, void* from, void* to)
{
	return agedge(g, from, to, NULL, 1);
}

void*
makeNode(void* g, char* name)
{
	return agnode(g, name, 1);
}

pointf
pos(void* node)
{
    return ND_coord((Agnode_t*) node);
}

*/
import "C"

import (
	"unsafe"
)

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
	gvc := C.gvContext() // does an implicit aginit()
	return Graph{
		G:     G{graph: C.makeGraph()},
		gvc:   gvc,
		nodes: map[string]unsafe.Pointer{},
	}
}

// destroy
func (g *Graph) Close() {
	C.gvFreeLayout(g.gvc, (*C.graph_t)(g.graph))
	C.closeGraph(g.graph)
	C.gvFreeContext(g.gvc)

	g.graph = nil
	g.gvc = nil
}

// Node adds a named node. Name should be unique.
func (g *Graph) Node(id string) Node {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))

	node := C.makeNode(g.graph, cid)
	g.nodes[id] = node
	return Node{G: G{graph: node}}
}

// Note: you'll probably want to add it to the main graph as well, otherwise you
// won't get a position for it back from Layout().
func (subg *Subgraph) Node(id string) Node {
	cid := C.CString(id)
	defer C.free(unsafe.Pointer(cid))

	node := C.makeNode(subg.graph, cid)
	return Node{G: G{graph: node}}
}

// Edge adds a directed edge.
// The endpoints have to have been added with Node() before.
func (g *Graph) Edge(fromID, toID string) Edge {
	from := g.nodes[fromID]
	to := g.nodes[toID]
	if from == nil || to == nil {
		panic("Unknown node id")
	}
	edge := C.makeEdge(g.graph, from, to)
	return Edge{
		G: G{graph: edge},
	}
}

// Subgraph creates a subgraph
func (g *Graph) Subgraph(name string) Subgraph {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	sub := C.makeSubgraph(g.graph, cname)

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
	C.agsafeset(g.graph, cattr, cvalue, cempty)
}

// Layout does all the calculations. Pos() will be ready after this.
func (g *Graph) Layout() map[string]Pos {
	ctype := C.CString("dot")
	defer C.free(unsafe.Pointer(ctype))
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
