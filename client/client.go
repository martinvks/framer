package client

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"time"

	"github.com/martinvks/framer/client/http2"
	"github.com/martinvks/framer/client/http3"
	"github.com/martinvks/framer/types"
)

func DoRequest(
	proto int,
	target *url.URL,
	timeout time.Duration,
	keyLogWriter io.Writer,
	ip net.IP,
	request *types.HttpRequest,
) (*types.HttpResponse, error) {
	switch proto {
	case types.H2:
		return http2.SendHTTP2Request(ip, target, timeout, keyLogWriter, request)
	case types.H3:
		return http3.SendHTTP3Request(ip, target, timeout, keyLogWriter, request)
	default:
		return nil, fmt.Errorf("unknown proto %d", proto)
	}
}
