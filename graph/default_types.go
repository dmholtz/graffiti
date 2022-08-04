package graph

// Node types

// Implementation of a blank node
type Node struct{}

// Implementation of a node representing a point in the geographic corrdinate system (latitude / longitude)
type GeoPoint struct {
	Lat float64
	Lon float64
}

// Edge types

// Simple implementation of a weighted half edge without any additional metadata
type WeightedHalfEdge[W Weight] struct {
	to     NodeId
	weight W
}

// Constructor method
func NewWeightedHalfEdge[W Weight](to int, weight W) WeightedHalfEdge[W] {
	return WeightedHalfEdge[W]{to: to, weight: weight}
}

// To implements IHalfEdge.To
func (e WeightedHalfEdge[W]) To() NodeId {
	return e.to
}

// Weight implements IWeightedHalfEdge.Weight
func (e WeightedHalfEdge[W]) Weight() W {
	return e.weight
}
