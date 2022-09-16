package http2

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/hpack"
	"io"
	"net"
	"net/url"
	"time"

	"github.com/Martinvks/httptestrunner/types"
	"github.com/Martinvks/httptestrunner/utils"
)

func SendHTTP2Request(target *url.URL, timeout time.Duration, keyLogWriter io.Writer, request *types.HttpRequest) (*types.HttpResponse, error) {
	ip, err := utils.LookUp(target.Hostname())
	if err != nil {
		return nil, err
	}

	port := target.Port()
	if port == "" {
		port = "443"
	}

	tcpConn, err := net.DialTimeout("tcp", net.JoinHostPort(ip.String(), port), timeout)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = tcpConn.Close()
	}()
	_ = tcpConn.SetDeadline(time.Now().Add(timeout))

	c := tls.Client(tcpConn, &tls.Config{
		NextProtos:         []string{"h2"},
		ServerName:         target.Hostname(),
		InsecureSkipVerify: true,
		KeyLogWriter:       keyLogWriter,
	})

	if _, err := c.Write(prepareHTTP2Request(request)); err != nil {
		return nil, err
	}

	response := types.HttpResponse{}
	headersDecoder := hpack.NewDecoder(^uint32(0), func(f hpack.HeaderField) {
		response.Headers = append(
			response.Headers,
			types.Header{Name: f.Name, Value: f.Value},
		)
	})

	framer := http2.NewFramer(nil, c)

	hasBody := false
	bodyRead := false
	headersDone := false
	for !headersDone || (hasBody && !bodyRead) {
		var f http2.Frame
		f, err = framer.ReadFrame()
		if err != nil {
			return nil, err
		}

		if ga, ok := f.(*http2.GoAwayFrame); ok {
			return nil, fmt.Errorf("received GOAWAY frame: error code %v", ga.ErrCode)
		}

		if f.Header().StreamID != 1 {
			continue
		}

		switch f := f.(type) {
		case *http2.HeadersFrame:
			if _, err := headersDecoder.Write(f.HeaderBlockFragment()); err != nil {
				return nil, err
			}
			headersDone = f.HeadersEnded()
			hasBody = !f.StreamEnded()

		case *http2.ContinuationFrame:
			if _, err := headersDecoder.Write(f.HeaderBlockFragment()); err != nil {
				return nil, err
			}
			headersDone = f.HeadersEnded()

		case *http2.DataFrame:
			// we should send window update, but who cares
			response.Body = append(response.Body, f.Data()...)
			bodyRead = f.StreamEnded()

		case *http2.RSTStreamFrame:
			return nil, fmt.Errorf("received RST_STREAM frame: error code %v", f.ErrCode)
		}
	}

	return &response, nil
}

func prepareHTTP2Request(request *types.HttpRequest) []byte {
	var hpackBuf []byte
	for i := range request.Headers {
		hpackBuf = hpackAppendHeader(hpackBuf, &request.Headers[i])
	}

	requestBuf := bytes.NewBuffer(nil)
	requestBuf.Write([]byte(http2.ClientPreface))

	framer := http2.NewFramer(requestBuf, nil)

	_ = framer.WriteSettings(http2.Setting{
		ID:  http2.SettingInitialWindowSize,
		Val: (1 << 30) - 1,
	})

	_ = framer.WriteWindowUpdate(0, (1<<30)-(1<<16)-1)

	_ = framer.WriteHeaders(http2.HeadersFrameParam{
		StreamID:      1,
		BlockFragment: hpackBuf,
		EndStream:     !request.HasBody(),
		EndHeaders:    true,
	})

	start := 0
	for start < len(request.Body) {
		end := start + 65536
		if end > len(request.Body) {
			end = len(request.Body)
		}
		_ = framer.WriteData(1, end == len(request.Body), request.Body[start:end])
		start = end
	}

	_ = framer.WriteSettingsAck()

	return requestBuf.Bytes()
}

func hpackAppendHeader(dst []byte, h *types.Header) []byte {
	dst = append(dst, 0x10)
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
