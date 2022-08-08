package graph

// Type alias for partition identifiers.
type PartitionId = uint16

// Capability description of a node in a partitioned graph
type Partitioner interface {
	Partition() PartitionId
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
	//AddFlag(partitionId PartitionId)
}
