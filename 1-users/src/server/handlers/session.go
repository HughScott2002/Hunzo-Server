package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"example.com/m/v2/src/db"
	"example.com/m/v2/src/models"
	"example.com/m/v2/src/utils"
)

func HandlerListActiveSessions(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the request body
	var request struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(body, &request); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate that email is provided
	if request.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Get all sessions for the user
	sessions, err := db.GetUserSessions(request.Email)
	if err != nil {
		http.Error(w, "Failed to retrieve user sessions", http.StatusInternalServerError)
		return
	}

	// Format the sessions for the response
	var activeSessions []map[string]string
	for _, session := range sessions {
		activeSessions = append(activeSessions, map[string]string{
			"id":          session.ID,
			"deviceInfo":  session.DeviceInfo,
			"createdAt":   session.CreatedAt.String(),
			"lastLoginAt": session.LastLoginAt.String(),
		})
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"activeSessions": activeSessions,
	})
}

func HandlerCheckSession(w http.ResponseWriter, r *http.Request) {
	var user *models.User
	var err error

	// Try to get the access token
	accessTokenCookie, err := r.Cookie("access_token")
	if err == nil {
		// Access token exists, try to validate it
		claims, err := utils.ValidateAccessToken(accessTokenCookie.Value)
		if err == nil {
			// Access token is valid, get the user
			user, err = db.GetUser(claims.Subject)
			fmt.Println("user1")
		}
	}

	// If we don't have a valid user at this point, try the refresh token
	if user == nil {
		refreshTokenCookie, err := r.Cookie("refresh_token")
		if err != nil {
			http.Error(w, "No valid session found", http.StatusUnauthorized)
			fmt.Println("user2")
			return
		}

		// Get the refresh token info
		tokenInfo, err := db.GetRefreshToken(refreshTokenCookie.Value)
		if err != nil {
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			fmt.Println("user3")
			
			return
		}

		// Get the user associated with this refresh token
		user, err = db.GetUser(tokenInfo.UserEmail)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		// Generate new access token
		newAccessToken, err := utils.GenerateAccessToken(user.Email)
		if err != nil {
			http.Error(w, "Error generating new access token", http.StatusInternalServerError)
			return
		}

		// Set new access token cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    newAccessToken,
			HttpOnly: true,
			Secure:   true,
			SameSite: http.SameSiteStrictMode,
			MaxAge:   900, // 15 minutes
		})

		// Update the session's last login time
		sessions, err := db.GetUserSessions(user.Email)
		if err == nil {
			for _, session := range sessions {
				if session.Token == refreshTokenCookie.Value {
					db.UpdateSessionLastLogin(session.ID)
					break
				}
			}
		}
	}

	if user == nil {
		http.Error(w, "No valid session found", http.StatusUnauthorized)
		return
	}

	// Prepare the response
	userData := map[string]interface{}{
		"user": map[string]string{
			"id":        user.AccountId,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"kycStatus": user.KYCStatus.String(),
		},
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userData)
}
