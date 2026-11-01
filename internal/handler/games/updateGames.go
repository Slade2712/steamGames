package games

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"

	models "steam_games/internal/models/games"
	"steam_games/internal/models/responses"
	"steam_games/internal/utils"
)

// UpdateGame modifie un jeu existant.
//
// @Summary Modifier un jeu
// @Description Remplace les informations d'un jeu existant à partir de son identifiant.
// @Description Le champ title est obligatoire.
// @Description Si _id est fourni dans le corps, il doit être identique à l'identifiant présent dans l'URL.
// @Description Les associations développeurs, éditeurs, plateformes, catégories et genres sont remplacées par celles fournies.
// @Description Un seul objet JSON est accepté dans le corps de la requête.
// @Description Cette route est réservée aux administrateurs.
// @ID update-game
// @Tags Games
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param id path int true "Identifiant du jeu" example(570)
// @Param game body models.Game true "Nouvelles informations du jeu"
//
// @Success 200 {object} models.Game "Jeu modifié avec succès"
// @Failure 400 {object} responses.GenericResponse "ID invalide, JSON invalide, plusieurs objets fournis, title manquant ou ID du corps différent de l'URL"
// @Failure 401 {object} responses.GenericResponse "Token Bearer manquant, invalide ou expiré"
// @Failure 403 {object} responses.GenericResponse "Accès réservé aux administrateurs"
// @Failure 404 {object} responses.GenericResponse "Jeu introuvable"
// @Failure 500 {object} responses.GenericResponse "Erreur interne lors de la modification du jeu"
// @Router /games/{id} [put]
func (h *Handler) UpdateGame(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil || id <= 0 {
		utils.SendJSONResponse(w, 400, responses.GenericResponse{Success: false, Message: "ID invalide"})
		return
	}
	var input models.Game
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&input); err != nil {
		utils.SendJSONResponse(w, 400, responses.GenericResponse{Success: false, Message: "JSON invalide"})
		return
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		utils.SendJSONResponse(w, 400, responses.GenericResponse{Success: false, Message: "Un seul objet JSON attendu"})
		return
	}
	if input.ID != 0 && input.ID != id {
		utils.SendJSONResponse(w, 400, responses.GenericResponse{Success: false, Message: "ID du corps différent de l’URL"})
		return
	}
	if strings.TrimSpace(input.Title) == "" {
		utils.SendJSONResponse(w, 400, responses.GenericResponse{Success: false, Message: "title obligatoire"})
		return
	}
	input.ID = id
	if err := h.service.UpdateGame(r.Context(), input); err != nil {
		h.gameWriteError(w, err)
		return
	}
	game, err := h.service.GetGameByID(r.Context(), id)
	if err != nil {
		h.gameWriteError(w, err)
		return
	}
	utils.SendJSONResponse(w, 200, game)
}
