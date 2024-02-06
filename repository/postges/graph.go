package postges

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"graphs/entity/postgre"
)

type GraphRepo struct {
	db *sqlx.DB
}

func NewGraphRepo(db *sqlx.DB) *GraphRepo {
	return &GraphRepo{
		db: db,
	}
}

// UpsertGraph - Rewrite graph
func (g *GraphRepo) UpsertGraph(ctx context.Context, graph *postgre.Graph) error {
	err := g.runInTransaction(ctx, func(tx *sqlx.Tx) error {
		// Rewrite graph tables in transaction
		err := truncateTable(ctx, tx, "edges")
		if err != nil {
			return fmt.Errorf("table edges: %w", err)
		}

		err = truncateTable(ctx, tx, "nodes")
		if err != nil {
			return fmt.Errorf("table nodes: %w", err)
		}

		err = truncateTable(ctx, tx, "graphs")
		if err != nil {
			return fmt.Errorf("table graphs: %w", err)
		}

		err = g.InsertGraph(ctx, tx, graph)
		if err != nil {
			return err
		}

		err = g.InsertNodes(ctx, tx, graph.Nodes)
		if err != nil {
			return err
		}

		err = g.InsertEdges(ctx, tx, graph.Edges)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

// InsertGraph ...
func (g *GraphRepo) InsertGraph(ctx context.Context, tx *sqlx.Tx, graph *postgre.Graph) error {
	q := `INSERT INTO graphs (id, name) VALUES (:id, :name)`
	_, err := tx.NamedExecContext(ctx, q, graph)
	if err != nil {
		return fmt.Errorf("failed to insert graph: %w", err)
	}

	return nil
}

// InsertNodes ...
func (g *GraphRepo) InsertNodes(ctx context.Context, tx *sqlx.Tx, nodes []postgre.Node) error {
	q := `INSERT INTO nodes (id, name, graph_id) VALUES (:id, :name, :graph_id)`
	_, err := tx.NamedExecContext(ctx, q, nodes)
	if err != nil {
		return fmt.Errorf("failed to insert nodes: %w", err)
	}

	return nil
}

// InsertEdges ...
func (g *GraphRepo) InsertEdges(ctx context.Context, tx *sqlx.Tx, edges []postgre.Edge) error {
	q := `INSERT INTO edges (id, previous_node, next_node, cost) VALUES (:id, :previous_node, :next_node, :cost)`
	_, err := tx.NamedExecContext(ctx, q, edges)
	if err != nil {
		return fmt.Errorf("failed to insert edges: %w", err)
	}

	return nil
}

func (g *GraphRepo) GetGraphCycle(ctx context.Context) ([]string, error) {
	var res = make([]string, 0)

	query := `WITH RECURSIVE search_graph(previous_node, next_node, id, depth, path, cycle)
                     AS (
          SELECT e.previous_node, e.next_node, e.id, 1,
                 ARRAY[e.id]::varchar[],
                 false
          FROM edges e
          UNION ALL
          SELECT e.previous_node, e.next_node, e.id, sg.depth + 1,
                 path || e.id,
                 e.id = ANY(path)
          FROM edges e, search_graph sg
          WHERE e.previous_node = sg.next_node AND NOT cycle
            and depth < 100
      )
	  SELECT path
	  FROM search_graph
	  WHERE cycle
	  limit 1;`

	// Execute the query
	err := g.db.GetContext(ctx, pq.Array(&res), query)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("failed to GetGraphCycle: %w", err)
	}

	return res, nil

}

func (g *GraphRepo) GetEdges(ctx context.Context) ([]postgre.Edge, error) {
	var res = make([]postgre.Edge, 0)

	query := `select id, previous_node, next_node, cost
		from edges;`
	// Execute the query
	rows, err := g.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to GetEdges: %w", err)
	}

	for rows.Next() {
		var edge postgre.Edge
		err := rows.Scan(&edge.ID,
			&edge.PreviousNode,
			&edge.NextNode,
			&edge.Cost,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, edge)
	}

	return res, nil
}
func (g *GraphRepo) GetNodes(ctx context.Context) ([]postgre.Node, error) {
	var res = make([]postgre.Node, 0)

	query := `select id,name, graph_id
		from nodes;`
	// Execute the query
	rows, err := g.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to GetNodes: %w", err)
	}

	for rows.Next() {
		var node postgre.Node
		err := rows.Scan(&node.ID,
			&node.Name,
			&node.GraphID,
		)
		if err != nil {
			return nil, err
		}

		res = append(res, node)
	}

	return res, nil
}

func (g *GraphRepo) GetGraph(ctx context.Context) (*postgre.Graph, error) {
	var (
		res postgre.Graph
	)

	query := `select id,name
		from graphs;`
	// Execute the query
	err := g.db.GetContext(ctx, &res, query)
	if err != nil {
		return nil, fmt.Errorf("failed to GetGraph: %w", err)
	}

	res.Nodes, err = g.GetNodes(ctx)
	if err != nil {
		return nil, err
	}

	res.Edges, err = g.GetEdges(ctx)
	if err != nil {
		return nil, err
	}

	return &res, nil
}

func (g *GraphRepo) runInTransaction(ctx context.Context, exec func(tx *sqlx.Tx) error) error {
	tx, err := g.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if err := exec(tx); err != nil {
		rerr := tx.Rollback()
		if rerr != nil {
			return fmt.Errorf("transaction rollback error %w", rerr)
		}

		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("transaction commit error %w", err)
	}
	return nil

}

func truncateTable(ctx context.Context, tx *sqlx.Tx, tableName string) error {
	query := fmt.Sprintf("TRUNCATE TABLE %s CASCADE", tableName)

	// Execute the query
	_, err := tx.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to truncate table %s: %v", tableName, err)
	}

	return nil
}
