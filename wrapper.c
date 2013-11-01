#include "wrapper.h"

#ifdef __cplusplus
extern "C" {
#endif

Graph*
makeGraph()
{
	Graph* g = malloc(sizeof(Graph));
	g->gvc = gvContext();	// does an implicit aginit()
	g->graph = agopen("no name", AGDIGRAPH);
	return g;
}

void
freeGraph(Graph* g)
{
	gvFreeLayout(g->gvc, g->graph);
	agclose(g->graph);
	g->graph = NULL;
	gvFreeContext(g->gvc);
	g->gvc = NULL;
    free(g);
}

void*
subgraph(void* g, char* name)
{
    Agraph_t* sub = agsubg(g, name);
    return sub;
}

void
layout(Graph* g)
{
    gvLayout(g->gvc, g->graph, "dot");
}

void*
set(void* g, char* attr, char* value)
{
    agsafeset(g, attr, value, "");
}

void*
node(void* g, char* id)
{
	void* n = agnode(g, id);
	// Make sure all nodes are rendered with the same width.
	set(n, "label", "examplename");
    return n;
}

void
edge(void* g, void* from, void* to)
{
    // return value ignored
    agedge(g, (Agnode_t*)from, (Agnode_t*)to);
}


pointf
pos(void* node)
{
    return ND_coord((Agnode_t*) node);
}

#ifdef __cplusplus
}
#endif
