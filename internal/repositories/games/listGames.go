package games

import (
	"context"
	"fmt"

	gameModel "steam_games/internal/models/games"
)

func (r *Repository) ListGames(ctx context.Context, limit int, skip int) ([]gameModel.Game, int, error) {
	var total int
	if err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM games").Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.QueryContext(ctx, "SELECT id FROM games ORDER BY id LIMIT $1 OFFSET $2", limit, skip)
	if err != nil {
		return nil, 0, err
	}

	ids := []int{}
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			rows.Close()
			return nil, 0, err
		}
		ids = append(ids, id)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, 0, err
	}
	rows.Close()

	games := make([]gameModel.Game, 0, len(ids))
	for _, id := range ids {
		game, err := r.GetGameByID(ctx, id)
		if err != nil {
			return nil, 0, fmt.Errorf("lecture du jeu %d : %w", id, err)
		}
		games = append(games, *game)
	}
	return games, total, nil
}
