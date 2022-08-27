package shortest_path

import (
	"container/heap"

	g "github.com/dmholtz/graffiti/graph"
)

// ArcFlagAStarRouter implements combines A* search with bidirectional arc flags.
type ArcFlagAStarRouter[N g.Partitioner, E g.IFlaggedHalfEdge[W], W g.Weight] struct {
	Graph     g.Graph[N, E]
	Transpose g.Graph[N, E]

	Heuristic Heuristic[W]
}

// String implements fmt.Stringer
func (r ArcFlagAStarRouter[N, E, W]) String() string {
	return "A-Star with bidirectional arc flags"
}

// A* with feasible heruistic is a lower-bounding algorithm
func (r ArcFlagAStarRouter[N, E, W]) Route(source, target g.NodeId, recordSearchSpace bool) ShortestPathResult[W] {
	var searchSpace []g.NodeId = nil
	if recordSearchSpace {
		searchSpace = make([]g.NodeId, 0)
	}

	dijkstraItems := make([]*AStarPqItem[W], r.Graph.NodeCount(), r.Graph.NodeCount())
	dijkstraItems[source] = &AStarPqItem[W]{Id: source, Distance: 0, Priority: 0, Predecessor: -1}

	pq := make(AStarPriorityQueue[W], 0)
	heap.Init(&pq)
	heap.Push(&pq, dijkstraItems[source])

	sourcePart := r.Transpose.GetNode(source).Partition()
	targetPart := r.Graph.GetNode(target).Partition()

	r.Heuristic.Init(source, target)

	pqPops := 0
	for len(pq) > 0 {
		currentPqItem := heap.Pop(&pq).(*AStarPqItem[W])
		currentNodeId := currentPqItem.Id
		pqPops++

		if recordSearchSpace {
			searchSpace = append(searchSpace, currentNodeId)
		}

		for _, edge := range r.Graph.GetHalfEdgesFrom(currentNodeId) {
			successor := edge.To()

			if !edge.IsFlagged(targetPart) {
				continue
			}
			// find the reverse edge
			var revEdge E
			for _, e := range r.Transpose.GetHalfEdgesFrom(successor) {
				if e.To() == currentNodeId {
					revEdge = e
					break
				}
			}
			if !revEdge.IsFlagged(sourcePart) {
				continue
			}

			if dijkstraItems[successor] == nil {
				newDistance := currentPqItem.Distance + edge.Weight()
				newPriority := newDistance + r.Heuristic.Evaluate(successor)
				pqItem := AStarPqItem[W]{Id: successor, Priority: newPriority, Distance: newDistance, Predecessor: currentNodeId}
				dijkstraItems[successor] = &pqItem
				heap.Push(&pq, &pqItem)
			} else {
				if updatedPriority := currentPqItem.Distance + edge.Weight() + r.Heuristic.Evaluate(successor); updatedPriority < dijkstraItems[successor].Priority {
					dijkstraItems[successor].Distance = currentPqItem.Distance + edge.Weight()
					dijkstraItems[successor].Priority = updatedPriority
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
