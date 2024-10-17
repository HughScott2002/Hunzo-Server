package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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
	// TODO: Implement logic to change password
	json.NewEncoder(w).Encode(map[string]string{"message": "Change password"})
}
