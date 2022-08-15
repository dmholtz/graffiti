package shortest_path

import (
	g "github.com/dmholtz/graffiti/graph"
)

// Atomic element of the priority queue used in Dijkstra's algorithm
type DijkstraPqItem[W g.Weight] struct {
	// ID of the node this item refers to
	Id g.NodeId
	// priority of this item in the heap
	Priority W
	// predecessor node of this item in the search tree
	Predecessor g.NodeId
	// index of this item in the underlying slice
	// The index is required for implementing heap.Interface and managed by the interface's methods.
	index int
}

// A priority queue implementation for Dijkstra's algorithm and other shortest path algorithms.
// Implements heap.Interface (https://pkg.go.dev/container/heap)
type DijkstraPriorityQueue[W g.Weight] []*DijkstraPqItem[W]

// Len implements heap.Interface
func (pq DijkstraPriorityQueue[W]) Len() int {
	return len(pq)
}

// Less implements heap.Interface
func (pq DijkstraPriorityQueue[W]) Less(i, j int) bool {
	// Min-Heap
	return pq[i].Priority < pq[j].Priority
}

// Swap implements heap.Interface
func (pq DijkstraPriorityQueue[W]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

// Push implements heap.Interface
func (pq *DijkstraPriorityQueue[W]) Push(item interface{}) {
	n := len(*pq)
	pqItem := item.(*DijkstraPqItem[W])
	pqItem.index = n
	*pq = append(*pq, pqItem)
}

// Pop implements heap.Interface
func (pq *DijkstraPriorityQueue[W]) Pop() interface{} {
	old_pq := *pq
	n := len(old_pq)
	pqItem := old_pq[n-1]
	old_pq[n-1] = nil
	*pq = old_pq[0 : n-1]
	return pqItem
}

// Atomic element of the priority queue used in A* Search
type AStarPqItem[W g.Weight] struct {
	// ID of the node this item refers to
	Id g.NodeId
	// minimum distance from source node to this node (g-value)
	Distance W
	// priority of this item in the heap (f-value)
	Priority W
	// predecessor node of this item in the search tree
	Predecessor g.NodeId
	// index of this item in the underlying slice
	// The index is required for implementing heap.Interface and managed by the interface's methods.
	index int
}

// A priority queue implementation for A* Search and other algorithms incorporating a heuristic.
// The priority is sorted according to the f-value: f = g + h, where f is the priority, g the minimum distance and h a heurisitc.
// An admissible heuristic estimates the distance from a node to the target without overestimating the true distance.
// Implements heap.Interface (https://pkg.go.dev/container/heap)
type AStarPriorityQueue[W g.Weight] []*AStarPqItem[W]

// Len implements heap.Interface
func (pq AStarPriorityQueue[W]) Len() int {
	return len(pq)
}

// Less implements heap.Interface
func (pq AStarPriorityQueue[W]) Less(i, j int) bool {
	// Min-Heap
	return pq[i].Priority < pq[j].Priority
}

// Swap implements heap.Interface
func (pq AStarPriorityQueue[W]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

// Push implements heap.Interface
func (pq *AStarPriorityQueue[W]) Push(item interface{}) {
	n := len(*pq)
	pqItem := item.(*AStarPqItem[W])
	pqItem.index = n
	*pq = append(*pq, pqItem)
}

// Pop implements heap.Interface
func (pq *AStarPriorityQueue[W]) Pop() interface{} {
	old_pq := *pq
	n := len(old_pq)
	pqItem := old_pq[n-1]
	old_pq[n-1] = nil
	*pq = old_pq[0 : n-1]
	return pqItem
}

// Atomic element of the priority queue used in the shortest path tree algorithm.
// Unlike DijkstraPqItem, this PqItem allows to store multiple predecessors of a node in case multiple shortest paths exist.
type ShortestPathTreePqItem[W g.Weight] struct {
	// ID of the node this item refers to
	Id g.NodeId
	// priority of this item in the heap
	Priority W
	// list of predecessor nodes of this item in the search tree
	Predecessors []g.NodeId
	// index of this item in the underlying slice
	// The index is required for implementing heap.Interface and managed by the interface's methods.
	index int
}

// A priority queue implementation for the shortest path tree algorithm.
// Implements heap.Interface (https://pkg.go.dev/container/heap)
type ShortestPathTreePriorityQueue[W g.Weight] []*ShortestPathTreePqItem[W]

// Len implements heap.Interface
func (pq ShortestPathTreePriorityQueue[W]) Len() int {
	return len(pq)
}

// Less implements heap.Interface
func (pq ShortestPathTreePriorityQueue[W]) Less(i, j int) bool {
	// Min-Heap
	return pq[i].Priority < pq[j].Priority
}

// Swap implements heap.Interface
func (pq ShortestPathTreePriorityQueue[W]) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].index, pq[j].index = i, j
}

// Push implements heap.Interface
func (pq *ShortestPathTreePriorityQueue[W]) Push(item interface{}) {
	n := len(*pq)
	pqItem := item.(*ShortestPathTreePqItem[W])
	pqItem.index = n
	*pq = append(*pq, pqItem)
}

// Pop implements heap.Interface
func (pq *ShortestPathTreePriorityQueue[W]) Pop() interface{} {
	old := *pq
	n := len(old)
	pqItem := old[n-1]
	old[n-1] = nil
	*pq = old[0 : n-1]
	return pqItem
}
