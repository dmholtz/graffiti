# graffiti

Graffiti is generic graph library written in Go.

## Prerequisites

- An installation of Go 1.18 or later (graffiti uses generics)

## Features

### Data structures

The `Graph` interface provides an abstract capability description of a graph:

```go
type Graph[N any, E IHalfEdge] interface {
    NodeCount() int
    EdgeCount() int
    GetNode(id NodeId) N
    GetHalfEdgesFrom(id NodeId) []E
}
```

There are currently two implementations of the `Graph` interface:

- Adjacency List
- Adjacency Array

### Shortest path algorithms

Shortest path algorithms aim at finding the shortest path between a source and a target node in a weighted graph.
Graffiti implements [Dijkstra's algorithm](https://en.wikipedia.org/wiki/Dijkstra%27s_algorithm), which serves as a baseline.
Beside that, the following speed-up techniques as well as the required preprocessing steps or heuristics (if applicable) are implemented:

- Bidirectional Dijkstra's algorithm
- A\* search algorithm
  - Haversine Heuristic
  - ALT (A\*, landmarks and triangular inequalities)
- Bidirectional A\* search algorithm
- Dijkstra's algorithm with arc flags
- Bidirectional Dijkstra's algorithm with arc flags
- Dijkstra's algorithm with two-level arc flags

## Demo

The [osm-ship-routing repository](https://github.com/dmholtz/osm-ship-routing) features a REST-API for global ship navigation.
The underlying graph is represented by graffiti's [AdjacencyArrayGraph](graph/adjacency_array.go) data structure and the graph search calls graffiti's [BidirectionalArcFlagRouter](algorithms/shortest_path/arc_flag_bi_dijkstra.go).

## License

Graffiti is licensed under the [MIT License](LICENSE).

```yaml
SPDX-License-Identifier: MIT
```
