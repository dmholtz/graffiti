package shortest_path

import (
	"fmt"
	"sync"

	g "github.com/dmholtz/graffiti/graph"
)

type addFlagJob struct {
	from      g.NodeId
	to        g.NodeId
	partition g.PartitionId
}

const MAX_GOROUTINES = 8

// Parallel implementation of arc flag preprocessing
func ComputeArcFlags[N g.Partitioner, E g.IFlaggedHalfEdge[W], W g.Weight](forwardGraph, transposedGraph g.Graph[N, E], partitionCount int) *g.AdjacencyArrayGraph[N, E] {

	// create a copy of the (forward) graph
	faag := g.NewAdjacencyArrayFromGraph(forwardGraph)

	// remove any existing arc flags
	for i, halfEdge := range faag.Edges {
		faag.Edges[i] = halfEdge.ResetFlag().(E)
	}

	// determine the flag range based on the first edge in the graph
	if faag.EdgeCount() < 1 {
		panic(fmt.Sprintf("Cannot compute arc flags - the graph does not contain any edges."))
	}
	flagRange := faag.Edges[0].FlagRange()

	// check if every partition is within the respective flag range
	for _, node := range faag.Nodes {
		if node.Partition() >= flagRange {
			panic(fmt.Sprintf("Partition exceeds flag range: %d >= %d", node.Partition(), flagRange))
		}
	}

	// determine boundary nodes for each region
	boundaryNodeSets := make([](map[g.NodeId]struct{}), 0)
	for i := 0; i < partitionCount; i++ {
		boundaryNodeSets = append(boundaryNodeSets, make(map[g.NodeId]struct{}, 0))
	}
	for tailNodeId := range faag.Nodes {
		for _, halfEdge := range faag.GetHalfEdgesFrom(tailNodeId) {
			headPartition := faag.GetNode(halfEdge.To()).Partition()
			if faag.GetNode(tailNodeId).Partition() != headPartition {
				boundaryNodeSets[headPartition][halfEdge.To()] = struct{}{}
			}
		}
	}

	// parallel implementation of arcflag preprocessing following the (multiple) producer - (single) consumer pattern.
	// n producers: each consumer starts a backward search in the transposed graph from a boundary node and computes the arcflags (heavy workload)
	// 1 consumer: the consumer synchronizes the results to avoid race conditions and writes the computed arcflags into memory (low workload)

	jobs := make(chan addFlagJob, 1<<16) // buffered synchronization channel between producers and consumers
	done := make(chan bool)              // indicates that the single consumer is done

	guard := make(chan struct{}, MAX_GOROUTINES) // restrict the number of parallel processes to the number of CPU cores
	wg := sync.WaitGroup{}                       // synchronizes the producers

	// start consumer
	go addFlag[N, E, W](jobs, faag, done)

	// start producers in parallel
	for partition, set := range boundaryNodeSets {
		setSize := len(set)
		wg.Add(setSize)
		fmt.Printf("Partition: %d, size=%d\n", partition, setSize)
		for boundaryNodeId := range set {
			guard <- struct{}{} // reserve 1 producer
			go backwardSearch[N, E, W](jobs, forwardGraph, transposedGraph, g.PartitionId(partition), boundaryNodeId, &wg, guard)
		}
	}

	// revise edges within the same partition
	for i := 0; i < faag.NodeCount(); i++ {
		for _, halfEdge := range faag.GetHalfEdgesFrom(i) {
			if faag.GetNode(i).Partition() == faag.GetNode(halfEdge.To()).Partition() {
				jobs <- addFlagJob{from: i, to: halfEdge.To(), partition: faag.GetNode(i).Partition()}
			}
		}
	}

	wg.Wait()   // await the producers
	close(jobs) // close the job channel
	<-done      // wait for the single consumer to terminate

	return faag
}

// producer function
func backwardSearch[N g.Partitioner, E g.IFlaggedHalfEdge[W], W g.Weight](jobs chan<- addFlagJob, forwardGraph, transposedGraph g.Graph[N, E], partition g.PartitionId, boundaryNodeId g.NodeId, wg *sync.WaitGroup, guard <-chan struct{}) {
	// calculate in reverse graph
	tree := ShortestPathTree[N, E, W](transposedGraph, boundaryNodeId)

	stack := make([]*ShortestPathTreeNode, 0)
	stack = append(stack, &tree)

	for len(stack) > 0 {
		// pop
		node := stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]

		for _, child := range node.Children {
			tailRev := node.Id  // tail node in the reversed graph
			headRev := child.Id // head node in the reversed graph
			jobs <- addFlagJob{from: headRev, to: tailRev, partition: partition}
			if forwardGraph.GetNode(child.Id).Partition() != partition && !child.Visited {
				stack = append(stack, child)
			}
			child.Visited = true
		}
	}
	<-guard // free resources for next producer
	wg.Done()
}

// (single) consumer
func addFlag[N g.Partitioner, E g.IFlaggedHalfEdge[W], W g.Weight](jobs <-chan addFlagJob, faag *g.AdjacencyArrayGraph[N, E], done chan<- bool) {
	// loop over jobs channel unti it is closed
	for job := range jobs {
		for i := faag.Offsets[job.from]; i < faag.Offsets[job.from+1]; i++ {
			edge := faag.Edges[i]
			if edge.To() == job.to {
				faag.Edges[i] = edge.AddFlag(job.partition).(E)
				break
			}
		}
	}
	done <- true // announce that consumer has terminated
}
