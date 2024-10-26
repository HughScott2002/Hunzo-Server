package implementations

import (
	"fmt"
	"sync"
	"time"

	"example.com/m/v2/src/models"
)

// MemoryImplementation implements the Database interface using in-memory storage
type MemoryImplementation struct {
	wallets map[string]*models.Wallet
	mu      sync.RWMutex
}

// NewMemoryImplementation creates a new memory implementation
func NewMemoryImplementation() *MemoryImplementation {
	return &MemoryImplementation{
		wallets: make(map[string]*models.Wallet),
	}
}

func (m *MemoryImplementation) AddWallet(wallet *models.Wallet) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.wallets[wallet.WalletId]; exists {
		return fmt.Errorf("wallet already exists")
	}

	m.wallets[wallet.WalletId] = wallet
	return nil
}

func (m *MemoryImplementation) GetWallet(id string) (*models.Wallet, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	wallet, exists := m.wallets[id]
	if !exists {
		return nil, fmt.Errorf("wallet not found")
	}

	return wallet, nil
}

func (m *MemoryImplementation) UpdateWallet(wallet *models.Wallet) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.wallets[wallet.WalletId]; !exists {
		return fmt.Errorf("wallet not found")
	}

	wallet.UpdatedAt = time.Now()
	m.wallets[wallet.WalletId] = wallet
	return nil
}

func (m *MemoryImplementation) WalletExists(id string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, exists := m.wallets[id]
	return exists, nil
}

func (m *MemoryImplementation) DeleteWallet(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.wallets[id]; !exists {
		return fmt.Errorf("wallet not found")
	}

	delete(m.wallets, id)
	return nil
}

func (m *MemoryImplementation) ListWallets(accountId string) ([]*models.Wallet, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var wallets []*models.Wallet
	for _, wallet := range m.wallets {
		if wallet.AccountId == accountId {
			wallets = append(wallets, wallet)
		}
	}

	return wallets, nil
}

func (m *MemoryImplementation) UpdateWalletStatus(id string, status models.WalletStatus) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	wallet, exists := m.wallets[id]
	if !exists {
		return fmt.Errorf("wallet not found")
	}

	wallet.Status = status
	wallet.UpdatedAt = time.Now()
	return nil
}

func (m *MemoryImplementation) UpdateWalletBalance(id string, balance float64) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	wallet, exists := m.wallets[id]
	if !exists {
		return fmt.Errorf("wallet not found")
	}

	wallet.Balance = balance
	wallet.UpdatedAt = time.Now()
	now := time.Now()
	wallet.LastActivity = &now
	return nil
}

func (m *MemoryImplementation) GetDefaultWallet(accountId string) (*models.Wallet, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, wallet := range m.wallets {
		if wallet.AccountId == accountId && wallet.IsDefault {
			return wallet, nil
		}
	}

	return nil, fmt.Errorf("no default wallet found")
}

func (m *MemoryImplementation) SetDefaultWallet(accountId string, walletId string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Verify wallet exists and belongs to account
	wallet, exists := m.wallets[walletId]
	if !exists {
		return fmt.Errorf("wallet not found")
	}
	if wallet.AccountId != accountId {
		return fmt.Errorf("wallet does not belong to account")
	}

	// Remove default status from all other wallets for this account
	for _, w := range m.wallets {
		if w.AccountId == accountId {
			w.IsDefault = false
			w.UpdatedAt = time.Now()
		}
	}

	// Set the new default wallet
	wallet.IsDefault = true
	wallet.UpdatedAt = time.Now()
	m.wallets[walletId] = wallet

	return nil
}

func (m *MemoryImplementation) FreezeWallet(accountId string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	updated := false
	for _, wallet := range m.wallets {
		if wallet.AccountId == accountId {
			wallet.Status = models.WalletStatusFrozen
			wallet.UpdatedAt = time.Now()
			updated = true
		}
	}

	if !updated {
		return fmt.Errorf("no wallets found for account")
	}

	return nil
}

func (m *MemoryImplementation) UnfreezeWallet(accountId string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	updated := false
	for _, wallet := range m.wallets {
		if wallet.AccountId == accountId {
			wallet.Status = models.WalletStatusActive
			wallet.UpdatedAt = time.Now()
			updated = true
		}
	}

	if !updated {
		return fmt.Errorf("no wallets found for account")
	}

	return nil
}
