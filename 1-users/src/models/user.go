package models

import (
	"encoding/json"
	"fmt"
)

type KYCStatus int

const (
	KYCStatusPending  KYCStatus = iota // 0
	KYCStatusApproved                  // 1
	KYCStatusRejected                  // 2
)

type User struct {
	AccountId         string    `json:"accountId"`
	FirstName         string    `json:"firstName"`
	LastName          string    `json:"lastName"`
	Address           string    `json:"address"`
	City              string    `json:"city"`
	State             string    `json:"state"`
	Country           string    `json:"country"`
	Currency          string    `json:"currency"`
	PostalCode        string    `json:"postalCode"`
	DOB               string    `json:"dob"`
	GovId             string    `json:"govId"`
	Email             string    `json:"email"`
	HashedPassword    string    `json:"password"`
	KYCStatus         KYCStatus `json:"kycstatus"`
	DataAuthorization bool      `json:"dataAuthorization"`
}

func (s KYCStatus) String() string {
	switch s {
	case KYCStatusPending:
		return "pending"
	case KYCStatusApproved:
		return "approved"
	case KYCStatusRejected:
		return "rejected"
	default:
		return "unknown"
	}
}

// UnmarshalJSON implements the json.Unmarshaler interface
func (s *KYCStatus) UnmarshalJSON(data []byte) error {
	var status string
	if err := json.Unmarshal(data, &status); err != nil {
		return err
	}

	switch status {
	case "pending":
		*s = KYCStatusPending
	case "approved":
		*s = KYCStatusApproved
	case "rejected":
		*s = KYCStatusRejected
	default:
		return fmt.Errorf("invalid KYC status: %s", status)
	}

	return nil
}

// MarshalJSON implements the json.Marshaler interface
func (s KYCStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}
