package server

import (
	"net/http"

	"example.com/m/v2/src/server/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"

)
// TODO: GET A LIST OF THE WALLETS AND EACH ATTACHED CARD

func Router() http.Handler {
	r := chi.NewRouter()
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
	r.Route("/api/wallets", func(r chi.Router) {
		//Wallet by account
		r.Get("/{walletId}", handlers.GetWallet)
		r.Get("/list/{accountId}", handlers.ListWallets) //List all the wallets
		r.Get("/recover", func(http.ResponseWriter, *http.Request) { panic("foo") })
	})

	return r
}
