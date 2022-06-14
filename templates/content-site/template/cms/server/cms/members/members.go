package members

import (
	"golang.org/x/crypto/bcrypt"
)

type MemberData struct {
	ID, Name, Email, Password string
	IsSuperuser               bool
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
