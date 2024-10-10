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

	// Protected routes
	r.Group(func(r chi.Router) {
		// r.Use(middleware.RequireAuth)  // Assuming you have this middleware
		r.Post("/refresh", middleware.RateLimitMiddleware(handlers.HandlerRefreshToken))
		r.Get("/list-sessions", middleware.RateLimitMiddleware(handlers.HandlerListActiveSessions))
		r.Post("/logout", handlers.HandlerLogout)
	})

	return r
}
