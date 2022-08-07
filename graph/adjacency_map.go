package graph

import "fmt"

type AdjacencyMapGraph[N any, E IHalfEdge] struct {
	Nodes      map[NodeId]N
	Edges      map[NodeId][]E
	EdgeCount_ int
}

// NodeCount implements Graph.NodeCount
func (amg *AdjacencyMapGraph[N, E]) NodeCount() int {
	return len(amg.Nodes)
}

// EdgeCount implements Graph.EdgeCount
func (amg *AdjacencyMapGraph[N, E]) EdgeCount() int {
	return amg.EdgeCount_
}

// GetNode implements Graph.GetNode
func (amg *AdjacencyMapGraph[N, E]) GetNode(id NodeId) N {
	node, ok := amg.Nodes[id]
	if !ok {
		panic(fmt.Sprintf("AdjacencyMapGraph does not contain a node with ID=%d.\n", id))
	}
	return node
}

// GetHalfEdgesFrom implements Graph.GetHalfEdgesFrom
func (amg *AdjacencyMapGraph[N, E]) GetHalfEdgesFrom(id NodeId) []E {
	edges, ok := amg.Edges[id]
	if !ok {
		panic(fmt.Sprintf("AdjacencyMapGraph does not contain a node with ID=%d.\n", id))
	}
	return edges
}
