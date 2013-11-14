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
	pos := g.Layout()
	pos_foo := pos["foo"]
	if pos_foo.X < 0 || pos_foo.X > 200 {
		t.Fatalf("Wrong x for foo: %v", pos_foo.X)
	}
	if pos_foo.Y != 18 {
		t.Fatalf("Wrong y for foo: %v", pos_foo.Y)
	}

}

func TestSub(t *testing.T) {
	g := MakeGraph()
	defer g.Close()

	g.Node("foo")
	sub := g.Subgraph("go sub")
	sub.Node("bar")
	pos := g.Layout()
	pos_foo := pos["foo"]
	if pos_foo.X < 0 || pos_foo.X > 200 {
		t.Fatalf("Wrong x for foo: %v", pos_foo.X)
	}
}
