package graph

import "fmt"

// The adjacency map is a graph data structure which stores for each node a list of its adjacent nodes.
// Unlike the adjacency list, the node IDs do not necessarily have to be in the range from 0 to n-1 but may be arbitrary (but unique) integers.
// This flexibility comes at the expense of linear space overhead in the number of nodes due to the underlying map datastructure.
//
// Implements the Graph interface
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
