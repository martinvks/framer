package utils

import (
	"net/url"

	"github.com/Martinvks/httptestrunner/types"
)

func GetRequest(target *url.URL, testCase TestCase) types.HttpRequest {

	defaultHeaders := types.Headers{
		{":authority", target.Host},
		{":method", testCase.Method},
		{":path", target.RequestURI()},
		{":scheme", "https"},
		{"x-id", testCase.Id},
	}

	return types.HttpRequest{
		Headers:      append(defaultHeaders, testCase.Headers...),
		Continuation: testCase.Continuation,
		Body:         []byte(testCase.Body),
		Trailer:      testCase.Trailer,
	}
}
