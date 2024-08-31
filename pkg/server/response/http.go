package response

import (
	"github.com/go-chi/render"
	"net/http"
)

type Object struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

func OK(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, http.StatusOK)
	render.JSON(w, r, data)
}

func BadRequest(w http.ResponseWriter, r *http.Request, err error, data interface{}) {
	render.Status(r, http.StatusBadRequest)

	v := Object{
		Success: false,
		Data:    data,
		Message: err.Error(),
	}
	render.JSON(w, r, v)
}

func Conflict(w http.ResponseWriter, r *http.Request, err error, data interface{}) {
	render.Status(r, http.StatusConflict)

	v := Object{
		Success: false,
		Data:    data,
		Message: err.Error(),
	}
	render.JSON(w, r, v)
}

func Created(w http.ResponseWriter, r *http.Request, data interface{}) {
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, data)
}

func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	render.Status(r, http.StatusInternalServerError)

	v := Object{
		Success: false,
		Message: err.Error(),
	}
	render.JSON(w, r, v)
}
