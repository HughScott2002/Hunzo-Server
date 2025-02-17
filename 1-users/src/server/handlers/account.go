package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"example.com/m/v2/src/db"
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

// Define a struct to unmarshal the password change request
type PasswordChangeRequest struct {
	AccountId          string `json:"account-id"`
	Email              string `json:"email"`
	CurrentPassword    string `json:"current-password"`
	NewPassword        string `json:"new-password"`
	ConfirmNewPassword string `json:"confirm-new-password"`
}

func HandlerChangePassword(w http.ResponseWriter, r *http.Request) {
	// Get device information
	deviceInfo := r.Header.Get("User-Agent")
	
	fmt.Printf("%s", deviceInfo)
	
	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Parse the JSON body
	var passwordChangeReq PasswordChangeRequest
	if err := json.Unmarshal(body, &passwordChangeReq); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate passwords match
	if passwordChangeReq.NewPassword != passwordChangeReq.ConfirmNewPassword {
		http.Error(w, "New passwords do not match", http.StatusBadRequest)
		return
	}

	// Get the user from database
	storedUser, err := db.GetUser(passwordChangeReq.Email)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify account ID matches
	if storedUser.AccountId != passwordChangeReq.AccountId {
		http.Error(w, "Invalid account ID", http.StatusUnauthorized)
		return
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.HashedPassword), []byte(passwordChangeReq.CurrentPassword))
	if err != nil {
		http.Error(w, "Current password is incorrect", http.StatusUnauthorized)
		return
	}

	// Check if new password is same as current password
	err = bcrypt.CompareHashAndPassword([]byte(storedUser.HashedPassword), []byte(passwordChangeReq.NewPassword))
	if err == nil {
		http.Error(w, "New password cannot be the same as current password", http.StatusBadRequest)
		return
	}

	// Hash the new password
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(passwordChangeReq.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing new password", http.StatusInternalServerError)
		return
	}

	// Update user's password
	storedUser.HashedPassword = string(hashedNewPassword)
	err = db.UpdateUser(storedUser)
	if err != nil {
		http.Error(w, "Failed to update password", http.StatusInternalServerError)
		return
	}

	// Invalidate all refresh tokens for this user
	err = db.DeleteUserSessions(passwordChangeReq.Email)
	if err != nil {
		log.Printf("Error deleting user sessions: %v", err)
	}

	// TODO: Produce password changed event
	// passwordChangedEvent := events.PasswordChangedEvent{
	// 	AccountId:  storedUser.AccountId,
	// 	DeviceInfo: deviceInfo,
	// 	ChangedAt:  time.Now(),
	// }
	// err = producer.ProducePasswordChangedEvent(passwordChangedEvent)
	// if err != nil {
	// 	log.Printf("Failed to produce password changed event: %v", err)
	// }

	// Return success response
	response := map[string]string{
		"message":   "Password changed successfully",
		"email":     passwordChangeReq.Email,
		"accountId": passwordChangeReq.AccountId,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
