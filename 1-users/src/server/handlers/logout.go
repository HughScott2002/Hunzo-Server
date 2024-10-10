package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/m/v2/src/db"
)

func HandlerLogout(w http.ResponseWriter, r *http.Request) {
	// Get the refresh token from the cookie
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err == nil {
		// If the refresh token exists, remove it from the stored tokens
		delete(db.RefreshTokens, refreshTokenCookie.Value)
	}

	// Clear the access token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	// Clear the refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
