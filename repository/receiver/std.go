package receiver

import (
	"context"
	"encoding/json"
	"fmt"
	"graphs/entity"
	jsonentity "graphs/entity/json"
	"io"
	"os"
	"sync"
)

// Receive receives purchase Subscriptions.
func Receive(ctx context.Context, graph *entity.Graph) {
	for {
		select {
		// part of graceful shutdown. Do current and exit when receive context cancelled
		case <-ctx.Done():
			return
		default:
			var requestQuery = jsonentity.RequestQuery{}
			dec := json.NewDecoder(os.Stdin)
			fmt.Print("> ")

			err := dec.Decode(&requestQuery)
			if err == io.EOF {
				continue
			}
			if err != nil {
				fmt.Printf("%v\n", err)
				continue
			}

			answer := GetAnswer(graph, &requestQuery)

			janswer, err := json.MarshalIndent(answer, "", " ")
			if err != nil {
				fmt.Printf("%v\n", err)
				continue
			}
			// print answer
			fmt.Println(string(janswer))
		}
	}
}

func GetAnswer(graph *entity.Graph, query *jsonentity.RequestQuery) *jsonentity.Answer {
	var (
		res = jsonentity.Answer{Answers: make([]map[string]jsonentity.PathResponse, 0, len(query.Queries))}
		wg  sync.WaitGroup
		//make buffered channel for wait response in concurrency
		ch = make(chan map[string]jsonentity.PathResponse, len(query.Queries))
		pw = 0
	)

	if len(query.Queries) == 0 {
		return &res
	}

	for _, q := range query.Queries {
		if q.Cheapest != nil && q.Cheapest.Start != "" && q.Cheapest.End != "" {
			wg.Add(1)
			pw++
			// make goroutine with result in channel for concurrently search path in graph
			go func(resultChannel chan<- map[string]jsonentity.PathResponse, start, end string) {
				defer wg.Done()
				r := jsonentity.PathResponse{From: start, To: end, Path: false}

				pa := graph.GetCheapestPaths(start, end)

				if len(pa) > 0 {
					r.Path = pa
				}

				resultChannel <- map[string]jsonentity.PathResponse{"cheapest": r}
			}(ch, q.Cheapest.Start, q.Cheapest.End)
		}

		if q.Paths != nil && q.Paths.Start != "" && q.Paths.End != "" {
			wg.Add(1)
			pw++
			// make goroutine with result in channel for concurrently search path in graph
			go func(resultChannel chan<- map[string]jsonentity.PathResponse, start, end string) {
				defer wg.Done()
				r := jsonentity.PathResponse{From: start, To: end, Paths: make([]string, 0)}

				pa := graph.GetPaths(start, end)
				if len(pa) > 0 {
					r.Paths = pa
				}

				resultChannel <- map[string]jsonentity.PathResponse{"paths": r}
			}(ch, q.Paths.Start, q.Paths.End)
		}
	}

	wg.Wait()

	for i := 0; i < pw; i++ {
		res.Answers = append(res.Answers, <-ch)
	}

	return &res
}

func GetAnswerIterate(graph *entity.Graph, query *jsonentity.RequestQuery) *jsonentity.Answer {

	var (
		res = jsonentity.Answer{Answers: make([]map[string]jsonentity.PathResponse, 0, len(query.Queries))}
	)

	if len(query.Queries) == 0 {
		return &res
	}

	for _, q := range query.Queries {
		if q.Cheapest != nil && q.Cheapest.Start != "" && q.Cheapest.End != "" {
			r := jsonentity.PathResponse{From: q.Cheapest.Start, To: q.Cheapest.End, Path: false}

			pa := graph.GetCheapestPaths(q.Cheapest.Start, q.Cheapest.End)

			if len(pa) > 0 {
				r.Path = pa
			}

			res.Answers = append(res.Answers, map[string]jsonentity.PathResponse{"cheapest": r})
		}

		if q.Paths != nil && q.Paths.Start != "" && q.Paths.End != "" {
			r := jsonentity.PathResponse{From: q.Paths.Start, To: q.Paths.End, Paths: make([]string, 0)}

			pa := graph.GetPaths(q.Paths.Start, q.Paths.End)
			if len(pa) > 0 {
				r.Paths = pa
			}

			res.Answers = append(res.Answers, map[string]jsonentity.PathResponse{"paths": r})
		}
	}

	return &res
}
