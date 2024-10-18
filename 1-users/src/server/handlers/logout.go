package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/m/v2/src/db"
)

func HandlerLogout(w http.ResponseWriter, r *http.Request) {
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err == nil {
		refreshToken := refreshTokenCookie.Value

		// Get the refresh token info
		tokenInfo, err := db.GetRefreshToken(refreshToken)
		if err == nil {
			// Delete the session associated with this refresh token
			sessions, err := db.GetUserSessions(tokenInfo.UserEmail)
			if err == nil {
				for _, session := range sessions {
					if session.Token == refreshToken {
						db.DeleteSession(session.ID)
						break
					}
				}
			}
		}

		// Delete the refresh token
		db.DeleteRefreshToken(refreshToken)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		HttpOnly: false,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1,
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logged out successfully"})
}
