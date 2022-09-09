package utils

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"os"
	"strings"
)

type TestCase struct {
	Id       string
	FileName string
	Data     TestCaseData
}

type TestCaseData struct {
	Method  string
	Headers []Header
	Body    string
}

func GetSingleTestCase(fileName string) (TestCase, error) {
	data, err := unmarshalTestCaseData(fileName)
	if err != nil {
		return TestCase{}, err
	}

	return TestCase{
		Id:       uuid.NewString(),
		FileName: fileName,
		Data:     data,
	}, nil
}

func GetAllTestCases(directory string) ([]TestCase, error) {
	jsonEntries, err := getJsonEntries(directory)
	if err != nil {
		return nil, err
	}

	var testCases []TestCase
	for _, entry := range jsonEntries {
		data, err := unmarshalTestCaseData(directory + "/" + entry.Name())
		if err != nil {
			return nil, err
		}
		testCases = append(testCases, TestCase{
			Id:       uuid.NewString(),
			FileName: entry.Name(),
			Data:     data,
		})
	}

	return testCases, nil
}

func getJsonEntries(directory string) ([]os.DirEntry, error) {
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
		return nil, fmt.Errorf("directory %v does not contain any json request files", directory)
	}

	return jsonEntries, nil
}

func unmarshalTestCaseData(fileName string) (TestCaseData, error) {
	content, err := os.ReadFile(fileName)

	if err != nil {
		return TestCaseData{}, err
	}

	data := TestCaseData{}
	err = json.Unmarshal(content, &data)

	if err != nil {
		return TestCaseData{}, err
	}

	return data, nil
}
