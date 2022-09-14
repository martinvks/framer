package utils

import (
	"encoding/csv"
	"io"
	"os"
)

type CsvWriter struct {
	writer *csv.Writer
}

func GetCsvWriter(filename string) (CsvWriter, error) {
	var writer io.Writer

	if filename == "" {
		writer = os.Stdout
	} else {
		f, err := os.Create(filename)
		if err != nil {
			return CsvWriter{}, err
		}
		writer = f
	}

	return CsvWriter{
		writer: csv.NewWriter(writer),
	}, nil
}

func (w CsvWriter) WriteHeaders() error {
	return w.writer.Write([]string{
		"filename",
		"request_id",
		"response_code",
		"error",
	})
}

func (w CsvWriter) WriteData(
	testFilename string,
	requestId string,
	responseCode string,
	error string,
) error {
	return w.writer.Write([]string{
		testFilename,
		requestId,
		responseCode,
		error,
	})
}

func (w CsvWriter) Close() {
	w.writer.Flush()
}
