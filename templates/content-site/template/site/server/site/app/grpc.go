package app

import (
	"context"
	"fmt"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"
	"net"
	"os"
	"site/config"
	"site/proto"
)

func unixConnect(ctx context.Context, addr string) (net.Conn, error) {
	unix_addr, err := net.ResolveUnixAddr("unix", addr)
	if nil != err {
		return nil, fmt.Errorf("resolve unix addrs: %s", err)
	}
	conn, err := net.DialUnix("unix", nil, unix_addr)
	if nil != err {
		return nil, fmt.Errorf("dial unix: %s", err)
	}
	return conn, nil
}

func CreateCMSClient(logger logr.Logger) proto.CMSClient {
	if _, err := os.Stat(config.Api.SocketPath); nil != err {
		logger.Error(err, "can't find API socket file")
		os.Exit(1)
	}
	conn, err := grpc.Dial(config.Api.SocketPath, grpc.WithInsecure(), grpc.WithContextDialer(unixConnect))
	if nil != err {
		logger.Error(err, "Failed to connect to CMS")
		os.Exit(1)
	}
	return proto.NewCMSClient(conn)
}
