package app

import (
	"cms/api"
	"cms/config"
	"cms/proto"
	"database/sql"
	"github.com/go-logr/logr"
	"google.golang.org/grpc"
	"net"
	"os"
	"path"
)

func StartGRPC(logger logr.Logger, db *sql.DB) {
	if _, err := os.Stat(config.Api.SocketPath); !os.IsNotExist(err) {
		err := os.Remove(config.Api.SocketPath)
		if nil != err {
			logger.Error(err, "can't remove API socket")
			os.Exit(1)
		}
	}
	err := os.MkdirAll(path.Dir(config.Api.SocketPath), os.ModePerm)
	if nil != err {
		logger.Error(err, "can't create API socket dir")
		os.Exit(1)
	}
	addr, err := net.ResolveUnixAddr("unix", config.Api.SocketPath)
	if nil != err {
		logger.Error(err, "can't resolve unix addr")
		os.Exit(1)
	}
	lis, err := net.ListenUnix("unix", addr)
	if nil != err {
		logger.Error(err, "can't listen unix addr")
		os.Exit(1)
	}
	if err := os.Chmod(config.Api.SocketPath, os.ModePerm); nil != err {
		logger.Error(err, "can't change socket file permission")
		os.Exit(1)
	}
	srv := api.NewAPI(db)
	s := grpc.NewServer()
	proto.RegisterCMSServer(s, srv)

	logger.Info("API server start", "socket", config.Api.SocketPath)

	if err = s.Serve(lis); err != nil {
		logger.Error(err, "failed to serve")
		os.Exit(1)
	}
}
