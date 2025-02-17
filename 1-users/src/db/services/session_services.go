package services

import (
	"net/http"
	"strings"
	"time"

	"example.com/m/v2/src/db"
	"example.com/m/v2/src/models"
	"github.com/google/uuid"
)

// Parse browser info from User-Agent
func ParseBrowser(userAgent string) string {
	if strings.Contains(userAgent, "Chrome") {
		return "Chrome"
	} else if strings.Contains(userAgent, "Firefox") {
		return "Firefox"
	} else if strings.Contains(userAgent, "Safari") {
		return "Safari"
	} else if strings.Contains(userAgent, "Edge") {
		return "Edge"
	} else {
		return "Unknown Browser"
	}
}

// Format time for session display
func FormatSessionTime(t time.Time) string {
	now := time.Now()
	if t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day() {
		return "Current Session"
	}
	return t.Format("Jan 2 at 3:04 PM")
}

// Create and save a new session
func CreateSession(r *http.Request, email string) (*models.Session, error) {
	// Get the IP address
	ipAddress := r.Header.Get("X-Forwarded-For")
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}

	// Get the User-Agent
	userAgent := r.Header.Get("User-Agent")

	// Create a new session
	session := &models.Session{
		ID:          uuid.New().String(),
		UserEmail:   email,
		DeviceInfo:  userAgent,
		IPAddress:   ipAddress,
		Browser:     ParseBrowser(userAgent),
		Country:     "United States", // TODO: Use a geo-IP service here
		Token:       "",              //TODO:  Set this as needed
		LastLoginAt: time.Now(),
		CreatedAt:   time.Now(),
	}

	// Save the session
	err := db.AddSession(session)
	if err != nil {
		return nil, err
	}

	return session, nil
}
