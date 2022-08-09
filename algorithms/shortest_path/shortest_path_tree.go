package shortest_path

import (
	"container/heap"

	g "github.com/dmholtz/graffiti/graph"
)

// Dijkstra's algorithm spans a directed, acyclic search graph (tree) with the source node being the root (entry point) of this search graph.
// This is due to the fact that multiple shortest paths to the same node might exist.
//
// Note: The term 'search tree' refers to the output of the search and is used to avoid confusions with the input graph, i.e. the graph on which the search is conducted.
func ShortestPathTree[N any, E g.IWeightedHalfEdge[W], W g.Weight](graph g.Graph[N, E], source g.NodeId) ShortestPathTreeNode {
	dijkstraItems := make([]*ShortestPathTreePqItem[W], graph.NodeCount(), graph.NodeCount())
	dijkstraItems[source] = &ShortestPathTreePqItem[W]{Id: source, Priority: 0, Predecessors: make([]int, 0)}

	pq := make(ShortestPathTreePriorityQueue[W], 0)
	heap.Init(&pq)
	heap.Push(&pq, dijkstraItems[source])

	successors := make([]*ShortestPathTreeNode, graph.NodeCount(), graph.NodeCount())

	for len(pq) > 0 {
		currentPqItem := heap.Pop(&pq).(*ShortestPathTreePqItem[W])
		currentNodeId := currentPqItem.Id

		if currentNodeId != source {
			if successors[currentNodeId] == nil {
				successors[currentNodeId] = &ShortestPathTreeNode{Id: currentNodeId, Children: make([]*ShortestPathTreeNode, 0)}
			}
			for _, pred := range currentPqItem.Predecessors {
				if successors[pred] == nil {
					successors[pred] = &ShortestPathTreeNode{Id: pred, Children: make([]*ShortestPathTreeNode, 0)}
				}
				successors[pred].Children = append(successors[pred].Children, successors[currentNodeId])
			}
		}

		for _, edge := range graph.GetHalfEdgesFrom(currentNodeId) {
			successor := edge.To()

			if dijkstraItems[successor] == nil {
				newPriority := dijkstraItems[currentNodeId].Priority + edge.Weight()
				pqItem := ShortestPathTreePqItem[W]{Id: successor, Priority: newPriority, Predecessors: []int{currentNodeId}}
				dijkstraItems[successor] = &pqItem
				heap.Push(&pq, &pqItem)
			} else {
				if updatedDistance := dijkstraItems[currentNodeId].Priority + edge.Weight(); updatedDistance < dijkstraItems[successor].Priority {
					dijkstraItems[successor].Priority = updatedDistance
					heap.Fix(&pq, dijkstraItems[successor].index)
					// reset predecessors
					dijkstraItems[successor].Predecessors = []int{currentNodeId}
				} else if updatedDistance == dijkstraItems[successor].Priority {
					// add another predecessor
					dijkstraItems[successor].Predecessors = append(dijkstraItems[successor].Predecessors, currentNodeId)
				}
			}
		}
	}

	return *successors[source]
}
