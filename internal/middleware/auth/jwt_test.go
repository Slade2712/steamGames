package auth

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMissingJWTIsUnauthorized(t *testing.T) {
	middleware := &Middleware{}
	handler := middleware.RequireLogin(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/api/v1/games", nil)
	handler(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("status obtenu %d, attendu 401", w.Code)
	}
}
