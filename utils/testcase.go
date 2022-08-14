package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type TestCase struct {
	Method  string
	Headers []Header
	Body    string
}

func GetTestCase(directory string) (TestCase, error) {
	dirEntry, err := getDirEntry(directory)

	if err != nil {
		return TestCase{}, err
	}

	testCase, err := unmarshalTestCase(directory, dirEntry)

	if err != nil {
		return TestCase{}, err
	}

	return testCase, nil
}

func getDirEntry(directory string) (os.DirEntry, error) {
	dirEntries, err := os.ReadDir(directory)

	if err != nil {
		return nil, err
	}

	var jsonEntries []os.DirEntry
	for _, entry := range dirEntries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".json") {
			jsonEntries = append(jsonEntries, entry)
		}
	}

	if len(jsonEntries) == 0 {
		return nil, errors.New(fmt.Sprintf("Directory %v does not contain any json request files", directory))
	}

	fmt.Println("Requests")
	for index, entry := range jsonEntries {
		fmt.Printf("%v. %v\n", index, entry.Name())
	}

	fmt.Println("Choose request number:")
	var input string
	_, err = fmt.Scanln(&input)
	fmt.Println()

	if err != nil {
		return nil, err
	}

	index, err := strconv.Atoi(input)
	if err != nil || (index < 0 || index >= len(jsonEntries)) {
		return nil, errors.New(fmt.Sprintf("Invalid request number: %v", input))
	}

	return jsonEntries[index], nil
}

func unmarshalTestCase(directory string, dirEntry os.DirEntry) (TestCase, error) {
	content, err := os.ReadFile(directory + "/" + dirEntry.Name())

	if err != nil {
		return TestCase{}, err
	}

	request := TestCase{
		Method: "GET",
	}
	err = json.Unmarshal(content, &request)

	if err != nil {
		return TestCase{}, err
	}

	return request, nil
}
