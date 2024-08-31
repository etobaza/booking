package main

import (
	"conf_res/internal/handler"
	"conf_res/internal/repository"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// Инициализация репозиториев
	repositories, err := repository.New(
		repository.WithPostgresStore(os.Getenv("DATABASE_URL")),
	)
	if err != nil {
		log.Fatalf("Failed to initialize repositories: %v", err)
	}
	defer repositories.Close()

	// Инициализация обработчика
	h := handler.New(repositories)

	// Запуск сервера
	srv := &http.Server{
		Handler:      h.Router(),
		Addr:         "0.0.0.0:8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}
