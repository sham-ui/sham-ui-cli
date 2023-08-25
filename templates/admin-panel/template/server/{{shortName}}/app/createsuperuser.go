package app

import (
	"github.com/go-logr/logr"
    "gopkg.in/AlecAivazis/survey.v1"
    "os"
	"{{shortName}}/config"
	"{{shortName}}/core/database"
	"{{shortName}}/members"
)

func CreateSuperUser(logger logr.Logger) {
	config.LoadConfiguration(logger, "config.cfg")
	db, err := database.ConnToDB(config.DataBase.GetURL())
	if nil != err {
		logger.Error(err, "Fail connect to db")
		os.Exit(1)
	}
	err = members.CreateMemberStructure(db)
	if nil != err {
		logger.Error(err, "Fail create members table")
		os.Exit(1)
	}
	logger.Info("Create members table")

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
