package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/m/v2/src/db"
	"github.com/go-chi/chi"
)

func GetWallet(w http.ResponseWriter, r *http.Request) {
	walletId := chi.URLParam(r, "walletId")
	wallet, err := db.GetWallet(walletId)
	if err != nil {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wallet)
}
