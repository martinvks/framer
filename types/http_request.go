package types

type HttpRequest struct {
	Headers Headers
	Body    []byte
}

func (request HttpRequest) HasBody() bool {
	return len(request.Body) > 0
}
