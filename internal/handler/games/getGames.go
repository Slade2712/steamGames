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

// GetGames récupère une liste paginée de jeux.
//
// @Summary Lister les jeux
// @Description Retourne une liste paginée de jeux Steam.
// @Description Les résultats peuvent être filtrés par titre, année de sortie, plateforme et genre.
// @Description Le tri est effectué sur le titre du jeu.
// @Description Si limit est supérieur à 100, sa valeur est automatiquement limitée à 100.
// @ID get-games
// @Tags Games
// @Produce json
//
// @Param limit query int false "Nombre maximum de jeux retournés" default(10) example(10)
// @Param skip query int false "Nombre de jeux à ignorer avant de commencer la liste" default(0) example(0)
// @Param ordering query string false "Ordre de tri par titre" Enums(title,-title) default(title) example(title)
// @Param title query string false "Recherche partielle dans le titre du jeu, insensible à la casse" example(Subnautica)
// @Param release_year query int false "Filtrer les jeux par année de sortie" example(2013)
// @Param platforms query string false "Filtrer les jeux par nom de plateforme" example(linux)
// @Param genre query string false "Filtrer par un ou plusieurs IDs de genres séparés par des virgules" example(2,25)
//
// @Success 200 {object} models.GamesResponse "Liste des jeux récupérée avec succès"
// @Failure 400 {object} responses.GenericResponse "Paramètre de requête invalide"
// @Failure 500 {object} responses.GenericResponse "Erreur interne lors de la récupération des jeux"
// @Router /games [get]
func (h *Handler) GetGames(response http.ResponseWriter,
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

// GetGameByID récupère un jeu à partir de son identifiant Steam.
//
// @Summary Récupérer un jeu par son ID
// @Description Retourne toutes les informations disponibles pour un jeu Steam.
// @Description La réponse contient notamment les développeurs, éditeurs, plateformes, catégories et genres associés au jeu.
// @ID get-game-by-id
// @Tags Games
// @Produce json
//
// @Param id path int true "Identifiant Steam du jeu" example(570)
//
// @Success 200 {object} models.Game "Jeu récupéré avec succès"
// @Failure 400 {object} responses.GenericResponse "L'identifiant fourni n'est pas un entier valide"
// @Failure 404 {object} responses.GenericResponse "Aucun jeu correspondant à cet ID"
// @Failure 500 {object} responses.GenericResponse "Erreur interne lors de la récupération du jeu"
// @Router /games/{id} [get]
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
