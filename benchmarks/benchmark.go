package benchmarks

// Summary summarizes the result of a benchmark.
type Summary struct {
	Runs   int     // number of executions
	Time   float64 // average execution time [ms]
	PqPops int     // average number of Pop() operations on priority queue
}

// BenchmarkResult reports details of a benchmark, i.e. the distribution of runtimes and pq-Pops.
type BenchmarkResult struct {
	TimeDistribution   []float64 `json:"times"`
	PqPopsDistribution []int     `json:"pq-pops"`
}

func NewBenchmarkResult() *BenchmarkResult {
	return &BenchmarkResult{TimeDistribution: make([]float64, 0), PqPopsDistribution: make([]int, 0)}
}

// Add adds a new observation to the benchmark.
func (br *BenchmarkResult) Add(time float64, pqPops int) {
	br.TimeDistribution = append(br.TimeDistribution, time)
	br.PqPopsDistribution = append(br.PqPopsDistribution, pqPops)
}

// Summarize builds a summary of the benchmark.
func (br BenchmarkResult) Summarize() Summary {
	runs := len(br.TimeDistribution)
	time := mean(br.TimeDistribution)
	pqPops := mean(br.PqPopsDistribution)
	return Summary{Runs: runs, Time: time, PqPops: pqPops}
}

// mean computest the mean of a slice.
func mean[N int | float64](slice []N) N {
	if len(slice) == 0 {
		return N(0)
	} else {
		cum := N(0)
		for _, sample := range slice {
			cum += sample
		}
		return cum / N(len(slice))
	}
}
