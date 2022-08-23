package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	fmi "github.com/dmholtz/graffiti/examples/io"
	g "github.com/dmholtz/graffiti/graph"
)

const defaultGraph = "graphs/ocean_equi_4.fmi"
const arcflag8 = "graphs/ocean_equi_4_grid_arcflags8_8.fmi"
const arcflag16 = "graphs/ocean_equi_4_grid_arcflags16_16.fmi"
const arcflag32 = "graphs/ocean_equi_4_grid_arcflags32_32.fmi"
const arcflag64 = "/Users/david/repos/osm-ship-routing/graphs/ocean_equi_4_grid_arcflags.fmi"
const arcflag128 = "graphs/ocean_equi_4_grid_arcflags128.fmi"

const NUMBER_OF_RUNS = 1000

type BenchmarkTask struct {
	Name       string
	ResultFile string
	Benchmark  *sp.Benchmarker[int]
}

func main() {
	//Baseline(true)
	//CompareArcFlagSize(true)
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

	falg128 := fmi.NewAdjacencyListFromFmi("graphs/ocean_equi_4_grid_arcflags128.fmi", fmi.ParsePartGeoPoint, fmi.ParseLargeFlaggedHalfEdge)
	faag128 := g.NewAdjacencyArrayFromGraph[g.PartGeoPoint, g.LargeFlaggedHalfEdge[int]](falg128)

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
		biArcflag128Benchmark},
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
