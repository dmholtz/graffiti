package shortest_path

import (
	"container/heap"

	g "github.com/dmholtz/graffiti/graph"
)

// BiDijkstraRouter implements the Router interface and features bidirectional search with Dijkstra's algorithm.
//
// Caveat: Always set the MaxInitializerValue to the maximum value of the generic type W, e.g. math.MaxInt in case of int.
type BiDijkstraRouter[N any, E g.IWeightedHalfEdge[W], W g.Weight] struct {
	Graph     g.Graph[N, E]
	Transpose g.Graph[N, E]

	MaxInitializerValue W
}

// String implements fmt.Stringer
func (r BiDijkstraRouter[N, E, W]) String() string {
	return "Bidirectional Dijkstra"
}

// Bidirectional Dijkstra runs a forward search from the source node to the target node in parallel with
// a backward search in the backward graph (transpose) from the target node to the source node.
//
// Bidirectional search requires both the forward graph and its transpose (backward graph) as input parameters.
// In case of undirected graphs, the same argument may be passed for both parameters.
//
// Reference: https://www.homepages.ucl.ac.uk/~ucahmto/math/2020/05/30/bidirectional-dijkstra.html
func (r BiDijkstraRouter[N, E, W]) Route(source, target g.NodeId, recordSearchSpace bool) ShortestPathResult[W] {
	var searchSpace []g.NodeId = nil
	if recordSearchSpace {
		searchSpace = make([]g.NodeId, 0)
	}

	// handle trivial search with source and target being the same node
	if source == target {
		return ShortestPathResult[W]{Length: W(0), Path: []g.NodeId{source}, PqPops: 0, SearchSpace: searchSpace}
	}

	dijkstraItemsForward := make([]*DijkstraPqItem[W], r.Graph.NodeCount(), r.Graph.NodeCount())
	dijkstraItemsForward[source] = &DijkstraPqItem[W]{Id: source, Priority: 0, Predecessor: -1}

	dijkstraItemsBackward := make([]*DijkstraPqItem[W], r.Transpose.NodeCount(), r.Transpose.NodeCount())
	dijkstraItemsBackward[target] = &DijkstraPqItem[W]{Id: target, Priority: 0, Predecessor: -1}

	pqForward := make(DijkstraPriorityQueue[W], 0)
	heap.Init(&pqForward)
	heap.Push(&pqForward, dijkstraItemsForward[source])

	pqBackward := make(DijkstraPriorityQueue[W], 0)
	heap.Init(&pqBackward)
	heap.Push(&pqBackward, dijkstraItemsBackward[target])

	// Once the algorithm terminates, mu contains the shortest path distance between source and target.
	mu := r.MaxInitializerValue // initialize with the largest representable number of weight type W

	middleNodeId := -1

	pqPops := 0
	for len(pqForward) > 0 && len(pqBackward) > 0 {
		forwardPqItem := heap.Pop(&pqForward).(*DijkstraPqItem[W])
		forwardNodeId := forwardPqItem.Id
		backwardPqItem := heap.Pop(&pqBackward).(*DijkstraPqItem[W])
		backwardNodeId := backwardPqItem.Id
		pqPops += 2

		if recordSearchSpace {
			searchSpace = append(searchSpace, forwardNodeId)
			searchSpace = append(searchSpace, backwardNodeId)
		}

		// forward search
		for _, edge := range r.Graph.GetHalfEdgesFrom(forwardNodeId) {
			successor := edge.To()

			if dijkstraItemsForward[successor] == nil {
				newPriority := dijkstraItemsForward[forwardNodeId].Priority + edge.Weight()
				pqItem := DijkstraPqItem[W]{Id: successor, Priority: newPriority, Predecessor: forwardNodeId}
				dijkstraItemsForward[successor] = &pqItem
				heap.Push(&pqForward, &pqItem)
			} else {
				if updatedDistance := dijkstraItemsForward[forwardNodeId].Priority + edge.Weight(); updatedDistance < dijkstraItemsForward[successor].Priority {
					dijkstraItemsForward[successor].Priority = updatedDistance
					dijkstraItemsForward[successor].Predecessor = forwardNodeId
					heap.Fix(&pqForward, dijkstraItemsForward[successor].index)
				}
			}

			if x := dijkstraItemsBackward[successor]; x != nil && dijkstraItemsForward[forwardNodeId].Priority+edge.Weight()+x.Priority < mu {
				mu = dijkstraItemsForward[forwardNodeId].Priority + edge.Weight() + x.Priority
				dijkstraItemsForward[successor].Predecessor = forwardNodeId
				middleNodeId = successor
			}
		}

		// backward search
		for _, edge := range r.Graph.GetHalfEdgesFrom(backwardNodeId) {
			successor := edge.To()

			if dijkstraItemsBackward[successor] == nil {
				newPriority := dijkstraItemsBackward[backwardNodeId].Priority + edge.Weight()
				pqItem := DijkstraPqItem[W]{Id: successor, Priority: newPriority, Predecessor: backwardNodeId}
				dijkstraItemsBackward[successor] = &pqItem
				heap.Push(&pqBackward, &pqItem)
			} else {
				if updatedDistance := dijkstraItemsBackward[backwardNodeId].Priority + edge.Weight(); updatedDistance < dijkstraItemsBackward[successor].Priority {
					dijkstraItemsBackward[successor].Priority = updatedDistance
					dijkstraItemsBackward[successor].Predecessor = backwardNodeId
					heap.Fix(&pqBackward, dijkstraItemsBackward[successor].index)
				}
			}

			if x := dijkstraItemsForward[successor]; x != nil && dijkstraItemsBackward[backwardNodeId].Priority+edge.Weight()+x.Priority < mu {
				mu = dijkstraItemsBackward[backwardNodeId].Priority + edge.Weight() + x.Priority
				dijkstraItemsBackward[successor].Predecessor = backwardNodeId
				middleNodeId = successor
			}
		}

		// stopping criterion
		if dijkstraItemsForward[forwardNodeId].Priority+dijkstraItemsBackward[backwardNodeId].Priority >= mu {
			break
		}
	}

	res := ShortestPathResult[W]{Length: W(-1), Path: make([]g.NodeId, 0), PqPops: pqPops, SearchSpace: searchSpace}

	// check if path exists
	if mu < r.MaxInitializerValue {
		res.Length = mu
		// sanity check: length == dijkstraItemsForward[middleNodeId].priority + dijkstraItemsBackward[middleNodeId].priority
		if dijkstraItemsForward[middleNodeId] != nil && dijkstraItemsBackward[middleNodeId] != nil {
			for nodeId := middleNodeId; nodeId != -1; nodeId = dijkstraItemsForward[nodeId].Predecessor {
				res.Path = append([]int{nodeId}, res.Path...)
			}
			if res.Path[len(res.Path)-1] == middleNodeId {
				res.Path = res.Path[0 : len(res.Path)-1]
			}
			for nodeId := middleNodeId; nodeId != -1; nodeId = dijkstraItemsBackward[nodeId].Predecessor {
				res.Path = append(res.Path, nodeId)
			}
		}
	}
	return res
}
