// cgo has problems with some of the selfreferential Graphviz structs. To work
// around that we wrap everything here with void pointers.

#ifndef _WRAPPER_H
#define _WRAPPER_H 1

#include <stdlib.h>
#include <gvc.h>

#ifdef __cplusplus
extern "C" {
#endif

struct Graph {
	GVC_t* gvc;
	void* graph; // That's an Agraph_t
};

typedef struct Graph Graph;

Graph* makeGraph();
void freeGraph(Graph*);
void layout(Graph*);
// subgraph, node, edge, set, and pos work on either Graph.graph or a subgraph.
void* subgraph(void*, char*);
void* node(void*, char*);
void edge(void*, void*, void*);
void set(void*, char*, char*);
pointf pos(void* node);

#ifdef __cplusplus
}
#endif
#endif				/* _WRAPPER_H */
