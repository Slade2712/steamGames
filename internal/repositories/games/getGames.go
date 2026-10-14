package games

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strconv"
	"strings"

	gameErrors "steam_games/internal/errors/games"
	models "steam_games/internal/models/games"
)

type Repository struct {
	db *sql.DB
}

func New(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) GetGames(
	context context.Context,
	limit int,
	skip int,
	ordering string,
	title string,
	releaseYear int,
	platform string,
	genre string,
) ([]models.Game, int, error) {

	query := `
		SELECT
			g.id,
			g.title,
			g.required_age,
			g.is_free,
			g.description,
			g.image,
			g.website,
			g.price,
			g.metacritic_score,
			g.recommendations,
			g.release_date,
			g.release_year
		FROM games g
		WHERE 1 = 1
	`

	args := []any{}
	argNumber := 1

	if title != "" {
		query += fmt.Sprintf(
			" AND g.title ILIKE $%d",
			argNumber,
		)

		args = append(
			args,
			"%"+title+"%",
		)

		argNumber++
	}

	if releaseYear != 0 {
		query += fmt.Sprintf(
			" AND g.release_year = $%d",
			argNumber,
		)

		args = append(
			args,
			releaseYear,
		)

		argNumber++
	}

	if platform != "" {
		query += fmt.Sprintf(`
			AND EXISTS (
				SELECT 1
				FROM game_platforms gp
				INNER JOIN platforms p
					ON p.id = gp.platform_id
				WHERE gp.game_id = g.id
				AND p.name = $%d
			)
		`, argNumber)

		args = append(
			args,
			platform,
		)

		argNumber++
	}

	if genre != "" {
		genreStrings := strings.Split(
			genre,
			",",
		)

		genreConditions := []string{}

		for _, genreString := range genreStrings {
			genreID, convertError := strconv.Atoi(
				genreString,
			)

			if convertError != nil {
				return nil, 0, fmt.Errorf(
					"genre invalide : %s",
					genreString,
				)
			}

			genreConditions = append(
				genreConditions,
				fmt.Sprintf(
					"gg.genre_id = $%d",
					argNumber,
				),
			)

			args = append(
				args,
				genreID,
			)

			argNumber++
		}

		query += `
			AND EXISTS (
				SELECT 1
				FROM game_genres gg
				WHERE gg.game_id = g.id
				AND (
					` + strings.Join(genreConditions, " OR ") + `
				)
			)
		`
	}

	countQuery := `
		SELECT COUNT(*)
		FROM (
			` + query + `
		) AS filtered_games
	`

	if ordering == "-title" {
		query += " ORDER BY g.title DESC"
	} else {
		query += " ORDER BY g.title ASC"
	}

	query += fmt.Sprintf(
		" LIMIT $%d OFFSET $%d",
		argNumber,
		argNumber+1,
	)

	queryArgs := append(
		[]any{},
		args...,
	)

	queryArgs = append(
		queryArgs,
		limit,
		skip,
	)

	var totalItems int

	countError := r.db.QueryRowContext(
		context,
		countQuery,
		args...,
	).Scan(&totalItems)

	if countError != nil {
		return nil, 0, fmt.Errorf(
			"erreur lors du comptage des jeux : %w",
			countError,
		)
	}

	rows, queryError := r.db.QueryContext(
		context,
		query,
		queryArgs...,
	)

	if queryError != nil {
		return nil, 0, fmt.Errorf(
			"erreur lors de la récupération des jeux : %w",
			queryError,
		)
	}

	defer rows.Close()

	games := []models.Game{}
	for rows.Next() {
		var game models.Game

		scanError := rows.Scan(
			&game.ID,
			&game.Title,
			&game.RequiredAge,
			&game.IsFree,
			&game.Description,
			&game.Image,
			&game.Website,
			&game.Price,
			&game.MetacriticScore,
			&game.Recommendations,
			&game.ReleaseDate,
			&game.ReleaseYear,
		)

		if scanError != nil {
			return nil, 0, fmt.Errorf(
				"erreur lors du scan d'un jeu : %w",
				scanError,
			)
		}

		developers, developersError := r.getDevelopers(
			context,
			game.ID,
		)

		if developersError != nil {
			return nil, 0, developersError
		}
		game.Developers = developers

		publishers, publishersError := r.getPublishers(
			context,
			game.ID,
		)

		if publishersError != nil {
			return nil, 0, publishersError
		}
		game.Publishers = publishers

		platforms, platformsError := r.getPlatforms(
			context,
			game.ID,
		)

		if platformsError != nil {
			return nil, 0, platformsError
		}
		game.Platforms = platforms

		categories, categoriesError := r.getCategories(
			context,
			game.ID,
		)

		if categoriesError != nil {
			return nil, 0, categoriesError
		}
		game.Categories = categories

		genres, genresError := r.getGenres(
			context,
			game.ID,
		)

		if genresError != nil {
			return nil, 0, genresError
		}
		game.Genres = genres

		games = append(
			games,
			game,
		)
	}

	if rowsError := rows.Err(); rowsError != nil {
		return nil, 0, rowsError
	}

	return games, totalItems, nil
}

func (r *Repository) GetGameByID(
	context context.Context,
	id int,
) (*models.Game, error) {

	query := `
		SELECT
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
		FROM games
		WHERE id = $1
	`

	var game models.Game

	scanError := r.db.QueryRowContext(
		context,
		query,
		id,
	).Scan(
		&game.ID,
		&game.Title,
		&game.RequiredAge,
		&game.IsFree,
		&game.Description,
		&game.Image,
		&game.Website,
		&game.Price,
		&game.MetacriticScore,
		&game.Recommendations,
		&game.ReleaseDate,
		&game.ReleaseYear,
	)

	if scanError != nil {
		if errors.Is(scanError, sql.ErrNoRows) {
			return nil, gameErrors.ErrGameNotFound
		}

		return nil, fmt.Errorf(
			"erreur lors de la récupération du jeu %d : %w",
			id,
			scanError,
		)
	}

	developers, developersError := r.getDevelopers(context, id)
	if developersError != nil {
		return nil, developersError
	}

	game.Developers = developers

	publishers, publishersError := r.getPublishers(context, id)
	if publishersError != nil {
		return nil, publishersError
	}

	game.Publishers = publishers

	platforms, platformsError := r.getPlatforms(context, id)
	if platformsError != nil {
		return nil, platformsError
	}

	game.Platforms = platforms

	categories, categoriesError := r.getCategories(context, id)
	if categoriesError != nil {
		return nil, categoriesError
	}

	game.Categories = categories

	genres, genresError := r.getGenres(context, id)
	if genresError != nil {
		return nil, genresError
	}

	game.Genres = genres

	return &game, nil
}

func (r *Repository) getDevelopers(
	context context.Context,
	gameID int,
) ([]string, error) {

	query := `
		SELECT developers.name
		FROM developers
		INNER JOIN game_developers
			ON developers.id = game_developers.developer_id
		WHERE game_developers.game_id = $1
		ORDER BY developers.name
	`

	rows, queryError := r.db.QueryContext(
		context,
		query,
		gameID,
	)

	if queryError != nil {
		return nil, fmt.Errorf(
			"erreur lors de la récupération des développeurs du jeu %d : %w",
			gameID,
			queryError,
		)
	}

	defer rows.Close()

	developers := []string{}

	for rows.Next() {
		var developer string

		scanError := rows.Scan(&developer)
		if scanError != nil {
			return nil, fmt.Errorf(
				"erreur lors du scan d'un développeur : %w",
				scanError,
			)
		}

		developers = append(developers, developer)
	}

	if rowsError := rows.Err(); rowsError != nil {
		return nil, rowsError
	}

	return developers, nil
}

func (r *Repository) getPublishers(
	context context.Context,
	gameID int,
) ([]string, error) {

	query := `
		SELECT publishers.name
		FROM publishers
		INNER JOIN game_publishers
			ON publishers.id = game_publishers.publisher_id
		WHERE game_publishers.game_id = $1
		ORDER BY publishers.name
	`

	rows, queryError := r.db.QueryContext(
		context,
		query,
		gameID,
	)

	if queryError != nil {
		return nil, fmt.Errorf(
			"erreur lors de la récupération des éditeurs du jeu %d : %w",
			gameID,
			queryError,
		)
	}

	defer rows.Close()

	publishers := []string{}

	for rows.Next() {
		var publisher string

		scanError := rows.Scan(&publisher)
		if scanError != nil {
			return nil, fmt.Errorf(
				"erreur lors du scan d'un éditeur : %w",
				scanError,
			)
		}

		publishers = append(publishers, publisher)
	}

	if rowsError := rows.Err(); rowsError != nil {
		return nil, rowsError
	}

	return publishers, nil
}

func (r *Repository) getPlatforms(
	context context.Context,
	gameID int,
) ([]string, error) {

	query := `
		SELECT platforms.name
		FROM platforms
		INNER JOIN game_platforms
			ON platforms.id = game_platforms.platform_id
		WHERE game_platforms.game_id = $1
		ORDER BY platforms.name
	`

	rows, queryError := r.db.QueryContext(
		context,
		query,
		gameID,
	)

	if queryError != nil {
		return nil, fmt.Errorf(
			"erreur lors de la récupération des plateformes du jeu %d : %w",
			gameID,
			queryError,
		)
	}

	defer rows.Close()

	platforms := []string{}

	for rows.Next() {
		var platform string

		scanError := rows.Scan(&platform)
		if scanError != nil {
			return nil, fmt.Errorf(
				"erreur lors du scan d'une plateforme : %w",
				scanError,
			)
		}

		platforms = append(platforms, platform)
	}

	if rowsError := rows.Err(); rowsError != nil {
		return nil, rowsError
	}

	return platforms, nil
}

func (r *Repository) getCategories(
	context context.Context,
	gameID int,
) ([]models.Category, error) {

	query := `
		SELECT
			categories.id,
			categories.description
		FROM categories
		INNER JOIN game_categories
			ON categories.id = game_categories.category_id
		WHERE game_categories.game_id = $1
		ORDER BY categories.id
	`

	rows, queryError := r.db.QueryContext(
		context,
		query,
		gameID,
	)

	if queryError != nil {
		return nil, fmt.Errorf(
			"erreur lors de la récupération des catégories du jeu %d : %w",
			gameID,
			queryError,
		)
	}

	defer rows.Close()

	categories := []models.Category{}

	for rows.Next() {
		var category models.Category

		scanError := rows.Scan(
			&category.ID,
			&category.Description,
		)

		if scanError != nil {
			return nil, fmt.Errorf(
				"erreur lors du scan d'une catégorie : %w",
				scanError,
			)
		}

		categories = append(categories, category)
	}

	if rowsError := rows.Err(); rowsError != nil {
		return nil, rowsError
	}

	return categories, nil
}

func (r *Repository) getGenres(
	context context.Context,
	gameID int,
) ([]models.Genre, error) {

	query := `
		SELECT
			genres.id,
			genres.description
		FROM genres
		INNER JOIN game_genres
			ON genres.id = game_genres.genre_id
		WHERE game_genres.game_id = $1
		ORDER BY genres.id
	`

	rows, queryError := r.db.QueryContext(
		context,
		query,
		gameID,
	)

	if queryError != nil {
		return nil, fmt.Errorf(
			"erreur lors de la récupération des genres du jeu %d : %w",
			gameID,
			queryError,
		)
	}

	defer rows.Close()

	genres := []models.Genre{}

	for rows.Next() {
		var genre models.Genre

		scanError := rows.Scan(
			&genre.ID,
			&genre.Description,
		)

		if scanError != nil {
			return nil, fmt.Errorf(
				"erreur lors du scan d'un genre : %w",
				scanError,
			)
		}

		genres = append(genres, genre)
	}

	if rowsError := rows.Err(); rowsError != nil {
		return nil, rowsError
	}

	return genres, nil
}
