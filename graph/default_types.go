package graph

// Node types

// Implementation of a blank node
type Node struct{}

// Implementation of a node representing a point in the geographic corrdinate system (latitude / longitude)
type GeoPoint struct {
	Lat float64
	Lon float64
}

// Implementation of a GeoPoint node for a partitioned graph
type PartGeoPoint struct {
	GeoPoint
	Partition_ PartitionId
}

// Partition implements Partitioner.Partition
func (pgp PartGeoPoint) Partition() PartitionId {
	return pgp.Partition_
}

// Edge types

// Simple implementation of a weighted half edge without any additional metadata
type WeightedHalfEdge[W Weight] struct {
	To_     NodeId
	Weight_ W
}

// Constructor method
func NewWeightedHalfEdge[W Weight](to int, weight W) WeightedHalfEdge[W] {
	return WeightedHalfEdge[W]{To_: to, Weight_: weight}
}

// To implements IHalfEdge.To
func (e WeightedHalfEdge[W]) To() NodeId {
	return e.To_
}

// Weight implements IWeightedHalfEdge.Weight
func (e WeightedHalfEdge[W]) Weight() W {
	return e.Weight_
}

// Simple implementation of a half edge with unsigned integer arc flag.
type FlaggedHalfEdge[W Weight, F FlagType] struct {
	WeightedHalfEdge[W]
	Flag F
}

// IsFlagged implements IFlaggedHalfEdge.IsFlagged
func (fhe FlaggedHalfEdge[W, F]) IsFlagged(p PartitionId) bool {
	return (fhe.Flag & (1 << p)) > 0
}

// Adds a flag for partition p to the half edge
func (fhe *FlaggedHalfEdge[W, F]) AddFlag(p PartitionId) {
	fhe.Flag = fhe.Flag | (1 << p)
}
