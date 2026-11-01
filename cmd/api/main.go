package main

import (
	"log"
	"net/http"
	"os"
	"steam_games/internal/database"
	"steam_games/internal/handler"

	gameHandler "steam_games/internal/handler/games"
	userHandler "steam_games/internal/handler/users"
	authMiddleware "steam_games/internal/middleware/auth"
	middlewareLogs "steam_games/internal/middleware/logs"
	gameRepository "steam_games/internal/repositories/games"
	userRepository "steam_games/internal/repositories/users"
	authService "steam_games/internal/services/auth"
	gameService "steam_games/internal/services/games"
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
	userRepo := userRepository.New(db)

	// ======================
	// SERVICES
	// ======================
	gameServ := gameService.New(gameRepo)
	jwtSecret := os.Getenv("JWT_SECRET")
	authServ, err := authService.New(userRepo, jwtSecret)
	if err != nil {
		log.Fatal(err)
	}

	// ======================
	// HANDLERS
	// ======================
	gameHand := gameHandler.New(gameServ)
	userHand := userHandler.New(authServ)
	authGuard := authMiddleware.New(authServ)

	// AUTHENTIFICATION
	router.HandleFunc("POST "+prefix+"register", userHand.Register)
	router.HandleFunc("POST "+prefix+"login", userHand.Login)
	router.HandleFunc("GET "+prefix+"users/me", authGuard.RequireLogin(userHand.Me))
	router.HandleFunc("GET "+prefix+"users", authGuard.RequireAdmin(userHand.GetAll))

	// HEALTH
	router.HandleFunc("GET "+prefix+"health", handler.Health)

	// GET GAMES
	router.HandleFunc("GET "+prefix+"games", gameHand.GetGames)

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
