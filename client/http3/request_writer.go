package http3

import (
	"bytes"
	"fmt"

	"github.com/martinvks/framer/types"
	"github.com/quic-go/qpack"
	"github.com/quic-go/quic-go"
)

const streamTypeControlStream = 0x00

const (
	maxVarInt1 = 63
	maxVarInt2 = 16383
	maxVarInt4 = 1073741823
	maxVarInt8 = 4611686018427387903
)

func writeRequest(connection quic.Connection, request *types.HttpRequest) (quic.Stream, error) {
	controlStream, err := connection.OpenUniStream()
	if err != nil {
		return nil, err
	}

	_, err = controlStream.Write([]byte{
		streamTypeControlStream,
		frameTypeSettings,
		0x0,
	})
	if err != nil {
		return nil, err
	}

	requestStream, err := connection.OpenStream()
	if err != nil {
		return nil, err
	}

	frames := [][]byte{getHeadersFrame(request.Headers)}

	if request.HasBody() {
		frames = append(frames, getDataFrame(request.Body))
	}

	if request.HasTrailerSection() {
		frames = append(frames, getHeadersFrame(request.Trailer))
	}

	for _, frame := range frames {
		_, _ = requestStream.Write(frame)
	}

	err = requestStream.Close()
	if err != nil {
		return nil, err
	}

	return requestStream, nil
}

func getHeadersFrame(headers types.Headers) []byte {
	headersFrame := bytes.NewBuffer(nil)
	qpackBuffer := bytes.NewBuffer(nil)
	qpackEncoder := qpack.NewEncoder(qpackBuffer)

	for _, h := range headers {
		_ = qpackEncoder.WriteField(qpack.HeaderField{Name: h.Name, Value: h.Value})
	}

	headersFrame.WriteByte(frameTypeHeaders)
	headersFrame.Write(getIntegerEncoding(uint64(qpackBuffer.Len())))
	headersFrame.Write(qpackBuffer.Bytes())
	return headersFrame.Bytes()
}

func getDataFrame(body []byte) []byte {
	dataFrame := bytes.NewBuffer(nil)
	dataFrame.WriteByte(frameTypeData)
	dataFrame.Write(getIntegerEncoding(uint64(len(body))))
	dataFrame.Write(body)
	return dataFrame.Bytes()
}

func getIntegerEncoding(i uint64) []byte {
	if i <= maxVarInt1 {
		return []byte{uint8(i)}
	} else if i <= maxVarInt2 {
		return []byte{uint8(i>>8) | 0x40, uint8(i)}
	} else if i <= maxVarInt4 {
		return []byte{uint8(i>>24) | 0x80, uint8(i >> 16), uint8(i >> 8), uint8(i)}
	} else if i <= maxVarInt8 {
		return []byte{
			uint8(i>>56) | 0xc0, uint8(i >> 48), uint8(i >> 40), uint8(i >> 32),
			uint8(i >> 24), uint8(i >> 16), uint8(i >> 8), uint8(i),
		}
	} else {
		panic(fmt.Sprintf("%#x doesn't fit into 62 bits", i))
	}
}
