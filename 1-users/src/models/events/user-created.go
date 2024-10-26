package events

import "example.com/m/v2/src/models"

type AccountCreatedEvent struct {
	AccountId string           `json:"accountId"`
	Currency  string           `json:"currency"`
	KYCStatus models.KYCStatus `json:"kycstatus"`
}
