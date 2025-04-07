package service

import (
	"fmt"
)

type Repository interface {
	Tables() ([]string, error)
	ExecuteQuery(query string) (string, error)
}

type Service struct {
	repo Repository
}

func NewDBService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Tables() ([]string, error) {
	tables, err := s.repo.Tables()
	if err != nil {
		return nil, fmt.Errorf("repo: %w", err)
	}
	return tables, nil
}

func (s *Service) ExecuteQuery(query string) (string, error) {
	result, err := s.repo.ExecuteQuery(query)
	if err != nil {
		return "", fmt.Errorf("repo: %w", err)
	}
	return result, nil
}
