package shortest_path

import (
	"container/heap"

	g "github.com/dmholtz/graffiti/graph"
)

// Heuristic interface for A* Search
type Heuristic[W g.Weight] interface {
	// Init initializes the heuristic with the source and the target node id of the next search.
	// The Init method must be called before any search by the search algorithm.
	Init(source g.NodeId, target g.NodeId)
	// Evaluate computes the value of the heuristic at node with ID=id.
	Evaluate(id g.NodeId) W
}

// AStarRouter implements the Router interface and provides A* search, a lower-bounding variant of Dijkstra's algorithm.
type AStarRouter[N any, E g.IWeightedHalfEdge[W], W g.Weight] struct {
	Graph     g.Graph[N, E]
	Heuristic Heuristic[W]
}

// String implements fmt.Stringer
func (r AStarRouter[N, E, W]) String() string {
	return "A-Star"
}

// A* with feasible heruistic is a lower-bounding algorithm
func (r AStarRouter[N, E, W]) Route(source, target g.NodeId, recordSearchSpace bool) ShortestPathResult[W] {
	var searchSpace []g.NodeId = nil
	if recordSearchSpace {
		searchSpace = make([]g.NodeId, 0)
	}

	dijkstraItems := make([]*AStarPqItem[W], r.Graph.NodeCount(), r.Graph.NodeCount())
	dijkstraItems[source] = &AStarPqItem[W]{Id: source, Distance: 0, Priority: 0, Predecessor: -1}

	pq := make(AStarPriorityQueue[W], 0)
	heap.Init(&pq)
	heap.Push(&pq, dijkstraItems[source])

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
