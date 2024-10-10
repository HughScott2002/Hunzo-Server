package events

type UserCreatedEvent struct {
	AccountId string `json:"accountId"`
	Email     string `json:"email"`
}
