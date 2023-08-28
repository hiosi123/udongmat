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

type UserSession struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
}

type User struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
	Password  string
}

func (s *SessionManager) GenerateSession(data UserSession) (string, error) {
	sessionId := uuid.NewString()
	jsonData, _ := json.Marshal(data)
	err := s.Rdb.Set(sessionId, string(jsonData), 24*time.Hour).Err()
	if err != nil {
		return "", err
	}

	return sessionId, nil
}

func (s *SessionManager) SignIn(email, password string) (string, error) {
	// check if the user exits
	var user User
	err := s.Conn.QueryRow(`SELECT id, first_name, last_name, email, password FROM users WHERE email = ?`,
		email).Scan(&user.Id, &user.FirstName, &user.LastName, &user.Email, &user.Password)
	if err != nil {
		return "", err
	}

	// check if the password matches
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	// create the session=
	sessionId := uuid.NewString()
	jsonData, _ := json.Marshal(UserSession{
		Id:        user.Id,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
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

func (s *SessionManager) GetSession(session string) (*UserSession, error) {
	data, err := s.Rdb.Get(session).Result()
	if err != nil {
		return nil, err
	}

	// unmarshal the data
	var userSession UserSession
	err = json.Unmarshal([]byte(data), &userSession)
	if err != nil {
		return nil, err
	}

	return &userSession, nil
}
