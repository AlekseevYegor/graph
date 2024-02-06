package postgre

import (
	"graphs/entity/xml"
)

type (
	Graph struct {
		ID    string `db:"id"`
		Name  string `db:"name"`
		Nodes []Node
		Edges []Edge
	}

	Node struct {
		ID      string `db:"id"`
		Name    string `db:"name"`
		GraphID string `db:"graph_id"`
	}

	Edge struct {
		ID           string  `db:"id"`
		PreviousNode string  `db:"previous_node"`
		NextNode     string  `db:"next_node"`
		Cost         float64 `db:"cost"`
	}
)

func NewGraph(graph xml.Graph) *Graph {
	var (
		nodes = make([]Node, 0, len(graph.Nodes.Nodes))
		edges = make([]Edge, 0, len(graph.Edges.Edges))
	)

	for _, node := range graph.Nodes.Nodes {
		nodes = append(nodes, makeNode(node, graph.ID))
	}

	for _, edge := range graph.Edges.Edges {
		edges = append(edges, makeEdge(edge))
	}

	return &Graph{
		ID:    graph.ID,
		Name:  graph.Name,
		Nodes: nodes,
		Edges: edges,
	}
}

func makeNode(node xml.Node, graphID string) Node {
	return Node{
		ID:      node.ID,
		Name:    node.Name,
		GraphID: graphID,
	}
}

func makeEdge(edge xml.Edge) Edge {
	return Edge{
		ID:           edge.ID,
		PreviousNode: edge.From,
		NextNode:     edge.To,
		Cost:         edge.Cost,
	}
}
