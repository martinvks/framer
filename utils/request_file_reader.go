package utils

import (
	"fmt"
	"os"
	"strings"

	"github.com/martinvks/framer/types"
)

type RequestFile struct {
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

func GetRequestFile(fileName string) (RequestFile, error) {
	data, err := unmarshalRequestFile(fileName)
	if err != nil {
		return RequestFile{}, fmt.Errorf("error unmarshalling %s: %w", fileName, err)
	}

	return RequestFile{
		FileName:    fileName,
		RequestData: data,
	}, nil
}

func GetRequestFiles(directory string) ([]RequestFile, error) {
	jsonEntries, err := getJsonEntries(directory)
	if err != nil {
		return nil, err
	}

	var requestFiles []RequestFile
	for _, entry := range jsonEntries {
		data, err := unmarshalRequestFile(directory + "/" + entry.Name())
		if err != nil {
			return nil, fmt.Errorf("error unmarshalling %s: %w", entry.Name(), err)
		}
		requestFiles = append(requestFiles, RequestFile{
			FileName:    entry.Name(),
			RequestData: data,
		})
	}

	return requestFiles, nil
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
