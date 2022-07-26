package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	"github.com/dmholtz/graffiti/examples/heuristics"
	fmi "github.com/dmholtz/graffiti/examples/io"
	g "github.com/dmholtz/graffiti/graph"
)

const defaultGraph = "graphs/ocean_equi_4.fmi"
const arcflag8 = "graphs/ocean_equi_4_grid_arcflags8_8.fmi"
const arcflag16 = "graphs/ocean_equi_4_grid_arcflags16_16.fmi"
const arcflag32 = "graphs/ocean_equi_4_grid_arcflags32_32.fmi"
const arcflag64 = "graphs/ocean_equi_4_grid_arcflags64.fmi"
const arcflag64_kd = "graphs/ocean_equi_4_kd_arcflags64.fmi"
const arcflag128 = "graphs/ocean_equi_4_grid_arcflags128.fmi"
const arcflag256 = "graphs/ocean_equi_4_grid_arcflags256.fmi"

const NUMBER_OF_RUNS = 1000

type BenchmarkTask struct {
	Name       string
	ResultFile string
	Benchmark  *sp.Benchmarker[int]
}

func main() {
	Baseline(false)
	CompareArcFlagSize(false)
	CompareGridType(false)
	CompareTwoLevelArcFlagSize(false)
	CompareAStar(false)
	CompareLandmarkCount(false)
	CompareLandmarkSelection(false)
	EvaluateArcflagAlt(false)
}

func Baseline(export bool) {

	// Load graphs

	alg := fmi.NewAdjacencyListFromFmi(defaultGraph, fmi.ParseGeoPoint, fmi.ParseWeightedHalfEdge)
	aag := g.NewAdjacencyArrayFromGraph[g.GeoPoint, g.WeightedHalfEdge[int]](alg)

	n := aag.NodeCount()

	// Build routers

	dijkstraRouter := sp.DijkstraRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag}
	dijkstraBenchmark := BenchmarkTask{Name: "Dijkstra's Algorithm", Benchmark: sp.NewBenchmarker[int](dijkstraRouter, n), ResultFile: "benchmarks/dijkstra.json"}

	biDijkstraRouter := sp.BiDijkstraRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Transpose: aag, MaxInitializerValue: math.MaxInt}
	biDijkstraBenchmark := BenchmarkTask{Name: "bidirectional Dijkstra's Algorithm", Benchmark: sp.NewBenchmarker[int](biDijkstraRouter, n), ResultFile: "benchmarks/bi-dijkstra.json"}

	RunBenchmarks([]BenchmarkTask{
		dijkstraBenchmark,
		biDijkstraBenchmark},
		NUMBER_OF_RUNS,
		export)

}

func CompareAStar(export bool) {
	// Load graphs

	alg := fmi.NewAdjacencyListFromFmi(defaultGraph, fmi.ParseGeoPoint, fmi.ParseWeightedHalfEdge)
	aag := g.NewAdjacencyArrayFromGraph[g.GeoPoint, g.WeightedHalfEdge[int]](alg)

	n := aag.NodeCount()

	// Build routers

	havHeuristic := heuristics.NewHaversineHeuristic[g.WeightedHalfEdge[int]](alg)
	havBackwardHeuristic := heuristics.NewHaversineHeuristic[g.WeightedHalfEdge[int]](alg)

	aStarRouter := sp.AStarRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Heuristic: havHeuristic}
	aStarBenchmark := BenchmarkTask{Name: "A* Search", Benchmark: sp.NewBenchmarker[int](aStarRouter, n), ResultFile: "benchmarks/astar.json"}

	biAStarRouter := sp.BidirectionalAStarRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Transpose: aag, ForwardHeuristic: havHeuristic, BackwardHeuristic: havBackwardHeuristic, MaxInitializerValue: math.MaxInt}
	biAStarBenchmark := BenchmarkTask{Name: "Bidirectional A* Search", Benchmark: sp.NewBenchmarker[int](biAStarRouter, n), ResultFile: "benchmarks/bi-astar.json"}

	RunBenchmarks([]BenchmarkTask{
		aStarBenchmark,
		biAStarBenchmark},
		NUMBER_OF_RUNS,
		export)

}

func CompareLandmarkCount(export bool) {
	// Load graphs

	alg := fmi.NewAdjacencyListFromFmi(defaultGraph, fmi.ParseGeoPoint, fmi.ParseWeightedHalfEdge)
	aag := g.NewAdjacencyArrayFromGraph[g.GeoPoint, g.WeightedHalfEdge[int]](alg)

	n := aag.NodeCount()

	// choose landmarks
	landmarks := sp.UniformLandmarks[g.GeoPoint, g.WeightedHalfEdge[int]](aag, 64)

	// Build routers
	alt2 := sp.NewAltHeurisitc[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, aag, landmarks[:2])
	alt4 := sp.NewAltHeurisitc[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, aag, landmarks[:4])
	alt8 := sp.NewAltHeurisitc[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, aag, landmarks[:8])
	alt16 := sp.NewAltHeurisitc[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, aag, landmarks[:16])
	alt32 := sp.NewAltHeurisitc[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, aag, landmarks[:32])
	alt64 := sp.NewAltHeurisitc[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, aag, landmarks[:64])

	alt2Router := sp.AStarRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Heuristic: alt2}
	alt2Benchmark := BenchmarkTask{Name: "ALT-2", Benchmark: sp.NewBenchmarker[int](alt2Router, n), ResultFile: "benchmarks/alt-2.json"}

	alt4Router := sp.AStarRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Heuristic: alt4}
	alt4Benchmark := BenchmarkTask{Name: "ALT-4", Benchmark: sp.NewBenchmarker[int](alt4Router, n), ResultFile: "benchmarks/alt-4.json"}

	alt8Router := sp.AStarRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Heuristic: alt8}
	alt8Benchmark := BenchmarkTask{Name: "ALT-8", Benchmark: sp.NewBenchmarker[int](alt8Router, n), ResultFile: "benchmarks/alt-8.json"}

	alt16Router := sp.AStarRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Heuristic: alt16}
	alt16Benchmark := BenchmarkTask{Name: "ALT-16", Benchmark: sp.NewBenchmarker[int](alt16Router, n), ResultFile: "benchmarks/alt-16.json"}

	alt32Router := sp.AStarRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Heuristic: alt32}
	alt32Benchmark := BenchmarkTask{Name: "ALT-32", Benchmark: sp.NewBenchmarker[int](alt32Router, n), ResultFile: "benchmarks/alt-32.json"}

	alt64Router := sp.AStarRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Heuristic: alt64}
	alt64Benchmark := BenchmarkTask{Name: "ALT-64", Benchmark: sp.NewBenchmarker[int](alt64Router, n), ResultFile: "benchmarks/alt-64.json"}

	RunBenchmarks([]BenchmarkTask{
		alt2Benchmark,
		alt4Benchmark,
		alt8Benchmark,
		alt16Benchmark,
		alt32Benchmark,
		alt64Benchmark},
		NUMBER_OF_RUNS,
		export)
}

func CompareLandmarkSelection(export bool) {
	// Load graphs

	alg := fmi.NewAdjacencyListFromFmi(defaultGraph, fmi.ParseGeoPoint, fmi.ParseWeightedHalfEdge)
	aag := g.NewAdjacencyArrayFromGraph[g.GeoPoint, g.WeightedHalfEdge[int]](alg)

	n := aag.NodeCount()

	// choose landmarks
	randomLandmarks := sp.UniformLandmarks[g.GeoPoint, g.WeightedHalfEdge[int]](aag, 8)
	oceanLandmarks := LoadLandmarkFile("graphs/landmarks/landmarks_ocean8.json")
	coastLandmarks := LoadLandmarkFile("graphs/landmarks/landmarks_coast8.json")

	// Build routers
	altRand := sp.NewAltHeurisitc[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, aag, randomLandmarks)
	altOcean := sp.NewAltHeurisitc[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, aag, oceanLandmarks)
	altCoast := sp.NewAltHeurisitc[g.GeoPoint, g.WeightedHalfEdge[int], int](aag, aag, coastLandmarks)

	altRandRouter := sp.AStarRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Heuristic: altRand}
	altRandBenchmark := BenchmarkTask{Name: "ALT (random landmarks)", Benchmark: sp.NewBenchmarker[int](altRandRouter, n), ResultFile: "benchmarks/alt-8-random.json"} // identical test is run in CompareLandmarkCount

	altOceanRouter := sp.AStarRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Heuristic: altOcean}
	altOceanBenchmark := BenchmarkTask{Name: "ALT (ocean landmarks)", Benchmark: sp.NewBenchmarker[int](altOceanRouter, n), ResultFile: "benchmarks/alt-8-ocean.json"}

	altCoastRouter := sp.AStarRouter[g.GeoPoint, g.WeightedHalfEdge[int], int]{Graph: aag, Heuristic: altCoast}
	altCoastBenchmark := BenchmarkTask{Name: "ALT (coast landmarks)", Benchmark: sp.NewBenchmarker[int](altCoastRouter, n), ResultFile: "benchmarks/alt-8-coast.json"}

	RunBenchmarks([]BenchmarkTask{
		altRandBenchmark,
		altOceanBenchmark,
		altCoastBenchmark},
		NUMBER_OF_RUNS,
		export)
}

func CompareArcFlagSize(export bool) {

	// Load graphs

	falg8 := fmi.NewAdjacencyListFromFmi(arcflag8, fmi.ParsePartGeoPoint, fmi.ParseFlaggedHalfEdge)
	faag8 := g.NewAdjacencyArrayFromGraph[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]](falg8)

	falg16 := fmi.NewAdjacencyListFromFmi(arcflag16, fmi.ParsePartGeoPoint, fmi.ParseFlaggedHalfEdge)
	faag16 := g.NewAdjacencyArrayFromGraph[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]](falg16)

	falg32 := fmi.NewAdjacencyListFromFmi(arcflag32, fmi.ParsePartGeoPoint, fmi.ParseFlaggedHalfEdge)
	faag32 := g.NewAdjacencyArrayFromGraph[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]](falg32)

	falg64 := fmi.NewAdjacencyListFromFmi(arcflag64, fmi.ParsePartGeoPoint, fmi.ParseFlaggedHalfEdge)
	faag64 := g.NewAdjacencyArrayFromGraph[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]](falg64)

	falg128 := fmi.NewAdjacencyListFromFmi(arcflag128, fmi.ParsePartGeoPoint, fmi.ParseLargeFlaggedHalfEdge)
	faag128 := g.NewAdjacencyArrayFromGraph[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int]](falg128)

	falg256 := fmi.NewAdjacencyListFromFmi(arcflag256, fmi.ParsePartGeoPoint, fmi.Parse256BitFlaggedHalfEdge)
	faag256 := g.NewAdjacencyArrayFromGraph[g.PartGeoPoint, g.B256FlaggedHalfEdge[int]](falg256)

	n := faag8.NodeCount()

	// Build routers

	arcflag8Router := sp.ArcFlagRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: faag8}
	arcflag8Benchmark := BenchmarkTask{Name: "8-bit arc flags", Benchmark: sp.NewBenchmarker[int](arcflag8Router, n), ResultFile: "benchmarks/arcflag8.json"}

	biArcflag8Router := sp.BidirectionalArcFlagRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: faag8, Transpose: faag8, MaxInitializerValue: math.MaxInt}
	biArcflag8Benchmark := BenchmarkTask{Name: "bidirectional 8-bit arc flags", Benchmark: sp.NewBenchmarker[int](biArcflag8Router, n), ResultFile: "benchmarks/bi-arcflag8.json"}

	arcflag16Router := sp.ArcFlagRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: faag16}
	arcflag16Benchmark := BenchmarkTask{Name: "16-bit arc flags", Benchmark: sp.NewBenchmarker[int](arcflag16Router, n), ResultFile: "benchmarks/arcflag16.json"}

	biArcflag16Router := sp.BidirectionalArcFlagRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: faag16, Transpose: faag16, MaxInitializerValue: math.MaxInt}
	biArcflag16Benchmark := BenchmarkTask{Name: "bidirectional 16-bit arc flags", Benchmark: sp.NewBenchmarker[int](biArcflag16Router, n), ResultFile: "benchmarks/bi-arcflag16.json"}

	arcflag32Router := sp.ArcFlagRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: faag32}
	arcflag32Benchmark := BenchmarkTask{Name: "32-bit arc flags", Benchmark: sp.NewBenchmarker[int](arcflag32Router, n), ResultFile: "benchmarks/arcflag32.json"}

	biArcflag32Router := sp.BidirectionalArcFlagRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: faag32, Transpose: faag32, MaxInitializerValue: math.MaxInt}
	biArcflag32Benchmark := BenchmarkTask{Name: "bidirectional 32-bit arc flags", Benchmark: sp.NewBenchmarker[int](biArcflag32Router, n), ResultFile: "benchmarks/bi-arcflag32.json"}

	arcflag64Router := sp.ArcFlagRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: faag64}
	arcflag64Benchmark := BenchmarkTask{Name: "64-bit arc flags", Benchmark: sp.NewBenchmarker[int](arcflag64Router, n), ResultFile: "benchmarks/arcflag64.json"}

	biArcflag64Router := sp.BidirectionalArcFlagRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: faag64, Transpose: faag64, MaxInitializerValue: math.MaxInt}
	biArcflag64Benchmark := BenchmarkTask{Name: "bidirectional 64-bit arc flags", Benchmark: sp.NewBenchmarker[int](biArcflag64Router, n), ResultFile: "benchmarks/bi-arcflag64.json"}

	arcflag128Router := sp.ArcFlagRouter[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int], int]{Graph: faag128}
	arcflag128Benchmark := BenchmarkTask{Name: "128-bit arc flags", Benchmark: sp.NewBenchmarker[int](arcflag128Router, n), ResultFile: "benchmarks/arcflag128.json"}

	biArcflag128Router := sp.BidirectionalArcFlagRouter[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int], int]{Graph: faag128, Transpose: faag128, MaxInitializerValue: math.MaxInt}
	biArcflag128Benchmark := BenchmarkTask{Name: "bidirectional 128-bit arc flags", Benchmark: sp.NewBenchmarker[int](biArcflag128Router, n), ResultFile: "benchmarks/bi-arcflag128.json"}

	arcflag256Router := sp.ArcFlagRouter[g.PartGeoPoint, g.B256FlaggedHalfEdge[int], int]{Graph: faag256}
	arcflag256Benchmark := BenchmarkTask{Name: "256-bit arc flags", Benchmark: sp.NewBenchmarker[int](arcflag256Router, n), ResultFile: "benchmarks/arcflag256.json"}

	biArcflag256Router := sp.BidirectionalArcFlagRouter[g.PartGeoPoint, g.B256FlaggedHalfEdge[int], int]{Graph: faag256, Transpose: faag256, MaxInitializerValue: math.MaxInt}
	biArcflag256Benchmark := BenchmarkTask{Name: "bidirectional 256-bit arc flags", Benchmark: sp.NewBenchmarker[int](biArcflag256Router, n), ResultFile: "benchmarks/bi-arcflag256.json"}

	RunBenchmarks([]BenchmarkTask{
		arcflag8Benchmark,
		biArcflag8Benchmark,
		arcflag16Benchmark,
		biArcflag16Benchmark,
		arcflag32Benchmark,
		biArcflag32Benchmark,
		arcflag64Benchmark,
		biArcflag64Benchmark,
		arcflag128Benchmark,
		biArcflag128Benchmark,
		arcflag256Benchmark,
		biArcflag256Benchmark},
		NUMBER_OF_RUNS,
		export)
}

func CompareTwoLevelArcFlagSize(export bool) {

	// Load graphs

	falg8 := fmi.NewAdjacencyListFromFmi(arcflag8, fmi.Parse2LPartGeoPoint, fmi.Parse2LFlaggedHalfEdge)
	faag8 := g.NewAdjacencyArrayFromGraph[g.TwoLevelPartGeoPoint, g.TwoLevelFlaggedHalfEdge[int, uint64, uint64]](falg8)

	falg16 := fmi.NewAdjacencyListFromFmi(arcflag16, fmi.Parse2LPartGeoPoint, fmi.Parse2LFlaggedHalfEdge)
	faag16 := g.NewAdjacencyArrayFromGraph[g.TwoLevelPartGeoPoint, g.TwoLevelFlaggedHalfEdge[int, uint64, uint64]](falg16)

	falg32 := fmi.NewAdjacencyListFromFmi(arcflag32, fmi.Parse2LPartGeoPoint, fmi.Parse2LFlaggedHalfEdge)
	faag32 := g.NewAdjacencyArrayFromGraph[g.TwoLevelPartGeoPoint, g.TwoLevelFlaggedHalfEdge[int, uint64, uint64]](falg32)

	n := faag8.NodeCount()

	// Build routers

	arcflag8Router := sp.TwoLevelArcFlagRouter[g.TwoLevelPartGeoPoint, g.TwoLevelFlaggedHalfEdge[int, uint64, uint64], int]{Graph: faag8}
	arcflag8Benchmark := BenchmarkTask{Name: "two-level 8-bit arc flags", Benchmark: sp.NewBenchmarker[int](arcflag8Router, n), ResultFile: "benchmarks/arcflag8-2level.json"}

	arcflag16Router := sp.TwoLevelArcFlagRouter[g.TwoLevelPartGeoPoint, g.TwoLevelFlaggedHalfEdge[int, uint64, uint64], int]{Graph: faag16}
	arcflag16Benchmark := BenchmarkTask{Name: "two-level 16-bit arc flags", Benchmark: sp.NewBenchmarker[int](arcflag16Router, n), ResultFile: "benchmarks/arcflag16-2level.json"}

	arcflag32Router := sp.TwoLevelArcFlagRouter[g.TwoLevelPartGeoPoint, g.TwoLevelFlaggedHalfEdge[int, uint64, uint64], int]{Graph: faag32}
	arcflag32Benchmark := BenchmarkTask{Name: "two-level 32-bit arc flags", Benchmark: sp.NewBenchmarker[int](arcflag32Router, n), ResultFile: "benchmarks/arcflag32-2level.json"}

	RunBenchmarks([]BenchmarkTask{
		arcflag8Benchmark,
		arcflag16Benchmark,
		arcflag32Benchmark},
		NUMBER_OF_RUNS,
		export)
}

func CompareGridType(export bool) {

	// Load graphs

	gridAlg64 := fmi.NewAdjacencyListFromFmi(arcflag64, fmi.ParsePartGeoPoint, fmi.ParseFlaggedHalfEdge)
	gridAag64 := g.NewAdjacencyArrayFromGraph[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]](gridAlg64)

	kdAlg64 := fmi.NewAdjacencyListFromFmi(arcflag64_kd, fmi.ParsePartGeoPoint, fmi.ParseFlaggedHalfEdge)
	kdAag64 := g.NewAdjacencyArrayFromGraph[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64]](kdAlg64)

	n := kdAag64.NodeCount()

	// Build routers

	arcflagGridRouter := sp.ArcFlagRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: gridAag64}
	arcflagGridBenchmark := BenchmarkTask{Name: "64-bit arc flags (grid)", Benchmark: sp.NewBenchmarker[int](arcflagGridRouter, n), ResultFile: "benchmarks/arcflag64_grid.json"}

	arcflagKdRouter := sp.ArcFlagRouter[g.PartGeoPoint, g.FlaggedHalfEdge[int, uint64], int]{Graph: kdAag64}
	arcflagkdBenchmark := BenchmarkTask{Name: "64-bit arc flags (kd)", Benchmark: sp.NewBenchmarker[int](arcflagKdRouter, n), ResultFile: "benchmarks/bi-arcflag64.json"}

	RunBenchmarks([]BenchmarkTask{
		arcflagGridBenchmark,
		arcflagkdBenchmark},
		NUMBER_OF_RUNS,
		export)
}

func EvaluateArcflagAlt(export bool) {

	// Load graphs
	falg256 := fmi.NewAdjacencyListFromFmi("graphs/ocean_equi_4_grid_arcflags256.fmi", fmi.ParsePartGeoPoint, fmi.Parse256BitFlaggedHalfEdge)
	faag256 := g.NewAdjacencyArrayFromGraph[g.PartGeoPoint, g.B256FlaggedHalfEdge[int]](falg256)

	n := faag256.NodeCount()

	// choose landmarks
	landmarks := sp.UniformLandmarks[g.PartGeoPoint, g.B256FlaggedHalfEdge[int]](faag256, 16)

	// precompute heuristic
	alt := sp.NewAltHeurisitc[g.PartGeoPoint, g.B256FlaggedHalfEdge[int], int](faag256, faag256, landmarks)

	// Build router
	arcflagAltRouter := sp.ArcFlagAStarRouter[g.PartGeoPoint, g.B256FlaggedHalfEdge[int], int]{Graph: faag256, Transpose: faag256, Heuristic: alt}
	arcflagAltBenchmark := BenchmarkTask{Name: "bidirectional 256-bit arc flags + ALT", Benchmark: sp.NewBenchmarker[int](arcflagAltRouter, n), ResultFile: "benchmarks/arcflag-alt.json"}

	RunBenchmarks([]BenchmarkTask{
		arcflagAltBenchmark},
		NUMBER_OF_RUNS,
		export)
}

func RunBenchmarks(tasks []BenchmarkTask, n int, export bool) {
	for _, task := range tasks {
		fmt.Printf("Run benchmark '%s'\n", task.Name)
		benchmark := task.Benchmark
		summary := benchmark.Run(n)
		fmt.Println(summary)
		if export {
			SaveBenchmark(task)
		}
	}
}

func SaveBenchmark(task BenchmarkTask) {
	file, _ := json.Marshal(task.Benchmark.Result)
	os.WriteFile(task.ResultFile, file, 0644)
}

func LoadLandmarkFile(filename string) []g.NodeId {
	var landmarks []g.NodeId
	file, _ := ioutil.ReadFile(filename)
	_ = json.Unmarshal([]byte(file), &landmarks)
	return landmarks
}
