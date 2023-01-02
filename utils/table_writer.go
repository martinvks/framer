package utils

import (
	"errors"
	"fmt"
)

const padding = 2

func WriteTable(tableHeaders []string, tableData [][]string) error {
	maxLengths, err := getMaxLengths(tableHeaders, tableData)
	if err != nil {
		return err
	}

	table := [][]string{tableHeaders}
	table = append(table, tableData...)

	for _, row := range table {
		for index := range row {
			fmt.Printf("%-*s", maxLengths[index]+padding, row[index])
		}
		fmt.Println()
	}

	return nil
}

func getMaxLengths(tableHeaders []string, tableData [][]string) ([]int, error) {
	var maxLengths = make([]int, len(tableHeaders))

	for i := range tableHeaders {
		maxLengths[i] = len(tableHeaders[i])
	}

	for _, data := range tableData {
		if len(data) != len(tableHeaders) {
			return nil, errors.New("data must be same size as headers")
		}

		for index := range data {
			length := len(data[index])
			if length > maxLengths[index] {
				maxLengths[index] = length
			}
		}
	}

	return maxLengths, nil
}
