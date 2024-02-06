package json

type (
	PathQuery struct {
		Start string `json:"start"`
		End   string `json:"end"`
	}

	Query struct {
		Paths    *PathQuery `json:"paths,omitempty"`
		Cheapest *PathQuery `json:"cheapest,omitempty"`
	}

	RequestQuery struct {
		Queries []Query
	}

	PathResponse struct {
		From  string      `json:"from"`
		To    string      `json:"to"`
		Paths interface{} `json:"paths,omitempty"` // [ [ "a", "b", "e" ], [ "a", "e" ] ]
		Path  interface{} `json:"path,omitempty"`  // [ "a", "e" ] - false

	}

	Answer struct {
		Answers []map[string]PathResponse `json:"answers"`
	}
)
