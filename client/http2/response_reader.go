package http2

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"

	"github.com/martinvks/framer/types"
)

func readResponse(tlsConn *tls.Conn) (*types.HttpResponse, error) {
	response := types.HttpResponse{}

	headersDecoder := hpack.NewDecoder(^uint32(0), func(f hpack.HeaderField) {
		response.Headers = append(
			response.Headers,
			types.Header{Name: f.Name, Value: f.Value},
		)
	})

	framer := http2.NewFramer(nil, tlsConn)

	hasBody := false
	bodyRead := false
	headersDone := false
	for !headersDone || (hasBody && !bodyRead) {
		frame, err := framer.ReadFrame()
		if err != nil {
			return nil, err
		}

		if ga, ok := frame.(*http2.GoAwayFrame); ok {
			return nil, fmt.Errorf("GOAWAY: error code %v", ga.ErrCode)
		}

		if frame.Header().StreamID != 1 {
			continue
		}

		switch frame := frame.(type) {
		case *http2.HeadersFrame:
			if _, err := headersDecoder.Write(frame.HeaderBlockFragment()); err != nil {
				return nil, err
			}
			headersDone = frame.HeadersEnded()
			hasBody = !frame.StreamEnded()

		case *http2.ContinuationFrame:
			if _, err := headersDecoder.Write(frame.HeaderBlockFragment()); err != nil {
				return nil, err
			}
			headersDone = frame.HeadersEnded()

		case *http2.DataFrame:
			response.Body = append(response.Body, frame.Data()...)
			bodyRead = frame.StreamEnded()

		case *http2.RSTStreamFrame:
			return nil, fmt.Errorf("RST_STREAM: error code %v", frame.ErrCode)
		}
	}

	return &response, nil
}
