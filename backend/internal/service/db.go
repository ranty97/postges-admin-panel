package service

import (
	"context"
	"errors"
	"fmt"
	"l6/internal/config"
	"l6/internal/domain"
	"os"
	"path/filepath"
)

type Repository interface {
	Tables(ctx context.Context) ([]string, error)
	ExecuteQuery(ctx context.Context, query string) (string, error)
	CreateBackup(ctx context.Context, dir string) (domain.BackupCreated, error)
	RestoreBackup(ctx context.Context, filename string, dir string) error
}

type Service struct {
	repo Repository
	cfg  *config.AppConfig
}

func NewDBService(repo Repository, cfg *config.AppConfig) *Service {
	return &Service{repo: repo, cfg: cfg}
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

func (s *Service) ListBackups(ctx context.Context) ([]domain.Backup, error) {

	dir, err := os.Open(s.cfg.BackupDir)
	if err != nil {
		return nil, fmt.Errorf("не удалось открыть директорию: %w", err)
	}
	defer dir.Close()

	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения директории: %w", err)
	}

	var backups []domain.Backup

	for _, fileInfo := range fileInfos {
		if !fileInfo.IsDir() {
			backup := domain.Backup{
				Filename:  fileInfo.Name(),
				CreatedAt: fileInfo.ModTime(),
				Size:      fileInfo.Size(),
			}
			backups = append(backups, backup)
		}
	}

	return backups, nil
}

func (s *Service) CreateBackup(ctx context.Context) (domain.BackupCreated, error) {
	backup, err := s.repo.CreateBackup(ctx, s.cfg.BackupDir)
	if err != nil {
		return domain.BackupCreated{}, fmt.Errorf("repo: %w", err)
	}
	return backup, nil
}

func (s *Service) DownloadBackup(ctx context.Context, filename string) ([]byte, error) {
	if filename == "" {
		return nil, errors.New("filename is required")
	}

	cleanFilename := filepath.Base(filename)
	fullPath := filepath.Join(s.cfg.BackupDir, cleanFilename)

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup file: %w", err)
	}

	return data, nil
}

func (s *Service) DeleteBackup(ctx context.Context, filename string) error {
	fullPath := filepath.Join(s.cfg.BackupDir, filename)
	err := os.Remove(fullPath)
	if err != nil {
		return fmt.Errorf("failed to delete backup file: %w", err)
	}
	return nil
}

func (s *Service) RestoreBackup(ctx context.Context, filename string) error {
	err := s.repo.RestoreBackup(ctx, filename, s.cfg.BackupDir)
	if err != nil {
		return fmt.Errorf("failed to restore backup file: %w", err)
	}
	return nil
}
