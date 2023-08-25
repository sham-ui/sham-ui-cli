package app

import (
	"github.com/go-logr/logr"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	articlesHandlers "site/articles/handlers"
	"site/assets"
	"site/config"
	"site/proto"
	"site/ssr"
)

func StartApplication(logger logr.Logger, n *negroni.Negroni, cmsClient proto.CMSClient, render ssr.Render) {
	logger.Info("Allowed domains", "domains", config.Server.AllowedDomains)

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
	r.HandleFunc("/api/articles", articlesHandlers.NewListHandler(logger, cmsClient)).Methods("GET")
	r.HandleFunc("/api/articles/{slug}", articlesHandlers.NewDetailHandler(logger, cmsClient)).Methods("GET")

	// Assets
	r.HandleFunc("/assets/{file}", articlesHandlers.NewAssetsHandler(logger, cmsClient)).Methods("GET")

	// Resources
	spaHandler := assets.NewHandler(logger, render)
	r.PathPrefix("/").Handler(spaHandler)

	// Middleware
	n.Use(c)
	n.UseHandler(r)
}
