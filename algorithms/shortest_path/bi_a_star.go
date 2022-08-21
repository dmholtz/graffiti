package shortest_path

import (
	"container/heap"

	g "github.com/dmholtz/graffiti/graph"
)

// BidirectionalAStarRouter implements the Router interface and provides bidirectional A* search, a lower-bounding variant of Dijkstra's algorithm.
type BidirectionalAStarRouter[N any, E g.IWeightedHalfEdge[W], W g.Weight] struct {
	Graph     g.Graph[N, E]
	Transpose g.Graph[N, E]

	ForwardHeuristic  Heuristic[W]
	BackwardHeuristic Heuristic[W]

	MaxInitializerValue W
}

// String implements fmt.Stringer
func (r BidirectionalAStarRouter[N, E, W]) String() string {
	return "Bidirectional A-Star"
}

// Bidirectional implementation of the lower-bounding A* algorithm following the symmetric approach by I. Pohl: "Bi-directional Search", 1971
// cf. Goldberg et al.: "Computing the Shortest Path: A* Search meets Graph Theory", 2004
func (r BidirectionalAStarRouter[N, E, W]) Route(source, target g.NodeId, recordSearchSpace bool) ShortestPathResult[W] {
	var searchSpace []g.NodeId = nil
	if recordSearchSpace {
		searchSpace = make([]g.NodeId, 0)
	}

	// handle trivial search with source and target being the same node
	if source == target {
		return ShortestPathResult[W]{Length: W(0), Path: []g.NodeId{source}, PqPops: 0, SearchSpace: searchSpace}
	}

	dijkstraItemsForward := make([]*AStarPqItem[W], r.Graph.NodeCount(), r.Graph.NodeCount())
	dijkstraItemsForward[source] = &AStarPqItem[W]{Id: source, Distance: 0, Priority: 0, Predecessor: -1}

	dijkstraItemsBackward := make([]*AStarPqItem[W], r.Graph.NodeCount(), r.Graph.NodeCount())
	dijkstraItemsBackward[target] = &AStarPqItem[W]{Id: target, Distance: 0, Priority: 0, Predecessor: -1}

	pqForward := make(AStarPriorityQueue[W], 0)
	heap.Init(&pqForward)
	heap.Push(&pqForward, dijkstraItemsForward[source])

	pqBackward := make(AStarPriorityQueue[W], 0)
	heap.Init(&pqBackward)
	heap.Push(&pqBackward, dijkstraItemsBackward[target])

	forwardSettled := make([]bool, r.Graph.NodeCount(), r.Graph.NodeCount())
	backwardSettled := make([]bool, r.Graph.NodeCount(), r.Graph.NodeCount())

	r.ForwardHeuristic.Init(source, target)
	r.BackwardHeuristic.Init(target, source)

	// Once the algorithm terminates, mu contains the shortest path distance between source and target.
	mu := r.MaxInitializerValue // initialize with the largest representable number of weight type W

	middleNodeId := -1

	pqPops := 0
	for len(pqForward) > 0 && len(pqBackward) > 0 {
		forwardPqItem := heap.Pop(&pqForward).(*AStarPqItem[W])
		forwardNodeId := forwardPqItem.Id
		forwardSettled[forwardNodeId] = true
		backwardPqItem := heap.Pop(&pqBackward).(*AStarPqItem[W])
		backwardNodeId := backwardPqItem.Id
		backwardSettled[backwardNodeId] = true
		pqPops += 2

		if recordSearchSpace {
			searchSpace = append(searchSpace, forwardNodeId)
			searchSpace = append(searchSpace, backwardNodeId)
		}

		// stopping criterion (Symmetric Approach, cf Pohl: Bi-directional Search, 1971)
		if dijkstraItemsForward[forwardNodeId].Priority >= mu {
			break
		}

		// forward search
		for _, edge := range r.Graph.GetHalfEdgesFrom(forwardNodeId) {
			successor := edge.To()
			// improvement by Kwa: An admissible bidirectional staged heuristic search algorithm
			if mu < r.MaxInitializerValue && dijkstraItemsBackward[successor] != nil && backwardSettled[successor] == true {
				// improvement by Kwa: An admissible bidirectional staged heuristic search algorithm
				if dijkstraItemsForward[successor] == nil {
					newDistance := forwardPqItem.Distance + edge.Weight()
					newPriority := newDistance + r.ForwardHeuristic.Evaluate(successor)
					pqItem := AStarPqItem[W]{Id: successor, Priority: newPriority, Distance: newDistance, Predecessor: forwardNodeId}
					dijkstraItemsForward[successor] = &pqItem
					// no put on priority queue
				}
				x := dijkstraItemsBackward[successor]
				if x != nil {
					if mu_new := dijkstraItemsForward[forwardNodeId].Distance + edge.Weight() + x.Distance; mu_new < mu {
						mu = mu_new
						dijkstraItemsForward[successor].Predecessor = forwardNodeId
						middleNodeId = successor
					}
				}
				continue
			}
			if dijkstraItemsForward[successor] == nil {
				newDistance := forwardPqItem.Distance + edge.Weight()
				newPriority := newDistance + r.ForwardHeuristic.Evaluate(successor)
				pqItem := AStarPqItem[W]{Id: successor, Priority: newPriority, Distance: newDistance, Predecessor: forwardNodeId}
				dijkstraItemsForward[successor] = &pqItem
				heap.Push(&pqForward, &pqItem)
			} else {
				if updatedPriority := forwardPqItem.Distance + edge.Weight() + r.ForwardHeuristic.Evaluate(successor); updatedPriority < dijkstraItemsForward[successor].Priority {
					dijkstraItemsForward[successor].Distance = forwardPqItem.Distance + edge.Weight()
					dijkstraItemsForward[successor].Priority = updatedPriority
					dijkstraItemsForward[successor].Predecessor = forwardNodeId
					heap.Fix(&pqForward, dijkstraItemsForward[successor].index)
				}
			}

			//// heuristic check
			//ld := Dijkstra[N, E, W](graph, successor, target, false).Length
			//h := r.ForwardHeuristic.Evaluate(successor)
			//if ld < h {
			//	fmt.Printf("Heuristic forward is not admissible: l < h: %d < %d\n", int(ld), int(h))
			//}

			x := dijkstraItemsBackward[successor]
			if x != nil {
				if mu_new := dijkstraItemsForward[forwardNodeId].Distance + edge.Weight() + x.Distance; mu_new < mu {
					mu = mu_new
					dijkstraItemsForward[successor].Predecessor = forwardNodeId
					middleNodeId = successor
				}
			}
		}

		// stopping criterion (Symmetric Approach, cf Pohl: Bi-directional Search, 1971)
		if dijkstraItemsBackward[backwardNodeId].Priority >= mu {
			break
		}

		// backward search
		for _, edge := range r.Transpose.GetHalfEdgesFrom(backwardNodeId) {
			successor := edge.To()
			// improvement by Kwa: An admissible bidirectional staged heuristic search algorithm
			if mu < r.MaxInitializerValue && dijkstraItemsForward[successor] != nil && forwardSettled[successor] == true {
				// improvement by Kwa: An admissible bidirectional staged heuristic search algorithm
				if dijkstraItemsBackward[successor] == nil {
					newDistance := backwardPqItem.Distance + edge.Weight()
					newPriority := newDistance + r.BackwardHeuristic.Evaluate(successor)
					pqItem := AStarPqItem[W]{Id: successor, Priority: newPriority, Distance: newDistance, Predecessor: backwardNodeId}
					dijkstraItemsBackward[successor] = &pqItem
					// no put on priority queue
				}
				x := dijkstraItemsForward[successor]
				if x != nil {
					if mu_new := dijkstraItemsBackward[backwardNodeId].Distance + edge.Weight() + x.Distance; mu_new < mu {
						mu = mu_new
						dijkstraItemsBackward[successor].Predecessor = backwardNodeId
						middleNodeId = successor
					}
				}
				continue
			}
			if dijkstraItemsBackward[successor] == nil {
				newDistance := backwardPqItem.Distance + edge.Weight()
				newPriority := newDistance + r.BackwardHeuristic.Evaluate(successor)
				pqItem := AStarPqItem[W]{Id: successor, Priority: newPriority, Distance: newDistance, Predecessor: backwardNodeId}
				dijkstraItemsBackward[successor] = &pqItem
				heap.Push(&pqBackward, &pqItem)
			} else {
				if updatedPriority := dijkstraItemsBackward[backwardNodeId].Distance + edge.Weight() + r.BackwardHeuristic.Evaluate(successor); updatedPriority < dijkstraItemsBackward[successor].Priority {
					dijkstraItemsBackward[successor].Distance = backwardPqItem.Distance + edge.Weight()
					dijkstraItemsBackward[successor].Priority = updatedPriority
					dijkstraItemsBackward[successor].Predecessor = backwardNodeId
					heap.Fix(&pqBackward, dijkstraItemsBackward[successor].index)
				}
			}
			x := dijkstraItemsForward[successor]
			if x != nil {
				if mu_new := dijkstraItemsBackward[backwardNodeId].Distance + edge.Weight() + x.Distance; mu_new < mu {
					mu = mu_new
					dijkstraItemsBackward[successor].Predecessor = backwardNodeId
					middleNodeId = successor
				}
			}
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
