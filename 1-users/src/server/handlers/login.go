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
	storedUser, err := db.GetUser(loginRequest.Email)
	if err != nil {
		http.Error(w, "User doesn't exist", http.StatusNotFound)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedUser.HashedPassword), []byte(loginRequest.UnHashedPassword))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	accessToken, err := utils.GenerateAccessToken(loginRequest.Email)
	if err != nil {
		http.Error(w, "Error generating access token", http.StatusInternalServerError)
		return
	}

	refreshToken, err := utils.GenerateRefreshToken(loginRequest.Email)
	if err != nil {
		http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
		return
	}

	accountId, err := utils.GenerateAccountId()
	if err != nil {
		http.Error(w, "Error generating AccountId", http.StatusInternalServerError)
		return
	}

	//Check if a session already exists
	//
	session := &models.Session{
		ID:          accountId,
		UserEmail:   storedUser.Email,
		DeviceInfo:  deviceInfo,
		Token:       refreshToken,
		LastLoginAt: time.Now(),
		CreatedAt:   time.Now(),
	}

	err = db.AddSession(session)
	if err != nil {
		http.Error(w, "Error creating session", http.StatusInternalServerError)
		return
	}

	err = db.AddRefreshToken(refreshToken, db.RefreshTokenInfo{
		UserEmail:  storedUser.Email,
		DeviceInfo: deviceInfo,
		CreatedAt:  time.Now(),
	})
	if err != nil {
		http.Error(w, "Error storing refresh token", http.StatusInternalServerError)
		return
	}
	// Set cookies
	utils.SetCookie(w, "access_token", accessToken, 15*60)        // 15 minutes
	utils.SetCookie(w, "refresh_token", refreshToken, 7*24*60*60) // 7 days
	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "access_token",
	// 	Value:    accessToken,
	// 	HttpOnly: false,
	// 	Secure:   false,
	// 	SameSite: http.SameSiteStrictMode,
	// 	Path:     "/",
	// 	MaxAge:   900, // 15 minutes
	// })

	// http.SetCookie(w, &http.Cookie{
	// 	Name:     "refresh_token",
	// 	Value:    refreshToken,
	// 	HttpOnly: false,
	// 	Secure:   false,
	// 	SameSite: http.SameSiteStrictMode,
	// 	Path:     "/",
	// 	MaxAge:   604800, // 7 days
	// })

	userData := map[string]interface{}{
		"user": map[string]string{
			"id":        storedUser.AccountId,
			"email":     storedUser.Email,
			"firstName": storedUser.FirstName,
			"lastName":  storedUser.LastName,
			"kycStatus": storedUser.KYCStatus.String(),
			// "access_token": accessToken,
			// "refresh_token": refreshToken,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userData)
}
