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
