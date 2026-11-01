package users

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	authMiddleware "steam_games/internal/middleware/auth"
	"steam_games/internal/models/responses"
	userModel "steam_games/internal/models/users"
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

// Register crée un nouveau compte utilisateur.
//
// @Summary Créer un compte utilisateur
// @Description Crée un nouveau compte utilisateur à partir d'une adresse e-mail et d'un mot de passe.
// @Description Le mot de passe doit contenir entre 10 et 72 caractères.
// @ID register-user
// @Tags Auth
// @Accept json
// @Produce json
//
// @Param request body users.RegisterRequest true "Informations nécessaires à la création du compte"
//
// @Success 201 {object} users.UserResponse "Compte créé avec succès"
// @Failure 400 {object} responses.GenericResponse "JSON invalide, adresse e-mail invalide ou mot de passe invalide"
// @Failure 409 {object} responses.GenericResponse "Adresse e-mail déjà utilisée"
// @Failure 500 {object} responses.GenericResponse "Erreur interne lors de la création du compte"
// @Router /auth/register [post]
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

// Login authentifie un utilisateur.
//
// @Summary Se connecter
// @Description Authentifie un utilisateur à partir de son adresse e-mail et de son mot de passe.
// @Description En cas de succès, retourne un token JWT de type Bearer ainsi que sa durée de validité.
// @ID login-user
// @Tags Auth
// @Accept json
// @Produce json
//
// @Param request body users.LoginRequest true "Identifiants de connexion"
//
// @Success 200 {object} users.LoginResponse "Authentification réussie"
// @Failure 400 {object} responses.GenericResponse "Corps JSON invalide"
// @Failure 401 {object} responses.GenericResponse "Adresse e-mail ou mot de passe incorrect"
// @Failure 500 {object} responses.GenericResponse "Erreur interne lors de l'authentification"
// @Router /auth/login [post]
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

// Me retourne l'utilisateur actuellement authentifié.
//
// @Summary Récupérer mon profil
// @Description Retourne les informations de l'utilisateur correspondant au token JWT envoyé dans l'en-tête Authorization.
// @Description Le token doit être envoyé sous la forme : Bearer <token>.
// @ID get-current-user
// @Tags Users
// @Produce json
// @Security BearerAuth
//
// @Success 200 {object} users.UserResponse "Utilisateur authentifié"
// @Failure 401 {object} responses.GenericResponse "Token absent, invalide ou expiré"
// @Router /users/me [get]
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	user := authMiddleware.CurrentUser(r)
	utils.SendJSONResponse(w, http.StatusOK, userModel.FromUser(user))
}

// GetAll retourne la liste des utilisateurs.
//
// @Summary Lister les utilisateurs
// @Description Retourne la liste de tous les utilisateurs enregistrés dans l'application.
// @Description Le token doit être envoyé sous la forme : Bearer <token>.
// @ID get-all-users
// @Tags Users
// @Produce json
// @Security BearerAuth
//
// @Success 200 {array} users.UserResponse "Liste des utilisateurs"
// @Failure 401 {object} responses.GenericResponse "Token absent, invalide ou expiré"
// @Failure 500 {object} responses.GenericResponse "Erreur interne lors de la récupération des utilisateurs"
// @Router /users [get]
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
