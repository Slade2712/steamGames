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
	userHandler "steam_games/internal/handler/users"
	authMiddleware "steam_games/internal/middleware/auth"
	gameRepository "steam_games/internal/repositories/games"
	gameService "steam_games/internal/services/games"
	middlewareCors "steam_games/internal/middleware/cors"
	middlewareLogs "steam_games/internal/middleware/logs"
	userRepository "steam_games/internal/repositories/users"
	authService "steam_games/internal/services/auth"
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

	// ======================
	// MIDDLEWARES
	// ======================
	corsHandler := middlewareCors.Middleware(router)
	loggedHandler := middlewareLogs.Logger(corsHandler)
	authGuard := authMiddleware.New(authServ)

	// AUTHENTIFICATION
	router.HandleFunc("POST "+prefix+"register", userHand.Register)
	router.HandleFunc("POST "+prefix+"login", userHand.Login)
	router.HandleFunc("GET "+prefix+"users/me", authGuard.RequireLogin(userHand.Me))
	router.HandleFunc("GET "+prefix+"users", authGuard.RequireAdmin(userHand.GetAll))
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
