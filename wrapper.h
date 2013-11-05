// cgo has problems with some of the selfreferential Graphviz structs. To work
// around that we wrap everything here with void pointers.

#ifndef _WRAPPER_H
#define _WRAPPER_H 1

#include <stdlib.h>
#include <gvc.h>

#ifdef __cplusplus
extern "C" {
#endif

void* makeGraph();
// subgraph, node, edge, set, and pos work on either Graph.graph or a subgraph.
pointf pos(void*);

#ifdef __cplusplus
}
#endif
#endif				/* _WRAPPER_H */
