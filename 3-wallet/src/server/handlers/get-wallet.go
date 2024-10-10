package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/m/v2/src/db"
	"github.com/go-chi/chi"
)

func GetWallet(w http.ResponseWriter, r *http.Request) {
	accountId := chi.URLParam(r, "accountId")
	wallet, exists := db.Wallets[accountId]
	if !exists {
		http.Error(w, "Wallet not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wallet)
}
