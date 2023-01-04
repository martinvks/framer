package http2

import (
	"bytes"
	"crypto/tls"
	"golang.org/x/net/http2"

	"github.com/martinvks/framer/types"
)

const (
	maxFrameSizeOctets                    int    = 16384
	maxWindowSize                         uint32 = (1 << 31) - 1
	initialFlowControlWindow              uint32 = (1 << 16) - 1
	windowSizeIncrement                          = maxWindowSize - initialFlowControlWindow
	literalHeaderFieldNeverIndexedNewName byte   = 0x10
)

func writeRequest(tlsConn *tls.Conn, request *types.HttpRequest) error {
	requestBuf := bytes.NewBuffer(nil)
	requestBuf.Write([]byte(http2.ClientPreface))

	framer := http2.NewFramer(requestBuf, nil)

	_ = framer.WriteSettings(http2.Setting{
		ID:  http2.SettingInitialWindowSize,
		Val: maxWindowSize,
	})

	_ = framer.WriteWindowUpdate(
		0,
		windowSizeIncrement,
	)

	_ = framer.WriteHeaders(http2.HeadersFrameParam{
		StreamID:      1,
		BlockFragment: hpackEncodeHeaders(request.Headers),
		EndStream:     !request.HasBody() && !request.HasTrailerSection(),
		EndHeaders:    !request.HasContinuationHeaders(),
	})

	if request.HasContinuationHeaders() {
		_ = framer.WriteContinuation(
			1,
			true,
			hpackEncodeHeaders(request.Continuation),
		)
	}

	if request.HasBody() {
		start := 0
		for start < len(request.Body) {
			end := start + maxFrameSizeOctets
			if end > len(request.Body) {
				end = len(request.Body)
			}
			_ = framer.WriteData(
				1,
				end == len(request.Body) && !request.HasTrailerSection(),
				request.Body[start:end],
			)
			start = end
		}
	}

	if request.HasTrailerSection() {
		_ = framer.WriteHeaders(http2.HeadersFrameParam{
			StreamID:      1,
			BlockFragment: hpackEncodeHeaders(request.Trailer),
			EndStream:     true,
			EndHeaders:    true,
		})
	}

	_ = framer.WriteSettingsAck()

	_, err := tlsConn.Write(requestBuf.Bytes())
	return err
}

func hpackEncodeHeaders(headers types.Headers) []byte {
	var hpackBuf []byte
	for _, header := range headers {
		hpackBuf = hpackAppendHeader(hpackBuf, &header)
	}
	return hpackBuf
}

func hpackAppendHeader(dst []byte, h *types.Header) []byte {
	dst = append(dst, literalHeaderFieldNeverIndexedNewName)
	dst = hpackAppendVarInt(dst, 7, uint64(len(h.Name)))
	dst = append(dst, h.Name...)
	dst = hpackAppendVarInt(dst, 7, uint64(len(h.Value)))
	dst = append(dst, h.Value...)
	return dst
}

func hpackAppendVarInt(dst []byte, n byte, val uint64) []byte {
	k := uint64((1 << n) - 1)
	if val < k {
		return append(dst, byte(val))
	}
	dst = append(dst, byte(k))
	val -= k
	for ; val >= 128; val >>= 7 {
		dst = append(dst, byte(0x80|(val&0x7f)))
	}
	return append(dst, byte(val))
}
