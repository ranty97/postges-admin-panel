package service

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DBService struct {
	db *sqlx.DB
}

func NewDBService(db *sqlx.DB) *DBService {
	return &DBService{db: db}
}

func (s *DBService) Tables() ([]string, error) {
	var tables []string
	err := s.db.Select(&tables, "SELECT table_name FROM information_schema.tables WHERE table_schema='public'")
	if err != nil {
		return nil, err
	}
	return tables, nil
}

func (s *DBService) ExecuteQuery(query string) (string, error) {
	result, err := s.db.Exec(query)
	if err != nil {
		return "", err
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Query executed successfully. %d rows affected.", rows), nil
}
