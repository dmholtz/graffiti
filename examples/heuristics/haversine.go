package heuristics

import (
	"math"

	g "github.com/dmholtz/graffiti/graph"
)

const EARTH_RADIUS = 6371e3 // unit meter

// Heuristic for A* search that estimates the distance between to GeoPoint using the haversine distance.
type HaversineHeuristic[E g.IWeightedHalfEdge[int]] struct {
	Graph  g.Graph[g.GeoPoint, E]
	Target g.GeoPoint
}

func NewHaversineHeuristic[E g.IWeightedHalfEdge[int]](graph g.Graph[g.GeoPoint, E]) *HaversineHeuristic[E] {
	return &HaversineHeuristic[E]{Graph: graph}
}

// Init implements Heuristic.Init
func (ah *HaversineHeuristic[E]) Init(source g.NodeId, target g.NodeId) {
	ah.Target = ah.Graph.GetNode(target)
}

// Evaluate implements Heuristic.Evaluate
func (ah HaversineHeuristic[E]) Evaluate(id g.NodeId) int {
	source := ah.Graph.GetNode(id)
	return Haversine(source, ah.Target)
}

func Haversine(first, second g.GeoPoint) int {
	phi1, phi2 := Phi(first), Phi(second)
	lambda1, lambda2 := Lambda(first), Lambda(second)

	deltaPhi := phi2 - phi1
	deltaLambda := lambda2 - lambda1

	a := math.Pow(math.Sin(deltaPhi/2), 2) + math.Cos(phi1)*math.Cos(phi2)*math.Pow(math.Sin(deltaLambda/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return int(EARTH_RADIUS * c * 0.9999)
}

// Latitude in radian
func Phi(p g.GeoPoint) float64 {
	return p.Lat * ratio // degree to radian
}

// Longitude in radian
func Lambda(p g.GeoPoint) float64 {
	return p.Lon * ratio // degree to radian
}

const ratio = math.Pi / 180
