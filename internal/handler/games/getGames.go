package games

import (
	"net/http"
	"strconv"

	responseModel "steam_games/internal/models/responses"
	gameService "steam_games/internal/services/games"
	"steam_games/internal/utils"
)

type Handler struct {
	service *gameService.Service
}

func New(service *gameService.Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GetGameByID(
	response http.ResponseWriter,
	request *http.Request,
) {
	idString := request.PathValue("id")
	id, convertError := strconv.Atoi(idString)

	if convertError != nil {
		utils.SendJSONResponse(
			response,
			http.StatusBadRequest,

			responseModel.GenericResponse{
				Success: false,
				Message: "ID invalide",
			},
		)
		return
	}

	game, getError := h.service.GetGameByID(
		request.Context(),
		id,
	)

	if getError != nil {
		utils.SendJSONResponse(
			response,
			http.StatusInternalServerError,
			responseModel.GenericResponse{
				Success: false,
				Message: getError.Error(),
			},
		)
		return
	}

	utils.SendJSONResponse(
		response,
		http.StatusOK,
		game,
	)
}
