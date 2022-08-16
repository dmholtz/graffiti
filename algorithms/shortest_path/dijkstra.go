package shortest_path

import (
	"container/heap"

	g "github.com/dmholtz/graffiti/graph"
)

// Efficient implementation of Dijkstra's Algorithm for finding the shortest path between a source and a target node in the graph.
// Implementation is based on a priority queue.
func Dijkstra[N any, E g.IWeightedHalfEdge[W], W g.Weight](graph g.Graph[N, E], source, target g.NodeId, recordSearchSpace bool) ShortestPathResult[W] {
	var searchSpace []g.NodeId = nil
	if recordSearchSpace {
		searchSpace = make([]g.NodeId, 0)
	}

	dijkstraItems := make([]*DijkstraPqItem[W], graph.NodeCount(), graph.NodeCount())
	dijkstraItems[source] = &DijkstraPqItem[W]{Id: source, Priority: 0, Predecessor: -1}

	pq := make(DijkstraPriorityQueue[W], 0)
	heap.Init(&pq)
	heap.Push(&pq, dijkstraItems[source])

	pqPops := 0
	for len(pq) > 0 {
		currentPqItem := heap.Pop(&pq).(*DijkstraPqItem[W])
		currentNodeId := currentPqItem.Id
		pqPops++

		if recordSearchSpace {
			searchSpace = append(searchSpace, currentNodeId)
		}

		for _, edge := range graph.GetHalfEdgesFrom(currentNodeId) {
			successor := edge.To()

			if dijkstraItems[successor] == nil {
				newPriority := dijkstraItems[currentNodeId].Priority + edge.Weight()
				pqItem := DijkstraPqItem[W]{Id: successor, Priority: newPriority, Predecessor: currentNodeId}
				dijkstraItems[successor] = &pqItem
				heap.Push(&pq, &pqItem)
			} else {
				if updatedDistance := dijkstraItems[currentNodeId].Priority + edge.Weight(); updatedDistance < dijkstraItems[successor].Priority {
					dijkstraItems[successor].Priority = updatedDistance
					dijkstraItems[successor].Predecessor = currentNodeId
					heap.Fix(&pq, dijkstraItems[successor].index)
				}
			}
		}

		if currentNodeId == target {
			break
		}
	}

	res := ShortestPathResult[W]{Length: W(-1), Path: make([]g.NodeId, 0), PqPops: pqPops, SearchSpace: searchSpace}
	if dijkstraItems[target] != nil {
		res.Length = dijkstraItems[target].Priority
		for nodeId := target; nodeId != -1; nodeId = dijkstraItems[nodeId].Predecessor {
			res.Path = append([]int{nodeId}, res.Path...)
		}
	}
	return res
}

// Efficient implementation of Dijkstra's Algorithm for finding the shortest path from a source to every other node in the graph.
// Implementation is based on a priority queue.
func DijkstraOneToAll[N any, E g.IWeightedHalfEdge[W], W g.Weight](graph g.Graph[N, E], source g.NodeId) ShortestPathToAllResult[W] {
	dijkstraItems := make([]*DijkstraPqItem[W], graph.NodeCount(), graph.NodeCount())
	dijkstraItems[source] = &DijkstraPqItem[W]{Id: source, Priority: 0, Predecessor: -1}

	pq := make(DijkstraPriorityQueue[W], 0)
	heap.Init(&pq)
	heap.Push(&pq, dijkstraItems[source])

	pqPops := 0
	for len(pq) > 0 {
		currentPqItem := heap.Pop(&pq).(*DijkstraPqItem[W])
		currentNodeId := currentPqItem.Id
		pqPops++

		for _, edge := range graph.GetHalfEdgesFrom(currentNodeId) {
			successor := edge.To()

			if dijkstraItems[successor] == nil {
				newPriority := dijkstraItems[currentNodeId].Priority + edge.Weight()
				pqItem := DijkstraPqItem[W]{Id: successor, Priority: newPriority, Predecessor: currentNodeId}
				dijkstraItems[successor] = &pqItem
				heap.Push(&pq, &pqItem)
			} else {
				if updatedDistance := dijkstraItems[currentNodeId].Priority + edge.Weight(); updatedDistance < dijkstraItems[successor].Priority {
					dijkstraItems[successor].Priority = updatedDistance
					dijkstraItems[successor].Predecessor = currentNodeId
					heap.Fix(&pq, dijkstraItems[successor].index)
				}
			}
		}
	}

	// assemble result item
	result := ShortestPathToAllResult[W]{Lengths: make([]W, graph.NodeCount(), graph.NodeCount()), Predecessors: make([]g.NodeId, graph.NodeCount(), graph.NodeCount()), PqPops: pqPops}
	for nodeId, pqItem := range dijkstraItems {
		if pqItem != nil {
			result.Lengths[nodeId] = pqItem.Priority
			result.Predecessors[nodeId] = pqItem.Predecessor
		} else {
			result.Lengths[nodeId] = -1
			result.Predecessors[nodeId] = -1
		}
	}
	return result
}
