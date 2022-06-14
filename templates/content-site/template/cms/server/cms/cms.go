package main

import (
	"cms/app"
	"cms/config"
	"flag"
	"github.com/gorilla/context"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"net/http"
	"strconv"
)

func main() {
	createSuperuserFlag := flag.Bool("createsuperuser", false, "create superuser member")
	flag.Parse()
	if *createSuperuserFlag {
		app.CreateSuperUser()
		return
	}
	n := negroni.New(negroni.NewRecovery(), negroni.NewLogger())
	db := app.StartApplication("config.cfg", n)
	go app.StartGRPC(db)
	port := strconv.Itoa(config.Server.Port)
	log.Infof("Server start on port :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, context.ClearHandler(n)))
}
