package app

import (
	articles "cms/articles/db"
	"cms/config"
	"cms/core/database"
	"cms/core/migrations"
	"cms/members"
	log "github.com/sirupsen/logrus"
	"gopkg.in/AlecAivazis/survey.v1"
)

func CreateSuperUser() {
	config.LoadConfiguration("config.cfg")
	db, err := database.ConnToDB(config.DataBase.GetURL())
	if nil != err {
		log.Fatalf("Fail connect to db: %s", err)
	}
	migrator, err := migrations.NewMigrator(db)
	if nil != err {
		log.Fatalf("Fail create migrator: %s", err)
	}
	for _, migrationByModule := range [][]migrations.Migration{
		members.Migrations(db),
		articles.Migrations(db),
	} {
		err = migrator.Apply(migrationByModule...)
		if nil != err {
			log.Fatalf("Fail apply migrations: %s", err)
		}
	}

	var email string
	err = survey.AskOne(&survey.Input{
		Message: "Email:",
	}, &email, nil)
	if nil != err {
		log.WithError(err).Fatal("can't get email")
	}
	var name string
	err = survey.AskOne(&survey.Input{
		Message: "Name:",
	}, &name, nil)
	if nil != err {
		log.WithError(err).Fatal("can't get name")
	}
	var password string
	err = survey.AskOne(&survey.Password{
		Message: "Password:",
	}, &password, nil)
	if nil != err {
		log.WithError(err).Fatal("can't get password")
	}

	hashedPw, err := members.HashPassword(password)
	if nil != err {
		log.WithError(err).Fatal("can't hash password")
	}

	err = members.CreateMember(db, &members.MemberData{
		Name:        name,
		Email:       email,
		Password:    hashedPw,
		IsSuperuser: true,
	})
	if nil == err {
		log.Info("Superuser created: ", name)
	} else {
		log.WithError(err).Fatal("can't create superuser")
	}
}
