package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/lib/pq"
)

type Client struct {
	DB *sql.DB
}

func NewClient(conn string) (*Client, error) {
	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("could not connect to Postgres: %w", err)
	}

	return &Client{ DB: db }, nil
}

func (c *Client) Close() error {
	return c.DB.Close()
}

func (c *Client) SearchTable(table, queryArgs string) ([]map[string]interface{}, error) {
	cols, err := getColumnNames(c.DB, table)
	if err != nil {
		return nil, err
	}

	concatExpr := fmt.Sprintf("concat_ws(' ', %s)", strings.Join(cols, ", "))
	query := fmt.Sprintf(`
		SELECT * FROM %s
		WHERE to_tsvector('english', %s) @@ plainto_tsquery($1)
	`, pq.QuoteIdentifier(table), concatExpr)
	rows, err := c.DB.Query(query, queryArgs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err = rows.Columns()
	if err != nil {
		return nil, err
	}

	results := []map[string]interface{}{}
	for rows.Next() {
		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}
		if err := rows.Scan(columnPointers...); err != nil {
			return nil, err
		}

		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			m[colName] = *val
		}

		results = append(results, m)
	}

	return results, nil
}

func getColumnNames(db *sql.DB, table string) ([]string, error) {
	rows, err := db.Query(`
		SELECT column_name
		FROM information_schema.columns
		WHERE table_name = $1
	`, table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []string
	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {
			return nil, err
		}
		cols = append(cols, fmt.Sprintf(`COALESCE(%s, '')`, pq.QuoteIdentifier(col))) // safely quote
	}
	return cols, nil
}