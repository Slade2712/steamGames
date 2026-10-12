package database

import (
	"database/sql"
	"fmt"
)

func Initialize(db *sql.DB) error {
	queries := []string{
		`
		CREATE TABLE IF NOT EXISTS games (
			id INTEGER PRIMARY KEY,
			title TEXT NOT NULL,
			required_age VARCHAR(5),
			is_free BOOLEAN,
			description TEXT,
			image TEXT,
			website TEXT,
			price NUMERIC(10,2),
			metacritic_score INTEGER,
			recommendations INTEGER,
			release_date DATE,
			release_year INTEGER
		);
		`,

		`
		CREATE TABLE IF NOT EXISTS developers (
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);
		`,

		`
		CREATE TABLE IF NOT EXISTS publishers (
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);
		`,

		`
		CREATE TABLE IF NOT EXISTS platforms (
			id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
			name TEXT NOT NULL UNIQUE
		);
		`,

		`
		CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY,
			description TEXT NOT NULL
		);
		`,

		`
		CREATE TABLE IF NOT EXISTS genres (
			id INTEGER PRIMARY KEY,
			description TEXT NOT NULL
		);
		`,

		`
		CREATE TABLE IF NOT EXISTS game_developers (
			game_id INTEGER REFERENCES games(id),
			developer_id INTEGER REFERENCES developers(id),
			PRIMARY KEY (game_id, developer_id)
		);
		`,

		`
		CREATE TABLE IF NOT EXISTS game_publishers (
			game_id INTEGER REFERENCES games(id),
			publisher_id INTEGER REFERENCES publishers(id),
			PRIMARY KEY (game_id, publisher_id)
		);
		`,

		`
		CREATE TABLE IF NOT EXISTS game_platforms (
			game_id INTEGER REFERENCES games(id),
			platform_id INTEGER REFERENCES platforms(id),
			PRIMARY KEY (game_id, platform_id)
		);
		`,

		`
		CREATE TABLE IF NOT EXISTS game_categories (
			game_id INTEGER REFERENCES games(id),
			category_id INTEGER REFERENCES categories(id),
			PRIMARY KEY (game_id, category_id)
		);
		`,

		`
		CREATE TABLE IF NOT EXISTS game_genres (
			game_id INTEGER REFERENCES games(id),
			genre_id INTEGER REFERENCES genres(id),
			PRIMARY KEY (game_id, genre_id)
		);
		`,
	}

	for _, query := range queries {
		_, execError := db.Exec(query)
		if execError != nil {
			return fmt.Errorf("Erreur initialisation database : %w", execError)
		}
	}

	return nil
}
