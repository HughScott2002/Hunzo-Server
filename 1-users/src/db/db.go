// db/db.go

package db

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"example.com/m/v2/src/models"
	"github.com/go-redis/redis/v8"
)

var db Database

type Database interface {
	AddUser(user *models.User) error
	GetUser(email string) (*models.User, error)
	GetUserByAccountId(accountId string) (*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(email string) error
	UserExists(email string) (bool, error)
	AddSession(session *models.Session) error
	GetSession(id string) (*models.Session, error)
	GetUserSessions(email string) ([]*models.Session, error)
	DeleteSession(id string) error
	DeleteUserSessions(email string) error
	UpdateSessionLastLogin(id string) error
	AddRefreshToken(token string, info RefreshTokenInfo) error
	GetRefreshToken(token string) (*RefreshTokenInfo, error)
	DeleteRefreshToken(token string) error
}

type RefreshTokenInfo struct {
	UserEmail  string
	DeviceInfo string
	CreatedAt  time.Time
}

type MemoryDB struct {
	users         map[string]models.User
	sessions      map[string]models.Session
	refreshTokens map[string]RefreshTokenInfo
	mu            sync.RWMutex
}

type RedisDB struct {
	client *redis.Client
}

type FutureDB struct {
	// Placeholder for future database implementation
}

func Init() error {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	mode := strings.ToLower(os.Getenv("MODE"))

	switch {
	case env == "local" && mode == "memcached":
		db = &MemoryDB{
			users:         make(map[string]models.User),
			sessions:      make(map[string]models.Session),
			refreshTokens: make(map[string]RefreshTokenInfo),
		}
	case env == "local" && mode != "memcached":
		redisPassword := os.Getenv("USER_REDIS_PASSWORD")
		redisAddr := fmt.Sprintf("user-redis:%s", os.Getenv("USER_REDIS_PORT"))
		redisClient := redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       0,
		})
		db = &RedisDB{client: redisClient}
		log.Println("USING DB & REDIS IN USER SERIVCE")
	case env == "local" && mode == "db":
		// Placeholder for future database implementation
		db = &FutureDB{}
	case env == "prod" || env == "production":
		// Placeholder for future database implementation
		db = &FutureDB{}
	default:
		return fmt.Errorf("unsupported environment or mode")
	}

	return nil
}

// Helper functions to call the database interface methods
func AddUser(user *models.User) error {
	exists, err := db.UserExists(user.Email)
	if err != nil {
		return fmt.Errorf("error checking user existence: %v", err)
	}
	if exists {
		return fmt.Errorf("user already exists")
	}
	return db.AddUser(user)
}

func GetUser(email string) (*models.User, error) { return db.GetUser(email) }
func GetUserByAccountId(accountId string) (*models.User, error) {
	return db.GetUserByAccountId(accountId)
}
func UpdateUser(user *models.User) error                      { return db.UpdateUser(user) }
func DeleteUser(email string) error                           { return db.DeleteUser(email) }
func UserExists(email string) (bool, error)                   { return db.UserExists(email) }
func AddSession(session *models.Session) error                { return db.AddSession(session) }
func GetSession(id string) (*models.Session, error)           { return db.GetSession(id) }
func GetUserSessions(email string) ([]*models.Session, error) { return db.GetUserSessions(email) }
func DeleteSession(id string) error                           { return db.DeleteSession(id) }
func DeleteUserSessions(email string) error                   { return db.DeleteUserSessions(email) }
func UpdateSessionLastLogin(id string) error                  { return db.UpdateSessionLastLogin(id) }
func AddRefreshToken(token string, info RefreshTokenInfo) error {
	return db.AddRefreshToken(token, info)
}
func GetRefreshToken(token string) (*RefreshTokenInfo, error) { return db.GetRefreshToken(token) }
func DeleteRefreshToken(token string) error                   { return db.DeleteRefreshToken(token) }

// MemoryDB implementations

func (m *MemoryDB) AddUser(user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.users[user.Email] = *user
	return nil
}

func (m *MemoryDB) GetUser(email string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	user, ok := m.users[email]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}
	return &user, nil
}
func (m *MemoryDB) GetUserByAccountId(accountId string) (*models.User, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, user := range m.users {
		if user.AccountId == accountId {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (m *MemoryDB) UpdateUser(user *models.User) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[user.Email]; !ok {
		return fmt.Errorf("user not found")
	}
	m.users[user.Email] = *user
	return nil
}

func (m *MemoryDB) DeleteUser(email string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[email]; !ok {
		return fmt.Errorf("user not found")
	}
	delete(m.users, email)
	return nil
}

func (m *MemoryDB) UserExists(email string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	_, exists := m.users[email]
	return exists, nil
}

func (m *MemoryDB) AddSession(session *models.Session) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.sessions[session.ID] = *session
	return nil
}

func (m *MemoryDB) GetSession(id string) (*models.Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	session, ok := m.sessions[id]
	if !ok {
		return nil, fmt.Errorf("session not found")
	}
	return &session, nil
}

func (m *MemoryDB) GetUserSessions(email string) ([]*models.Session, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	var userSessions []*models.Session
	for _, session := range m.sessions {
		if session.UserEmail == email {
			sessionCopy := session
			userSessions = append(userSessions, &sessionCopy)
		}
	}
	return userSessions, nil
}

func (m *MemoryDB) DeleteSession(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.sessions[id]; !ok {
		return fmt.Errorf("session not found")
	}
	delete(m.sessions, id)
	return nil
}

func (m *MemoryDB) DeleteUserSessions(email string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for id, session := range m.sessions {
		if session.UserEmail == email {
			delete(m.sessions, id)
		}
	}
	return nil
}

func (m *MemoryDB) UpdateSessionLastLogin(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	session, ok := m.sessions[id]
	if !ok {
		return fmt.Errorf("session not found")
	}
	session.LastLoginAt = time.Now()
	m.sessions[id] = session
	return nil
}

func (m *MemoryDB) AddRefreshToken(token string, info RefreshTokenInfo) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.refreshTokens[token] = info
	return nil
}

func (m *MemoryDB) GetRefreshToken(token string) (*RefreshTokenInfo, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	info, ok := m.refreshTokens[token]
	if !ok {
		return nil, fmt.Errorf("refresh token not found")
	}
	return &info, nil
}

func (m *MemoryDB) DeleteRefreshToken(token string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.refreshTokens[token]; !ok {
		return fmt.Errorf("refresh token not found")
	}
	delete(m.refreshTokens, token)
	return nil
}

// RedisDB implementations

func (r *RedisDB) AddUser(user *models.User) error {
	ctx := context.Background()
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, "user:"+user.Email, userJSON, 0).Err()
}

func (r *RedisDB) GetUser(email string) (*models.User, error) {
	ctx := context.Background()
	userJSON, err := r.client.Get(ctx, "user:"+email).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	var user models.User
	err = json.Unmarshal(userJSON, &user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (r *RedisDB) GetUserByAccountId(accountId string) (*models.User, error) {
	ctx := context.Background()
	pattern := "user:*"
	iter := r.client.Scan(ctx, 0, pattern, 0).Iterator()
	for iter.Next(ctx) {
		userJSON, err := r.client.Get(ctx, iter.Val()).Bytes()
		if err != nil {
			continue
		}
		var user models.User
		if err := json.Unmarshal(userJSON, &user); err != nil {
			continue
		}
		if user.AccountId == accountId {
			return &user, nil
		}
	}
	return nil, fmt.Errorf("user not found")
}

func (r *RedisDB) UpdateUser(user *models.User) error {
	ctx := context.Background()
	exists, err := r.UserExists(user.Email)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user not found")
	}
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, "user:"+user.Email, userJSON, 0).Err()
}

func (r *RedisDB) DeleteUser(email string) error {
	ctx := context.Background()
	exists, err := r.UserExists(email)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("user not found")
	}
	return r.client.Del(ctx, "user:"+email).Err()
}

func (r *RedisDB) UserExists(email string) (bool, error) {
	ctx := context.Background()
	exists, err := r.client.Exists(ctx, "user:"+email).Result()
	if err != nil {
		return false, err
	}
	return exists == 1, nil
}

func (r *RedisDB) AddSession(session *models.Session) error {
	ctx := context.Background()
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	}
	pipe := r.client.Pipeline()
	pipe.Set(ctx, "session:"+session.ID, sessionJSON, 24*time.Hour)
	pipe.SAdd(ctx, "user_sessions:"+session.UserEmail, session.ID)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisDB) GetSession(id string) (*models.Session, error) {
	ctx := context.Background()
	sessionJSON, err := r.client.Get(ctx, "session:"+id).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("session not found")
		}
		return nil, err
	}
	var session models.Session
	err = json.Unmarshal(sessionJSON, &session)
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *RedisDB) GetUserSessions(email string) ([]*models.Session, error) {
	ctx := context.Background()
	sessionIDs, err := r.client.SMembers(ctx, "user_sessions:"+email).Result()
	if err != nil {
		return nil, err
	}
	var userSessions []*models.Session
	for _, id := range sessionIDs {
		session, err := r.GetSession(id)
		if err != nil {
			continue // Skip sessions that can't be retrieved
		}
		userSessions = append(userSessions, session)
	}
	return userSessions, nil
}

func (r *RedisDB) DeleteSession(id string) error {
	ctx := context.Background()
	session, err := r.GetSession(id)
	if err != nil {
		return err
	}
	pipe := r.client.Pipeline()
	pipe.Del(ctx, "session:"+id)
	pipe.SRem(ctx, "user_sessions:"+session.UserEmail, id)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisDB) DeleteUserSessions(email string) error {
	ctx := context.Background()
	sessionIDs, err := r.client.SMembers(ctx, "user_sessions:"+email).Result()
	if err != nil {
		return err
	}
	pipe := r.client.Pipeline()
	for _, id := range sessionIDs {
		pipe.Del(ctx, "session:"+id)
	}
	pipe.Del(ctx, "user_sessions:"+email)
	_, err = pipe.Exec(ctx)
	return err
}

func (r *RedisDB) UpdateSessionLastLogin(id string) error {
	// _ := context.Background()
	session, err := r.GetSession(id)
	if err != nil {
		return err
	}
	session.LastLoginAt = time.Now()
	return r.AddSession(session) // This will overwrite the existing session
}

func (r *RedisDB) AddRefreshToken(token string, info RefreshTokenInfo) error {
	ctx := context.Background()
	infoJSON, err := json.Marshal(info)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, "refresh_token:"+token, infoJSON, 7*24*time.Hour).Err()
}

func (r *RedisDB) GetRefreshToken(token string) (*RefreshTokenInfo, error) {
	ctx := context.Background()
	infoJSON, err := r.client.Get(ctx, "refresh_token:"+token).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("refresh token not found")
		}
		return nil, err
	}
	var info RefreshTokenInfo
	err = json.Unmarshal(infoJSON, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}

func (r *RedisDB) DeleteRefreshToken(token string) error {
	ctx := context.Background()
	return r.client.Del(ctx, "refresh_token:"+token).Err()
}

// FutureDB implementations (placeholders)

func (f *FutureDB) AddUser(user *models.User) error {
	return fmt.Errorf("FutureDB: AddUser not implemented")
}

func (f *FutureDB) GetUser(email string) (*models.User, error) {
	return nil, fmt.Errorf("FutureDB: GetUser not implemented")
}
func (f *FutureDB) GetUserByAccountId(accountId string) (*models.User, error) {
	return nil, fmt.Errorf("FutureDB: GetUserByAccountId not implemented")
}

func (f *FutureDB) UpdateUser(user *models.User) error {
	return fmt.Errorf("FutureDB: UpdateUser not implemented")
}

func (f *FutureDB) DeleteUser(email string) error {
	return fmt.Errorf("FutureDB: DeleteUser not implemented")
}

func (f *FutureDB) UserExists(email string) (bool, error) {
	return false, fmt.Errorf("FutureDB: UserExists not implemented")
}

func (f *FutureDB) AddSession(session *models.Session) error {
	return fmt.Errorf("FutureDB: AddSession not implemented")
}

func (f *FutureDB) GetSession(id string) (*models.Session, error) {
	return nil, fmt.Errorf("FutureDB: GetSession not implemented")
}

func (f *FutureDB) GetUserSessions(email string) ([]*models.Session, error) {
	return nil, fmt.Errorf("FutureDB: GetUserSessions not implemented")
}

func (f *FutureDB) DeleteSession(id string) error {
	return fmt.Errorf("FutureDB: DeleteSession not implemented")
}

func (f *FutureDB) DeleteUserSessions(email string) error {
	return fmt.Errorf("FutureDB: DeleteUserSessions not implemented")
}

func (f *FutureDB) UpdateSessionLastLogin(id string) error {
	return fmt.Errorf("FutureDB: UpdateSessionLastLogin not implemented")
}

func (f *FutureDB) AddRefreshToken(token string, info RefreshTokenInfo) error {
	return fmt.Errorf("FutureDB: AddRefreshToken not implemented")
}

func (f *FutureDB) GetRefreshToken(token string) (*RefreshTokenInfo, error) {
	return nil, fmt.Errorf("FutureDB: GetRefreshToken not implemented")
}

func (f *FutureDB) DeleteRefreshToken(token string) error {
	return fmt.Errorf("FutureDB: DeleteRefreshToken not implemented")
}
