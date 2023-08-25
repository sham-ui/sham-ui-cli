package members

import (
	"database/sql"
	"fmt"
)

func CreateMemberStructure(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS members(
			"id" SERIAL UNIQUE PRIMARY KEY,
			"name" varchar(100),
			"email" varchar(100),
			"password" varchar(255),
			"is_superuser" boolean NOT NULL DEFAULT false
		);
		ALTER TABLE ONLY members DROP CONSTRAINT IF EXISTS member_email_unique;
		ALTER TABLE ONLY members ADD CONSTRAINT member_email_unique UNIQUE (email);
	`)
	if err != nil {
		return fmt.Errorf("create table: %s", err)
	}
	return nil
}

func IsUniqueEmailForMemberID(db *sql.DB, id, email string) (bool, error) {
	var existingEmail string
	err := db.QueryRow("SELECT email FROM members WHERE email = $1 AND id != $2", email, id).Scan(&existingEmail)
	if err == sql.ErrNoRows {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("select email: %s", err)
	}
	return false, nil
}

func IsUniqueEmail(db *sql.DB, email string) (bool, error) {
	var existingEmail string
	err := db.QueryRow("SELECT email FROM members WHERE email = $1", email).Scan(&existingEmail)
	if err == sql.ErrNoRows {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("select email: %s", err)
	}
	return false, nil
}

// CreateMember creates the new member record
func CreateMember(db *sql.DB, m *MemberData) error {
	_, err := db.Exec("INSERT INTO members(name, email, password, is_superuser) VALUES ($1,$2, $3, $4)", m.Name, m.Email, m.Password, m.IsSuperuser)
	if nil != err {
		return fmt.Errorf("insert into members: %s", err)
	}
	return nil
}

// UpdateMemberName uses the member ID to update a name
func UpdateMemberName(db *sql.DB, id string, name string) error {
	_, err := db.Exec("UPDATE members SET name = $2 WHERE id = $1", id, name)
	return err
}

// UpdateMemberEmail uses the member ID to update a email
func UpdateMemberEmail(db *sql.DB, id string, email string) error {
	_, err := db.Exec("UPDATE members SET email = $2 WHERE id = $1", id, email)
	return err
}

// UpdateMemberPassword uses the member ID to update a password
func UpdateMemberPassword(db *sql.DB, id string, password string) error {
	_, err := db.Exec("UPDATE members SET password = $2 WHERE id = $1", id, password)
	return err
}

// UpdateMemberData uses the member ID to update a name, email, is_superuser
func UpdateMemberData(db *sql.DB, m *MemberData) error {
	_, err := db.Exec("UPDATE members SET name = $2, email = $3, is_superuser = $4  WHERE id = $1", m.ID, m.Name, m.Email, m.IsSuperuser)
	if nil != err {
		return fmt.Errorf("update member: %s", err)
	}
	return err
}

// HasMemberForID check member for id
func HasMemberForID(db *sql.DB, id string) (bool, error) {
	var existingId string
	err := db.QueryRow("SELECT id FROM members WHERE id = $1", id).Scan(&existingId)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("select id: %s", err)
	}
	return true, nil
}

// DeleteMember uses the member ID to delete
func DeleteMember(db *sql.DB, id string) error {
	_, err := db.Exec("DELETE FROM members WHERE id = $1", id)
	return err
}
