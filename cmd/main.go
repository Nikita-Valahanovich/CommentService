package main

import (
	"CommentService/pkg/handler"
	"CommentService/pkg/middleware"
	"CommentService/pkg/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	db, err := storage.NewDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	r := mux.NewRouter()

	// Подключаем middleware
	r.Use(middleware.RequestIDMiddleware)
	r.Use(middleware.LoggingMiddleware)

	h := handlers.NewHandler(db)

	r.HandleFunc("/comments", h.AddComment).Methods("POST")
	r.HandleFunc("/comments/{news_id}", h.GetComments).Methods("GET")

	log.Println("Сервер запущен на :8081")
	err = http.ListenAndServe(":8081", r)
	if err != nil {
		log.Fatal(err)
	}
}
