package types

type HttpRequest struct {
	Headers      Headers
	Continuation Headers
	Trailer      Headers
	Body         []byte
}

func (request *HttpRequest) HasContinuationHeaders() bool {
	return len(request.Continuation) > 0
}

func (request *HttpRequest) HasBody() bool {
	return len(request.Body) > 0
}

func (request *HttpRequest) HasTrailerSection() bool {
	return len(request.Trailer) > 0
}
