package main

import (
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"net/http"
	"site/app"
	"site/config"
	"site/logger"
	"site/ssr"
	"strconv"
)

func main() {
	loggerMiddleware := negroni.NewLogger()
	loggerMiddleware.ALogger = logger.Logger
	config.LoadConfiguration("config.cfg")
	n := negroni.New(negroni.NewRecovery(), loggerMiddleware)
	render := ssr.NewServerSideRender()
	cmsClient := app.CreateCMSClient()
	app.StartApplication(n, cmsClient, render)
	port := strconv.Itoa(config.Server.Port)
	log.Infof("Server start on port :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, context.ClearHandler(n)))
}
