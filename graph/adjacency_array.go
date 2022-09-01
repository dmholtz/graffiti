package graph

import "fmt"

// The adjacency array is a graph data structure that stores for each node a list of its adjacent nodes.
// Unlike AdjacencyListGraph, the adjacent nodes are kept in a single, flat slice, which is segmented by an additional offset slice.
// Therefore, AdjacencyArrayGraph reduces the overhead induced by nested slices at the expense of being a non-expandable static graph.
//
// Please note: This implementation requires node indices from 0 to n-1, where n is the node count of the graph.
// Graphs that do not fulfill this condition should be wrapped by a wrapper that maps nodes to node ids.
//
// Implements the Graph interface
type AdjacencyArrayGraph[N any, E IHalfEdge] struct {
	Nodes   []N   // stores the nodes
	Edges   []E   // adjacency array: flat representation of leaving edges
	Offsets []int // use values at index i, i+1 to obtain the segment of adjacent edges for the i-th node
}

// Create new AdjacencyArrayGraph as a snapshot from another Graph interface type
func NewAdjacencyArrayFromGraph[N any, E IHalfEdge](g Graph[N, E]) *AdjacencyArrayGraph[N, E] {
	nodes := make([]N, 0)
	edges := make([]E, 0)
	offsets := make([]int, g.NodeCount()+1, g.NodeCount()+1)

	for i := 0; i < g.NodeCount(); i++ {
		// copy node
		nodes = append(nodes, g.GetNode(i))

		// copy all leaving edges from node
		for _, halfEdge := range g.GetHalfEdgesFrom(i) {
			edges = append(edges, halfEdge)
		}

		// set end-segment (non-inclusive) offset
		// = set start-segment (inclusive) offset of the next node
		offsets[i+1] = len(edges)
	}

	return &AdjacencyArrayGraph[N, E]{Nodes: nodes, Edges: edges, Offsets: offsets}
}

// NodeCount implements Graph.NodeCount
func (aag *AdjacencyArrayGraph[N, E]) NodeCount() int {
	return len(aag.Nodes)
}

// EdgeCount implements Graph.EdgeCount
func (aag *AdjacencyArrayGraph[N, E]) EdgeCount() int {
	return len(aag.Edges)
}

// GetNode implements Graph.GetNode
func (alg *AdjacencyArrayGraph[N, E]) GetNode(id NodeId) N {
	if id < 0 || id >= alg.NodeCount() {
		panic(fmt.Sprintf("AdjacencyArrayGraph does not contain a node with ID=%d.\n", id))
	}
	return alg.Nodes[id]
}

// GetHalfEdgesFrom implements Graph.GetHalfEdgesFrom
func (aag *AdjacencyArrayGraph[N, E]) GetHalfEdgesFrom(id NodeId) []E {
	if id < 0 || id >= aag.NodeCount() {
		panic(fmt.Sprintf("AdjacencyListGraph does not contain a node with ID=%d.\n", id))
	}
	return aag.Edges[aag.Offsets[id]:aag.Offsets[id+1]]
}
