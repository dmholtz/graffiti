package shortest_path

import (
	g "github.com/dmholtz/graffiti/graph"
)

// Efficient implementation of Dijkstra's Algorithm for finding shortest paths between nodes in a graph using a priority queue.
func Dijkstra[N any, E g.IWeightedHalfEdge[W], W g.Weight](graph g.Graph[N, E], source, target g.NodeId, recordSearchSpace bool) ShortestPathResult[W] {
	return ShortestPathResult[W]{}
}
