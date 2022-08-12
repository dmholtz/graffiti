package shortest_path

import (
	"container/heap"
	"fmt"
	"sync"

	g "github.com/dmholtz/graffiti/graph"
)

// Type of boundary node
const (
	// A l1-boundary node is a node at the boundary of a level 1 region.
	l1_BOUNDARY_NODE = iota
	// A l2-boundary node is a node at the boundary of a level 2 region. Every l1-boundary node is also a l2-boundary node.
	l2_BOUNDARY_NODE = iota
)

const (
	l1_JOB = iota
	l2_JOB = iota
)

type addTwoLevelFlagJob struct {
	from      g.NodeId
	to        g.NodeId
	partition g.PartitionId
	kind      uint8
}

// Implementation with restricted types due to syntactic limitations of Golang
func ComputeTwoLevelArcFlags[N g.TwoLevelPartitioner, E g.ITwoLevelFlaggedHalfEdge[W], W g.Weight](forwardGraph, transposedGraph g.Graph[N, E]) *g.AdjacencyArrayGraph[N, E] {

	l1PartCount := g.PartitionId(32)
	l2PartCount := g.PartitionId(32)

	// create a copy of the (forward) graph
	faag := g.NewAdjacencyArrayFromGraph(forwardGraph)

	// remove any existing arc flags
	for i, halfEdge := range faag.Edges {
		faag.Edges[i] = halfEdge.ResetFlags().(E)
	}

	// Initialize l1/l2 boundary node maps
	l1BoundaryNodes := make(map[g.PartitionId]map[g.NodeId]struct{})
	l2BoundaryNodes := make(map[g.PartitionId]map[g.PartitionId]map[g.NodeId]struct{})
	for l1Part := g.PartitionId(0); l1Part < l1PartCount; l1Part++ {
		l1BoundaryNodes[l1Part] = make(map[g.NodeId]struct{})
		l2BoundaryNodes[l1Part] = make(map[g.PartitionId]map[g.NodeId]struct{})
		for l2Part := g.PartitionId(0); l2Part < l2PartCount; l2Part++ {
			l2BoundaryNodes[l1Part][l2Part] = make(map[g.NodeId]struct{})
		}
	}

	// determine l1 boundary nodes
	for tailNodeId, tailNode := range faag.Nodes {
		tailL1Part := tailNode.L1Part()
		tailL2Part := tailNode.L2Part()
		for _, edge := range faag.GetHalfEdgesFrom(tailNodeId) {
			headL1Part := faag.GetNode(edge.To()).L1Part()
			headL2Part := faag.GetNode(edge.To()).L2Part()
			if tailL1Part != headL1Part {
				l1BoundaryNodes[headL1Part][edge.To()] = struct{}{}
				l2BoundaryNodes[headL1Part][headL2Part][edge.To()] = struct{}{}
			} else {
				// same l1 partition
				if tailL2Part != headL2Part {
					l2BoundaryNodes[headL1Part][headL2Part][edge.To()] = struct{}{}
				}
			}
		}
	}

	l1Jobs := make(chan addFlagJob, 1<<16)
	done := make(chan bool)
	guard := make(chan struct{}, MAX_GOROUTINES)
	wg := sync.WaitGroup{}

	// start consumer
	go addLevel1Flag[N, E, W](l1Jobs, faag, done)

	for l1Part, nodeSet := range l1BoundaryNodes {
		setSize := len(nodeSet)
		wg.Add(setSize)
		fmt.Printf("Partition: %d, size=%d\n", l1Part, setSize)
		for boundaryNodeId := range nodeSet {
			guard <- struct{}{}
			go l1BoundaryBackwardSearch[N, E, W](l1Jobs, forwardGraph, transposedGraph, boundaryNodeId, &wg, guard)
		}
	}

	//revise edges within the same l1 partition
	for i := 0; i < faag.NodeCount(); i++ {
		tailL1Part := faag.GetNode(i).L1Part()
		for _, halfEdge := range faag.GetHalfEdgesFrom(i) {
			if tailL1Part == faag.GetNode(halfEdge.To()).L1Part() {
				l1Jobs <- addFlagJob{from: i, to: halfEdge.To(), partition: tailL1Part}
			}
		}
	}

	wg.Wait()
	close(l1Jobs)
	<-done
	fmt.Println("Done with L1 boundary nodes.")
	// l1 partitions are now done

	jobs := make(chan addTwoLevelFlagJob, 1<<10)
	go addTwoLevelFlag[N, E, W](jobs, faag, done)

	// precompute the l1 partition size for each l1 partition
	l1PartSizes := make(map[g.PartitionId]int)
	for _, node := range faag.Nodes {
		l1Part := node.L1Part()
		if _, ok := l1PartSizes[l1Part]; !ok {
			l1PartSizes[l1Part] = 1
		} else {
			l1PartSizes[l1Part] += 1
		}
	}

	// l2 boundary nodes
	for l1Part, l2Map := range l2BoundaryNodes {
		// determine the l1 partition size, which serves as a pruning / early stopping criterion for the backward search
		l2BoundaryNodeSize := 0
		for _, nodeSet := range l2Map {
			l2BoundaryNodeSize += len(nodeSet)
		}
		fmt.Printf("Partition: %d, size=%d\n", l1Part, l2BoundaryNodeSize)
		wg.Add(l2BoundaryNodeSize)

		for _, nodeSet := range l2Map {
			for boundaryNodeId := range nodeSet {
				guard <- struct{}{}
				go l2BoundaryBackwardSearch[N, E, W](jobs, forwardGraph, transposedGraph, boundaryNodeId, l1PartSizes[l1Part], &wg, guard)
			}

		}
	}

	// revise edges within the same l2 partition
	for i := 0; i < faag.NodeCount(); i++ {
		tailL1Part := faag.GetNode(i).L1Part()
		tailL2Part := faag.GetNode(i).L2Part()
		for _, halfEdge := range faag.GetHalfEdgesFrom(i) {
			if tailL1Part == (faag.GetNode(halfEdge.To()).L1Part()) {
				if tailL2Part == (faag.GetNode(halfEdge.To()).L2Part()) {
					jobs <- addTwoLevelFlagJob{from: i, to: halfEdge.To(), partition: tailL2Part, kind: l2_JOB}
				}
			}
		}
	}

	wg.Wait()
	close(jobs)
	<-done

	return faag
}

// (single) consumer
func addTwoLevelFlag[N g.TwoLevelPartitioner, E g.ITwoLevelFlaggedHalfEdge[W], W g.Weight](jobs <-chan addTwoLevelFlagJob, faag *g.AdjacencyArrayGraph[N, E], done chan<- bool) {
	for job := range jobs {
		for i := faag.Offsets[job.from]; i < faag.Offsets[job.from+1]; i++ {
			edge := faag.Edges[i]
			if edge.To() == job.to {
				if job.kind == l1_JOB {
					faag.Edges[i] = edge.AddL1Flag(job.partition).(E)
					break
				} else if job.kind == l2_JOB {
					faag.Edges[i] = edge.AddL2Flag(job.partition).(E)
					break
				} else {
					panic("Invalid job kind")
				}
			}
		}
	}
	done <- true
}

// (single) consumer
func addLevel1Flag[N g.TwoLevelPartitioner, E g.ITwoLevelFlaggedHalfEdge[W], W g.Weight](jobs <-chan addFlagJob, faag *g.AdjacencyArrayGraph[N, E], done chan<- bool) {
	for job := range jobs {
		for i := faag.Offsets[job.from]; i < faag.Offsets[job.from+1]; i++ {
			edge := faag.Edges[i]
			if edge.To() == job.to {
				faag.Edges[i] = edge.AddL1Flag(job.partition).(E)
				break
			}
		}
	}
	done <- true
}

// producer function
// TODO change edge type to IWeightedHalfEdge??
func l1BoundaryBackwardSearch[N g.TwoLevelPartitioner, E g.ITwoLevelFlaggedHalfEdge[W], W g.Weight](jobs chan<- addFlagJob, forwardGraph, transposedGraph g.Graph[N, E], boundaryNodeId g.NodeId, wg *sync.WaitGroup, guard <-chan struct{}) {
	// calculate in reverse graph
	tree := ShortestPathTree[N, E, W](transposedGraph, boundaryNodeId)
	l1Part := forwardGraph.GetNode(boundaryNodeId).L1Part()

	stack := make([]*ShortestPathTreeNode, 0)
	stack = append(stack, &tree)

	for len(stack) > 0 {
		// pop
		treeNode := stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]

		for _, child := range treeNode.Children {
			if child.Visited {
				continue
			}
			child.Visited = true

			tailRev := treeNode.Id
			headRev := child.Id
			jobs <- addFlagJob{from: headRev, to: tailRev, partition: l1Part}

			if forwardGraph.GetNode(child.Id).L1Part() != l1Part {
				stack = append(stack, child)
			}
		}
	}

	<-guard
	wg.Done()
}

// producer function
// TODO change edge type to IWeightedHalfEdge??
func l2BoundaryBackwardSearch[N g.TwoLevelPartitioner, E g.ITwoLevelFlaggedHalfEdge[W], W g.Weight](jobs chan<- addTwoLevelFlagJob, forwardGraph, transposedGraph g.Graph[N, E], boundaryNodeId g.NodeId, l1PartSize int, wg *sync.WaitGroup, guard <-chan struct{}) {
	l1Part := forwardGraph.GetNode(boundaryNodeId).L1Part()
	l2Part := forwardGraph.GetNode(boundaryNodeId).L2Part()

	// calculate in reverse graph
	tree := prundedShortestPathTree[N, E, W](transposedGraph, boundaryNodeId, l1Part, l1PartSize)

	stack := make([]*ShortestPathTreeNode, 0)
	stack = append(stack, &tree)

	for len(stack) > 0 {
		// pop
		treeNode := stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]

		for _, child := range treeNode.Children {
			if !child.Visited {
				stack = append(stack, child)
			}
			child.Visited = true

			tailRev := treeNode.Id
			headRev := child.Id
			if forwardGraph.GetNode(headRev).L1Part() == l1Part && forwardGraph.GetNode(tailRev).L1Part() == l1Part {
				jobs <- addTwoLevelFlagJob{from: headRev, to: tailRev, partition: l2Part, kind: l2_JOB}
			}
		}
	}

	<-guard
	wg.Done()
}

// Pruned variant of ShortestPathTree method: Search stops once every node of the specified l1-partition has been settled.
func prundedShortestPathTree[N g.TwoLevelPartitioner, E g.IWeightedHalfEdge[W], W g.Weight](graph g.Graph[N, E], source g.NodeId, l1Part g.PartitionId, l1PartSize int) ShortestPathTreeNode {
	dijkstraItems := make([]*ShortestPathTreePqItem[W], graph.NodeCount(), graph.NodeCount())
	dijkstraItems[source] = &ShortestPathTreePqItem[W]{Id: source, Priority: 0, Predecessors: make([]int, 0)}

	pq := make(ShortestPathTreePriorityQueue[W], 0)
	heap.Init(&pq)
	heap.Push(&pq, dijkstraItems[source])

	successors := make([]*ShortestPathTreeNode, graph.NodeCount(), graph.NodeCount())

	l1SettledCount := 0
	for len(pq) > 0 {
		currentPqItem := heap.Pop(&pq).(*ShortestPathTreePqItem[W])
		currentNodeId := currentPqItem.Id

		if graph.GetNode(currentNodeId).L1Part() == l1Part {
			l1SettledCount += 1
		}

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

		// Pruning / stopping criterion
		if l1SettledCount >= l1PartSize {
			break
		}
	}

	return *successors[source]
}
