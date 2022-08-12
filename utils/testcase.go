package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
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

	if len(dirEntries) == 0 {
		return nil, errors.New(fmt.Sprintf("Directory %v is empty", directory))
	}

	fmt.Println("Requests")
	for index, value := range dirEntries {
		fmt.Printf("%v. %v\n", index, value.Name())
	}

	fmt.Println("Choose request number:")
	var input string
	_, err = fmt.Scanln(&input)

	if err != nil {
		return nil, err
	}

	index, err := strconv.Atoi(input)
	if err != nil || (index < 0 || index >= len(dirEntries)) {
		return nil, errors.New(fmt.Sprintf("Invalid request number: %v", input))
	}

	return dirEntries[index], nil
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
