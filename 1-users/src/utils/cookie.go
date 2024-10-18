package utils

import (
	"net/http"
	"os"
)

func setCookie(w http.ResponseWriter, name, value string, maxAge int) {
	isSecure := os.Getenv("ENVIRONMENT") == "production" || os.Getenv("ENVIRONMENT") == "prod"

	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		HttpOnly: true,
		Secure:   isSecure, // Only set to true when using HTTPS
		SameSite: http.SameSiteLaxMode,
		Path:     "/",
		MaxAge:   maxAge,
	})
}
