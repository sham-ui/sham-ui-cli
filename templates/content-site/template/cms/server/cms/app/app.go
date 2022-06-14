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
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	"net/http"
	"strings"
	"time"
)

func StartApplication(configPath string, n *negroni.Negroni) *sql.DB {
	config.LoadConfiguration(configPath)

	db, err := database.ConnToDB(config.DataBase.GetURL())
	if nil != err {
		log.Fatalf("Fail connect to db: %s", err)
	}

	sessionsStore, err := sessions.NewStore(db, config.Session.Secret)
	if nil == err {
		log.Info("Create pg session store")
	} else {
		log.WithError(err).Fatal("Fail create pg session store")
	}
	// Run a background goroutine to clean up expired sessions from the database.
	defer sessionsStore.StopCleanup(sessionsStore.Cleanup(time.Minute * 5))

	migrator, err := migrations.NewMigrator(db)
	if nil != err {
		log.Fatalf("Fail create migrator: %s", err)
	}
	for _, migrationByModule := range [][]migrations.Migration{
		members.Migrations(db),
		articlesDB.Migrations(db),
	} {
		err = migrator.Apply(migrationByModule...)
		if nil != err {
			log.Fatalf("Fail apply migrations: %s", err)
		}
	}

	log.Infof("Allowed domains: %s", strings.Join(config.Server.AllowedDomains, ", "))

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
	r.HandleFunc("/api/csrftoken", sessionHandlers.NewCsrfTokenHandler()).Methods(http.MethodGet)
	r.HandleFunc("/api/login", authenticationHandlers.NewLoginHandler(sessionsStore, db)).Methods(http.MethodPost)
	r.HandleFunc("/api/logout", authenticationHandlers.NewLogoutHandler(sessionsStore)).Methods(http.MethodPost)

	// Session Routes
	r.HandleFunc("/api/validsession", sessionHandlers.NewValidSessionHandler(sessionsStore)).Methods(http.MethodGet)

	// Member CRUD routes
	r.HandleFunc("/api/members/email", membersHandlers.NewUpdateEmailHandler(db, sessionsStore)).Methods(http.MethodPut)
	r.HandleFunc("/api/members/name", membersHandlers.NewUpdateNameHandler(db, sessionsStore)).Methods(http.MethodPut)
	r.HandleFunc("/api/members/password", membersHandlers.NewUpdatePasswordHandler(db, sessionsStore)).Methods(http.MethodPut)

	// Superuser routes
	r.HandleFunc("/api/admin/members", membersHandlers.NewListHandler(db, sessionsStore)).Methods(http.MethodGet)
	r.HandleFunc("/api/admin/members", membersHandlers.NewCreateHandler(db, sessionsStore)).Methods(http.MethodPost)
	r.HandleFunc("/api/admin/members/{id:[0-9]+}", membersHandlers.NewUpdateHandler(db, sessionsStore)).Methods(http.MethodPut)
	r.HandleFunc("/api/admin/members/{id:[0-9]+}/password", membersHandlers.NewResetPasswordHandler(db, sessionsStore)).Methods(http.MethodPut)
	r.HandleFunc("/api/admin/members/{id:[0-9]+}", membersHandlers.NewDeleteHandler(db, sessionsStore)).Methods(http.MethodDelete)
	r.HandleFunc("/api/admin/server-info", serverHandlers.NewInfoHandler(sessionsStore)).Methods(http.MethodGet)

	// Articles category routes
	r.HandleFunc("/api/categories", category.NewCreateHandler(db, sessionsStore)).Methods(http.MethodPost)
	r.HandleFunc("/api/categories", category.NewListHandler(db, sessionsStore)).Methods(http.MethodGet)
	r.HandleFunc("/api/categories/{id:[0-9]+}", category.NewUpdateHandler(db, sessionsStore)).Methods(http.MethodPut)
	r.HandleFunc("/api/categories/{id:[0-9]+}", category.NewDeleteHandler(db, sessionsStore)).Methods(http.MethodDelete)

	// Articles tag routes
	r.HandleFunc("/api/tags", tag.NewCreateHandler(db, sessionsStore)).Methods(http.MethodPost)
	r.HandleFunc("/api/tags", tag.NewListHandler(db, sessionsStore)).Methods(http.MethodGet)
	r.HandleFunc("/api/tags/{id:[0-9]+}", tag.NewUpdateHandler(db, sessionsStore)).Methods(http.MethodPut)
	r.HandleFunc("/api/tags/{id:[0-9]+}", tag.NewDeleteHandler(db, sessionsStore)).Methods(http.MethodDelete)

	// Articles routes
	r.HandleFunc("/api/articles", article.NewListHandler(db, sessionsStore)).Methods(http.MethodGet)
	r.HandleFunc("/api/articles", article.NewCreateHandler(db, sessionsStore)).Methods(http.MethodPost)
	r.HandleFunc("/api/articles/{id:[0-9]+}", article.NewUpdateHandler(db, sessionsStore)).Methods(http.MethodPut)
	r.HandleFunc("/api/articles/{id:[0-9]+}", article.NewDetailHandler(db, sessionsStore)).Methods(http.MethodGet)
	r.HandleFunc("/api/articles/{id:[0-9]+}", article.NewDeleteHandler(db, sessionsStore)).Methods(http.MethodDelete)

	// Resources
	spaHandler := assets.NewHandler(sessionsStore)
	r.PathPrefix("/").Handler(spaHandler)

	// Middleware
	n.Use(c)
	n.UseHandler(CSRF(r))

	return db
}
