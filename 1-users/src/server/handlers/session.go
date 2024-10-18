package handlers

import (
	"encoding/json"
	"io"
	"log"
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

	log.Println("Starting session check")

	// Try to get the access token
	accessTokenCookie, err := r.Cookie("access_token")
	if err != nil {
		log.Println("Access token not found:", err)
	} else {
		// Access token exists, try to validate it
		claims, err := utils.ValidateAccessToken(accessTokenCookie.Value)
		if err != nil {
			log.Println("Access token validation failed:", err)
		} else {
			// Access token is valid, get the user
			user, err = db.GetUser(claims.Subject)
			if err != nil {
				log.Println("Failed to get user from access token:", err)
			} else {
				log.Println("User retrieved from access token")
			}
		}
	}

	// If we don't have a valid user at this point, try the refresh token
	if user == nil {
		log.Println("Attempting to use refresh token")
		refreshTokenCookie, err := r.Cookie("refresh_token")
		if err != nil {
			log.Println("Refresh token not found:", err)
			http.Error(w, "No valid session found", http.StatusUnauthorized)
			return
		}

		log.Printf("Refresh token found: %s", refreshTokenCookie.Value)

		// Get the refresh token info
		tokenInfo, err := db.GetRefreshToken(refreshTokenCookie.Value)
		if err != nil {
			log.Println("Invalid refresh token:", err)
			http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
			return
		}

		log.Printf("Refresh token info retrieved: %+v", tokenInfo)

		// Get the user associated with this refresh token
		user, err = db.GetUser(tokenInfo.UserEmail)
		if err != nil {
			log.Println("User not found from refresh token:", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		log.Println("User retrieved from refresh token")

		// Generate new access token
		newAccessToken, err := utils.GenerateAccessToken(user.Email)
		if err != nil {
			log.Println("Error generating new access token:", err)
			http.Error(w, "Error generating new access token", http.StatusInternalServerError)
			return
		}

		// Set new access token cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "access_token",
			Value:    newAccessToken,
			HttpOnly: false,
			Secure:   false,
			Path:     "/",

			SameSite: http.SameSiteStrictMode,
			MaxAge:   900, // 15 minutes
		})

		log.Println("New access token set")

		// Update the session's last login time
		sessions, err := db.GetUserSessions(user.Email)
		if err == nil {
			for _, session := range sessions {
				if session.Token == refreshTokenCookie.Value {
					db.UpdateSessionLastLogin(session.ID)
					log.Println("Session last login time updated")
					break
				}
			}
		} else {
			log.Println("Failed to get user sessions:", err)
		}
	}

	if user == nil {
		log.Println("No valid user found after all checks")
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
	log.Println("Session check completed successfully")
}
