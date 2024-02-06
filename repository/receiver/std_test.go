package receiver

import (
	"graphs/entity"
	jsonentity "graphs/entity/json"
	"testing"
)

func BenchmarkGetAnswer(b *testing.B) {

	graph := entity.Graph{
		AdjacencyList: map[string][]entity.Edge{
			"a": {{Next: "e", Cost: 42}, {Next: "b", Cost: 10}},
			"e": {{Next: "c", Cost: 3}},
			"c": {{Next: "a", Cost: 42}, {Next: "d", Cost: 5}},
			"b": {{Next: "d", Cost: 20}, {Next: "f", Cost: 10}},
			"f": {{Next: "i", Cost: 10}},
			"i": {{Next: "h", Cost: 10}},
			"h": {{Next: "g", Cost: 10}},
			"d": {{Next: "g", Cost: 10}},
		},
	}
	query := jsonentity.RequestQuery{
		Queries: []jsonentity.Query{
			{Paths: &jsonentity.PathQuery{Start: "a", End: "e"}},
			{Paths: &jsonentity.PathQuery{Start: "a", End: "f"}},
			{Paths: &jsonentity.PathQuery{Start: "a", End: "d"}},
			{Paths: &jsonentity.PathQuery{Start: "a", End: "g"}},
			{Paths: &jsonentity.PathQuery{Start: "f", End: "g"}},
			{Paths: &jsonentity.PathQuery{Start: "b", End: "g"}},
			{Cheapest: &jsonentity.PathQuery{Start: "a", End: "d"}},
			{Cheapest: &jsonentity.PathQuery{Start: "a", End: "g"}},
			{Cheapest: &jsonentity.PathQuery{Start: "e", End: "g"}},
		},
	}

	b.ResetTimer()
	GetAnswer(&graph, &query)
}

func BenchmarkGetAnswerIterate(b *testing.B) {
	graph := entity.Graph{
		AdjacencyList: map[string][]entity.Edge{
			"a": {{Next: "e", Cost: 42}, {Next: "b", Cost: 10}},
			"e": {{Next: "c", Cost: 3}},
			"c": {{Next: "a", Cost: 42}, {Next: "d", Cost: 5}},
			"b": {{Next: "d", Cost: 20}, {Next: "f", Cost: 10}},
			"f": {{Next: "i", Cost: 10}},
			"i": {{Next: "h", Cost: 10}},
			"h": {{Next: "g", Cost: 10}},
			"d": {{Next: "g", Cost: 10}},
		},
	}
	query := jsonentity.RequestQuery{
		Queries: []jsonentity.Query{
			{Paths: &jsonentity.PathQuery{Start: "a", End: "e"}},
			{Paths: &jsonentity.PathQuery{Start: "a", End: "f"}},
			{Paths: &jsonentity.PathQuery{Start: "a", End: "d"}},
			{Paths: &jsonentity.PathQuery{Start: "a", End: "g"}},
			{Paths: &jsonentity.PathQuery{Start: "f", End: "g"}},
			{Paths: &jsonentity.PathQuery{Start: "b", End: "g"}},
			{Cheapest: &jsonentity.PathQuery{Start: "a", End: "d"}},
			{Cheapest: &jsonentity.PathQuery{Start: "a", End: "g"}},
			{Cheapest: &jsonentity.PathQuery{Start: "e", End: "g"}},
		},
	}

	b.ResetTimer()
	GetAnswerIterate(&graph, &query)
}
