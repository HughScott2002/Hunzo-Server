package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"example.com/m/v2/src/server"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

//TODO: Use bcrypt or Argon2 for password hashing and salting.
//TODO: Implement rate limiting on login attempts to prevent brute-force attacks.
//TODO: Use secure session management techniques.

//JWT
//TODO: Use strong, randomly generated secrets for signing JWTs
//TODO: Include essential claims like 'exp' (expiration), 'iat' (issued at), and 'jti' (JWT ID).
//TODO: Keep JWT payload minimal to reduce token size.

func main() {
	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow your frontend origin
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	//Routing
	r.Mount("/api/users/auth", server.Router())
	r.Post("/api/users/dump", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		fmt.Printf("%s\n", body)
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})
	r.Get("/api/users/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("()"))
	})

	fmt.Println("User server is running on Port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
