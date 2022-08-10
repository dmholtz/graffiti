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
func (fhe FlaggedHalfEdge[W, F]) AddFlag(p PartitionId) IFlaggedHalfEdge[W] {
	fhe.Flag = fhe.Flag | (1 << p)
	return fhe
}

// Adds a flag for partition p to the half edge
func (fhe FlaggedHalfEdge[W, F]) ResetFlag() IFlaggedHalfEdge[W] {
	fhe.Flag = 0
	return fhe
}

// Simple implementation of a half edge with two level arc flags.
type TwoLevelFlaggedHalfEdge[W Weight, F1, F2 FlagType] struct {
	WeightedHalfEdge[W]
	L1Flag F1
	L2Flag F2
}

// IsL1Flagged implements ITwoLevelFlaggedHalfEdge.IsL1Flagged
func (fhe TwoLevelFlaggedHalfEdge[W, F1, F2]) IsL1Flagged(p PartitionId) bool {
	return (fhe.L1Flag & (1 << p)) > 0
}

// IsL2Flagged implements ITwoLevelFlaggedHalfEdge.IsL2Flagged
func (fhe TwoLevelFlaggedHalfEdge[W, F1, F2]) IsL2Flagged(p PartitionId) bool {
	return (fhe.L2Flag & (1 << p)) > 0
}

func (fhe TwoLevelFlaggedHalfEdge[W, F1, F2]) AddL1Flag(p PartitionId) ITwoLevelFlaggedHalfEdge[W] {
	fhe.L1Flag = fhe.L1Flag | (1 << p)
	return fhe
}

func (fhe TwoLevelFlaggedHalfEdge[W, F1, F2]) AddL2Flag(p PartitionId) ITwoLevelFlaggedHalfEdge[W] {
	fhe.L1Flag = fhe.L1Flag | (1 << p)
	return fhe
}

// Adds a flag for partition p to the half edge
func (fhe TwoLevelFlaggedHalfEdge[W, F1, F2]) ResetFlags() ITwoLevelFlaggedHalfEdge[W] {
	fhe.L1Flag = 0
	fhe.L2Flag = 0
	return fhe
}
