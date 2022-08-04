package graph

import "fmt"

// Implementation of the Graph interface
type AdjacencyListGraph[N any, E IHalfEdge] struct {
	nodes     []N   // stores the nodes
	edges     [][]E // adjacency list: stores the leaving edges for each node
	edgeCount int   // caches the number of edges in the graph
}

// NodeCount implements Graph.NodeCount
func (alg *AdjacencyListGraph[N, E]) NodeCount() int {
	return len(alg.nodes)
}

// EdgeCount implements Graph.EdgeCount
func (alg *AdjacencyListGraph[N, E]) EdgeCount() int {
	return alg.edgeCount
}

// GetNode implements Graph.GetNode
func (alg *AdjacencyListGraph[N, E]) GetNode(id NodeId) N {
	if id < 0 || id >= alg.NodeCount() {
		panic(fmt.Sprintf("AdjacencyListGraph does not contain a node with ID=%d.\n", id))
	}
	return alg.nodes[id]
}

// GetHalfEdgesFrom implements Graph.GetHalfEdgesFrom
func (alg *AdjacencyListGraph[N, E]) GetHalfEdgesFrom(id NodeId) []E {
	if id < 0 || id >= alg.NodeCount() {
		panic(fmt.Sprintf("AdjacencyListGraph does not contain a node with ID=%d.\n", id))
	}
	return alg.edges[id]
}

// AppendNode(n) adds node 'n' to the graph and assigns the next unused ID (i.e. the previous node count) to it.
// Additionally, the assigned ID is returned.
func (alg *AdjacencyListGraph[N, E]) AppendNode(n N) int {
	nodeId := alg.NodeCount()                   // cache the ID of the new node
	alg.nodes = append(alg.nodes, n)            // append node n
	alg.edges = append(alg.edges, make([]E, 0)) // initialize a list of leaving edges
	return nodeId                               // return cached ID
}

// InsertHalfEdge(tail, e) inserts a new half edge e from tail node with ID='tail' to the graph.
// The method fails iff either tail or head node of the edge do not exist.
// If the same edge already exists, nothing is changed, i.e. duplicate edges are ignored.
func (alg *AdjacencyListGraph[N, E]) InsertHalfEdge(tail NodeId, e E) {
	if tail < 0 || tail >= alg.NodeCount() {
		panic(fmt.Sprintf("AdjacencyListGraph does not contain the tail node with ID=%d.\n", tail))
	}
	if e.To() < 0 || e.To() >= alg.NodeCount() {
		panic(fmt.Sprintf("AdjacencyListGraph does not contain the head node with ID=%d of the edge %v.\n", e.To(), e))
	}
	// check for duplicates
	for _, leavingEdge := range alg.edges[tail] {
		if e.To() == leavingEdge.To() {
			return // ignore duplicate edges
		}
	}
	alg.edges[tail] = append(alg.edges[tail], e)
	alg.edgeCount++
}
