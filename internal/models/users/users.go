package users

import "time"

// User correspond à un utilisateur enregistré dans PostgreSQL.
// PasswordHash ne doit jamais être envoyé dans une réponse HTTP.
type User struct {
	ID           int
	Email        string
	PasswordHash string
	Role         string
	CreatedAt    time.Time
}

// RegisterRequest est le JSON reçu pour créer un compte.
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest est le JSON reçu pour se connecter.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse est la réponse HTTP après une connexion réussie.
type LoginResponse struct {
	Token     string `json:"token"`
	TokenType string `json:"token_type"`
	ExpiresIn int    `json:"expires_in"`
}

// UserResponse expose seulement les informations publiques du compte.
type UserResponse struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}

// FromUser transforme un utilisateur de la base en réponse publique.
func FromUser(user *User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}
}
