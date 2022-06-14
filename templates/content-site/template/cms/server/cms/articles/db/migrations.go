package db

import (
	"cms/core/migrations"
	"database/sql"
)

func Migrations(db *sql.DB) []migrations.Migration {
	return []migrations.Migration{
		migrations.NewDBMigration(db, "0001_articles_initial", `
			CREATE TABLE category(
				"id" SERIAL UNIQUE PRIMARY KEY,
				"name" varchar(100),
				"slug" varchar(100)
			);
			ALTER TABLE ONLY category ADD CONSTRAINT category_name_unique UNIQUE ("name");
			ALTER TABLE ONLY category ADD CONSTRAINT category_slug_unique UNIQUE ("slug");
			CREATE UNIQUE INDEX category_slug_index ON category (slug);

			CREATE TABLE tag(
				"id" SERIAL UNIQUE PRIMARY KEY,
				"name" varchar(100),
				"slug" varchar(100)
			);
			ALTER TABLE ONLY tag ADD CONSTRAINT tag_name_unique UNIQUE ("name");
			ALTER TABLE ONLY tag ADD CONSTRAINT tag_slug_unique UNIQUE ("slug");
			CREATE UNIQUE INDEX tag_slug_index ON tag (slug);

			CREATE TABLE article(
				"id" SERIAL UNIQUE PRIMARY KEY,
				"title" varchar(100),
				"slug" varchar(100),
				"category_id" INTEGER REFERENCES category(id) ON DELETE RESTRICT,
				"short_body" text,
				"body" text,
				"published_at" TIMESTAMPTZ NOT NULL DEFAULT NOW()
			);
			ALTER TABLE ONLY article ADD CONSTRAINT article_name_unique UNIQUE ("title");
			ALTER TABLE ONLY article ADD CONSTRAINT article_slug_unique UNIQUE ("slug");
			CREATE UNIQUE INDEX article_slug_index ON article (slug);
			CREATE INDEX article_category_id_index ON article (category_id);

			CREATE TABLE article_tag(
				"id" SERIAL UNIQUE PRIMARY KEY,
				"tag_id" INTEGER REFERENCES tag(id) ON DELETE CASCADE,
				"article_id" INTEGER REFERENCES article(id) ON DELETE CASCADE
			);
			ALTER TABLE ONLY article_tag ADD CONSTRAINT article_tag_tag_id_article_id_unique UNIQUE ("tag_id", "article_id");
			CREATE INDEX article_tag_tag_id ON article_tag (tag_id);
			CREATE INDEX article_tag_article_id ON article_tag (article_id);
		`),
	}
}
