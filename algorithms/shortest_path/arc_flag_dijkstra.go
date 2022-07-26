package shortest_path

import (
	"container/heap"

	g "github.com/dmholtz/graffiti/graph"
)

// ArcFlagRouter implements the Router interface and improves Dijkstra's algorithm by incorporating arc flags.
type ArcFlagRouter[N g.Partitioner, E g.IFlaggedHalfEdge[W], W g.Weight] struct {
	Graph g.Graph[N, E]
}

// String implements fmt.Stringer
func (r ArcFlagRouter[N, E, W]) String() string {
	return "ArcFlag Dijkstra"
}

// Implementation of Dijkstra's Algorithm with arc flags
func (r ArcFlagRouter[N, E, W]) Route(source, target g.NodeId, recordSearchSpace bool) ShortestPathResult[W] {
	var searchSpace []g.NodeId = nil
	if recordSearchSpace {
		searchSpace = make([]g.NodeId, 0)
	}

	dijkstraItems := make([]*DijkstraPqItem[W], r.Graph.NodeCount(), r.Graph.NodeCount())
	dijkstraItems[source] = &DijkstraPqItem[W]{Id: source, Priority: 0, Predecessor: -1}

	pq := make(DijkstraPriorityQueue[W], 0)
	heap.Init(&pq)
	heap.Push(&pq, dijkstraItems[source])

	targetPartition := r.Graph.GetNode(target).Partition()

	pqPops := 0
	for len(pq) > 0 {
		currentPqItem := heap.Pop(&pq).(*DijkstraPqItem[W])
		currentNodeId := currentPqItem.Id
		pqPops++

		if recordSearchSpace {
			searchSpace = append(searchSpace, currentNodeId)
		}

		for _, edge := range r.Graph.GetHalfEdgesFrom(currentNodeId) {
			if !edge.IsFlagged(targetPartition) {
				continue
			}

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
