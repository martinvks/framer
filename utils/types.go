package utils

type Header struct {
	Name  string
	Value string
}

type HTTPMessage struct {
	Headers []Header
	Body    []byte
}
