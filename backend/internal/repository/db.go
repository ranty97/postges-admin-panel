package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	cfg "l6/internal/config"
	"l6/internal/domain"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	db  *sqlx.DB
	cfg *cfg.PostgresConfig
}

func NewDB(db *sqlx.DB, cfg *cfg.PostgresConfig) *DB {
	return &DB{db: db, cfg: cfg}
}

func (d *DB) Tables(ctx context.Context) ([]string, error) {

	query := `
		SELECT table_name 
		FROM information_schema.tables 
		WHERE table_schema='public'
	`

	var tables []string
	err := d.db.SelectContext(ctx, &tables, query)
	if err != nil {
		return nil, fmt.Errorf("postgres: %w", err)
	}
	return tables, nil
}

func (d *DB) ExecuteQuery(ctx context.Context, query string) (string, error) {
	if len(query) >= 6 && strings.ToUpper(query[:6]) == "SELECT" {
		rows, err := d.db.QueryxContext(ctx, query)
		if err != nil {
			return "", fmt.Errorf("postgres: %w", err)
		}
		defer rows.Close()

		columnTypes, err := rows.ColumnTypes()
		if err != nil {
			return "", fmt.Errorf("postgres: get column types: %w", err)
		}

		var results []map[string]interface{}
		for rows.Next() {
			row := make(map[string]interface{})
			err := rows.MapScan(row)
			if err != nil {
				return "", fmt.Errorf("postgres: %w", err)
			}

			for key, value := range row {
				if byteArray, ok := value.([]byte); ok {
					dataType := ""
					for _, col := range columnTypes {
						if col.Name() == key {
							dataType = col.DatabaseTypeName()
							break
						}
					}

					switch dataType {
					case "NUMERIC", "DECIMAL", "FLOAT", "REAL", "DOUBLE PRECISION":
						floatVal, err := strconv.ParseFloat(string(byteArray), 64)
						if err == nil {
							row[key] = floatVal
						} else {
							row[key] = string(byteArray)
						}
					case "INTEGER", "BIGINT", "SMALLINT":
						intVal, err := strconv.ParseInt(string(byteArray), 10, 64)
						if err == nil {
							row[key] = intVal
						} else {
							row[key] = string(byteArray)
						}
					default:
						row[key] = string(byteArray)
					}
				}
			}

			results = append(results, row)
		}

		jsonData, err := json.Marshal(results)
		if err != nil {
			return "", fmt.Errorf("json marshal: %w", err)
		}
		return string(jsonData), nil
	}

	result, err := d.db.ExecContext(ctx, query)
	if err != nil {
		return "", fmt.Errorf("postgres: %w", err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("postgres: %w", err)
	}
	return fmt.Sprintf("rows affected: %d", rows), nil
}

func (d *DB) CreateBackup(ctx context.Context, dir string) (domain.BackupCreated, error) {
	if dir == "" {
		dir = os.Getenv("BACKUP_DIR")
		if dir == "" {
			return domain.BackupCreated{}, errors.New("backup directory not provided and BACKUP_DIR is not set")
		}
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return domain.BackupCreated{}, fmt.Errorf("failed to create backup directory: %w", err)
	}

	timestamp := time.Now().Format("2006-01-02_15-04-05")
	filename := fmt.Sprintf("%s_backup_%s.sql", d.cfg.Database, timestamp)
	filePath := filepath.Join(dir, filename)

	cmd := exec.CommandContext(
		ctx,
		"pg_dump",
		"-h", d.cfg.Host,
		"-p", d.cfg.Port,
		"-U", d.cfg.Username,
		"-f", filePath,
		d.cfg.Database,
	)

	if d.cfg.Password != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", d.cfg.Password))
	}

	if output, err := cmd.CombinedOutput(); err != nil {
		return domain.BackupCreated{}, fmt.Errorf("pg_dump failed: %w, output: %s", err, string(output))
	}

	return domain.BackupCreated{
		Filename: filename,
		Message:  "Backup created successfully",
		Success:  true,
	}, nil
}

func (d *DB) RestoreBackup(ctx context.Context, filename string, dir string) error {
	if filename == "" {
		return errors.New("filename is required")
	}

	cleanFilename := filepath.Base(filename)
	fullPath := filepath.Join(dir, cleanFilename)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file does not exist: %w", err)
	}

	cmd := exec.CommandContext(
		ctx,
		"psql",
		"-h", d.cfg.Host,
		"-p", d.cfg.Port,
		"-U", d.cfg.Username,
		"-d", d.cfg.Database,
		"-f", fullPath,
	)

	if d.cfg.Password != "" {
		cmd.Env = append(os.Environ(), fmt.Sprintf("PGPASSWORD=%s", d.cfg.Password))
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("pg_restore failed: %w, output: %s", err, string(output))
	}

	return nil
}

func (d *DB) DeleteAllTables(ctx context.Context) error {
	query := `
		DO $$ 
		DECLARE 
    		r RECORD;
		BEGIN 
    		EXECUTE 'SET session_replication_role = replica';

    		FOR r IN (SELECT tablename FROM pg_tables WHERE schemaname = 'public') 
    		LOOP
        		EXECUTE 'DROP TABLE IF EXISTS public.' || r.tablename || ' CASCADE';
    		END LOOP;

    		EXECUTE 'SET session_replication_role = DEFAULT';
		END $$;`

	_, err := d.db.ExecContext(ctx, query)
	if err != nil {
		return fmt.Errorf("failed to delete all tables: %w", err)
	}
	return nil
}
