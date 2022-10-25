package utils

import (
	"fmt"
	"net/url"

	"github.com/Martinvks/httptestrunner/types"
)

func GetRequest(
	addIdHeader bool,
	addIdQuery bool,
	proto int,
	target *url.URL,
	testCase TestCase,
) types.HttpRequest {
	headers := getHeaders(target, addIdHeader, addIdQuery, testCase)
	continuation := getContinuationHeaders(proto, testCase)

	return types.HttpRequest{
		Headers:      headers,
		Continuation: continuation,
		Body:         []byte(testCase.Body),
		Trailer:      testCase.Trailer,
	}
}

func getHeaders(target *url.URL, addIdHeader bool, addIdQuery bool, testCase TestCase) types.Headers {
	var headers types.Headers

	if testCase.AddDefaultHeaders {
		if addIdQuery {
			query := target.Query()
			query.Set("id", testCase.Id)
			target.RawQuery = query.Encode()
		}

		headers = types.Headers{
			{":authority", target.Host},
			{":method", "GET"},
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
	} else {
		headers = testCase.Headers
	}

	if addIdHeader {
		headers = append(headers, types.Header{Name: "x-id", Value: testCase.Id})
	}

	return headers
}

func getContinuationHeaders(proto int, testCase TestCase) types.Headers {
	var continuation types.Headers

	switch proto {
	case types.H2:
		continuation = testCase.Continuation
	case types.H3:
		if len(testCase.Continuation) > 0 {
			fmt.Printf("WARN: continuation headers not supported for HTTP/3 and will be ignored\n")
		}
	}

	return continuation
}
