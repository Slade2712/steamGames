package database

import (
	"database/sql"
	"fmt"
)

func Initialize(db *sql.DB) error {
	queries := []string{
		``,
	}

	for _, query := range queries {
		_, execError := db.Exec(query)
		if execError != nil {
			return fmt.Errorf("Initialisation de la base de données : %w", execError)
		}
	}

	return nil
}
