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

// CreateGame ajoute un jeu. L'ID est attribué par le repository.
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
