package games

import (
	"net/http"
	"strconv"

	gameErrors "steam_games/internal/errors/games"
	models "steam_games/internal/models/games"
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

func (h *Handler) GetGames(
	response http.ResponseWriter,
	request *http.Request,
) {
	params := request.URL.Query()

	limit := 10
	skip := 0
	releaseYear := 0

	if params.Get("limit") != "" {
		value, error := strconv.Atoi(params.Get("limit"))
		if error != nil {
			utils.SendJSONResponse(
				response,
				http.StatusBadRequest,
				responseModel.GenericResponse{
					Success: false,
					Message: gameErrors.ErrInvalidLimit.Error(),
				},
			)
			return
		}

		if value > 100 {
			limit = 100
		} else {
			limit = value
		}
	}

	if params.Get("skip") != "" {
		value, error := strconv.Atoi(params.Get("skip"))

		if error != nil {
			utils.SendJSONResponse(
				response,
				http.StatusBadRequest,
				responseModel.GenericResponse{
					Success: false,
					Message: gameErrors.ErrInvalidSkip.Error(),
				},
			)
			return
		}

		skip = value
	}

	if params.Get("release_year") != "" {
		value, error := strconv.Atoi(
			params.Get("release_year"),
		)

		if error != nil {
			utils.SendJSONResponse(
				response,
				http.StatusBadRequest,
				responseModel.GenericResponse{
					Success: false,
					Message: gameErrors.ErrInvalidReleaseYear.Error(),
				},
			)
			return
		}

		releaseYear = value
	}

	ordering := params.Get("ordering")

	if ordering == "" {
		ordering = "title"
	}

	if ordering != "title" && ordering != "-title" {
		utils.SendJSONResponse(
			response,
			http.StatusBadRequest,
			responseModel.GenericResponse{
				Success: false,
				Message: gameErrors.ErrInvalidOrdering.Error(),
			},
		)
		return
	}

	games, totalItems, getError := h.service.GetGames(
		request.Context(),
		limit,
		skip,
		ordering,
		params.Get("title"),
		releaseYear,
		params.Get("platforms"),
		params.Get("genre"),
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

	totalPages := (totalItems + limit - 1) / limit
	page := (skip / limit) + 1

	params.Set("limit", strconv.Itoa(limit))
	params.Set("skip", strconv.Itoa(skip))
	self := request.URL.Path + "?" + params.Encode()
	var next *string
	if skip+limit < totalItems {
		params.Set("skip", strconv.Itoa(skip+limit))
		nextURL := request.URL.Path + "?" + params.Encode()
		next = &nextURL
	}

	result := models.GamesResponse{
		Items:      games,
		Page:       page,
		PageSize:   limit,
		TotalItems: totalItems,
		TotalPages: totalPages,
		Links: models.GamesLinks{
			Self: self,
			Next: next,
		},
	}

	utils.SendJSONResponse(
		response,
		http.StatusOK,
		result,
	)
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
				Message: gameErrors.ErrInvalidGameID.Error(),
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
