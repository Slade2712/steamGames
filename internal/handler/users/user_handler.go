package users

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	authMiddleware "steam_games/internal/middleware/auth"
	userModel "steam_games/internal/models/users"
	"steam_games/internal/models/responses"
	authService "steam_games/internal/services/auth"
	"steam_games/internal/utils"
)

type Handler struct {
	service *authService.Service
}

func New(service *authService.Service) *Handler {
	return &Handler{service: service}
}

func readJSON(w http.ResponseWriter, r *http.Request, value any) error {
	r.Body = http.MaxBytesReader(w, r.Body, 16*1024)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(value); err != nil {
		return err
	}
	if decoder.Decode(new(any)) != io.EOF {
		return errors.New("un seul objet JSON est autorisé")
	}
	return nil
}

func sendError(w http.ResponseWriter, status int, message string) {
	utils.SendJSONResponse(w, status, responses.GenericResponse{Success: false, Message: message})
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var data userModel.RegisterRequest
	if err := readJSON(w, r, &data); err != nil {
		sendError(w, http.StatusBadRequest, "JSON invalide")
		return
	}

	user, err := h.service.Register(r.Context(), data.Email, data.Password)
	if err != nil {
		switch {
		case errors.Is(err, authService.ErrInvalidInput):
			sendError(w, http.StatusBadRequest, "email invalide ou mot de passe de 10 à 72 octets requis")
		case errors.Is(err, authService.ErrEmailTaken):
			sendError(w, http.StatusConflict, "email déjà utilisé")
		default:
			sendError(w, http.StatusInternalServerError, "erreur interne")
		}
		return
	}

	utils.SendJSONResponse(w, http.StatusCreated, userModel.FromUser(user))
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var data userModel.LoginRequest
	if err := readJSON(w, r, &data); err != nil {
		sendError(w, http.StatusBadRequest, "JSON invalide")
		return
	}

	token, err := h.service.Login(r.Context(), data.Email, data.Password)
	if err != nil {
		if errors.Is(err, authService.ErrInvalidCredentials) {
			sendError(w, http.StatusUnauthorized, "identifiants invalides")
		} else {
			sendError(w, http.StatusInternalServerError, "erreur interne")
		}
		return
	}

	utils.SendJSONResponse(w, http.StatusOK, userModel.LoginResponse{
		Token: token, TokenType: "Bearer", ExpiresIn: authService.TokenExpiresIn(),
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user := authMiddleware.CurrentUser(r)
	utils.SendJSONResponse(w, http.StatusOK, userModel.FromUser(user))
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.ListUsers(r.Context())
	if err != nil {
		sendError(w, http.StatusInternalServerError, "erreur interne")
		return
	}

	result := make([]userModel.UserResponse, 0, len(users))
	for i := range users {
		result = append(result, userModel.FromUser(&users[i]))
	}
	utils.SendJSONResponse(w, http.StatusOK, result)
}
