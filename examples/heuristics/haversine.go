package heuristics

import (
	"math"

	g "github.com/dmholtz/graffiti/graph"
)

const EARTH_RADIUS = 6371e3 // unit meter

// Heuristic for A* search that estimates the distance between to GeoPoint using the haversine distance.
type HaversineHeuristic[E g.IWeightedHalfEdge[W], W g.Weight] struct {
	Graph  g.Graph[g.GeoPoint, E]
	Target g.NodeId
}

func NewHaversineHeuristic[E g.IWeightedHalfEdge[W], W g.Weight](graph g.Graph[g.GeoPoint, E]) *HaversineHeuristic[E, W] {
	return &HaversineHeuristic[E, W]{Graph: graph}
}

// Init implements Heuristic.Init
func (ah *HaversineHeuristic[E, W]) Init(source g.NodeId, target g.NodeId) {
	ah.Target = target
}

// Evaluate implements Heuristic.Evaluate
func (ah HaversineHeuristic[E, W]) Evaluate(id g.NodeId) W {

	first := ah.Graph.GetNode(id)
	second := ah.Graph.GetNode(ah.Target)

	return W(Haversine(first, second))
}

func Haversine(first, second g.GeoPoint) int {
	deltaPhi := Phi(second) - Phi(first)
	deltaLambda := Lambda(second) - Lambda(first)

	a := math.Pow(math.Sin(deltaPhi/2), 2) + math.Cos(Phi(first))*math.Cos(Phi(second))*math.Pow(math.Sin(deltaLambda/2), 2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return int(EARTH_RADIUS * c * 0.9999)
}

// Latitude in radian
func Phi(p g.GeoPoint) float64 {
	return Deg2Rad(p.Lat)
}

// Longitude in radian
func Lambda(p g.GeoPoint) float64 {
	return Deg2Rad(p.Lon)
}

const ratio = math.Pi / 180

// Convert degree to radian
func Deg2Rad(degree float64) float64 {
	return degree * ratio
}

// Convert radian to degree
func Rad2Deg(radian float64) float64 {
	return radian / ratio
}
