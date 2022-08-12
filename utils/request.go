package utils

import "net/url"

func GetRequest(target *url.URL, testCase TestCase) HTTPMessage {

	defaultHeaders := []Header{
		{":authority", target.Host},
		{":method", testCase.Method},
		{":path", target.RequestURI()},
		{":scheme", "https"},
	}

	return HTTPMessage{
		Body:    []byte(testCase.Body),
		Headers: append(defaultHeaders, testCase.Headers...),
	}
}
