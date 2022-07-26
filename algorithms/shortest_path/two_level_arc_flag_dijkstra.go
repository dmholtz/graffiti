package shortest_path

import (
	"container/heap"

	g "github.com/dmholtz/graffiti/graph"
)

// TwoLevelArcFlagRouter implements the Router interface and improves Dijkstra's algorithm by incorporating two-level arc flags.
type TwoLevelArcFlagRouter[N g.TwoLevelPartitioner, E g.ITwoLevelFlaggedHalfEdge[W], W g.Weight] struct {
	Graph g.Graph[N, E]
}

// String implements fmt.Stringer
func (r TwoLevelArcFlagRouter[N, E, W]) String() string {
	return "Two-level ArcFlag Dijkstra"
}

// Implementation of Dijkstra's Algorithm with two-level arc flags
func (r TwoLevelArcFlagRouter[N, E, W]) Route(source, target g.NodeId, recordSearchSpace bool) ShortestPathResult[W] {
	var searchSpace []g.NodeId = nil
	if recordSearchSpace {
		searchSpace = make([]g.NodeId, 0)
	}

	dijkstraItems := make([]*DijkstraPqItem[W], r.Graph.NodeCount(), r.Graph.NodeCount())
	dijkstraItems[source] = &DijkstraPqItem[W]{Id: source, Priority: 0, Predecessor: -1}

	pq := make(DijkstraPriorityQueue[W], 0)
	heap.Init(&pq)
	heap.Push(&pq, dijkstraItems[source])

	l1TargetPartition := r.Graph.GetNode(target).L1Part()
	l2TargetPartition := r.Graph.GetNode(target).L2Part()

	pqPops := 0
	for len(pq) > 0 {
		currentPqItem := heap.Pop(&pq).(*DijkstraPqItem[W])
		currentNodeId := currentPqItem.Id
		pqPops++

		if recordSearchSpace {
			searchSpace = append(searchSpace, currentNodeId)
		}

		currentL1Part := r.Graph.GetNode(currentNodeId).L1Part()
		for _, edge := range r.Graph.GetHalfEdgesFrom(currentNodeId) {
			// restrict the search space to the edges that are flagged with the l1-target-partition
			if !edge.IsL1Flagged(l1TargetPartition) {
				continue
			}

			successor := edge.To()

			// if the current edge is within the l1-target-partition, check whether the l2-partition flag is set
			if currentL1Part == l1TargetPartition && r.Graph.GetNode(successor).L1Part() == l1TargetPartition {
				if !edge.IsL2Flagged(l2TargetPartition) {
					continue
				}
			}

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
