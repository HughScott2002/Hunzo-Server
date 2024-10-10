package models

type Wallet struct {
	AccountId string  `json:"accountId"`
	Balance   float64 `json:"balance"`
}
