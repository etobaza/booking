package repository

import (
	"conf_res/internal/repository/postgres"
	"context"
	"fmt"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
	"path/filepath"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repositories struct {
	PostgresDB      *pgxpool.Pool
	ReservationRepo *postgres.ReservationRepository
	// Другие репозитории можно добавить сюда
}

type Option func(*Repositories) error

func New(opts ...Option) (repo *Repositories, err error) {
	r := &Repositories{}

	for _, opt := range opts {
		if err := opt(r); err != nil {
			return nil, err
		}
	}

	return r, nil
}

func WithPostgresStore(dsn string) (opt Option) {
	return func(r *Repositories) (err error) {
		// Подключение к базе данных
		r.PostgresDB, err = pgxpool.Connect(context.Background(), dsn)
		if err != nil {
			return
		}

		// Применение миграций
		if err = applyMigrations(r.PostgresDB); err != nil {
			return
		}

		// Инициализация репозиториев
		r.ReservationRepo = postgres.NewReservationRepository(r.PostgresDB)

		return
	}
}

func (r *Repositories) Close() {
	if r.PostgresDB != nil {
		r.PostgresDB.Close()
	}
	log.Println("Database connection closed.")
}

func applyMigrations(dbpool *pgxpool.Pool) (err error) {
	// Определение директории миграций
	migrationDir := "migrations"

	// Чтение файлов миграций
	files, err := os.ReadDir(migrationDir)
	if err != nil {
		return fmt.Errorf("failed to read migration directory: %w", err)
	}

	// Применение миграций
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			err := executeMigration(dbpool, filepath.Join(migrationDir, file.Name()))
			if err != nil {
				return fmt.Errorf("failed to execute migration %s: %w", file.Name(), err)
			}
		}
	}

	log.Println("Migrations applied successfully")
	return
}

func executeMigration(dbpool *pgxpool.Pool, filepath string) (err error) {
	// Чтение содержимого SQL файла
	sqlBytes, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read migration file: %w", err)
	}
	sql := string(sqlBytes)

	// Выполнение миграции в транзакции
	tx, err := dbpool.Begin(context.Background())
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}(tx, context.Background()) // Обеспечение отката транзакции в случае ошибки

	_, err = tx.Exec(context.Background(), sql)
	if err != nil {
		return fmt.Errorf("failed to execute migration: %w", err)
	}

	// Фиксация транзакции
	err = tx.Commit(context.Background())
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return
}
