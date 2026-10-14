package handler

import (
	"net/http"
	"steam_games/internal/utils"
)

type HealthCheck struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Health(response http.ResponseWriter, request *http.Request) {
	utils.SendJSONResponse(response, http.StatusOK, HealthCheck{
		Success: true,
		Message: "OK",
	})
}
