package partitioning

import (
	"fmt"
	"sort"

	g "github.com/dmholtz/graffiti/graph"
)

type KDNode struct {
	node g.PartGeoPoint
	id   g.NodeId
}

// KDPartitioning is a preprocessing step for arc flags and computes a partitioning by constructing a kD-tree of the graph
func KdPartitioning[E g.IHalfEdge](graph *g.AdjacencyArrayGraph[g.PartGeoPoint, E], depth int) *g.AdjacencyArrayGraph[g.PartGeoPoint, E] {
	if depth > 8 {
		panic(fmt.Sprintf("256 bit are reserved for partitions. Got: depth=%d, 2^%d > 256", depth, depth))
	}

	kdNodes := make([]KDNode, 0, graph.NodeCount())
	// reset existing partitions
	for i := 0; i < graph.NodeCount(); i++ {
		point := graph.GetNode(i)
		point.Partition_ = 0
		kdNodes = append(kdNodes, KDNode{node: point, id: i})
	}

	queue := make([][]KDNode, 0)
	queue = append(queue, kdNodes)

	for d := 0; d < depth; d++ {
		end := len(queue)
		fmt.Printf("Depth = %d, length = %d\n", d, end)

		for i := 0; i < end; i++ {
			kdNodes = queue[0]
			queue = queue[1:]

			if d%2 != 0 {
				// north-south split (lat)
				sort.Slice(kdNodes, func(i, j int) bool {
					return kdNodes[i].node.Lat < kdNodes[j].node.Lat
				})
			} else {
				// east-west split (lon)
				sort.Slice(kdNodes, func(i, j int) bool {
					return kdNodes[i].node.Lon < kdNodes[j].node.Lon
				})
			}

			first := kdNodes[:len(kdNodes)/2]
			for j := 0; j < len(first); j++ {
				first[j].node.Partition_ = first[j].node.Partition() << 1
			}

			second := kdNodes[len(kdNodes)/2:]
			for j := 0; j < len(second); j++ {
				second[j].node.Partition_ = (second[j].node.Partition() << 1) + 1
			}

			queue = append(queue, first)
			queue = append(queue, second)
		}
	}

	kdNodes = make([]KDNode, 0)
	for _, s := range queue {
		kdNodes = append(kdNodes, s...)
	}

	for _, kdNode := range kdNodes {
		graph.Nodes[kdNode.id] = kdNode.node
	}

	return graph
}
