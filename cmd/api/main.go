package main

import (
	"log"
	"net/http"
	"os"
	"steam_games/internal/database"
	"steam_games/internal/handler"

	gameHandler "steam_games/internal/handler/games"
	gameRepository "steam_games/internal/repositories/games"
	gameService "steam_games/internal/services/games"
	middlewareLogs "steam_games/internal/middleware/logs"
)

func main() {

	// DB
	db, connectError := database.Connect()
	if connectError != nil {
		log.Fatal(connectError)
	}
	defer db.Close()

	initError := database.Initialize(db)
	if initError != nil {
		log.Fatal(initError)
	}

	prefix := os.Getenv("API_PREFIX")
	log.Println(prefix)
	router := http.NewServeMux()

	// ======================
	// REPOSITORIES
	// ======================
	gameRepo := gameRepository.New(db)

	// ======================
	// SERVICES
	// ======================
	gameServ := gameService.New(gameRepo)

	// ======================
	// HANDLERS
	// ======================
	gameHand := gameHandler.New(gameServ)

	// HEALTH
	router.HandleFunc("GET "+prefix+"health", handler.Health)

	// GET GAME BY ID
	router.HandleFunc("GET "+prefix+"games/{id}", gameHand.GetGameByID)

	// SERVER
	internalPort := os.Getenv("PORT")
	externalPort := os.Getenv("API_PORT")

	log.Printf("API démarrée sur http://localhost:%s", externalPort)

	serveError := http.ListenAndServe(":"+internalPort, middlewareLogs.Logger(router))
	if serveError != nil {
		log.Fatal(serveError)
	}
}
