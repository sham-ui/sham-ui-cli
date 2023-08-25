package main

import (
	"cms/app"
	"cms/config"
	"cms/logger"
	"flag"
	"github.com/gorilla/context"
	"github.com/urfave/negroni"
	"net/http"
	"strconv"
)

func main() {
	createSuperuserFlag := flag.Bool("createsuperuser", false, "create superuser member")
	flag.Parse()

	log := logger.NewLogger(0)

	if *createSuperuserFlag {
		app.CreateSuperUser(log)
		return
	}
	n := negroni.New(negroni.NewRecovery(), logger.CreateNegroniLogger(log))
	db := app.StartApplication(log, "config.cfg", n)
	go app.StartGRPC(log, db)
	port := strconv.Itoa(config.Server.Port)
	log.Info("Server started", "port", port)
	if err := http.ListenAndServe(":"+port, context.ClearHandler(n)); err != nil {
		log.Error(err, "Server stopped")
	}
}
