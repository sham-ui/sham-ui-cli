package app

import (
	articles "cms/articles/db"
	"cms/config"
	"cms/core/database"
	"cms/core/migrations"
	"cms/members"
	"github.com/go-logr/logr"
	"gopkg.in/AlecAivazis/survey.v1"
	"os"
)

func CreateSuperUser(logger logr.Logger) {
	config.LoadConfiguration(logger, "config.cfg")
	db, err := database.ConnToDB(config.DataBase.GetURL())
	if nil != err {
		logger.Error(err, "Fail connect to db")
		os.Exit(1)
	}
	migrator, err := migrations.NewMigrator(logger, db)
	if nil != err {
		logger.Error(err, "Fail create migrator")
		os.Exit(1)
	}
	for _, migrationByModule := range [][]migrations.Migration{
		members.Migrations(db),
		articles.Migrations(db),
	} {
		err = migrator.Apply(migrationByModule...)
		if nil != err {
			logger.Error(err, "Fail apply migrations")
			os.Exit(1)
		}
	}

	var email string
	err = survey.AskOne(&survey.Input{
		Message: "Email:",
	}, &email, nil)
	if nil != err {
		logger.Error(err, "Can't get email")
		os.Exit(1)
	}
	var name string
	err = survey.AskOne(&survey.Input{
		Message: "Name:",
	}, &name, nil)
	if nil != err {
		logger.Error(err, "Can't get name")
		os.Exit(1)
	}
	var password string
	err = survey.AskOne(&survey.Password{
		Message: "Password:",
	}, &password, nil)
	if nil != err {
		logger.Error(err, "Can't get password")
		os.Exit(1)
	}

	hashedPw, err := members.HashPassword(password)
	if nil != err {
		logger.Error(err, "Can't hash password")
		os.Exit(1)
	}

	err = members.CreateMember(db, &members.MemberData{
		Name:        name,
		Email:       email,
		Password:    hashedPw,
		IsSuperuser: true,
	})
	if err != nil {
		logger.Error(err, "Can't create superuser")
		os.Exit(1)
	}
	logger.Info("Superuser created", "email", email, "name", name)
}
