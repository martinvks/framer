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

	for index := range tableHeaders {
		fmt.Printf("%-*s", maxLengths[index]+padding, tableHeaders[index])
	}
	fmt.Println()

	for _, data := range tableData {
		for index := range data {
			fmt.Printf("%-*s", maxLengths[index]+padding, data[index])
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
