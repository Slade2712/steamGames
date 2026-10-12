package main

import (
	"log"
	"net/http"
	"os"
	"steam_games/internal/database"
	"steam_games/internal/handler"
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

	// HANDLER //
	// HEALTH
	router.HandleFunc("GET "+prefix+"health", handler.Health)

	// SERVER
	internalPort := os.Getenv("PORT")
	externalPort := os.Getenv("API_PORT")

	log.Printf("API démarrée sur http://localhost:%s", externalPort)

	serveError := http.ListenAndServe(":"+internalPort, router)
	if serveError != nil {
		log.Fatal(serveError)
	}
}
