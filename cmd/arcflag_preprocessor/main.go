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
	falg := fmi.NewAdjacencyListFromFmi(inputGraphFile, fmi.ParsePartGeoPoint, fmi.ParseFlaggedHalfEdge)
	faag := g.NewAdjacencyArrayFromGraph[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]](falg)
	elapsed := time.Since(start)
	fmt.Printf("[TIME-FileReader] = %s\n", elapsed)

	start = time.Now()
	faag = partitioning.GridPartitioning(faag, 8, 8) // 64 partitions
	//faag = partitioning.KdPartitioning(faag, 6)      // 64 partitions
	elapsed = time.Since(start)
	fmt.Printf("[TIME-Partitioning] = %s\n", elapsed)

	start = time.Now()
	faag = sp.ComputeArcFlags[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int](faag, faag, 64)
	elapsed = time.Since(start)
	fmt.Printf("[TIME-ArcFlagComputation] = %s\n", elapsed)

	fmi.WriteFmi[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]](faag, outputGraphFile, fmi.PartGeoPoint2FmiLine, fmi.FlaggedHalfEdge2FmiLine)

	falg = fmi.NewAdjacencyListFromFmi(outputGraphFile, fmi.ParsePartGeoPoint, fmi.ParseFlaggedHalfEdge)

	testedRouter := sp.ArcFlagRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: falg}
	baselineRouter := sp.DijkstraRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: falg}

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
