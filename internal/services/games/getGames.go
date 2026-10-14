package games

import (
	"context"
	gameErrors "steam_games/internal/errors/games"
	models "steam_games/internal/models/games"
	gameRepository "steam_games/internal/repositories/games"
)

type Service struct {
	repository *gameRepository.Repository
}

func New(repository *gameRepository.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) GetGameByID(
	context context.Context,
	id int,
) (*models.Game, error) {

	if id <= 0 {
		return nil, gameErrors.ErrInvalidGameID
	}

	game, getError := s.repository.GetGameByID(
		context,
		id,
	)

	if getError != nil {
		return nil, getError
	}

	return game, nil
}
