package postgres

import (
	"conf_res/internal/models"
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
)

type ReservationRepository struct {
	db *pgxpool.Pool
}

func NewReservationRepository(db *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{db: db}
}

func (s *ReservationRepository) CreateReservation(ctx context.Context, reservation *models.Reservation) (err error) {
	var count int
	err = s.db.QueryRow(ctx,
		`SELECT COUNT(*) FROM reservations 
         WHERE room_id=$1 AND 
         ($2 < end_time AND $3 > start_time)`,
		reservation.RoomID, reservation.StartTime, reservation.EndTime).
		Scan(&count)
	if err != nil {
		return
	}

	if count > 0 {
		return errors.New("time slot conflict")
	}

	// Выполняем вставку новой записи
	err = s.db.QueryRow(ctx,
		`INSERT INTO reservations (room_id, start_time, end_time) 
         VALUES ($1, $2, $3) RETURNING id`,
		reservation.RoomID, reservation.StartTime, reservation.EndTime).
		Scan(&reservation.ID) // Получаем сгенерированный ID
	if err != nil {
		return
	}

	return
}

func (s *ReservationRepository) GetReservations(ctx context.Context, roomID string) (res []models.Reservation, err error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, room_id, start_time, end_time 
         FROM reservations 
         WHERE room_id=$1`, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []models.Reservation
	for rows.Next() {
		var r models.Reservation
		if err = rows.Scan(&r.ID, &r.RoomID, &r.StartTime, &r.EndTime); err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}

	return reservations, nil
}

func (s *ReservationRepository) GetAllReservations(ctx context.Context) (res []models.Reservation, err error) {
	rows, err := s.db.Query(ctx,
		`SELECT id, room_id, start_time, end_time 
         FROM reservations`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reservations []models.Reservation
	for rows.Next() {
		var r models.Reservation
		if err = rows.Scan(&r.ID, &r.RoomID, &r.StartTime, &r.EndTime); err != nil {
			return nil, err
		}
		reservations = append(reservations, r)
	}

	return reservations, nil
}
