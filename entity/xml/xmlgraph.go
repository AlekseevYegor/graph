package xml

import (
	"encoding/xml"
	"fmt"
)

type (
	Graph struct {
		XMLName xml.Name `xml:"graph"`
		ID      string   `xml:"id"`
		Name    string   `xml:"name"`
		Nodes   Nodes    `xml:"nodes"`
		Edges   Edges    `xml:"edges"`
	}

	Nodes struct {
		XMLName xml.Name `xml:"nodes"`
		Nodes   []Node   `xml:"node"`
	}

	Node struct {
		XMLName xml.Name `xml:"node"`
		ID      string   `xml:"id"`
		Name    string   `xml:"name"`
	}

	Edges struct {
		XMLName xml.Name `xml:"edges"`
		Edges   []Edge   `xml:"node"`
	}

	Edge struct {
		XMLName xml.Name `xml:"node"`
		ID      string   `xml:"id"`
		From    string   `xml:"from"`
		To      string   `xml:"to"`
		Cost    float64  `xml:"cost"`
	}
)

func (g *Graph) Validate() error {
	if g.ID == "" || g.Name == "" {
		return fmt.Errorf("graph must have both <id> and <name>")
	}

	// Validate at least one <node> in the <nodes> group
	if len(g.Nodes.Nodes) == 0 {
		return fmt.Errorf("at least one <node> must be present in the <nodes> group")
	}

	// Validate unique <id> tags for nodes
	nodeIDs := make(map[string]bool)
	for _, node := range g.Nodes.Nodes {
		if nodeIDs[node.ID] {
			return fmt.Errorf("duplicate <id> tags for nodes are not allowed")
		}
		nodeIDs[node.ID] = true
	}

	for _, edge := range g.Edges.Edges {
		// Validate <from> and <to> tags in edges correspond to defined nodes
		if !nodeIDs[edge.From] || !nodeIDs[edge.To] {
			return fmt.Errorf("undefined nodes in <from> or <to> tags in edges")
		}

		// Validate <from> and <to> tags in edges not equals
		if edge.From == edge.To {
			return fmt.Errorf(fmt.Sprintf("edge id: %s pointed to itself", edge.ID))
		}

		// Validate cost must be greater than 0
		if edge.Cost < 0 {
			return fmt.Errorf(fmt.Sprintf("cost must be greather than 0 for edge id: %s", edge.ID))
		}
	}

	return nil
}
