package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"example.com/m/v2/src/db"
	"example.com/m/v2/src/models"
	"example.com/m/v2/src/utils"
	"golang.org/x/crypto/bcrypt"
)

func HandlerLogin(w http.ResponseWriter, r *http.Request) {
	deviceInfo := r.Header.Get("User-Agent")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Printf("%s\n", body)

	var loginRequest models.User
	// Decode the JSON body into the User struct
	if err := json.Unmarshal(body, &loginRequest); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}
	// Check if the user exists in memory
	storedUser, exists := db.Users[loginRequest.Email]
	if !exists {
		http.Error(w, "User doesn't exist", http.StatusNotFound)
		return
	}
	// Compare the stored hashed password with the hash of the provided password
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.HashedPassword), []byte(loginRequest.HashedPassword))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Generate access token
	accessToken, err := utils.GenerateAccessToken(loginRequest.Email)
	if err != nil {
		http.Error(w, "Error generating access token", http.StatusInternalServerError)
		return
	}
	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(loginRequest.Email)
	if err != nil {
		http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
		return
	}

	db.RefreshTokens[refreshToken] = db.RefreshTokenInfo{
		UserEmail:  storedUser.Email,
		DeviceInfo: deviceInfo,
		CreatedAt:  time.Now(),
	}

	// Set cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    accessToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   900, // 15 minutes
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   604800, // 7 days
	})

	// Return the session token as part of the response
	// response := map[string]string{
	// 	"message":   "Login successful",
	// 	"kycstatus": storedUser.KYCStatus.String(),
	// }

	// After successful login
	userData := map[string]interface{}{
		"user": map[string]string{
			"id":        storedUser.AccountId,
			"email":     storedUser.Email,
			"firstName": storedUser.FirstName,
			"lastName":  storedUser.LastName,
			"kycStatus": storedUser.KYCStatus.String(),
		},
	}

	// Set headers and return the response
	w.Header().Set("Content-Type", "application/json")
	// json.NewEncoder(w).Encode(response)
	json.NewEncoder(w).Encode(userData)

}
