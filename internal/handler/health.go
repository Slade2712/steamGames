package handler

import (
	"context"
	"database/sql"
	"net/http"
	"time"

	"steam_games/internal/utils"
)

type HealthCheck struct {
	Database bool   `json:"database"`
	Message  string `json:"message"`
}

// Documentation
type HealthCheckSuccess struct {
	Database bool   `json:"database" example:"true"`
	Message  string `json:"message" example:"OK"`
}

type HealthCheckUnavailable struct {
	Database bool   `json:"database" example:"false"`
	Message  string `json:"message" example:"Database unavailable"`
}

type HealthHandler struct {
	db *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Health vérifie l'état de santé de l'API.
//
// @Summary Vérifier l'état de santé de l'API
// @Description Vérifie que l'API fonctionne et que la connexion à la base de données PostgreSQL est disponible.
// @Description La connexion à PostgreSQL est testée avec un délai maximum de 2 secondes.
// @ID health-check
// @Tags Health
// @Produce json
// @Success 200 {object} HealthCheckSuccess "API opérationnelle et base de données disponible"
// @Failure 503 {object} HealthCheckUnavailable "Base de données indisponible"
// @Router /health [get]
func (handler *HealthHandler) Health(
	response http.ResponseWriter,
	request *http.Request,
) {
	context, cancel := context.WithTimeout(request.Context(), 2*time.Second)
	defer cancel()

	pingError := handler.db.PingContext(context)
	if pingError != nil {
		utils.SendJSONResponse(
			response,
			http.StatusServiceUnavailable,
			HealthCheck{
				Database: false,
				Message:  "Database unavailable",
			},
		)
		return
	}

	utils.SendJSONResponse(
		response,
		http.StatusOK,
		HealthCheck{
			Database: true,
			Message:  "OK",
		},
	)
}
