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

// Input: partitioned graph (max 64 partitions) with zero arc-flags
// Output: partitioned graph with non-trivial arc-flags

func b() {
	var faag g.AdjacencyArrayGraph[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]]
	ComputeArcFlags[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int](&faag, &faag, 64)
}

// Implementation with restricted types due to syntactic limitations of Golang
func ComputeArcFlags[N g.Partitioner, E g.IFlaggedHalfEdge[W], W g.Weight](forwardGraph, transposedGraph g.Graph[N, E], partitionCount int) *g.AdjacencyArrayGraph[N, E] {

	// create a copy of the (forward) graph
	faag := g.NewAdjacencyArrayFromGraph(forwardGraph)

	// remove any existing arc flags
	for i, halfEdge := range faag.Edges {
		faag.Edges[i] = halfEdge.ResetFlag().(E)
	}

	// determine boundary nodes for each region
	boundaryNodeSets := make([](map[g.NodeId]bool), 0) // TODO change bool to struct{}
	for i := 0; i < partitionCount; i++ {
		boundaryNodeSets = append(boundaryNodeSets, make(map[g.NodeId]bool, 0))
	}
	for tailNodeId := range faag.Nodes {
		for _, halfEdge := range faag.GetHalfEdgesFrom(tailNodeId) {
			headPartition := faag.GetNode(halfEdge.To()).Partition()
			if faag.GetNode(tailNodeId).Partition() != headPartition {
				boundaryNodeSets[headPartition][halfEdge.To()] = true
			}
		}
	}

	jobs := make(chan addFlagJob, 1<<16)
	done := make(chan bool)
	guard := make(chan struct{}, MAX_GOROUTINES)
	wg := sync.WaitGroup{}

	// start consumer
	go addFlag[N, E, W](jobs, faag, done)

	for partition, set := range boundaryNodeSets {
		setSize := len(set)
		wg.Add(setSize)
		fmt.Printf("Partition: %d, size=%d\n", partition, setSize)
		for boundaryNodeId := range set {
			guard <- struct{}{}
			go backwardSearch[N, E, W](jobs, forwardGraph, transposedGraph, g.PartitionId(partition), boundaryNodeId, &wg, guard)
		}
	}

	wg.Wait()
	close(jobs)
	<-done

	// revise edges within the same partition
	for i := 0; i < faag.NodeCount(); i++ {
		for _, halfEdge := range faag.GetHalfEdgesFrom(i) {
			if faag.GetNode(i).Partition() == faag.GetNode(halfEdge.To()).Partition() {
				addFlag1[N, E, W](faag, i, halfEdge.To(), faag.GetNode(i).Partition())
			}
		}
	}

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
			tailRev := node.Id
			headRev := child.Id
			jobs <- addFlagJob{from: headRev, to: tailRev, partition: partition}
			if forwardGraph.GetNode(child.Id).Partition() != partition && !child.Visited {
				stack = append(stack, child)
			}
			child.Visited = true
		}
	}
	<-guard
	wg.Done()
}

// (single) consumer
func addFlag[N g.Partitioner, E g.IFlaggedHalfEdge[W], W g.Weight](jobs <-chan addFlagJob, faag *g.AdjacencyArrayGraph[N, E], done chan<- bool) {
	for job := range jobs {
		for i := faag.Offsets[job.from]; i < faag.Offsets[job.from+1]; i++ {
			edge := faag.Edges[i]
			if edge.To() == job.to {
				faag.Edges[i] = edge.AddFlag(job.partition).(E)
				break
			}
		}
	}
	done <- true
}

func addFlag1[N g.Partitioner, E g.IFlaggedHalfEdge[W], W g.Weight](faag *g.AdjacencyArrayGraph[N, E], from g.NodeId, to g.NodeId, partition g.PartitionId) {
	for i := faag.Offsets[from]; i < faag.Offsets[from+1]; i++ {
		edge := faag.Edges[i]
		if edge.To() == to {
			faag.Edges[i] = edge.AddFlag(partition).(E)
			break
		}
	}
}
