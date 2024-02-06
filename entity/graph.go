package entity

import (
	"graphs/entity/postgre"
	"sort"
)

type (
	Edge struct {
		Next string
		Cost float64
	}

	Graph struct {
		AdjacencyList map[string][]Edge
	}

	PathsCost struct {
		Path      []string
		TotalCost float64
	}
)

func NewGraph(graphDB postgre.Graph) *Graph {
	graph := Graph{AdjacencyList: make(map[string][]Edge, len(graphDB.Nodes))}

	for _, n := range graphDB.Nodes {
		graph.AdjacencyList[n.ID] = nil
	}

	for _, e := range graphDB.Edges {
		f, ok := graph.AdjacencyList[e.PreviousNode]
		if ok {
			if f == nil {
				f = []Edge{{Next: e.NextNode, Cost: e.Cost}}
			} else {
				f = append(f, Edge{Next: e.NextNode, Cost: e.Cost})
			}

			graph.AdjacencyList[e.PreviousNode] = f
		}
	}

	return &graph
}

func (g Graph) GetPaths(start, end string) [][]string {
	var (
		cost, totalCost float64
		visited         = make(map[string]int)
		path            = make([]string, 0)
		allPaths        = make([]PathsCost, 0)
		response        = make([][]string, 0)
	)

	g.dfsAllPathsWithCost(start, end, visited, path, cost, totalCost, &allPaths)

	for _, p := range allPaths {
		response = append(response, p.Path)
	}

	return response
}

func (g Graph) GetCheapestPaths(start, end string) []string {
	var (
		cost, totalCost float64
		visited         = make(map[string]int)
		path            = make([]string, 0)
		allPaths        = make([]PathsCost, 0)
	)

	g.dfsAllPathsWithCost(start, end, visited, path, cost, totalCost, &allPaths)

	// If it has no path return nil
	if len(allPaths) == 0 {
		return nil
	}

	// Sort slice from min to max TotalCost
	sort.Slice(allPaths, func(i, j int) bool {
		return allPaths[i].TotalCost < allPaths[j].TotalCost
	})

	// return minimum coast path
	return allPaths[0].Path
}

func (g Graph) dfsAllPathsWithCost(current, finish string, visited map[string]int, path []string, cost, totalCost float64, allPaths *[]PathsCost) {
	visited[current] = 1

	path = append(path, current)
	totalCost += cost

	// If the current node is the target, return the path
	if current == finish {
		*allPaths = append(*allPaths, PathsCost{Path: path, TotalCost: totalCost})
	} else {
		for _, next := range g.AdjacencyList[current] {
			if visited[next.Next] != 1 {
				g.dfsAllPathsWithCost(next.Next, finish, visited, path, next.Cost, totalCost, allPaths)
			} else if visited[next.Next] == 1 {
				// path has cycle skip the path
				continue
			}
		}
	}

	visited[current] = 2
}
