package graph

// Type alias for partition identifiers.
type PartitionId = uint16

// Capability description of a node in a partitioned graph
type Partitioner interface {
	Partition() PartitionId
}

type TwoLevelPartitioner interface {
	L1Part() PartitionId
	L2Part() PartitionId
}

// Arc flags are unsigned integers.
type FlagType interface {
	uint64 | uint32 | uint16 | uint8
}

// Capability description of a weighted half edge with arc flags
type IFlaggedHalfEdge[W Weight] interface {
	// IFlaggedHalfEdge inherits all capabilities of IWeightedHalfEdge.
	IWeightedHalfEdge[W]

	// IsFlagged returns true iff the arc flag for the given partitionId is 1.
	IsFlagged(partitionId PartitionId) bool
	// AddFlag sets the arc flag for the given partitionId to 1.
	AddFlag(partitionId PartitionId) IFlaggedHalfEdge[W]
	ResetFlag() IFlaggedHalfEdge[W]
}

// Capability description of a weighted half edge with two level arc flags
type ITwoLevelFlaggedHalfEdge[W Weight] interface {
	// IFlaggedHalfEdge inherits all capabilities of IWeightedHalfEdge.
	IWeightedHalfEdge[W]

	// IsL1Flagged returns true iff the level 1 arc flag for the given partitionId is set.
	IsL1Flagged(partitionId PartitionId) bool
	// IsL2Flagged returns true iff the level 2 arc flag for the given partitionId is set.
	IsL2Flagged(partitionId PartitionId) bool
	// AddFlag sets the level 1 arc flag for the given partitionId to 1.
	AddL1Flag(partitionId PartitionId) ITwoLevelFlaggedHalfEdge[W]

	AddL2Flag(partitionId PartitionId) ITwoLevelFlaggedHalfEdge[W]
	ResetFlags() ITwoLevelFlaggedHalfEdge[W]
}
