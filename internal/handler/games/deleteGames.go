package games

import (
	"errors"
	"net/http"
	gameErrors "steam_games/internal/errors/games"
	"steam_games/internal/models/responses"
	"steam_games/internal/utils"
	"strconv"
)

func (h *Handler) DeleteGame(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		utils.SendJSONResponse(w, 400, responses.GenericResponse{Success: false, Message: "ID invalide"})
		return
	}
	if err := h.service.DeleteGame(r.Context(), id); err != nil {
		h.gameWriteError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) gameWriteError(w http.ResponseWriter, err error) {
	status := 500
	message := "Erreur interne"
	if errors.Is(err, gameErrors.ErrGameNotFound) {
		status = 404
		message = "Jeu introuvable"
	}
	if errors.Is(err, gameErrors.ErrInvalidGameID) {
		status = 400
		message = "ID invalide"
	}
	if errors.Is(err, gameErrors.ErrInvalidGameData) {
		status = 400
		message = "Données invalides"
	}
	utils.SendJSONResponse(w, status, responses.GenericResponse{Success: false, Message: message})
}
