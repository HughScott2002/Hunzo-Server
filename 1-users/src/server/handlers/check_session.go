package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/m/v2/src/db"
	"example.com/m/v2/src/utils"
)

func HandlerCheckSession(w http.ResponseWriter, r *http.Request) {
	// Get the access token from the cookie
	accessTokenCookie, err := r.Cookie("access_token")
	if err != nil {
		http.Error(w, "Access token not found", http.StatusUnauthorized)
		return
	}

	// Validate the access token
	claims, err := utils.ValidateAccessToken(accessTokenCookie.Value)
	if err != nil {
		http.Error(w, "Invalid access token: "+err.Error(), http.StatusUnauthorized)
		return
	}

	// Get the user email from the token claims
	email := claims.Subject

	// Fetch the user from the database
	user, exists := db.Users[email]
	if !exists {
		http.Error(w, "User not found", http.StatusNotFound)
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
			// Add any other user fields you want to include
		},
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userData)
}
