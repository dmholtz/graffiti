package partitioning

import (
	"math"

	g "github.com/dmholtz/graffiti/graph"
)

// GridPartitioning is a preprocessing step for arc flags and computes the grid partitioning of a graph of GeoPoints.
func GridPartitioning[E g.IHalfEdge](graph *g.AdjacencyArrayGraph[g.PartGeoPoint, E], lat, lon int) *g.AdjacencyArrayGraph[g.PartGeoPoint, E] {

	for i := 0; i < graph.NodeCount(); i++ {
		geoPoint := graph.GetNode(i)

		// determine column (lon) and row (lat) index of geoPoint in the grid
		col := int(math.Min(((geoPoint.Lon + 180) / 360 * float64(lon)), float64(lon-1)))
		row := int(math.Min(((geoPoint.Lat + 90) / 180 * float64(lat)), float64(lat-1)))

		partition := row*lon + col

		geoPoint.Partition_ = g.PartitionId(partition)

		graph.Nodes[i] = geoPoint
	}

	return graph
}
