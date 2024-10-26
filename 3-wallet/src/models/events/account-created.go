package events

type AccountCreatedEvent struct {
	AccountId string `json:"accountId"`
	// Email     string `json:"email"`
	KYCStatus string `json:"kycstatus"`
}
