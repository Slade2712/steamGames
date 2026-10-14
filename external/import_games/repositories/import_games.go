package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"steam_games/internal/models"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) UpsertGames(
	context context.Context,
	games []models.GameImport,
) error {
	transaction, beginError := r.db.BeginTx(context, nil)
	if beginError != nil {
		return fmt.Errorf(
			"impossible de démarrer la transaction : %w",
			beginError,
		)
	}

	defer transaction.Rollback()

	for _, game := range games {
		if err := upsertGame(context, transaction, game); err != nil {
			return err
		}

		if err := clearGameRelations(
			context,
			transaction,
			game.ID,
		); err != nil {
			return err
		}

		if err := saveDevelopers(
			context,
			transaction,
			game,
		); err != nil {
			return err
		}

		if err := savePublishers(
			context,
			transaction,
			game,
		); err != nil {
			return err
		}

		if err := savePlatforms(
			context,
			transaction,
			game,
		); err != nil {
			return err
		}

		if err := upsertCategories(
			context,
			transaction,
			game,
		); err != nil {
			return err
		}

		if err := upsertGenres(
			context,
			transaction,
			game,
		); err != nil {
			return err
		}
	}

	if commitError := transaction.Commit(); commitError != nil {
		return fmt.Errorf(
			"impossible de valider la transaction : %w",
			commitError,
		)
	}

	return nil
}

func upsertGame(
	context context.Context,
	transaction *sql.Tx,
	game models.GameImport,
) error {
	query := `
		SELECT id
		FROM games
		WHERE id = $1
	`

	var gameID int

	scanError := transaction.QueryRowContext(
		context,
		query,
		int(game.ID),
	).Scan(&gameID)

	if scanError != nil && scanError != sql.ErrNoRows {
		return fmt.Errorf(
			"erreur lors de la recherche du jeu %d : %w",
			int(game.ID),
			scanError,
		)
	}

	if scanError == sql.ErrNoRows {
		query = `
			INSERT INTO games (
				id,
				title,
				required_age,
				is_free,
				description,
				image,
				website,
				price,
				metacritic_score,
				recommendations,
				release_date,
				release_year
			)
			VALUES (
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8,
				$9,
				$10,
				$11,
				$12
			)
		`

		_, execError := transaction.ExecContext(
			context,
			query,
			int(game.ID),
			game.Title,
			requiredAgeToString(game.RequiredAge),
			bool(game.IsFree),
			game.Description,
			game.Image,
			game.Website,
			game.Price,
			game.MetacriticScore,
			game.Recommendations,
			game.ReleaseDate,
			game.ReleaseYear,
		)

		if execError != nil {
			return fmt.Errorf(
				"erreur lors de l'insertion du jeu %d : %w",
				int(game.ID),
				execError,
			)
		}

		return nil
	}

	query = `
		UPDATE games
		SET
			title = $2,
			required_age = $3,
			is_free = $4,
			description = $5,
			image = $6,
			website = $7,
			price = $8,
			metacritic_score = $9,
			recommendations = $10,
			release_date = $11,
			release_year = $12
		WHERE id = $1
	`

	_, execError := transaction.ExecContext(
		context,
		query,
		int(game.ID),
		game.Title,
		requiredAgeToString(game.RequiredAge),
		bool(game.IsFree),
		game.Description,
		game.Image,
		game.Website,
		game.Price,
		game.MetacriticScore,
		game.Recommendations,
		game.ReleaseDate,
		game.ReleaseYear,
	)

	if execError != nil {
		return fmt.Errorf(
			"erreur lors de la mise à jour du jeu %d : %w",
			int(game.ID),
			execError,
		)
	}

	return nil
}

func clearGameRelations(
	context context.Context,
	transaction *sql.Tx,
	gameID int,
) error {
	tables := []string{
		"game_developers",
		"game_publishers",
		"game_platforms",
		"game_categories",
		"game_genres",
	}

	for _, table := range tables {
		query := "DELETE FROM " + table + " WHERE game_id = $1"

		_, execError := transaction.ExecContext(
			context,
			query,
			gameID,
		)

		if execError != nil {
			return fmt.Errorf(
				"impossible de supprimer les anciennes relations %s du jeu %d : %w",
				table,
				gameID,
				execError,
			)
		}
	}

	return nil
}

func saveDevelopers(
	context context.Context,
	transaction *sql.Tx,
	game models.GameImport,
) error {
	for _, developer := range game.Developers {
		var developerID int

		query := `
			SELECT id
			FROM developers
			WHERE name = $1
		`

		scanError := transaction.QueryRowContext(
			context,
			query,
			developer,
		).Scan(&developerID)

		if scanError != nil && scanError != sql.ErrNoRows {
			return fmt.Errorf(
				"erreur lors de la recherche du développeur %s : %w",
				developer,
				scanError,
			)
		}

		if scanError == sql.ErrNoRows {
			query = `
				INSERT INTO developers (name)
				VALUES ($1)
			`

			_, execError := transaction.ExecContext(
				context,
				query,
				developer,
			)

			if execError != nil {
				return fmt.Errorf(
					"erreur lors de l'insertion du développeur %s : %w",
					developer,
					execError,
				)
			}

			query = `
				SELECT id
				FROM developers
				WHERE name = $1
			`

			scanError = transaction.QueryRowContext(
				context,
				query,
				developer,
			).Scan(&developerID)

			if scanError != nil {
				return fmt.Errorf(
					"erreur lors de la récupération de l'id du développeur %s : %w",
					developer,
					scanError,
				)
			}
		}

		query = `
			INSERT INTO game_developers (
				game_id,
				developer_id
			)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`

		_, execError := transaction.ExecContext(
			context,
			query,
			int(game.ID),
			developerID,
		)

		if execError != nil {
			return fmt.Errorf(
				"erreur lors de l'association du développeur %s au jeu %d : %w",
				developer,
				int(game.ID),
				execError,
			)
		}
	}

	return nil
}

func savePublishers(
	context context.Context,
	transaction *sql.Tx,
	game models.GameImport,
) error {
	for _, publisher := range game.Publishers {
		var publisherID int

		query := `
			SELECT id
			FROM publishers
			WHERE name = $1
		`

		scanError := transaction.QueryRowContext(
			context,
			query,
			publisher,
		).Scan(&publisherID)

		if scanError != nil && scanError != sql.ErrNoRows {
			return fmt.Errorf(
				"erreur lors de la recherche de l'éditeur %s : %w",
				publisher,
				scanError,
			)
		}

		if scanError == sql.ErrNoRows {
			query = `
				INSERT INTO publishers (name)
				VALUES ($1)
			`

			_, execError := transaction.ExecContext(
				context,
				query,
				publisher,
			)

			if execError != nil {
				return fmt.Errorf(
					"erreur lors de l'insertion de l'éditeur %s : %w",
					publisher,
					execError,
				)
			}

			query = `
				SELECT id
				FROM publishers
				WHERE name = $1
			`

			scanError = transaction.QueryRowContext(
				context,
				query,
				publisher,
			).Scan(&publisherID)

			if scanError != nil {
				return fmt.Errorf(
					"erreur lors de la récupération de l'id de l'éditeur %s : %w",
					publisher,
					scanError,
				)
			}
		}

		query = `
			INSERT INTO game_publishers (
				game_id,
				publisher_id
			)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`

		_, execError := transaction.ExecContext(
			context,
			query,
			int(game.ID),
			publisherID,
		)

		if execError != nil {
			return fmt.Errorf(
				"erreur lors de l'association de l'éditeur %s au jeu %d : %w",
				publisher,
				int(game.ID),
				execError,
			)
		}
	}

	return nil
}

func savePlatforms(
	context context.Context,
	transaction *sql.Tx,
	game models.GameImport,
) error {
	for _, platform := range game.Platforms {
		var platformID int

		query := `
			SELECT id
			FROM platforms
			WHERE name = $1
		`

		scanError := transaction.QueryRowContext(
			context,
			query,
			platform,
		).Scan(&platformID)

		if scanError != nil && scanError != sql.ErrNoRows {
			return fmt.Errorf(
				"erreur lors de la recherche de la plateforme %s : %w",
				platform,
				scanError,
			)
		}

		if scanError == sql.ErrNoRows {
			query = `
				INSERT INTO platforms (name)
				VALUES ($1)
			`

			_, execError := transaction.ExecContext(
				context,
				query,
				platform,
			)

			if execError != nil {
				return fmt.Errorf(
					"erreur lors de l'insertion de la plateforme %s : %w",
					platform,
					execError,
				)
			}

			query = `
				SELECT id
				FROM platforms
				WHERE name = $1
			`

			scanError = transaction.QueryRowContext(
				context,
				query,
				platform,
			).Scan(&platformID)

			if scanError != nil {
				return fmt.Errorf(
					"erreur lors de la récupération de l'id de la plateforme %s : %w",
					platform,
					scanError,
				)
			}
		}

		query = `
			INSERT INTO game_platforms (
				game_id,
				platform_id
			)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`

		_, execError := transaction.ExecContext(
			context,
			query,
			int(game.ID),
			platformID,
		)

		if execError != nil {
			return fmt.Errorf(
				"erreur lors de l'association de la plateforme %s au jeu %d : %w",
				platform,
				int(game.ID),
				execError,
			)
		}
	}

	return nil
}

func upsertCategories(
	context context.Context,
	transaction *sql.Tx,
	game models.GameImport,
) error {
	for _, category := range game.Categories {
		query := `
			INSERT INTO categories (
				id,
				description
			)
			VALUES ($1, $2)
			ON CONFLICT (id)
			DO UPDATE SET
				description = $2
		`

		_, execError := transaction.ExecContext(
			context,
			query,
			int(category.ID),
			category.Description,
		)

		if execError != nil {
			return fmt.Errorf(
				"impossible d'enregistrer la catégorie %d (%s) : %w",
				int(category.ID),
				category.Description,
				execError,
			)
		}

		query = `
			INSERT INTO game_categories (
				game_id,
				category_id
			)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`

		_, execError = transaction.ExecContext(
			context,
			query,
			int(game.ID),
			int(category.ID),
		)

		if execError != nil {
			return fmt.Errorf(
				"impossible d'associer la catégorie %d au jeu %d : %w",
				int(category.ID),
				int(game.ID),
				execError,
			)
		}
	}

	return nil
}

func upsertGenres(
	context context.Context,
	transaction *sql.Tx,
	game models.GameImport,
) error {
	for _, genre := range game.Genres {
		query := `
			INSERT INTO genres (
				id,
				description
			)
			VALUES ($1, $2)
			ON CONFLICT (id)
			DO UPDATE SET
				description = $2
		`

		_, execError := transaction.ExecContext(
			context,
			query,
			int(genre.ID),
			genre.Description,
		)

		if execError != nil {
			return fmt.Errorf(
				"impossible d'enregistrer le genre %d (%s) : %w",
				int(genre.ID),
				genre.Description,
				execError,
			)
		}

		query = `
			INSERT INTO game_genres (
				game_id,
				genre_id
			)
			VALUES ($1, $2)
			ON CONFLICT DO NOTHING
		`

		_, execError = transaction.ExecContext(
			context,
			query,
			int(game.ID),
			int(genre.ID),
		)

		if execError != nil {
			return fmt.Errorf(
				"impossible d'associer le genre %d au jeu %d : %w",
				int(genre.ID),
				int(game.ID),
				execError,
			)
		}
	}

	return nil
}

func requiredAgeToString(requiredAge any) string {
	switch value := requiredAge.(type) {
	case string:
		return value

	case float64:
		return fmt.Sprintf("%.0f", value)

	case nil:
		return ""

	default:
		return fmt.Sprint(value)
	}
}
