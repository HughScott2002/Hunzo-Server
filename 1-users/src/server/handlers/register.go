package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"example.com/m/v2/src/db"
	"example.com/m/v2/src/events/producer"
	"example.com/m/v2/src/models"
	"example.com/m/v2/src/models/events"
	"example.com/m/v2/src/utils"
	"golang.org/x/crypto/bcrypt"
)

func HandlerRegister(w http.ResponseWriter, r *http.Request) {
	//Change
	deviceInfo := r.Header.Get("User-Agent")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Println(deviceInfo)
	fmt.Printf("%s\n", body)
	//

	var user models.User
	// Decode the JSON body into the User struct
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Check if the user already exists
	if _, exists := db.Users[user.Email]; exists {
		utils.ErrorResponse(w, "User already exists", 500)
		// http.Error(w, "User already exists", http.StatusConflict)
		return
	}
	//Create the Account Id
	user.AccountId, err = utils.GenerateAccountId()
	if err != nil {
		http.Error(w, "Error generating account ID", http.StatusInternalServerError)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.HashedPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	user.HashedPassword = string(hashedPassword)

	// Set the initial KYC status to Pending
	user.KYCStatus = models.KYCStatusPending

	// Save the user in memory
	db.Users[user.Email] = user

	userCreatedEvent := events.UserCreatedEvent{
		AccountId: user.AccountId,
		Email:     user.Email,
	}

	err = producer.ProduceUserCreatedEvent(userCreatedEvent)
	if err != nil {
		log.Printf("failed to produce user created event: %v", err)
	}

	//TODO: DEAL WITH SESSION TOKENS (Probably should be on the user object)
	// Generate access token
	accessToken, err := utils.GenerateAccessToken(user.Email)
	if err != nil {
		http.Error(w, "Error generating access token", http.StatusInternalServerError)
		return
	}
	// Generate refresh token
	refreshToken, err := utils.GenerateRefreshToken(user.Email)
	if err != nil {
		http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
		return
	}

	db.RefreshTokens[refreshToken] = db.RefreshTokenInfo{
		UserEmail:  user.Email,
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

	// // Create a response object
	// response := map[string]string{
	// 	"message":   "User registered successfully",
	// 	"kycstatus": user.KYCStatus.String(),
	// 	"accountId": user.AccountId,
	// }
	// // Set headers and return the response
	// w.Header().Set("Content-Type", "application/json")
	// w.WriteHeader(http.StatusCreated)
	// // w.Write(users[user.Email])
	// json.NewEncoder(w).Encode(response)
	// After successful registration
	userData := map[string]interface{}{
		"user": map[string]string{
			"id":        user.AccountId,
			"email":     user.Email,
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"kycStatus": user.KYCStatus.String(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userData)
}