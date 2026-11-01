package games

import (
	"context"
	gameErrors "steam_games/internal/errors/games"
)

func (s *Service) DeleteGame(ctx context.Context, id int) error {
	if id <= 0 {
		return gameErrors.ErrInvalidGameID
	}
	return s.repository.DeleteGame(ctx, id)
}
