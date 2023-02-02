package http3

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/context"
	"io"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/martinvks/framer/types"
	"github.com/quic-go/quic-go"
)

const (
	frameTypeData     = 0x00
	frameTypeHeaders  = 0x01
	frameTypeSettings = 0x04
)

func SendHTTP3Request(
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

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %w", err)
	}

	udpConn, err := net.ListenPacket("udp", ":0")
	if err != nil {
		return nil, err
	}

	defer func() { _ = udpConn.Close() }()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	go func() {
		<-ctx.Done()
		_ = udpConn.Close()
	}()

	udpAddr := &net.UDPAddr{
		IP:   ip,
		Port: portInt,
	}

	tlsConfig := &tls.Config{
		NextProtos:         []string{"h3", "h3-29"},
		ServerName:         target.Hostname(),
		InsecureSkipVerify: true,
		KeyLogWriter:       keyLogWriter,
	}

	quicConfig := &quic.Config{
		Versions:           []quic.VersionNumber{quic.Version1, quic.VersionDraft29},
		MaxIncomingStreams: -1,
	}

	connection, err := quic.DialEarlyContext(
		ctx,
		udpConn,
		udpAddr,
		target.Hostname(),
		tlsConfig,
		quicConfig,
	)
	if err != nil {
		return nil, err
	}
	defer func() { _ = connection.CloseWithError(0, "") }()

	requestStream, err := writeRequest(connection, request)
	if err != nil {
		return nil, err
	}

	return readResponse(ctx, requestStream)
}
