package games

import (
	"context"
	"database/sql"
	gameErrors "steam_games/internal/errors/games"
	models "steam_games/internal/models/games"
)

func (r *Repository) UpdateGame(ctx context.Context, g models.Game) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, `UPDATE games SET title=$2,required_age=$3,is_free=$4,description=$5,image=$6,website=$7,price=$8,metacritic_score=$9,recommendations=$10,release_date=$11,release_year=$12 WHERE id=$1`, g.ID, g.Title, g.RequiredAge, g.IsFree, g.Description, g.Image, g.Website, g.Price, g.MetacriticScore, g.Recommendations, g.ReleaseDate, g.ReleaseYear)
	if err != nil {
		return err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return gameErrors.ErrGameNotFound
	}

	for _, table := range []string{"game_developers", "game_publishers", "game_platforms", "game_categories", "game_genres"} {
		if _, err := tx.ExecContext(ctx, "DELETE FROM "+table+" WHERE game_id=$1", g.ID); err != nil {
			return err
		}
	}

	if err := saveGameAssociations(ctx, tx, g.ID, g); err != nil {
		return err
	}
	return tx.Commit()
}

func saveNames(ctx context.Context, tx *sql.Tx, gameID int, table string, linkTable string, idColumn string, names []string) error {
	for _, name := range names {
		var id int
		err := tx.QueryRowContext(ctx, "INSERT INTO "+table+"(name) VALUES($1) ON CONFLICT(name) DO UPDATE SET name=EXCLUDED.name RETURNING id", name).Scan(&id)
		if err != nil {
			return err
		}

		_, err = tx.ExecContext(ctx, "INSERT INTO "+linkTable+"(game_id,"+idColumn+") VALUES($1,$2) ON CONFLICT DO NOTHING", gameID, id)
		if err != nil {
			return err
		}
	}
	return nil
}

func saveGameAssociations(ctx context.Context, tx *sql.Tx, gameID int, g models.Game) error {
	if err := saveNames(ctx, tx, gameID, "developers", "game_developers", "developer_id", g.Developers); err != nil {
		return err
	}

	if err := saveNames(ctx, tx, gameID, "publishers", "game_publishers", "publisher_id", g.Publishers); err != nil {
		return err
	}

	if err := saveNames(ctx, tx, gameID, "platforms", "game_platforms", "platform_id", g.Platforms); err != nil {
		return err
	}

	for _, c := range g.Categories {
		if _, err := tx.ExecContext(ctx, "INSERT INTO categories(id,description) VALUES($1,$2) ON CONFLICT(id) DO NOTHING", c.ID, c.Description); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO game_categories(game_id,category_id) VALUES($1,$2) ON CONFLICT DO NOTHING", gameID, c.ID); err != nil {
			return err
		}
	}
	for _, genre := range g.Genres {
		if _, err := tx.ExecContext(ctx, "INSERT INTO genres(id,description) VALUES($1,$2) ON CONFLICT(id) DO NOTHING", genre.ID, genre.Description); err != nil {
			return err
		}
		if _, err := tx.ExecContext(ctx, "INSERT INTO game_genres(game_id,genre_id) VALUES($1,$2) ON CONFLICT DO NOTHING", gameID, genre.ID); err != nil {
			return err
		}
	}
	return nil
}
