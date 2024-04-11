package net_addr

import (
	"strings"
)

const (
	UnixNetwork = "unix"
	TcpNetwork  = "tcp"
)

func Resolve(addr string) (string, string) {
	switch {
	case strings.HasPrefix(addr, "unix:"):
		return UnixNetwork, addr[5:]
	case strings.HasPrefix(addr, "/"):
		return UnixNetwork, addr
	default:
		return TcpNetwork, addr
	}
}
