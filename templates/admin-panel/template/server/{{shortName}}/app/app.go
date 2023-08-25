package app

import (
	"database/sql"
	"github.com/go-logr/logr"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/urfave/negroni"
	"os"
	"{{shortName}}/assets"
	authenticationHandlers "{{shortName}}/authentication/handlers"
	"{{shortName}}/config"
	"{{shortName}}/core/database"
	"{{shortName}}/core/sessions"
	"{{shortName}}/members"
	membersHandlers "{{shortName}}/members/handlers"
	serverHandlers "{{shortName}}/server/handlers"
	sessionHandlers "{{shortName}}/session/handlers"
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

	err = members.CreateMemberStructure(db)
	if nil != err {
		logger.Error(err, "Fail create members table")
		os.Exit(1)
	} else {
		logger.Info("Create members table")
	}

	logger.Info("Allowed domains", "domains", config.Server.AllowedDomains)

	CSRF := csrf.Protect(
		[]byte("32-byte-long-auth-key"),
		csrf.RequestHeader("X-CSRF-Token"),
		csrf.CookieName("{{ shortName }}_csrf"),
		csrf.Secure(false), // Disabled for localhost non-https debugging
	)

	c := cors.New(cors.Options{
		AllowedOrigins:   config.Server.AllowedDomains,
		AllowedMethods:   []string{"PUT", "POST", "GET", "DELETE"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"X-CSRF-Token", "Content-Type"},
		ExposedHeaders:   []string{"X-CSRF-Token", "Content-Type"},
		// Enable Debugging for testing, consider disabling in production
		Debug: false,
	})

	r := mux.NewRouter()

	// Authentication Routes
	r.HandleFunc("/api/csrftoken", sessionHandlers.NewCsrfTokenHandler(logger)).Methods("GET")
	r.HandleFunc("/api/login", authenticationHandlers.NewLoginHandler(logger, sessionsStore, db)).Methods("POST")
	r.HandleFunc("/api/logout", authenticationHandlers.NewLogoutHandler(logger, sessionsStore)).Methods("POST")

	// Session Routes
	r.HandleFunc("/api/validsession", sessionHandlers.NewValidSessionHandler(logger, sessionsStore)).Methods("GET")

	// Member CRUD routes
	{{#if signupEnabled}}
	r.HandleFunc("/api/members", membersHandlers.NewSignupHandler(logger, db, sessionsStore)).Methods("POST")
	{{/if}}
	r.HandleFunc("/api/members/email", membersHandlers.NewUpdateEmailHandler(logger, db, sessionsStore)).Methods("PUT")
	r.HandleFunc("/api/members/name", membersHandlers.NewUpdateNameHandler(logger, db, sessionsStore)).Methods("PUT")
	r.HandleFunc("/api/members/password", membersHandlers.NewUpdatePasswordHandler(logger, db, sessionsStore)).Methods("PUT")

	// Superuser sections
	r.HandleFunc("/api/admin/members", membersHandlers.NewListHandler(logger, db, sessionsStore)).Methods("GET")
	r.HandleFunc("/api/admin/members", membersHandlers.NewCreateHandler(logger, db, sessionsStore)).Methods("POST")
	r.HandleFunc("/api/admin/members/{id:[0-9]+}", membersHandlers.NewUpdateHandler(logger, db, sessionsStore)).Methods("PUT")
	r.HandleFunc("/api/admin/members/{id:[0-9]+}/password", membersHandlers.NewResetPasswordHandler(logger, db, sessionsStore)).Methods("PUT")
	r.HandleFunc("/api/admin/members/{id:[0-9]+}", membersHandlers.NewDeleteHandler(logger, db, sessionsStore)).Methods("DELETE")
	r.HandleFunc("/api/admin/server-info", serverHandlers.NewInfoHandler(logger, sessionsStore)).Methods("GET")

	// Resources
	spaHandler := assets.NewHandler(logger, sessionsStore)
	r.PathPrefix("/").Handler(spaHandler)

	// Middleware
	n.Use(c)
	n.UseHandler(CSRF(r))

	return db
}
