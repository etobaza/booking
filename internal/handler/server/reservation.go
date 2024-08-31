package server

import (
	"conf_res/internal/models"
	"conf_res/internal/repository/postgres"
	"conf_res/pkg/server/response"
	"encoding/json"
	"github.com/go-chi/chi"
	"net/http"
)

type Handler struct {
	ReservationRepository *postgres.ReservationRepository
}

// New создает новый экземпляр обработчика для бронирований
func New(repo *postgres.ReservationRepository) *Handler {
	return &Handler{
		ReservationRepository: repo,
	}
}

// Routes возвращает маршрутизатор с маршрутами для бронирований
func (h *Handler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Post("/", h.CreateReservation)
	r.Get("/{room_id}", h.GetReservations)
	r.Get("/", h.GetAllReservations)
	return r
}

func (h *Handler) CreateReservation(w http.ResponseWriter, r *http.Request) {
	var reservation models.Reservation
	if err := json.NewDecoder(r.Body).Decode(&reservation); err != nil {
		response.BadRequest(w, r, err, nil)
		return
	}

	// Создаем бронирование
	err := h.ReservationRepository.CreateReservation(r.Context(), &reservation)
	if err != nil {
		if err.Error() == "time slot conflict" {
			response.Conflict(w, r, err, nil)
		} else {
			response.InternalServerError(w, r, err)
		}
		return
	}

	// Возвращаем 201 Created с созданным бронированием
	response.Created(w, r, reservation)
}

// GetReservations возвращает список бронирований для комнаты
func (h *Handler) GetReservations(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room_id")
	reservations, err := h.ReservationRepository.GetReservations(r.Context(), roomID)
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	// Возвращаем список бронирований
	response.OK(w, r, reservations)
}

// GetAllReservations возвращает список всех бронирований
func (h *Handler) GetAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := h.ReservationRepository.GetAllReservations(r.Context())
	if err != nil {
		response.InternalServerError(w, r, err)
		return
	}

	// Возвращаем список всех бронирований
	response.OK(w, r, reservations)
}
