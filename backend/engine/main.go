package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/nadmax/open-cbe-search/core/engine/indexer"
	"github.com/nadmax/open-cbe-search/core/engine/postgres"
)

const dataDir = "data"

func main() {
	conn := os.Getenv("POSTGRES_URL")
	if conn == "" {
		log.Fatalf("No Postgres URL defined")
	}

	c, err := postgres.NewClient(conn)
	if err != nil {
		log.Fatalf("Error connecting to Postgres: %s", err)
	}
	defer c.Close()

	start := time.Now()
	err = filepath.WalkDir(dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(d.Name(), ".csv") {
			tableName := strings.TrimSuffix(d.Name(), ".csv")
			fmt.Printf("Inserting %s into table '%s'\n", d.Name(), tableName)
			if err := indexer.BulkInsertCSV(c.DB, path, tableName); err != nil {
				log.Printf("Error inserting %s: %s", d.Name(), err)
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking data directory: %s", err)
	}

	fmt.Printf("Finished indexing datas in %s^n", time.Since(start))
}