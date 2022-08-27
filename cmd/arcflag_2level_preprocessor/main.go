package main

import (
	"fmt"
	"math/rand"
	"time"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	fmi "github.com/dmholtz/graffiti/examples/io"
	"github.com/dmholtz/graffiti/examples/partitioning"
	g "github.com/dmholtz/graffiti/graph"
)

//const inputGraphFile = "graphs/ocean_10k.fmi"
const inputGraphFile = "graphs/ocean_equi_4.fmi"

const outputGraphFile = "out.fmi"

func main() {

	start := time.Now()
	falg := fmi.NewAdjacencyListFromFmi(inputGraphFile, fmi.Parse2LPartGeoPoint, fmi.Parse2LFlaggedHalfEdge)
	faag := g.NewAdjacencyArrayFromGraph[g.TwoLevelPartGeoPoint, g.TwoLevelFlaggedHalfEdge[int, uint64, uint64]](falg)
	elapsed := time.Since(start)
	fmt.Printf("[TIME-FileReader] = %s\n", elapsed)

	start = time.Now()
	faag = partitioning.TwoLevelGridPartitioning(faag, 4, 8, 4, 8) // 32x32 partitions
	//faag = partitioning.TwoLevelGridPartitioning(faag, 4, 4, 4, 4) // 16x16 partitions
	//faag = partitioning.TwoLevelGridPartitioning(faag, 2, 4, 2, 4) // 8x8 partitions
	elapsed = time.Since(start)
	fmt.Printf("[TIME-Partitioning] = %s\n", elapsed)

	start = time.Now()
	faag = sp.ComputeTwoLevelArcFlags[g.TwoLevelPartGeoPoint, g.TwoLevelFlaggedHalfEdge[int, uint64, uint64], int](faag, faag)
	elapsed = time.Since(start)
	fmt.Printf("[TIME-ArcFlagComputation] = %s\n", elapsed)

	fmi.WriteFmi[g.TwoLevelPartGeoPoint, g.TwoLevelFlaggedHalfEdge[int, uint64, uint64]](faag, outputGraphFile, fmi.TwoLevelPartGeoPoint2FmiLine, fmi.TwoLevelFlaggedHalfEdge2FmiLine)

	falg = fmi.NewAdjacencyListFromFmi(outputGraphFile, fmi.Parse2LPartGeoPoint, fmi.Parse2LFlaggedHalfEdge)

	testedRouter := sp.TwoLevelArcFlagRouter[g.TwoLevelPartGeoPoint, g.TwoLevelFlaggedHalfEdge[int, uint64, uint64], int]{Graph: falg}
	baselineRouter := sp.DijkstraRouter[g.TwoLevelPartGeoPoint, g.TwoLevelFlaggedHalfEdge[int, uint64, uint64], int]{Graph: falg}

	n := 100 // number of random tests
	fmt.Printf("Compare %d random searches of bidirectional Dijkstra against textbook Dijkstra.\n", n)
	dijkstraPqPops, arcFlagDijkstraPqPops := 0, 0
	for i := 0; i < n; i++ {
		source := rand.Intn(faag.NodeCount())
		target := rand.Intn(faag.NodeCount())

		dijkstraResult := baselineRouter.Route(source, target, false)
		arcFlagDijkstraResult := testedRouter.Route(source, target, false)

		if dijkstraResult.Length != arcFlagDijkstraResult.Length {
			fmt.Printf("[Path(source=%d, target=%d)]: Different lengths found: Dijkstra=%d, ArcFlagDijkstra=%d\n", source, target, dijkstraResult.Length, arcFlagDijkstraResult.Length)
		}

		// maintain performance indicators
		dijkstraPqPops += dijkstraResult.PqPops
		arcFlagDijkstraPqPops += arcFlagDijkstraResult.PqPops
	}
	fmt.Printf("Avgerage number of Pop() operations on priority queue: %d (Dijkstra), %d (BiDijkstra)\n", dijkstraPqPops/n, arcFlagDijkstraPqPops/n)
}
