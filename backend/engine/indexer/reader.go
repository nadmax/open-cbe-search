package indexer

import (
	"encoding/csv"
	"io"
)

func ReadBatch(reader *csv.Reader, size int) ([][]string, error) {
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
