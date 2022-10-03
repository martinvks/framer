package utils

import (
	"net/url"

	"github.com/Martinvks/httptestrunner/types"
)

func GetRequest(target *url.URL, testCase TestCase) types.HttpRequest {
	var headers types.Headers

	if testCase.AddDefaultHeaders {
		headers = getHeadersWithDefault(target, testCase)
	} else {
		headers = testCase.Headers
	}

	headers = append(headers, types.Header{Name: "x-id", Value: testCase.Id})

	return types.HttpRequest{
		Headers:      headers,
		Continuation: testCase.Continuation,
		Body:         []byte(testCase.Body),
		Trailer:      testCase.Trailer,
	}
}

func getHeadersWithDefault(target *url.URL, testCase TestCase) types.Headers {
	headers := types.Headers{
		{":authority", target.Host},
		{":method", testCase.Method},
		{":path", target.RequestURI()},
		{":scheme", "https"},
	}

	var toSkip = make(map[string]struct{})

	for i := range headers {
		header := &headers[i]
		if val, ok := testCase.Headers.Get(header.Name); ok {
			header.Value = val
			toSkip[header.Name] = struct{}{}
		}
	}

	for _, header := range testCase.Headers {
		if _, ok := toSkip[header.Name]; ok {
			delete(toSkip, header.Name)
			continue
		}
		headers = append(headers, header)
	}

	return headers
}
