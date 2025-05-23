package indexer

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"
)

const bulkSize = 5000

func BulkInsertCSV(db *sql.DB, filePath, tableName string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	headers, err := reader.Read()
	if err != nil {
		return fmt.Errorf("could not read headers: %w", err)
	}

	var columns []string
	for _, h := range headers {
		col := fmt.Sprintf("%q TEXT", h)
		columns = append(columns, col)
	}
	createStmt := fmt.Sprintf(`CREATE TABLE IF NOT EXISTS "%s" (%s);`, tableName, strings.Join(columns, ", "))
	if _, err := db.Exec(createStmt); err != nil {
		return fmt.Errorf("table creation failed: %w", err)
	}

	total := 0
	batchCount := 0

	for {
		rows, err := ReadBatch(reader, bulkSize)
		if err != nil && err != io.EOF {
			return fmt.Errorf("read batch error: %w", err)
		}

		if len(rows) == 0 {
			break
		}

		valueArgs := []interface{}{}
		var rowsSQL []string

		for _, row := range rows {
			rowVals := []string{}
			for _, val := range row {
				valueArgs = append(valueArgs, val)
				rowVals = append(rowVals, fmt.Sprintf("$%d", len(valueArgs)))
			}
			rowsSQL = append(rowsSQL, fmt.Sprintf("(%s)", strings.Join(rowVals, ",")))
		}

		stmt := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES %s;`,
			tableName,
			strings.Join(quoteIdentifiers(headers), ","),
			strings.Join(rowsSQL, ","))
		if _, err := db.Exec(stmt, valueArgs...); err != nil {
			return fmt.Errorf("batch insert failed: %w", err)
		}

		total += len(rows)
		batchCount++
		fmt.Printf("Batch %d inserted (%d records)\n", batchCount, len(rows))

		if err == io.EOF {
			break
		}
	}

	fmt.Printf("Finished inserting %d records into table '%s'\n", total, tableName)	

	return nil
}

func quoteIdentifiers(ids []string) []string {
	var out []string
	for _, id := range ids {
		out = append(out, fmt.Sprintf(`"%s"`, id))
	}

	return out
}