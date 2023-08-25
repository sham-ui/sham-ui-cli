package main

import (
	"github.com/gorilla/context"
	"github.com/urfave/negroni"
	"net/http"
	"site/app"
	"site/config"
	"site/logger"
	"site/ssr"
	"strconv"
)

func main() {
	logr := logger.NewLogger(0)
	config.LoadConfiguration(logr, "config.cfg")
	n := negroni.New(negroni.NewRecovery(), logger.CreateNegroniLogger(logr))
	render := ssr.NewServerSideRender(logr)
	render.Start()
	cmsClient := app.CreateCMSClient(logr)
	app.StartApplication(logr, n, cmsClient, render)
	port := strconv.Itoa(config.Server.Port)

	logr.Info("Server started", "port", port)
	logr.Error(http.ListenAndServe(":"+port, context.ClearHandler(n)), "Server stopped")
}
