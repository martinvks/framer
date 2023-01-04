package http2

import (
	"crypto/tls"
	"io"
	"net"
	"net/url"
	"time"

	"github.com/martinvks/framer/types"
)

func SendHTTP2Request(
	ip net.IP,
	target *url.URL,
	timeout time.Duration,
	keyLogWriter io.Writer,
	request *types.HttpRequest,
) (*types.HttpResponse, error) {
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

	tlsConn := tls.Client(tcpConn, &tls.Config{
		NextProtos:         []string{"h2"},
		ServerName:         target.Hostname(),
		InsecureSkipVerify: true,
		KeyLogWriter:       keyLogWriter,
	})

	err = writeRequest(tlsConn, request)
	if err != nil {
		return nil, err
	}

	return readResponse(tlsConn)
}
