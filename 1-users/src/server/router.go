package server

import (
	"net/http"

	"example.com/m/v2/src/server/handlers"
	"example.com/m/v2/src/server/middleware"
	"github.com/go-chi/chi/v5"
)

// TODO: FIX AUTH REFRESH
// TODO: SEND BACK RELIVANT USER DATA TO THE CLIENT
// TODO: Fix the rate limiter
func Router() http.Handler {
	r := chi.NewRouter()

	r.Post("/login", handlers.HandlerLogin)
	r.Post("/register", handlers.HandlerRegister)
	// r.Post("/update", handlers.HandlerUpdateUser)

	// Protected routes
	r.Group(func(r chi.Router) {
		// r.Use(middleware.RequireAuth)  // Assuming you have this middleware
		r.Get("/check-session", handlers.HandlerCheckSession)
		r.Post("/refresh", middleware.RateLimitMiddleware(handlers.HandlerRefreshToken))
		// r.Get("/list-sessions", middleware.RateLimitMiddleware(handlers.HandlerListActiveSessions))
		r.Post("/logout", handlers.HandlerLogout)

		// User account management
		r.Route("/account", func(r chi.Router) {
			r.Get("/", handlers.HandlerGetUserProfile)
			r.Put("/", handlers.HandlerUpdateUserProfile)
			r.Delete("/", handlers.HandlerDeleteUserAccount)
			r.Put("/change-password", handlers.HandlerChangePassword)
		})

		// Security settings
		r.Route("/security", func(r chi.Router) {
			r.Get("/sessions", handlers.HandlerListActiveSessions)
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
