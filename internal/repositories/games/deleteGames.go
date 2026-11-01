package games

import (
	"context"
	gameErrors "steam_games/internal/errors/games"
)

func (r *Repository) DeleteGame(ctx context.Context, id int) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	for _, table := range []string{"game_developers", "game_publishers", "game_platforms", "game_categories", "game_genres"} {
		if _, err := tx.ExecContext(ctx, "DELETE FROM "+table+" WHERE game_id=$1", id); err != nil {
			return err
		}
	}
	result, err := tx.ExecContext(ctx, "DELETE FROM games WHERE id=$1", id)
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
	return tx.Commit()
}
