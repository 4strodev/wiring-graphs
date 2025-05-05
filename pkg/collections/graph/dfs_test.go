package graph

import "testing"

func TestDfsVisit(t *testing.T) {
	nodeA := NewNode("A")
	nodeB := NewNode("B")
	nodeC := NewNode("C")

	nodeA.Connect(nodeB, OUT)
	nodeB.Connect(nodeC, OUT)
	nodeC.Connect(nodeA, OUT)

	state := make(map[*Node[string]]dfsStates)
	parent := make(map[*Node[string]]*Node[string])
	cycle := []*Node[string]{}

	cycleDetected := dfsVisit(nodeA, state, parent, &cycle)
	if !cycleDetected {
		t.Fatalf("Cycle should be detected\n")
	}
}
