package games

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	models "steam_games/internal/models/games"
	"steam_games/internal/models/responses"
	"steam_games/internal/utils"
)

// CreateGame crée un nouveau jeu.
//
// @Summary Créer un jeu
// @Description Crée un nouveau jeu Steam.
// @Description L'identifiant _id ne doit pas être fourni : il est généré automatiquement.
// @Description Le champ title est obligatoire.
// @Description Un seul objet JSON est accepté dans le corps de la requête.
// @Description Cette route est réservée aux administrateurs.
// @ID create-game
// @Tags Games
// @Accept json
// @Produce json
// @Security BearerAuth
//
// @Param game body models.Game true "Informations du jeu à créer"
//
// @Success 201 {object} models.Game "Jeu créé avec succès"
// @Failure 400 {object} responses.GenericResponse "JSON invalide, plusieurs objets fournis, title manquant ou _id fourni"
// @Failure 401 {object} responses.GenericResponse "Token Bearer manquant, invalide ou expiré"
// @Failure 403 {object} responses.GenericResponse "Accès réservé aux administrateurs"
// @Failure 500 {object} responses.GenericResponse "Erreur interne lors de la création du jeu"
// @Router /games [post]
func (h *Handler) CreateGame(w http.ResponseWriter, r *http.Request) {
	var game models.Game
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&game); err != nil {
		utils.SendJSONResponse(w, 400, responses.GenericResponse{Success: false, Message: "JSON invalide"})
		return
	}

	if err := decoder.Decode(new(any)); err != io.EOF {
		utils.SendJSONResponse(w, 400, responses.GenericResponse{Success: false, Message: "Un seul objet JSON attendu"})
		return
	}

	if game.ID != 0 || strings.TrimSpace(game.Title) == "" {
		utils.SendJSONResponse(w, 400, responses.GenericResponse{Success: false, Message: "Ne pas fournir _id ; title est obligatoire"})
		return
	}

	id, err := h.service.CreateGame(r.Context(), game)
	if err != nil {
		h.gameWriteError(w, err)
		return
	}

	created, err := h.service.GetGameByID(r.Context(), id)
	if err != nil {
		h.gameWriteError(w, err)
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, created)
}
