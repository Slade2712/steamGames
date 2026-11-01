package main

import (
	"log"
	"net/http"
	"os"
	"steam_games/internal/database"
	"steam_games/internal/handler"

	httpSwagger "github.com/swaggo/http-swagger/v2"

	_ "steam_games/docs"

	gameHandler "steam_games/internal/handler/games"
	gameRepository "steam_games/internal/repositories/games"
	gameService "steam_games/internal/services/games"

	middlewareCors "steam_games/internal/middleware/cors"
	middlewareLogs "steam_games/internal/middleware/logs"
)

// @title Steam Games API
// @version 1.0
// @description API REST pour consulter les jeux Steam.
// @BasePath /api/v1

// @tag.name Health
// @tag.description Vérification de l'état de santé de l'API.

// @tag.name Games
// @tag.description Consultation, recherche, filtrage et pagination des jeux Steam.
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
	healthHand := handler.NewHealthHandler(db)

	// HEALTH
	router.HandleFunc("GET "+prefix+"health", healthHand.Health)

	// SWAGGER
	router.HandleFunc("GET /swagger", func(response http.ResponseWriter, request *http.Request) {
		http.Redirect(response, request, "/swagger/index.html", http.StatusMovedPermanently)
	})
	router.Handle(
		"GET /swagger/",
		httpSwagger.Handler(
			httpSwagger.UIConfig(map[string]string{
				"tagsSorter": `(a, b) => {
				if (a === "Health") return -1;
				if (b === "Health") return 1;
				return a.localeCompare(b);
			}`,
			}),
		),
	)

	// GET GAMES
	router.HandleFunc("GET "+prefix+"games", gameHand.GetGames)

	// GET GAME BY ID
	router.HandleFunc("GET "+prefix+"games/{id}", gameHand.GetGameByID)

	// ======================
	// MIDDLEWARES
	// ======================
	corsHandler := middlewareCors.Middleware(router)
	loggedHandler := middlewareLogs.Logger(corsHandler)

	// ======================
	// SERVER
	// ======================
	internalPort := os.Getenv("PORT")
	externalPort := os.Getenv("API_PORT")

	log.Printf("API démarrée sur http://localhost:%s", externalPort)

	serveError := http.ListenAndServe(":"+internalPort, loggedHandler)
	if serveError != nil {
		log.Fatal(serveError)
	}
}
