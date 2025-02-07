package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"example.com/m/v2/src/db"
	"example.com/m/v2/src/models"
	"golang.org/x/crypto/bcrypt"
)

func HandlerGetUserProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement logic to get user profile
	json.NewEncoder(w).Encode(map[string]string{"message": "Get user profile"})
}

func HandlerUpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement logic to update user profile
	deviceInfo := r.Header.Get("User-Agent")
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	fmt.Printf("%s\n", body)

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Update user profile",
		"device":  deviceInfo,
		"body":    string(body),
	})
}

func HandlerDeleteUserAccount(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement logic to delete user account
	json.NewEncoder(w).Encode(map[string]string{"message": "Delete user account"})
}

func HandlerChangePassword(w http.ResponseWriter, r *http.Request) {
	// Get device information
	deviceInfo := r.Header.Get("User-Agent")

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Define a struct to unmarshal the password change request
	type PasswordChangeRequest struct {
		AccountId          string `json:"account-id"`
		Email              string `json:"email"`
		CurrentPassword    string `json:"current-password"`
		NewPassword        string `json:"new-password"`
		ConfirmNewPassword string `json:"confirm-new-password"`
	}

	// Parse the JSON body
	var passwordChangeReq PasswordChangeRequest
	if err := json.Unmarshal(body, &passwordChangeReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate input
	if passwordChangeReq.NewPassword != passwordChangeReq.ConfirmNewPassword {
		http.Error(w, "New passwords do not match", http.StatusBadRequest)
		return
	}

	var loginRequest models.User
	// Get the user email from the access token
	// Assuming you have a function to extract email from the token
	storedUser, err := db.GetUser(loginRequest.Email)
	if err != nil {
		http.Error(w, "User doesn't exist", http.StatusNotFound)
		return
	}
	
	email, err := utils.GetEmailFromAccessToken(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Retrieve user from database
	user, err := db.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(passwordChangeReq.CurrentPassword))
	if err != nil {
		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	// Hash the new password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(passwordChangeReq.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing new password", http.StatusInternalServerError)
		return
	}

	// Update user's password in the database
	err = db.UpdateUserPassword(email, string(hashedNewPassword))
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	// Optional: Produce an event for password change
	passwordChangedEvent := events.PasswordChangedEvent{
		AccountId:  user.AccountId,
		DeviceInfo: deviceInfo,
		ChangedAt:  time.Now(),
	}
	err = producer.ProducePasswordChangedEvent(passwordChangedEvent)
	if err != nil {
		fmt.Printf("failed to produce password changed event: %v", err)
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Password changed successfully",
	})
}
