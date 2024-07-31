package main

import (
	"errors"
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSetupRouter(t *testing.T) {
	router := setupRouter()

	// Test if router is not nil
	assert.NotNil(t, router, "Router should not be nil")

	// Test if specific routes are set up
	routes := []string{
		"/swagger/doc.json",
		"/swagger/*",
		"/api/register",
		"/api/login",
		"/api/address/search",
		"/api/address/geocode",
	}

	for _, route := range routes {
		assert.True(t, routeExists(router, route), "Route %s should exist", route)
	}
}

func TestAuthMiddleware(t *testing.T) {
	handler := AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Test with no token
	req := httptest.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code, "Should return 403 when no token is provided")

	// Test with invalid token
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer invalidtoken")
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusForbidden, rr.Code, "Should return 403 when invalid token is provided")

	// Test with valid token
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{"user_id": "123"})
	req = httptest.NewRequest("GET", "/", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code, "Should return 200 when valid token is provided")
}

// Helper function to check if a route exists in the router
func routeExists(router *chi.Mux, path string) bool {
	found := false
	_ = chi.Walk(router, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		if route == path {
			found = true
			return errors.New("route found")
		}
		return nil
	})
	return found
}
