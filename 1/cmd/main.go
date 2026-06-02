package main

import (
	"log"
	"main/internal/handler"
	"main/internal/repository"
	"main/internal/service"
	"net/http"
)

func main() {
	repo := repository.NewInMemoryRepository()
	svs := service.NewUserService(repo)
	h := handler.NewHttpHandler(svs)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", h.Create)
	mux.HandleFunc("PUT /users/{id}", h.Update)
	mux.HandleFunc("DELETE /users/{id}", h.Delete)
	mux.HandleFunc("GET /users/{id}", h.GetByID)
	mux.HandleFunc("GET /users", h.GetAll)

	log.Println("Starting HTTP server on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
