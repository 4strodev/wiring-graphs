package main

import (
	"log"

	"github.com/4strodev/wiring/pkg/collections/graph"
)

func main() {
	var numbers = []int{1, 2, 3, 4, 5}
	var nodes = []*graph.Node[int]{}
	var g = graph.NewGraph[int]()

	for _, n := range numbers {
		node := graph.NewNode(n)
		nodes = append(nodes, node)
		g.Add(node)
	}

	g.Connect(nodes[0], nodes[1], graph.OUT)
	g.Connect(nodes[1], nodes[2], graph.OUT)
	g.Connect(nodes[2], nodes[0], graph.OUT)

	cicle, found := g.DetectCircularRelations()
	if !found {
		log.Fatal("No cicle detected")
	} else {
		log.Println("Cicle detected", cicle)
	}

}
