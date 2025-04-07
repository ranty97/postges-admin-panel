package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
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

		// Получаем типы колонок
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

			// Обработка значений перед добавлением в результат
			for key, value := range row {
				// Если значение - байтовый массив, преобразуем его в число или строку
				if byteArray, ok := value.([]byte); ok {
					// Определяем тип колонки
					dataType := ""
					for _, col := range columnTypes {
						if col.Name() == key {
							dataType = col.DatabaseTypeName()
							break
						}
					}

					// В зависимости от типа данных в PostgreSQL выполняем разное преобразование
					switch dataType {
					case "NUMERIC", "DECIMAL", "FLOAT", "REAL", "DOUBLE PRECISION":
						// Преобразуем в число
						floatVal, err := strconv.ParseFloat(string(byteArray), 64)
						if err == nil {
							row[key] = floatVal
						} else {
							// Если не удалось преобразовать в число, оставляем как строку
							row[key] = string(byteArray)
						}
					case "INTEGER", "BIGINT", "SMALLINT":
						// Преобразуем в целое число
						intVal, err := strconv.ParseInt(string(byteArray), 10, 64)
						if err == nil {
							row[key] = intVal
						} else {
							row[key] = string(byteArray)
						}
					default:
						// Для других типов данных преобразуем в строку
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
