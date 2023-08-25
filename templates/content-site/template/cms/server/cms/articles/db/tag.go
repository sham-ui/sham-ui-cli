package db

import (
	"database/sql"
	"fmt"
)

type Tag struct {
	Name string
	Slug string
}

type TagRepository struct {
	db *sql.DB
}

func (r *TagRepository) IsUniqueTag(slug string) (bool, error) {
	var existingName string
	row := r.db.QueryRow("SELECT id FROM tag WHERE slug = $1", slug)
	err := row.Scan(&existingName)
	if err == sql.ErrNoRows {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("select id: %s", err)
	}
	return false, nil
}

// CreateTag creates the new tag record
func (r *TagRepository) CreateTag(d *Tag) error {
	_, err := r.db.Exec("INSERT INTO tag(name, slug) VALUES ($1,$2)", d.Name, d.Slug)
	if nil != err {
		return fmt.Errorf("insert into tag: %s", err)
	}
	return nil
}

// UpdateTag update tag record
func (r *TagRepository) UpdateTag(id string, d *Tag) error {
	_, err := r.db.Exec("UPDATE tag SET name = $2, slug = $3 WHERE id = $1", id, d.Name, d.Slug)
	if nil != err {
		return fmt.Errorf("update tag: %s", err)
	}
	return nil
}

// DeleteTag delete tag record
func (r *TagRepository) DeleteTag(id string) error {
	_, err := r.db.Exec("DELETE FROM tag WHERE id = $1", id)
	if nil != err {
		return fmt.Errorf("delete from tag: %s", err)
	}
	return nil
}

func (r *TagRepository) GetOrCreateTag(tx *sql.Tx, tag Tag) (int, error) {
	var id int
	err := tx.QueryRow("SELECT id FROM tag WHERE slug = $1", tag.Slug).Scan(&id)
	if nil != err {
		if err == sql.ErrNoRows {
			err := tx.QueryRow("INSERT INTO tag(name, slug) VALUES ($1,$2) RETURNING id", tag.Name, tag.Slug).Scan(&id)
			if nil != err {
				return 0, fmt.Errorf("insert tag: %s", err)
			}
		} else {
			return 0, fmt.Errorf("query tag: %s", err)
		}
	}
	return id, nil
}

func NewTagRepository(db *sql.DB) *TagRepository {
	return &TagRepository{db: db}
}
