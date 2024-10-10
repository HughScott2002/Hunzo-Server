package db

import (
	"sync"
	"time"

	"example.com/m/v2/src/models"
)

type RefreshTokenInfo struct {
	UserEmail  string
	DeviceInfo string
	CreatedAt  time.Time
}

var Users = make(map[string]models.User)
var Sessions = make(map[string]models.Session)
var RefreshTokens = make(map[string]RefreshTokenInfo)

// Change the map to hold values of the type visitor.
var Visitors = make(map[string]*models.Visitor)
var Mu sync.Mutex
