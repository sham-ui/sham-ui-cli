package db

import (
	"database/sql"
	"fmt"
)

type Category struct {
	Name string
	Slug string
}

type CategoryRepository struct {
	db *sql.DB
}

func (r *CategoryRepository) IsUniqueCategory(slug string) (bool, error) {
	var existingName string
	err := r.db.QueryRow("SELECT id FROM category WHERE slug = $1", slug).Scan(&existingName)
	if err == sql.ErrNoRows {
		return true, nil
	}
	if err != nil {
		return false, fmt.Errorf("select id: %s", err)
	}
	return false, nil
}

// CreateCategory creates the new category record
func (r *CategoryRepository) CreateCategory(d *Category) error {
	_, err := r.db.Exec("INSERT INTO category(name, slug) VALUES ($1,$2)", d.Name, d.Slug)
	if nil != err {
		return fmt.Errorf("insert into category: %s", err)
	}
	return nil
}

// UpdateCategory update category record
func (r *CategoryRepository) UpdateCategory(id string, d *Category) error {
	_, err := r.db.Exec("UPDATE category SET name = $2, slug = $3 WHERE id = $1", id, d.Name, d.Slug)
	if nil != err {
		return fmt.Errorf("update category: %s", err)
	}
	return nil
}

// DeleteCategory delete category record
func (r *CategoryRepository) DeleteCategory(id string) error {
	_, err := r.db.Exec("DELETE FROM category WHERE id = $1", id)
	if nil != err {
		return fmt.Errorf("delete from category: %s", err)
	}
	return nil
}

func NewCategoryRepository(db *sql.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}
