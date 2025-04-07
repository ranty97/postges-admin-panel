package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type DB struct {
	db *sqlx.DB
}

func NewDB(db *sqlx.DB) *DB {
	return &DB{db: db}
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

		var results []map[string]interface{}
		for rows.Next() {
			row := make(map[string]interface{})
			err := rows.MapScan(row)
			if err != nil {
				return "", fmt.Errorf("postgres: %w", err)
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
