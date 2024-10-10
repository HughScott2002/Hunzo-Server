package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/m/v2/src/db"
	"example.com/m/v2/src/utils"
)

func HandlerListActiveSessions(w http.ResponseWriter, r *http.Request) {
	// Get the current user's email from the access token
	accessToken, err := r.Cookie("access_token")
	if err != nil {
		http.Error(w, "Access token not found", http.StatusUnauthorized)
		return
	}

	claims, err := utils.ValidateAccessToken(accessToken.Value)
	if err != nil {
		http.Error(w, "Invalid access token", http.StatusUnauthorized)
		return
	}

	userEmail := claims.Subject

	// Collect active sessions for this user
	var activeSessions []map[string]string
	for _, tokenInfo := range db.RefreshTokens {
		if tokenInfo.UserEmail == userEmail {
			activeSessions = append(activeSessions, map[string]string{
				"deviceInfo": tokenInfo.DeviceInfo,
				"createdAt":  tokenInfo.CreatedAt.String(),
			})
		}
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"activeSessions": activeSessions,
	})
}
