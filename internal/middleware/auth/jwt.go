package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"steam_games/internal/models/responses"
	userModel "steam_games/internal/models/users"
	authService "steam_games/internal/services/auth"
	"steam_games/internal/utils"
)

type userContextKey struct{}

type Middleware struct {
	service *authService.Service
}

func New(service *authService.Service) *Middleware {
	return &Middleware{service: service}
}

func CurrentUser(request *http.Request) *userModel.User {
	user, _ := request.Context().Value(userContextKey{}).(*userModel.User)
	return user
}

func (m *Middleware) RequireLogin(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		parts := strings.Fields(authorization)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			utils.SendJSONResponse(w, http.StatusUnauthorized, responses.GenericResponse{
				Success: false, Message: "token Bearer manquant",
			})
			return
		}

		user, err := m.service.UserFromToken(r.Context(), parts[1])
		if err != nil {
			if errors.Is(err, authService.ErrInvalidToken) {
				utils.SendJSONResponse(w, http.StatusUnauthorized, responses.GenericResponse{
					Success: false, Message: "token invalide ou expiré",
				})
			} else {
				utils.SendJSONResponse(w, http.StatusInternalServerError, responses.GenericResponse{
					Success: false, Message: "erreur interne",
				})
			}
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey{}, user)
		next(w, r.WithContext(ctx))
	}
}

func (m *Middleware) RequireAdmin(next http.HandlerFunc) http.HandlerFunc {
	return m.RequireLogin(func(w http.ResponseWriter, r *http.Request) {
		user := CurrentUser(r)
		if user.Role != "ADMIN" {
			utils.SendJSONResponse(w, http.StatusForbidden, responses.GenericResponse{
				Success: false, Message: "accès réservé aux administrateurs",
			})
			return
		}
		next(w, r)
	})
}
