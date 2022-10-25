package http3

import (
	"bufio"
	"errors"
	"golang.org/x/net/context"
	"io"

	"github.com/Martinvks/httptestrunner/types"
	"github.com/lucas-clemente/quic-go"
	"github.com/marten-seemann/qpack"
)

type http3Frame struct {
	frameType int
	length    uint64
	data      []byte
}

func readResponse(ctx context.Context, requestStream quic.Stream) (*types.HttpResponse, error) {
	response := types.HttpResponse{}

	headersDecoder := qpack.NewDecoder(func(f qpack.HeaderField) {
		response.Headers = append(
			response.Headers,
			types.Header{Name: f.Name, Value: f.Value})
	})

	reader := bufio.NewReader(requestStream)

	for {
		frame, err := readFrame(reader)

		if err != nil {
			if ctx.Err() != nil {
				return nil, errors.New("timeout")
			}

			if err == io.EOF {
				break
			}

			return nil, err
		}

		switch frame.frameType {
		case frameTypeData:
			response.Body = append(response.Body, frame.data...)
		case frameTypeHeaders:
			if _, err := headersDecoder.Write(frame.data); err != nil {
				return nil, err
			}
		}
	}

	return &response, nil
}

func readFrame(reader *bufio.Reader) (*http3Frame, error) {
	frameType, err := readVarInt(reader)
	if err != nil {
		return nil, err
	}
	length, err := readVarInt(reader)
	if err != nil {
		return nil, err
	}
	data := make([]byte, length)
	if _, err := io.ReadFull(reader, data); err != nil {
		return nil, err
	}

	return &http3Frame{
		frameType: int(frameType),
		length:    length,
		data:      data,
	}, nil
}

func readVarInt(b io.ByteReader) (uint64, error) {
	firstByte, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	// the first two bits of the first byte encode the length
	intLen := 1 << ((firstByte & 0xc0) >> 6)
	b1 := firstByte & (0xff - 0xc0)
	if intLen == 1 {
		return uint64(b1), nil
	}
	b2, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	if intLen == 2 {
		return uint64(b2) + uint64(b1)<<8, nil
	}
	b3, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	b4, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	if intLen == 4 {
		return uint64(b4) + uint64(b3)<<8 + uint64(b2)<<16 + uint64(b1)<<24, nil
	}
	b5, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	b6, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	b7, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	b8, err := b.ReadByte()
	if err != nil {
		return 0, err
	}
	return uint64(b8) + uint64(b7)<<8 + uint64(b6)<<16 + uint64(b5)<<24 + uint64(b4)<<32 + uint64(b3)<<40 + uint64(b2)<<48 + uint64(b1)<<56, nil
}
