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
	node_baz := g.Node("baz")
	node_baz.Set("label", "Baz!")
	// fmt.Printf("Graph: %+v\n", g)
	// g.Node("baz") // breaks?
	g.Edge("foo", "bar")
	edge := g.Edge("baz", "bar")
	edge.Set("shape", "record")
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
	if x < 0 || x > 200 {
		t.Fatalf("Wrong x for foo: %v", x)
	}
	if y != 18 {
		t.Fatalf("Wrong y for foo: %v", y)
	}

}
