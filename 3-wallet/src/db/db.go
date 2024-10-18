package db

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"example.com/m/v2/src/models"
	"github.com/go-redis/redis/v8"
)

var db Database

type Database interface {
	AddWallet(wallet *models.Wallet) error
	GetWallet(wallet *models.Wallet) error
	WalletExists(wallet *models.Wallet) error
}

type MemoryDB struct {
	wallets map[string]*models.Wallet
	mu      sync.RWMutex
}

type RedisDB struct {
	client *redis.Client
}

func Init() error {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	mode := strings.ToLower(os.Getenv("MODE"))

	switch {
	case env == "local" && mode == "memcached":
		db = &MemoryDB{
			wallets: make(map[string]*models.Wallet),
		}
	case env == "local" && mode != "memcashed":
		redisPassword := os.Getenv("WALLET_REDIS_PASSWORD")
		redisAddr := fmt.Sprintf("wallet-redis:%s", os.Getenv("WALLET_REDIS_PORT"))
		redisClient := redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       0,
		})
		db = &RedisDB{client: redisClient}
		log.Println("USING DB & REDIS IN USER SERVICE")
	case env == "prod" || env == "production":
		db = &FutureDB{}
	default:
		return fmt.Errorf("unsupported environment or mode")
	}

	return nil
}

//Helper functions
func AddWallet(wallet models.Wallet) error{
	exists, err := 
}


// MemoryDB inplementations
func (m *MemoryDB) AddWallet(wallet *models.Wallet) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.wallets[wallet.AccountId] = wallet
	return nil
}

// GetWallet implements Database.
func (m *MemoryDB) GetWallet(wallet *models.Wallet) error {
	// panic("unimplemented")
	return nil
}

// DB & Redis Implementations
// AddWallet implements Database.
func (r *RedisDB) AddWallet(wallet *models.Wallet) error {
	panic("unimplemented")
}

// GetWallet implements Database.
func (r *RedisDB) GetWallet(wallet *models.Wallet) error {
	panic("unimplemented")
}

type FutureDB struct {
	// Placeholder for future database implementation
}

// AddWallet implements Database.
func (f *FutureDB) AddWallet(wallet *models.Wallet) error {
	panic("unimplemented")
}

// GetWallet implements Database.
func (f *FutureDB) GetWallet(wallet *models.Wallet) error {
	panic("unimplemented")
}
