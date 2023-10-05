package storage

import "github.com/jmoiron/sqlx"

type AccountStorage struct {
	Conn *sqlx.DB
}

func NewAccountStorage(conn *sqlx.DB) *AccountStorage {
	return &AccountStorage{Conn: conn}
}

type NewAccount struct {
	Name     string
	Email    string
	Password string
}

func (s *AccountStorage) CreateNewAccount(data NewAccount) (int, error) {
	res, err := s.Conn.Exec("insert into account (name, email, password) values (?, ?, ?)", data.Name, data.Email, data.Password)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}
