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

func NewAdjacencyMap[N any, E IHalfEdge]() *AdjacencyMapGraph[N, E] {
	amg := AdjacencyMapGraph[N, E]{}
	amg.Nodes = make(map[NodeId]N)
	amg.Edges = make(map[NodeId][]E)
	return &amg
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

// AddNode(id, node) adds or overwrites a node with ID=id.
func (amg *AdjacencyMapGraph[N, E]) AddNode(id NodeId, n N) {
	amg.Nodes[id] = n
	if _, ok := amg.Edges[id]; !ok {
		amg.Edges[id] = make([]E, 0)
	}
}

// InsertHalfEdge(tail, e) inserts a new half edge e from tail node with ID='tail' to the graph.
// The method fails iff either tail or head node of the edge do not exist.
// If the same edge already exists, nothing is changed, i.e. duplicate edges are ignored.
func (amg *AdjacencyMapGraph[N, E]) InsertHalfEdge(tail NodeId, e E) {
	if _, ok := amg.Nodes[tail]; !ok {
		panic(fmt.Sprintf("AdjacencyMapGraph does not contain the tail node with ID=%d.\n", tail))
	}
	if _, ok := amg.Nodes[e.To()]; !ok {
		panic(fmt.Sprintf("AdjacencyMapGraph does not contain the head node with ID=%d of the edge %v.\n", e.To(), e))
	}
	// check for duplicates
	for _, leavingEdge := range amg.Edges[tail] {
		if e.To() == leavingEdge.To() {
			return // ignore duplicate edges
		}
	}
	amg.Edges[tail] = append(amg.Edges[tail], e)
	amg.EdgeCount_++
}
