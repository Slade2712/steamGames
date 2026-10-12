package database

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=Europe/Paris",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	db, openError := sql.Open("postgres", dsn)
	if openError != nil {
		return nil, fmt.Errorf("Ouverture de la base de données : %w", openError)
	}

	pingError := db.Ping()
	if pingError != nil {
		db.Close()
		return nil, fmt.Errorf("Connexion à la base de données : %w", pingError)
	}

	return db, nil
}
