package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/mail"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	userModel "steam_games/internal/models/users"
	userRepository "steam_games/internal/repositories/users"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidInput = errors.New("email ou mot de passe invalide")
var ErrEmailTaken = errors.New("cet email est déjà utilisé")
var ErrInvalidCredentials = errors.New("identifiants invalides")
var ErrInvalidToken = errors.New("token invalide ou expiré")

const tokenDuration = time.Hour

type Service struct {
	repository *userRepository.Repository
	secret     []byte
}

func New(repository *userRepository.Repository, secret string) (*Service, error) {
	if len(secret) < 32 {
		return nil, errors.New("JWT_SECRET doit contenir au moins 32 caractères")
	}
	return &Service{repository: repository, secret: []byte(secret)}, nil
}

func (s *Service) Register(ctx context.Context, email string, password string) (*userModel.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email || len(email) > 255 || utf8.RuneCountInString(password) < 10 || len([]byte(password)) > 72 {
		return nil, ErrInvalidInput
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.repository.Create(ctx, email, string(hash))
	if err != nil {
		var pgError *pq.Error
		if errors.As(err, &pgError) && pgError.Code == "23505" {
			return nil, ErrEmailTaken
		}
		return nil, err
	}
	return user, nil
}

func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := s.repository.FindByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return "", ErrInvalidCredentials
	}
	if err != nil {
		return "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		return "", ErrInvalidCredentials
	}

	now := time.Now()
	claims := jwt.RegisteredClaims{
		Issuer:    "steam-games",
		Subject:   strconv.Itoa(user.ID),
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(now.Add(tokenDuration)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

func (s *Service) UserFromToken(ctx context.Context, tokenString string) (*userModel.User, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (any, error) {
			return s.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer("steam-games"),
		jwt.WithExpirationRequired(),
	)
	if err != nil || !token.Valid {
		return nil, ErrInvalidToken
	}

	id, err := strconv.Atoi(claims.Subject)
	if err != nil || id <= 0 {
		return nil, ErrInvalidToken
	}

	user, err := s.repository.FindByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrInvalidToken
	}
	if err != nil {
		return nil, fmt.Errorf("recherche de l'utilisateur : %w", err)
	}
	return user, nil
}

func (s *Service) ListUsers(ctx context.Context) ([]userModel.User, error) {
	return s.repository.FindAll(ctx)
}

func TokenExpiresIn() int {
	return int(tokenDuration.Seconds())
}
