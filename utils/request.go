package utils

import (
	"fmt"
	"net/url"

	"github.com/martinvks/framer/types"
)

func GetRequest(id string, addIdQuery bool, requestData RequestData, commonArgs types.CommonArguments) types.HttpRequest {
	headers := getHeaders(
		id,
		addIdQuery,
		requestData,
		commonArgs,
	)

	continuation := getContinuationHeaders(
		commonArgs.Proto,
		requestData,
	)

	return types.HttpRequest{
		Headers:      headers,
		Continuation: continuation,
		Body:         []byte(requestData.Body),
		Trailer:      requestData.Trailer,
	}
}

func getHeaders(id string, addIdQuery bool, requestData RequestData, commonArgs types.CommonArguments) types.Headers {
	var headers types.Headers

	if requestData.AddDefaultHeaders {
		headers = getHeadersWithDefaults(id, addIdQuery, commonArgs.Target, requestData)
	} else {
		headers = requestData.Headers
	}

	if commonArgs.AddIdHeader {
		headers = append(headers, types.Header{Name: commonArgs.IdHeaderName, Value: id})
	}

	headers = append(headers, commonArgs.CommonHeaders...)

	return headers
}

func getHeadersWithDefaults(
	id string,
	addIdQuery bool,
	target *url.URL,
	requestData RequestData,
) types.Headers {
	var headers types.Headers

	if addIdQuery {
		query := target.Query()
		query.Set("id", id)
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
		if val, ok := requestData.Headers.Get(header.Name); ok {
			header.Value = val
			toSkip[header.Name] = struct{}{}
		}
	}

	for _, header := range requestData.Headers {
		if _, ok := toSkip[header.Name]; ok {
			delete(toSkip, header.Name)
			continue
		}
		headers = append(headers, header)
	}

	return headers
}

func getContinuationHeaders(proto int, requestData RequestData) types.Headers {
	var continuation types.Headers

	switch proto {
	case types.H2:
		continuation = requestData.Continuation
	case types.H3:
		if len(requestData.Continuation) > 0 {
			fmt.Printf("WARN: continuation headers not supported for HTTP/3 and will be ignored\n")
		}
	}

	return continuation
}
