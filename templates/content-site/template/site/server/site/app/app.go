package app

import (
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	articlesHandlers "site/articles/handlers"
	"site/assets"
	"site/config"
	"site/proto"
	"site/ssr"
	"strings"
)

func StartApplication(n *negroni.Negroni, cmsClient proto.CMSClient, render ssr.Render) {
	log.Infof("Allowed domains: %s", strings.Join(config.Server.AllowedDomains, ", "))

	c := cors.New(cors.Options{
		AllowedOrigins:   config.Server.AllowedDomains,
		AllowedMethods:   []string{"PUT", "POST", "GET"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-CSRF-Token", "Content-Type"},
		ExposedHeaders:   []string{"X-CSRF-Token", "Content-Type"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	r := mux.NewRouter()

	// Articles
	r.HandleFunc("/api/articles", articlesHandlers.NewListHandler(cmsClient)).Methods("GET")
	r.HandleFunc("/api/articles/{slug}", articlesHandlers.NewDetailHandler(cmsClient)).Methods("GET")

	// Assets
	r.HandleFunc("/assets/{file}", articlesHandlers.NewAssetsHandler(cmsClient)).Methods("GET")

	// Resources
	spaHandler := assets.NewHandler(render)
	r.PathPrefix("/").Handler(spaHandler)

	// Middleware
	n.Use(c)
	n.UseHandler(r)
}
