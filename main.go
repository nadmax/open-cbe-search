package main

import (
	"bytes"
	"io"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/elastic/go-elasticsearch/v9"
)

const (
	bulkSize = 500
	dataDir = "data"
)

func readBatch(reader *csv.Reader, size int) ([][]string, error) {
	var batch [][]string
	for len(batch) < size {
		row, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				return batch, io.EOF
			}
			return batch, err
		}
		batch = append(batch, row)
	}
	return batch, nil
}

func bulkIndexCSV(es *elasticsearch.Client, filePath, indexName string) error {
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

	var (
		bulkBuffer bytes.Buffer
		count int
		batchCount int
	)

		for {
		rows, err := readBatch(reader, bulkSize)
		if err != nil && err != io.EOF {
			return fmt.Errorf("read batch error: %w", err)
		}

		for _, row := range rows {
			doc := make(map[string]string)
			for j, val := range row {
				doc[headers[j]] = val
			}

			meta := map[string]map[string]string{
				"index": {"_index": indexName},
			}
			metaLine, _ := json.Marshal(meta)
			docLine, _ := json.Marshal(doc)

			bulkBuffer.Write(metaLine)
			bulkBuffer.WriteByte('\n')
			bulkBuffer.Write(docLine)
			bulkBuffer.WriteByte('\n')

			count++
		}

		if bulkBuffer.Len() > 0 {
			res, err := es.Bulk(bytes.NewReader(bulkBuffer.Bytes()))
			if err != nil {
				return fmt.Errorf("bulk request error: %w", err)
			}
			defer res.Body.Close()

			if res.IsError() {
				body, _ := io.ReadAll(res.Body)
				return fmt.Errorf("bulk response error: %s", body)
			}

			bulkBuffer.Reset()
			batchCount++
			fmt.Printf("Batch %d indexed (%d records)\n", batchCount, len(rows))
		}

		if err == io.EOF {
			break
		}
	}

	fmt.Printf("Finished indexing %d documents into '%s'\n", count, indexName)
	return nil
}

func main() {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}

	res, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	log.Println(res)

	err = filepath.WalkDir(dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() && strings.HasSuffix(d.Name(), ".csv") {
			indexName := strings.TrimSuffix(d.Name(), ".csv")
			fmt.Printf("Indexing %s into index '%s'\n", d.Name(), indexName)
			if err := bulkIndexCSV(es, path, indexName); err != nil {
				log.Printf("Error indexing %s: %s", d.Name(), err)
			}
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Error walking data directory: %s", err)
	}
}
