package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"

	"github.com/martinvks/framer/types"
)

type jsonFileData struct {
	AddDefaultHeaders bool
	Headers           map[string]jsonHeaderValue
	Continuation      map[string]jsonHeaderValue
	Trailer           map[string]jsonHeaderValue
	Body              string
}

type jsonHeaderValue []string

func (v *jsonHeaderValue) UnmarshalJSON(input []byte) error {
	var header string
	err := json.Unmarshal(input, &header)
	if err == nil {
		*v = append(*v, header)
		return nil
	}

	var headers []string
	err = json.Unmarshal(input, &headers)
	if err == nil {
		*v = append(*v, headers...)
		return nil
	}

	return fmt.Errorf("invalid json header value \"%s\", must be a string or array of strings", string(input))
}

func unmarshalRequestFile(fileName string) (RequestData, error) {
	content, err := os.ReadFile(fileName)
	if err != nil {
		return RequestData{}, err
	}

	contentWithEnv := replaceWithEnvironmentVariables(content)

	fileData := jsonFileData{
		AddDefaultHeaders: true,
	}
	err = json.Unmarshal(contentWithEnv, &fileData)

	if err != nil {
		return RequestData{}, err
	}

	return RequestData{
		AddDefaultHeaders: fileData.AddDefaultHeaders,
		Headers:           mapToRequestHeaders(fileData.Headers),
		Continuation:      mapToRequestHeaders(fileData.Continuation),
		Trailer:           mapToRequestHeaders(fileData.Trailer),
		Body:              fileData.Body,
	}, nil
}

func mapToRequestHeaders(values map[string]jsonHeaderValue) types.Headers {
	var headers types.Headers

	for key, jsonHV := range values {
		for _, value := range jsonHV {
			headers = append(headers, types.Header{
				Name:  key,
				Value: value,
			})
		}
	}

	return headers
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
