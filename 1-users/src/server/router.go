package server

import (
	"net/http"

	"example.com/m/v2/src/server/handlers"
	"example.com/m/v2/src/server/middleware"
	"github.com/go-chi/chi/v5"
)

// TODO: FIX THE AUTH, it needs to accept tokens on each request and make sure its vaild
// TODO: Fix the rate limiter
func Router() http.Handler {
	r := chi.NewRouter()

	// r.Post("/update", handlers.HandlerUpdateUser)

	// Protected routes
	r.Group(func(r chi.Router) {
		// r.Use(middleware.RequireAuth)  // Assuming you have this middleware
		r.Get("/check-session", handlers.HandlerCheckSession)
		r.Post("/refresh", middleware.RateLimitMiddleware(handlers.HandlerRefreshToken))
		// r.Get("/list-sessions", middleware.RateLimitMiddleware(handlers.HandlerListActiveSessions))

		// User account management
		r.Route("/account", func(r chi.Router) {
			r.Post("/login", handlers.HandlerLogin)
			r.Post("/register", handlers.HandlerRegister)
			r.Post("/logout", handlers.HandlerLogout)
			r.Get("/{accountid}", handlers.HandlerGetUserProfile)
			r.Put("/", handlers.HandlerUpdateUserProfile)
			r.Delete("/", handlers.HandlerDeleteUserAccount)
			r.Put("/change-password", handlers.HandlerChangePassword)
		})

		// Security settings
		r.Route("/security", func(r chi.Router) {
			// All sessions currently logged in 
			r.Post("/sessions", handlers.HandlerListActiveSessions)
			
			//TODO: ADD Logout for individual sessions
			//TODO: ADD Logout for all sessions
			
			// Two-factor authentication on these devices has been remembered for 30 days.
			//TODO: ADD Remembered Devices
			//TODO: ADD A Activity history logger
			// The last 30 days of activity on your account:
			// Event	       Source	        IP address	 Date and time	    Country
			// Log in failure  Chrome (Linux)	127.0.0.2	 Jan 2, 12:22 PM	United States
			
			r.Get("/sessions/{sessionId}", handlers.HandlerListActiveSessions)
			// Two-factor authentication (2FA) helps accounts secure by adding an extra layer of protection beyond a password. 
			// By default, we require you to set up a 2FA app that can generate 2FA codes, 
			// but you can add a security key to log in even quicker.
			r.Post("/enable-2fa", handlers.HandlerEnable2FA)
			r.Post("/disable-2fa", handlers.HandlerDisable2FA)

			
		})

		// Device and session management
		r.Route("/devices", func(r chi.Router) {
			r.Get("/", handlers.HandlerListDevices)
			r.Delete("/{deviceId}", handlers.HandlerRemoveDevice)
			r.Delete("/", handlers.HandlerRemoveAllDevices)
		})
	})

	return r
}

func HandlerPlaceHolder(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte{})
}
