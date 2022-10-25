package utils

import (
	"encoding/csv"
	"io"
	"os"
)

type CsvWriter struct {
	addIdField bool
	writer     *csv.Writer
}

func GetCsvWriter(filename string, addIdField bool) (CsvWriter, error) {
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
		addIdField: addIdField,
		writer:     csv.NewWriter(writer),
	}, nil
}

func (w CsvWriter) WriteHeaders() error {
	var headers []string
	if w.addIdField {
		headers = append(headers, "request_id")
	}

	headers = append(
		headers,
		"filename",
		"response_code",
		"error",
	)

	return w.writer.Write(headers)
}

func (w CsvWriter) WriteData(
	testFilename string,
	requestId string,
	responseCode string,
	error string,
) error {
	var row []string
	if w.addIdField {
		row = append(row, requestId)
	}

	row = append(
		row,
		testFilename,
		responseCode,
		error,
	)

	return w.writer.Write(row)
}

func (w CsvWriter) Close() {
	w.writer.Flush()
}
