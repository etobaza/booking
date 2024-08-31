package tests

import (
	"conf_res/internal/models"
	"conf_res/internal/repository/postgres"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupDB() (*postgres.ReservationRepository, func()) {
	// Подключаемся к тестовой базе данных
	dbpool, err := pgxpool.Connect(context.Background(), "postgres://user:password@localhost:5432/booking_test?sslmode=disable")
	if err != nil {
		panic(err)
	}

	repo := postgres.NewReservationRepository(dbpool)

	// Функция очистки базы данных после тестов
	cleanup := func() {
		_, err := dbpool.Exec(context.Background(), "TRUNCATE TABLE reservations RESTART IDENTITY CASCADE")
		if err != nil {
			panic(err)
		}
		dbpool.Close()
	}

	return repo, cleanup
}

func TestCreateReservation_Success(t *testing.T) {
	repo, cleanup := setupDB()
	defer cleanup()

	ctx := context.Background()

	// Добавляем тестовые данные
	reservation := &models.Reservation{
		RoomID:    "room_1",
		StartTime: time.Now().Add(1 * time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
	}

	err := repo.CreateReservation(ctx, reservation)
	require.NoError(t, err)
	assert.NotEmpty(t, reservation.ID)
}

func TestCreateReservation_Conflict(t *testing.T) {
	repo, cleanup := setupDB()
	defer cleanup()

	ctx := context.Background()

	// Добавляем первое бронирование
	firstReservation := &models.Reservation{
		RoomID:    "room_1",
		StartTime: time.Now().Add(1 * time.Hour),
		EndTime:   time.Now().Add(2 * time.Hour),
	}

	err := repo.CreateReservation(ctx, firstReservation)
	require.NoError(t, err)

	// Пытаемся добавить пересекающееся бронирование
	conflictingReservation := &models.Reservation{
		RoomID:    "room_1",
		StartTime: time.Now().Add(90 * time.Minute), // Пересекается по времени
		EndTime:   time.Now().Add(3 * time.Hour),
	}

	err = repo.CreateReservation(ctx, conflictingReservation)
	require.Error(t, err)
	assert.Equal(t, "time slot conflict", err.Error())
}

func TestCreateReservation_Concurrent(t *testing.T) {
	repo, cleanup := setupDB()
	defer cleanup()

	ctx := context.Background()

	startTime := time.Now().Add(1 * time.Hour)
	endTime := time.Now().Add(2 * time.Hour)

	// Параллельно создаем бронирования
	done := make(chan error, 2)
	for i := 0; i < 2; i++ {
		go func() {
			reservation := &models.Reservation{
				RoomID:    "room_1",
				StartTime: startTime,
				EndTime:   endTime,
			}
			err := repo.CreateReservation(ctx, reservation)
			done <- err
		}()
	}

	// Ожидаем завершения обеих горутин
	firstErr := <-done
	secondErr := <-done

	// Один из запросов должен пройти успешно, второй должен вернуть ошибку
	assert.NoError(t, firstErr)
	assert.Error(t, secondErr)
	assert.Equal(t, "time slot conflict", secondErr.Error())
}
