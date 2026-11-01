package games

import (
	"context"
	gameErrors "steam_games/internal/errors/games"
	models "steam_games/internal/models/games"
	"strings"
)

func (s *Service) UpdateGame(ctx context.Context, game models.Game) error {
	if game.ID <= 0 {
		return gameErrors.ErrInvalidGameID
	}
	if strings.TrimSpace(game.Title) == "" {
		return gameErrors.ErrInvalidGameData
	}
	return s.repository.UpdateGame(ctx, game)
}
