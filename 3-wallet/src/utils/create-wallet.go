package utils

import (
	"log"

	"example.com/m/v2/src/db"
	"example.com/m/v2/src/models"
)

func CreateWallet(accountId string) {
	wallet := &models.Wallet{
		AccountId: accountId,
		Balance:   0,
	}
	db.Wallets[accountId] = wallet
	log.Printf("Wallet created for account: %s", accountId)
}
