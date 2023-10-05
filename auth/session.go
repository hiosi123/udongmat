package auth

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"golang.org/x/crypto/bcrypt"
)

type SessionManager struct {
	Rdb  *redis.Client
	Conn *sqlx.DB
}

func NewSessionManager(rdb *redis.Client, conn *sqlx.DB) *SessionManager {
	return &SessionManager{Rdb: rdb, Conn: conn}
}

type AccountSession struct {
	Id    int
	Name  string
	Email string
}

type Account struct {
	Id       int
	Name     string
	Email    string
	Password string
}

func (s *SessionManager) GenerateSession(data AccountSession) (string, error) {
	sessionId := uuid.NewString()
	jsonData, _ := json.Marshal(data)
	err := s.Rdb.Set(sessionId, string(jsonData), 24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (s *SessionManager) SignIn(email, password string) (string, error) {
	// check if the account exits
	var account Account
	err := s.Conn.QueryRow(`SELECT id, name, email, password FROM account WHERE email = ?`,
		email).Scan(&account.Id, &account.Name, &account.Email, &account.Password)
	if err != nil {
		return "", err
	}

	// check if the password matches
	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil {
		return "", err
	}

	// create the session=
	sessionId := uuid.NewString()
	jsonData, _ := json.Marshal(AccountSession{
		Id:    account.Id,
		Name:  account.Name,
		Email: account.Email,
	})
	err = s.Rdb.Set(sessionId, string(jsonData), 24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (s *SessionManager) SignOut(sessionId string) error {
	return s.Rdb.Del(sessionId).Err()
}

func (s *SessionManager) GetSession(session string) (*AccountSession, error) {
	data, err := s.Rdb.Get(session).Result()
	if err != nil {
		return nil, err
	}

	// unmarshal the data
	var accountSession AccountSession
	err = json.Unmarshal([]byte(data), &accountSession)
	if err != nil {
		return nil, err
	}

	return &accountSession, nil
}
