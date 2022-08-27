package partitioning

import (
	"fmt"
	"math"

	g "github.com/dmholtz/graffiti/graph"
)

// GridPartitioning is a preprocessing step for two-level arc flags and computes the two-level grid partitioning of a graph of GeoPoints.
func TwoLevelGridPartitioning[E g.IHalfEdge](graph *g.AdjacencyArrayGraph[g.TwoLevelPartGeoPoint, E], l1_lat, l1_lon, l2_lat, l2_lon int) *g.AdjacencyArrayGraph[g.TwoLevelPartGeoPoint, E] {
	if l1_lat*l1_lon > 32 {
		panic(fmt.Sprintf("32 bit are reserved for level 1 partitions. Got: l1_lat * l1_lon > 32"))
	}
	if l2_lat*l2_lon > 32 {
		panic(fmt.Sprintf("32 bit are reserved for level 2 partitions. Got: l2_lat * l2_lon > 32"))
	}

	l_lat := l1_lat * l2_lat
	l_lon := l1_lon * l2_lon

	for i := 0; i < graph.NodeCount(); i++ {
		geoPoint := graph.GetNode(i)

		// determine column (lon) and row (lat) index of geoPoint in the grid
		col := int(math.Min(((geoPoint.Lon + 180) / 360 * float64(l_lon)), float64(l_lon-1)))
		row := int(math.Min(((geoPoint.Lat + 90) / 180 * float64(l_lat)), float64(l_lat-1)))

		level1_partition := (row/l1_lat)*l1_lon + (col / l1_lon)
		level2_partition := (row%l1_lat)*l2_lon + (col % l1_lon)

		geoPoint.L1Part_ = g.PartitionId(level1_partition)
		geoPoint.L2Part_ = g.PartitionId(level2_partition)

		graph.Nodes[i] = geoPoint
	}

	return graph
}
