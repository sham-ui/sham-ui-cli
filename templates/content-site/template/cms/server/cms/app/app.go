package app

import (
	articlesDB "cms/articles/db"
	"cms/articles/handlers/article"
	"cms/articles/handlers/category"
	"cms/articles/handlers/tag"
	"cms/assets"
	authenticationHandlers "cms/authentication/handlers"
	"cms/config"
	"cms/core/database"
	"cms/core/migrations"
	"cms/core/sessions"
	"cms/members"
	membersHandlers "cms/members/handlers"
	serverHandlers "cms/server/handlers"
	sessionHandlers "cms/session/handlers"
	"database/sql"
	"github.com/go-logr/logr"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"net/http"
	"os"
	"time"
)

func StartApplication(logger logr.Logger, configPath string, n *negroni.Negroni) *sql.DB {
	config.LoadConfiguration(logger, configPath)

	db, err := database.ConnToDB(config.DataBase.GetURL())
	if nil != err {
		logger.Error(err, "Fail connect to db")
		os.Exit(1)
	}

	sessionsStore, err := sessions.NewStore(db, config.Session.Secret)
	if nil == err {
		logger.Info("Create pg session store")
	} else {
		logger.Error(err, "Fail create pg session store")
		os.Exit(1)
	}
	// Run a background goroutine to clean up expired sessions from the database.
	defer sessionsStore.StopCleanup(sessionsStore.Cleanup(time.Minute * 5))

	migrator, err := migrations.NewMigrator(logger, db)
	if nil != err {
		logger.Error(err, "Fail create migrator")
		os.Exit(1)
	}
	for _, migrationByModule := range [][]migrations.Migration{
		members.Migrations(db),
		articlesDB.Migrations(db),
	} {
		err = migrator.Apply(migrationByModule...)
		if nil != err {
			logger.Error(err, "Fail apply migrations")
			os.Exit(1)
		}
	}

	logger.Info("Allowed domains", "domains", config.Server.AllowedDomains)

	CSRF := csrf.Protect(
		[]byte("32-byte-long-auth-key"),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.CookieName("cms_csrf"),
		csrf.Secure(false), // Disabled for localhost non-https debugging
	)

	c := cors.New(cors.Options{
		AllowedOrigins:   config.Server.AllowedDomains,
		AllowedMethods:   []string{http.MethodPut, http.MethodPost, http.MethodGet, http.MethodDelete},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-CSRF-Token", "Content-Type"},
		ExposedHeaders:   []string{"X-CSRF-Token", "Content-Type"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	r := mux.NewRouter()

	// Authentication routes
	r.HandleFunc("/api/csrftoken", sessionHandlers.NewCsrfTokenHandler(logger)).Methods(http.MethodGet)
	r.HandleFunc("/api/login", authenticationHandlers.NewLoginHandler(logger, sessionsStore, db)).Methods(http.MethodPost)
	r.HandleFunc("/api/logout", authenticationHandlers.NewLogoutHandler(logger, sessionsStore)).Methods(http.MethodPost)

	// Session Routes
	r.HandleFunc("/api/validsession", sessionHandlers.NewValidSessionHandler(logger, sessionsStore)).Methods(http.MethodGet)

	// Member CRUD routes
	r.HandleFunc("/api/members/email", membersHandlers.NewUpdateEmailHandler(logger, db, sessionsStore)).Methods(http.MethodPut)
	r.HandleFunc("/api/members/name", membersHandlers.NewUpdateNameHandler(logger, db, sessionsStore)).Methods(http.MethodPut)
	r.HandleFunc("/api/members/password", membersHandlers.NewUpdatePasswordHandler(logger, db, sessionsStore)).Methods(http.MethodPut)

	// Superuser routes
	r.HandleFunc("/api/admin/members", membersHandlers.NewListHandler(logger, db, sessionsStore)).Methods(http.MethodGet)
	r.HandleFunc("/api/admin/members", membersHandlers.NewCreateHandler(logger, db, sessionsStore)).Methods(http.MethodPost)
	r.HandleFunc("/api/admin/members/{id:[0-9]+}", membersHandlers.NewUpdateHandler(logger, db, sessionsStore)).Methods(http.MethodPut)
	r.HandleFunc("/api/admin/members/{id:[0-9]+}/password", membersHandlers.NewResetPasswordHandler(logger, db, sessionsStore)).Methods(http.MethodPut)
	r.HandleFunc("/api/admin/members/{id:[0-9]+}", membersHandlers.NewDeleteHandler(logger, db, sessionsStore)).Methods(http.MethodDelete)
	r.HandleFunc("/api/admin/server-info", serverHandlers.NewInfoHandler(logger, sessionsStore)).Methods(http.MethodGet)

	// Articles category routes
	r.HandleFunc("/api/categories", category.NewCreateHandler(logger, db, sessionsStore)).Methods(http.MethodPost)
	r.HandleFunc("/api/categories", category.NewListHandler(logger, db, sessionsStore)).Methods(http.MethodGet)
	r.HandleFunc("/api/categories/{id:[0-9]+}", category.NewUpdateHandler(logger, db, sessionsStore)).Methods(http.MethodPut)
	r.HandleFunc("/api/categories/{id:[0-9]+}", category.NewDeleteHandler(logger, db, sessionsStore)).Methods(http.MethodDelete)

	// Articles tag routes
	r.HandleFunc("/api/tags", tag.NewCreateHandler(logger, db, sessionsStore)).Methods(http.MethodPost)
	r.HandleFunc("/api/tags", tag.NewListHandler(logger, db, sessionsStore)).Methods(http.MethodGet)
	r.HandleFunc("/api/tags/{id:[0-9]+}", tag.NewUpdateHandler(logger, db, sessionsStore)).Methods(http.MethodPut)
	r.HandleFunc("/api/tags/{id:[0-9]+}", tag.NewDeleteHandler(logger, db, sessionsStore)).Methods(http.MethodDelete)

	// Articles routes
	r.HandleFunc("/api/articles", article.NewListHandler(logger, db, sessionsStore)).Methods(http.MethodGet)
	r.HandleFunc("/api/articles", article.NewCreateHandler(logger, db, sessionsStore)).Methods(http.MethodPost)
	r.HandleFunc("/api/articles/{id:[0-9]+}", article.NewUpdateHandler(logger, db, sessionsStore)).Methods(http.MethodPut)
	r.HandleFunc("/api/articles/{id:[0-9]+}", article.NewDetailHandler(logger, db, sessionsStore)).Methods(http.MethodGet)
	r.HandleFunc("/api/articles/{id:[0-9]+}", article.NewDeleteHandler(logger, db, sessionsStore)).Methods(http.MethodDelete)

	// Upload routes
	r.HandleFunc("/api/upload-image", article.NewUploadHandler(logger, sessionsStore)).Methods(http.MethodPost)
	r.HandleFunc("/assets/{file}", article.NewImagesHandler(logger, sessionsStore)).Methods(http.MethodGet)

	// Resources
	spaHandler := assets.NewHandler(logger, sessionsStore)
	r.PathPrefix("/").Handler(spaHandler)

	// Middleware
	n.Use(c)
	n.UseHandler(CSRF(r))

	return db
}
