package users

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	userModel "steam_games/internal/models/users"
)

func TestReadJSON(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"email":"test@example.com","password":"motdepasse"}`))
	w := httptest.NewRecorder()
	var data userModel.RegisterRequest
	if err := readJSON(w, r, &data); err != nil {
		t.Fatal(err)
	}
	if data.Email != "test@example.com" {
		t.Fatalf("email inattendu : %s", data.Email)
	}
}

func TestReadJSONRejectsUnknownFields(t *testing.T) {
	r := httptest.NewRequest(http.MethodPost, "/register", strings.NewReader(`{"email":"test@example.com","admin":true}`))
	w := httptest.NewRecorder()
	var data userModel.RegisterRequest
	if err := readJSON(w, r, &data); err == nil {
		t.Fatal("un champ inconnu devrait être refusé")
	}
}
