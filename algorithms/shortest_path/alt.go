package shortest_path

import (
	"math/rand"

	g "github.com/dmholtz/graffiti/graph"
)

// Stores the lenghts of shortest path from/to a landmark L to/from every node
type LandmarkDistances[W g.Weight] struct {
	Landmark g.NodeId // node ID of the landmark
	From     []W      // stores the length of shortest paths from L to every node
	To       []W      // stores the length of shortest paths from every node to L
}

// Heuristic for ALT algorithm (A*, Landmarks and Triangular Inequalities)
type AltHeuristic[W g.Weight] struct {
	LandmarkDistancesCollection map[g.NodeId]LandmarkDistances[W]

	// dynamic attributes
	Source g.NodeId // source node of the current search: updated via Init
	Target g.NodeId // target node of the current search: updated via Init

	ActiveLandmarks []LandmarkDistances[W]
}

func NewAltHeurisitc[N any, E g.IWeightedHalfEdge[W], W g.Weight](graph, transpose g.Graph[N, E], landmarks []g.NodeId) *AltHeuristic[W] {

	// preprocessing
	landmarkDistancesCollection := make(map[g.NodeId]LandmarkDistances[W], 0)
	// TODO implement parallelization using go-routines
	for _, landmark := range landmarks {
		// compute distances from landmark l to every node: one-to-all-dijkstra in (forward) graph starting at l
		distancesFrom := DijkstraOneToAll[N, E, W](graph, landmark).Lengths

		// compute distances from every node to landmark l: one-to-all-dijsktra in transposed graph starting at l
		distancesTo := DijkstraOneToAll[N, E, W](transpose, landmark).Lengths

		landmarkDistances := LandmarkDistances[W]{Landmark: landmark, From: distancesFrom, To: distancesTo}
		landmarkDistancesCollection[landmark] = landmarkDistances
	}

	ah := AltHeuristic[W]{LandmarkDistancesCollection: landmarkDistancesCollection, ActiveLandmarks: make([]LandmarkDistances[W], 0)}

	// default implementation for setting active landmarks: set all landmarks to active landmarks
	// Caveat: The query time might suffer from a too large number of landmarks due to the huge computational overhead.
	for _, landmarkDistances := range landmarkDistancesCollection {
		ah.ActiveLandmarks = append(ah.ActiveLandmarks, landmarkDistances)
	}
	return &ah
}

// Init implements Heuristic.Init
func (ah *AltHeuristic[W]) Init(source g.NodeId, target g.NodeId) {
	ah.Source = source
	ah.Target = target
}

// Evaluate implements Heuristic.Evaluate
func (ah AltHeuristic[W]) Evaluate(id g.NodeId) W {
	upper_bound := W(0)
	for _, landmark := range ah.ActiveLandmarks {
		upper_bound = max(upper_bound, landmark.From[ah.Target]-landmark.From[id])
		upper_bound = max(upper_bound, landmark.To[id]-landmark.To[ah.Target])
	}
	return upper_bound
}

// Maximum Implementation for generic (weight) number types
// max(a, b) returns a iff a is greater or equal than b.
func max[W g.Weight](a, b W) W {
	if a >= b {
		return a
	} else {
		return b
	}
}

// UniformLandmarks chooses n nodes uniformly and at random from the graph.
func UniformLandmarks[N any, E g.IHalfEdge](graph g.Graph[N, E], n int) []g.NodeId {
	landmarks := make([]g.NodeId, 0, n)
	// choose a seed that is different from the seed generating the test sequences
	// otherwise, the landmarks would be severely biased
	rand.Seed(rand.Int63())
	for i := 0; i < n; i++ {
		landmark := rand.Intn(graph.NodeCount())
		landmarks = append(landmarks, landmark)
	}
	return landmarks
}
