package benchmarks

import (
	"math/rand"
	"time"

	sp "github.com/dmholtz/graffiti/algorithms/shortest_path"
	g "github.com/dmholtz/graffiti/graph"
)

const DEFAULT_SEED = 314159265359

type Benchmarker[W g.Weight] struct {
	NodeRange g.NodeId
	Router    sp.Router[W]
	Result    BenchmarkResult
}

func NewBenchmarker[W g.Weight](router sp.Router[W], nodeCount int) *Benchmarker[W] {
	return &Benchmarker[W]{NodeRange: nodeCount, Router: router, Result: *NewBenchmarkResult()}
}

func (b Benchmarker[W]) Run(n int) Summary {
	rand.Seed(DEFAULT_SEED)

	for i := 0; i < n; i++ {
		source := rand.Intn(b.NodeRange)
		target := rand.Intn(b.NodeRange)

		start := time.Now()
		routingResult := b.Router.Route(source, target, false)
		time := float64(time.Since(start)) / 1000000 // ms

		b.Result.Add(time, routingResult.PqPops)
	}
	return b.Result.Summarize()
}
