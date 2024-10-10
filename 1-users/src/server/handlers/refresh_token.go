package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"example.com/m/v2/src/db"
	"example.com/m/v2/src/utils"
)

func HandlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	deviceInfo := r.Header.Get("User-Agent")
	// Get the refresh token from the cookie
	refreshTokenCookie, err := r.Cookie("refresh_token")
	if err != nil {
		http.Error(w, "Refresh token not found", http.StatusUnauthorized)
		return
	}

	refreshToken := refreshTokenCookie.Value

	// Validate the refresh token
	user, exists := db.RefreshTokens[refreshToken]
	if !exists || user.DeviceInfo != deviceInfo {
		http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
		return
	}

	// var user models.User
	// Generate new access token
	newAccessToken, err := utils.GenerateAccessToken(user.UserEmail)
	if err != nil {
		http.Error(w, "Error generating new access token", http.StatusInternalServerError)
		return
	}

	// Optionally, rotate the refresh token
	newRefreshToken, err := utils.GenerateRefreshToken(user.UserEmail)
	if err != nil {
		http.Error(w, "Error generating new refresh token", http.StatusInternalServerError)
		return
	}

	// Update stored refresh token
	delete(db.RefreshTokens, refreshToken)
	db.RefreshTokens[newRefreshToken] = user

	// Set new cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    newAccessToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   900, // 15 minutes
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   604800, // 7 days
	})

	sendUser, exists := db.Users[user.UserEmail]
	if !exists {
		http.Error(w, "User doesn't exist", http.StatusNotFound)
		return
	}
	userData := map[string]interface{}{
		"user": map[string]string{
			"id":        sendUser.AccountId,
			"email":     sendUser.Email,
			"firstName": sendUser.FirstName,
			"lastName":  sendUser.LastName,
			"kycStatus": sendUser.KYCStatus.String(),
		},
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userData)
	// json.NewEncoder(w).Encode(map[string]string{"message": "Tokens refreshed successfully"})
}
