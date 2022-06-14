package app

import (
	"cms/api"
	"cms/config"
	"cms/proto"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"path"
)

func StartGRPC(db *sql.DB) {
	if _, err := os.Stat(config.Api.SocketPath); !os.IsNotExist(err) {
		err := os.Remove(config.Api.SocketPath)
		if nil != err {
			log.WithError(err).Fatal("can't remove API socket")
		}
	}
	err := os.MkdirAll(path.Dir(config.Api.SocketPath), os.ModePerm)
	if nil != err {
		log.WithError(err).Fatal("can't create unix path dirs")
	}
	addr, err := net.ResolveUnixAddr("unix", config.Api.SocketPath)
	if nil != err {
		log.WithError(err).Fatal("can't resolve unix addr")
	}
	lis, err := net.ListenUnix("unix", addr)
	if nil != err {
		log.WithError(err).Fatal("can't listen addr")
	}
	if err := os.Chmod(config.Api.SocketPath, os.ModePerm); nil != err {
		log.WithError(err).Fatal("can't change socket file permission")
	}
	srv := api.NewAPI(db)
	s := grpc.NewServer()
	proto.RegisterCMSServer(s, srv)

	log.WithField("socket", config.Api.SocketPath).Info("API server start")

	err = s.Serve(lis)
	if nil != err {
		log.WithError(err).Fatal("failed to server")
	}
}
