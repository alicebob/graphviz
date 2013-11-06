Very selective Go bindings for graph layout with Graphviz.

Input via objects, no dot file.

Output is node locations only, no rendered pictures.

Only directed non-strict graph is supported, and only dot layout is supported.

Tested on graphviz 2.34.0.

Build:

    $ go get https://github.com/alicebob/graphviz

	On Mavericks you need to use apple's gcc4.2
	$ brew install apple-gcc42
	$ cd /usr/bin
	$ sudo mv gcc gcc_mavs
	$ sudo ln -s /usr/local/Cellar/apple-gcc42/4.2.1-5666.3/bin/gcc-4.2 gcc

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
	node_bar := sub1.Node("bar")
	node_bar.Set("label", "A bar")
	sub1.Node("baz")

	positions := g.Layout()

	fmt.Printf("Foo: %v,%v\n", positions["foo"].X, positions["foo"].Y)
    ... &c

