package utils

import "net/url"

func GetRequest(target *url.URL, testCase TestCase) HTTPMessage {

	defaultHeaders := []Header{
		{":authority", target.Host},
		{":method", testCase.Data.Method},
		{":path", target.RequestURI()},
		{":scheme", "https"},
		{"x-id", testCase.Id},
	}

	return HTTPMessage{
		Body:    []byte(testCase.Data.Body),
		Headers: append(defaultHeaders, testCase.Data.Headers...),
	}
}
