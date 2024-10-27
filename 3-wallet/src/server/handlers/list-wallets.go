package handlers

import (
	"encoding/json"
	"net/http"

	"example.com/m/v2/src/db"
	"github.com/go-chi/chi"
)

func ListWallets(w http.ResponseWriter, r *http.Request) {
	accountId := chi.URLParam(r, "accountId")
	wallet, err := db.ListWallets(accountId)
	if err != nil {
		http.Error(w, "No Wallets Found", http.StatusNotFound)
		return
	}
	

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(wallet)
}
