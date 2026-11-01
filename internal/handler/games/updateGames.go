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
