package service

import (
	"context"
	"fmt"
)

type Repository interface {
	Tables(ctx context.Context) ([]string, error)
	ExecuteQuery(ctx context.Context, query string) (string, error)
}

type Service struct {
	repo Repository
}

func NewDBService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) Tables(ctx context.Context) ([]string, error) {
	tables, err := s.repo.Tables(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo: %w", err)
	}
	return tables, nil
}

func (s *Service) ExecuteQuery(ctx context.Context, query string) (string, error) {
	result, err := s.repo.ExecuteQuery(ctx, query)
	if err != nil {
		return "", fmt.Errorf("repo: %w", err)
	}
	return result, nil
}
