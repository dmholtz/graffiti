package io

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"

	g "github.com/dmholtz/graffiti/graph"
)

// fmi parse states
const (
	PARSE_NODE_COUNT = iota
	PARSE_EDGE_COUNT = iota
	PARSE_NODES      = iota
	PARSE_EDGES      = iota
)

// Build an AdjacencyListGraph from an .fmi file.
// nodeParseFnc parses a line of the .fmi file and returns a (nodeId, node) tuple
// edgeParseFnc parses a line of the .fmi file and returns a (nodeId, halfEdge) tuple
func NewAdjacencyListFromFmi[N any, E g.IHalfEdge](filename string, nodeParseFnc func(line string) (int, N), edgeParseFnc func(line string) (int, E)) *g.AdjacencyListGraph[N, E] {

	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	numNodes := 0
	numParsedNodes := 0

	alg := g.AdjacencyListGraph[N, E]{}
	id2index := make(map[int]int)

	parseState := PARSE_NODE_COUNT
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) < 1 {
			// skip empty lines
			continue
		} else if line[0] == '#' {
			// skip comments
			continue
		}

		switch parseState {
		case PARSE_NODE_COUNT:
			if val, err := strconv.Atoi(line); err == nil {
				numNodes = val
				parseState = PARSE_EDGE_COUNT
			}
		case PARSE_EDGE_COUNT:
			parseState = PARSE_NODES
		case PARSE_NODES:
			id, node := nodeParseFnc(line)
			id2index[id] = alg.NodeCount()
			alg.AppendNode(node)
			numParsedNodes++
			if numParsedNodes == numNodes {
				parseState = PARSE_EDGES
			}
		case PARSE_EDGES:
			from, edge := edgeParseFnc(line)
			alg.InsertHalfEdge(id2index[from], edge)
		}
	}

	if alg.NodeCount() != numNodes {
		// cannot check edge count because ocean.fmi contains duplicates, which are removed during import
		panic("Invalid parsing result")
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return &alg
}

// Serialize a graph into the .fmi format
// node2Fmi outputs a string description of a node in the .fmi format
// edge2Fmi outputs a string description of a half edge in the .fmi format
func WriteFmi[N any, E g.IHalfEdge](graph g.Graph[N, E], filename string, node2Fmi func(id g.NodeId, node N) string, edge2Fmi func(from g.NodeId, halfEdge E) string) {
	file, cErr := os.Create(filename)

	if cErr != nil {
		log.Fatal(cErr)
	}
	writer := bufio.NewWriter(file)

	// write number of nodes and number of edges
	writer.WriteString(fmt.Sprintf("%d\n", graph.NodeCount()))
	writer.WriteString(fmt.Sprintf("%d\n", graph.EdgeCount()))

	// list all nodes structured as "id lat lon" and append additional information such as partitions
	for i := 0; i < graph.NodeCount(); i++ {
		node := graph.GetNode(i)
		writer.WriteString(node2Fmi(i, node))
	}

	// list all edges structured as "fromId targetId distance" and append additional information such as arcflags
	for id := 0; id < graph.NodeCount(); id++ {
		for _, halfEdge := range graph.GetHalfEdgesFrom(id) {
			writer.WriteString(edge2Fmi(id, halfEdge))
		}
	}

	writer.Flush()
}

// Parsing functions

func ParseGeoPoint(line string) (int, g.GeoPoint) {
	var id int
	var lat, lon float64
	fmt.Sscanf(line, "%d %f %f", &id, &lat, &lon)
	return id, g.GeoPoint{Lon: lon, Lat: lat}
}

func ParsePartGeoPoint(line string) (int, g.PartGeoPoint) {
	var id int
	var lat, lon float64
	var part g.PartitionId
	fmt.Sscanf(line, "%d %f %f %d", &id, &lat, &lon, &part)
	return id, g.PartGeoPoint{GeoPoint: g.GeoPoint{Lon: lon, Lat: lat}, Partition_: part}
}

func Parse2LPartGeoPoint(line string) (int, g.TwoLevelPartGeoPoint) {
	var id int
	var lat, lon float64
	var l1Part, l2Part g.PartitionId
	fmt.Sscanf(line, "%d %f %f %d %d", &id, &lat, &lon, &l1Part, &l2Part)
	return id, g.TwoLevelPartGeoPoint{GeoPoint: g.GeoPoint{Lon: lon, Lat: lat}, L1Part_: l1Part, L2Part_: l2Part}
}

func ParseWeightedHalfEdge(line string) (int, g.WeightedHalfEdge[int]) {
	var from, to, weight int
	fmt.Sscanf(line, "%d %d %d", &from, &to, &weight)
	return from, g.WeightedHalfEdge[int]{To_: to, Weight_: weight}
}

func ParseFlaggedHalfEdge(line string) (int, g.FlaggedHalfEdge[int, uint64]) {
	var from, to, weight int
	var flag uint64
	fmt.Sscanf(line, "%d %d %d %d", &from, &to, &weight, &flag)
	return from, g.FlaggedHalfEdge[int, uint64]{To_: to, Weight_: weight, Flag: flag}
}

func ParseLargeFlaggedHalfEdge(line string) (int, g.LargeFlaggedHalfEdge[int]) {
	var from, to, weight int
	var msbFlag, lsbFlag uint64
	fmt.Sscanf(line, "%d %d %d %d %d", &from, &to, &weight, &msbFlag, &lsbFlag)
	return from, g.LargeFlaggedHalfEdge[int]{To_: to, Weight_: weight, MsbFlag: msbFlag, LsbFlag: lsbFlag}
}

func Parse2LFlaggedHalfEdge(line string) (int, g.TwoLevelFlaggedHalfEdge[int, uint64, uint64]) {
	var from, to, weight int
	var l1Flag, l2Flag uint64
	fmt.Sscanf(line, "%d %d %d %d %d", &from, &to, &weight, &l1Flag, &l2Flag)
	return from, g.TwoLevelFlaggedHalfEdge[int, uint64, uint64]{To_: to, Weight_: weight, L1Flag: l1Flag, L2Flag: l2Flag}
}

// Printer functions

func PartGeoPoint2FmiLine(id g.NodeId, node g.PartGeoPoint) string {
	return fmt.Sprintf("%d %f %f %d\n", id, node.Lat, node.Lon, node.Partition_)
}

func FlaggedHalfEdge2FmiLine(from g.NodeId, edge g.FlaggedHalfEdge[int, uint64]) string {
	return fmt.Sprintf("%d %d %d %d\n", from, edge.To_, edge.Weight_, edge.Flag)
}

func TwoLevelPartGeoPoint2FmiLine(id g.NodeId, node g.TwoLevelPartGeoPoint) string {
	return fmt.Sprintf("%d %f %f %d %d\n", id, node.Lat, node.Lon, node.L1Part_, node.L2Part_)
}

func TwoLevelFlaggedHalfEdge2FmiLine(from g.NodeId, edge g.TwoLevelFlaggedHalfEdge[int, uint64, uint64]) string {
	return fmt.Sprintf("%d %d %d %d %d\n", from, edge.To_, edge.Weight_, edge.L1Flag, edge.L2Flag)
}
