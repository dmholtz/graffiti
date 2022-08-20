package graph

import "reflect"

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

// Implementation of a GeoPoint node for a two level partitioned graph
type TwoLevelPartGeoPoint struct {
	GeoPoint
	L1Part_ PartitionId
	L2Part_ PartitionId
}

// L1Part implements Partitioner.L1Part
func (pgp TwoLevelPartGeoPoint) L1Part() PartitionId {
	return pgp.L1Part_
}

// L2Part implements Partitioner.L2Part
func (pgp TwoLevelPartGeoPoint) L2Part() PartitionId {
	return pgp.L2Part_
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

// Simple implementation of a weighted half edge with unsigned integer arc flag.
type FlaggedHalfEdge[W Weight, F FlagType] struct {
	// TODO revert to nested struct once bug in golang has been fixed
	To_     int
	Weight_ W
	Flag    F
}

// To implements IFlaggedHalfEdge.To
func (fhe FlaggedHalfEdge[W, F]) To() NodeId {
	return fhe.To_
}

// Weight implements IFlaggedHalfEdge.Weight
func (fhe FlaggedHalfEdge[W, F]) Weight() W {
	return fhe.Weight_
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

// FlagRange implements IFlaggedHalfEdge.FlagRange
func (fhe FlaggedHalfEdge[W, F]) FlagRange() PartitionId {
	return PartitionId(reflect.TypeOf(fhe.Flag).Bits())
}

// Simple implementation of a half edge with two level arc flags.
type TwoLevelFlaggedHalfEdge[W Weight, F1, F2 FlagType] struct {
	// TODO revert to nested struct once bug in golang has been fixed
	To_     int
	Weight_ W
	L1Flag  F1
	L2Flag  F2
}

// To implements ITwoLevelFlaggedHalfEdge.To
func (fhe TwoLevelFlaggedHalfEdge[W, F1, F2]) To() NodeId {
	return fhe.To_
}

// Weight implements ITwoLevelFlaggedHalfEdge.Weight
func (fhe TwoLevelFlaggedHalfEdge[W, F1, F2]) Weight() W {
	return fhe.Weight_
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
	fhe.L2Flag = fhe.L2Flag | (1 << p)
	return fhe
}

// Adds a flag for partition p to the half edge
func (fhe TwoLevelFlaggedHalfEdge[W, F1, F2]) ResetFlags() ITwoLevelFlaggedHalfEdge[W] {
	fhe.L1Flag = 0
	fhe.L2Flag = 0
	return fhe
}

// L1FlagRange implements ITwoLevelFlaggedHalfEdge.L1FlagRange
func (fhe TwoLevelFlaggedHalfEdge[W, F1, F2]) L1FlagRange() PartitionId {
	return PartitionId(reflect.TypeOf(fhe.L1Flag).Bits())
}

// L2FlagRange implements ITwoLevelFlaggedHalfEdge.L2FlagRange
func (fhe TwoLevelFlaggedHalfEdge[W, F1, F2]) L2FlagRange() PartitionId {
	return PartitionId(reflect.TypeOf(fhe.L2Flag).Bits())
}

// Example of custom half edge with large arc flags (128 bit)
type LargeFlaggedHalfEdge[W Weight] struct {
	// TODO revert to nested struct once bug in golang has been fixed
	To_     int
	Weight_ W
	MsbFlag uint64
	LsbFlag uint64
}

// To implements IFlaggedHalfEdge.To
func (lfe LargeFlaggedHalfEdge[W]) To() NodeId {
	return lfe.To_
}

// Weight implements IFlaggedHalfEdge.Weight
func (lfe LargeFlaggedHalfEdge[W]) Weight() W {
	return lfe.Weight_
}

// IsFlagged implements IFlaggedHalfEdge.IsFlagged
func (lfe LargeFlaggedHalfEdge[W]) IsFlagged(p PartitionId) bool {
	if p < 64 {
		return (lfe.LsbFlag & (1 << p)) > 0
	} else {
		return (lfe.MsbFlag & (1 << (p - 64))) > 0
	}
}

// Adds a flag for partition p to the half edge
func (lfe LargeFlaggedHalfEdge[W]) AddFlag(p PartitionId) IFlaggedHalfEdge[W] {
	if p < 64 {
		lfe.LsbFlag = lfe.LsbFlag | (1 << p)
		return lfe
	} else {
		lfe.MsbFlag = lfe.MsbFlag | (1 << (p - 64))
		return lfe
	}
}

// Resets the flag vector of the half edeg
func (lfe LargeFlaggedHalfEdge[W]) ResetFlag() IFlaggedHalfEdge[W] {
	lfe.LsbFlag = 0
	lfe.MsbFlag = 0
	return lfe
}

// FlagRange implements IFlaggedHalfEdge.FlagRange
func (lfe LargeFlaggedHalfEdge[W]) FlagRange() PartitionId {
	return 128
}
