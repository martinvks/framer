package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/Martinvks/httptestrunner/types"
)

type TestCase struct {
	FileName string
	RequestData
}

type RequestData struct {
	AddDefaultHeaders bool
	Headers           types.Headers
	Continuation      types.Headers
	Trailer           types.Headers
	Body              string
}

func GetSingleTestCase(fileName string) (TestCase, error) {
	data, err := unmarshalTestCaseData(fileName)
	if err != nil {
		return TestCase{}, err
	}

	return TestCase{
		FileName:    fileName,
		RequestData: data,
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
			FileName:    entry.Name(),
			RequestData: data,
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

func unmarshalTestCaseData(fileName string) (RequestData, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return RequestData{}, err
	}

	contentWithEnv := replaceWithEnvironmentVariables(content)

	data := RequestData{
		AddDefaultHeaders: true,
	}
	err = json.Unmarshal(contentWithEnv, &data)

	if err != nil {
		return RequestData{}, err
	}

	return data, nil
}

// replaces all occurrences of ${ENVIRONMENT_VARIABLE_KEY} with ENVIRONMENT_VARIABLE_VALUE
func replaceWithEnvironmentVariables(fileContent []byte) []byte {
	re := regexp.MustCompile(`\$\{([^="]+)}`)

	return re.ReplaceAllFunc(fileContent, func(templateMatch []byte) []byte {
		parts := re.FindSubmatch(templateMatch)
		env := os.Getenv(string(parts[1]))
		return []byte(env)
	})
}
