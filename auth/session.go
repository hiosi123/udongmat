package auth

import (
	"github.com/go-redis/redis"
	"github.com/jmoiron/sqlx"
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

}
