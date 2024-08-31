package handler

import (
	"conf_res/internal/handler/server"
	"conf_res/internal/repository"
	"github.com/go-chi/chi"
	"net/http"
)

type Handler struct {
	router http.Handler
}

// New создает новый экземпляр обработчика
func New(repositories *repository.Repositories) *Handler {
	// Создаем обработчик для бронирований
	reservationHandler := server.New(repositories.ReservationRepo)

	// Устанавливаем маршруты
	r := chi.NewRouter()
	r.Mount("/reservations", reservationHandler.Routes())

	// Возвращаем обработчик
	return &Handler{
		router: r,
	}
}

// Router возвращает маршрутизатор
func (h *Handler) Router() http.Handler {
	return h.router
}
