Very selective Go bindings for graph layout with Graphviz.
Input via objects, no dot file.
Output is node locations only, no rendered pictures.
Only directed non-strict graph is supported, only dot layout is supported.

Usage:

	g := MakeGraph()
	defer g.Close()

	g.Node("foo")
	g.Node("bar")
	g.Node("baz")
	g.Node("bat")
	g.Edge("foo", "bar")
	sub1 := g.Subgraph("my sub")
	sub1.Rank("same")
	sub1.Node("bar")
	sub1.Node("baz")

	g.Layout()

	x, y, err := g.Pos("foo")
	x, y, err := g.Pos("bat")
    ... &c

