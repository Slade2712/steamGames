package services

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"steam_games/internal/models"
	repository "steam_games/external/import_games/repositories"
)

type Service struct {
    repository *repository.Repository
}

func New(repository *repository.Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) Import(
	context context.Context,
	filePath string,
) error {

	file, openError := os.Open(filePath)
	if openError != nil {
		return fmt.Errorf(
			"impossible d'ouvrir le fichier %s : %w",
			filePath,
			openError,
		)
	}
	defer file.Close()

	var games []models.GameImport
	decodeError := json.NewDecoder(file).Decode(&games)
	if decodeError != nil {
		return fmt.Errorf(
			"impossible de décoder le fichier JSON : %w",
			decodeError,
		)
	}

	const batchSize = 1500

	for start := 0; start < len(games); start += batchSize {
		end := start + batchSize

		if end > len(games) {
			end = len(games)
		}

		batch := games[start:end]
		upsertGamesError := s.repository.UpsertGames(
			context,
			batch,
		)
		if upsertGamesError != nil {
			return fmt.Errorf(
				"erreur import jeux %d à %d : %w",
				start+1,
				end,
				upsertGamesError,
			)
		}

		fmt.Printf(
			"Import : %d/%d jeux\n",
			end,
			len(games),
		)
	}

	return nil
}
