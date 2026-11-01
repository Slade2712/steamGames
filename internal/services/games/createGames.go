package games

import (
	"context"
	"strings"

	gameErrors "steam_games/internal/errors/games"
	models "steam_games/internal/models/games"
)

func (s *Service) CreateGame(ctx context.Context, game models.Game) (int, error) {
	if strings.TrimSpace(game.Title) == "" || game.ID != 0 {
		return 0, gameErrors.ErrInvalidGameData
	}
	return s.repository.CreateGame(ctx, game)
}
