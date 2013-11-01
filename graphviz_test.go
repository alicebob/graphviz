package graphviz

import (
	// "fmt"
	"testing"
)

func TestIt(t *testing.T) {
	g := MakeGraph()
	defer g.Close()

	g.Node("foo")
	g.Node("bar")
	g.Node("baz")
	// fmt.Printf("Graph: %+v\n", g)
	// g.Node("baz") // breaks?
	g.Edge("foo", "bar")
	g.Edge("baz", "bar")
	// g.Edge("baz", "baq") // breaks
	sub1 := g.Subgraph("my sub")
	sub1.Rank("same")
	sub1.Node("foo")
	sub1.Node("bar")
	sub2 := g.Subgraph("my 2nd sub")
	sub2.Rank("source")
	sub2.Node("baz")
	g.Layout()
	x, y, err := g.Pos("foo")
	if err != nil {
		t.Fatalf("%s", err)
	}
	if x != 78 || y != 18 {
		t.Fatalf("Wrong x/y for foo: %v,%v", x, y)
	}

}
