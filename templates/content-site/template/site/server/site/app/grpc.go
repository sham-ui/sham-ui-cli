package app

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
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

func CreateCMSClient() proto.CMSClient {
	if _, err := os.Stat(config.Api.SocketPath); nil != err {
		log.WithError(err).Fatal("can't find API socket file")
	}
	conn, err := grpc.Dial(config.Api.SocketPath, grpc.WithInsecure(), grpc.WithContextDialer(unixConnect))
	if nil != err {
		log.WithError(err).Fatal("can't dial with CMS")
	}
	return proto.NewCMSClient(conn)
}
