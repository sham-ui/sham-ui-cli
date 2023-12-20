package dialer

import (
	"context"
	"net"
	"strings"
	"time"
)

func dialUnix(ctx context.Context, addr string) (net.Conn, error) {
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "unix", addr)
	if err != nil {
		return nil, NewDialUnixSocketError(err)
	}
	return conn, nil
}

func dialTCP(ctx context.Context, addr string) (net.Conn, error) {
	var dialer net.Dialer
	conn, err := dialer.DialContext(ctx, "tcp", addr)
	if err != nil {
		return nil, NewDialTCPError(err)
	}
	return conn, nil
}

func Dial(ctx context.Context, addr string) (net.Conn, error) {
	switch {
	case strings.HasPrefix(addr, "unix:"):
		return dialUnix(ctx, addr[5:])
	case strings.HasPrefix(addr, "/"):
		return dialUnix(ctx, addr)
	default:
		return dialTCP(ctx, addr)
	}
}

func WithTimeout(timeout time.Duration) func(ctx context.Context, addr string) (net.Conn, error) {
	return func(ctx context.Context, addr string) (net.Conn, error) {
		ctxTimeout, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return Dial(ctxTimeout, addr)
	}
}
