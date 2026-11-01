package games

import (
	"errors"
	"net/http"
	gameErrors "steam_games/internal/errors/games"
	"steam_games/internal/models/responses"
	"steam_games/internal/utils"
	"strconv"
)

// DeleteGame supprime un jeu.
//
// @Summary Supprimer un jeu
// @Description Supprime définitivement un jeu à partir de son identifiant.
// @Description Les associations du jeu avec ses développeurs, éditeurs, plateformes, catégories et genres sont également supprimées.
// @Description Cette route est réservée aux administrateurs.
// @ID delete-game
// @Tags Games
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Identifiant du jeu à supprimer" example(570)
//
// @Success 204 "Jeu supprimé avec succès"
// @Failure 400 {object} responses.GenericResponse "Identifiant du jeu invalide"
// @Failure 401 {object} responses.GenericResponse "Token Bearer manquant, invalide ou expiré"
// @Failure 403 {object} responses.GenericResponse "Accès réservé aux administrateurs"
// @Failure 404 {object} responses.GenericResponse "Jeu introuvable"
// @Failure 500 {object} responses.GenericResponse "Erreur interne lors de la suppression du jeu"
// @Router /games/{id} [delete]
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
