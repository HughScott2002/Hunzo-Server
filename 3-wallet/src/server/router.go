package server

import (
	"net/http"

	"example.com/m/v2/src/server/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)
// TODO: GET A LIST OF THE WALLETS AND EACH ATTACHED CARD

func Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Route("/api/wallets", func(r chi.Router) {
		//Wallet by account
		r.Get("/{walletId}", handlers.GetWallet)
		r.Get("/list/{accountId}", handlers.ListWallets)
		r.Get("/recover", func(http.ResponseWriter, *http.Request) { panic("foo") })
	})

	return r
}
