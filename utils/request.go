package utils

import (
	"github.com/Martinvks/httptestrunner/types"
	"net/url"
)

func GetRequest(target *url.URL, testCase TestCase) types.HttpRequest {

	defaultHeaders := types.Headers{
		{":authority", target.Host},
		{":method", testCase.Data.Method},
		{":path", target.RequestURI()},
		{":scheme", "https"},
		{"x-id", testCase.Id},
	}

	return types.HttpRequest{
		Body:    []byte(testCase.Data.Body),
		Headers: append(defaultHeaders, testCase.Data.Headers...),
	}
}
