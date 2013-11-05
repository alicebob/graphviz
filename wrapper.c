#include "wrapper.h"

#ifdef __cplusplus
extern "C" {
#endif

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

#ifdef __cplusplus
}
#endif
