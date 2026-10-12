package main

import (
	"context"
	"log"
	"os"

	"steam_games/internal/database"
	repository "steam_games/internal/repositories/import_games"
	service "steam_games/internal/services/import_games"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal(
			"Utilisation : go run ./cmd/import <chemin-du-fichier.json>",
		)
	}

	filePath := os.Args[1]
	db, connectionError := database.Connect()
	if connectionError != nil {
		log.Fatal(
			"Erreur de connexion à la base de données : ",
			connectionError,
		)
	}
	defer db.Close()

	initializeError := database.Initialize(db)
	if initializeError != nil {
		log.Fatal(
			"Erreur lors de l'initialisation de la base de données : ",
			initializeError,
		)
	}

	gameRepository := repository.New(db)
	importService := service.New(gameRepository)
	importError := importService.Import(
		context.Background(),
		filePath,
	)

	if importError != nil {
		log.Fatal(
			"Erreur lors de l'import : ",
			importError,
		)
	}

	log.Println("Import terminé avec succès")
}
