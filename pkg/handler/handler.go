package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	_ "strings"

	"CommentService/pkg/models"
	"CommentService/pkg/storage"
	"github.com/gorilla/mux"
	"time"
)

type Handler struct {
	DB *storage.DB
}

func NewHandler(db *storage.DB) *Handler {
	return &Handler{DB: db}
}

func (h *Handler) AddComment(w http.ResponseWriter, r *http.Request) {
	var req struct {
		NewsID   int    `json:"news_id"`
		ParentID *int   `json:"parent_id"`
		Content  string `json:"content"`
		Author   string `json:"author"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	comment := &models.Comment{
		NewsID:    req.NewsID,
		ParentID:  req.ParentID,
		Content:   req.Content,
		Author:    req.Author,
		CreatedAt: time.Now(),
	}

	err = h.DB.InsertComment(comment)
	if err != nil {
		http.Error(w, "Ошибка сохранения комментария", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Комментарий принят на модерацию"})
}

func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	newsIDStr := vars["news_id"]
	newsID, err := strconv.Atoi(newsIDStr)
	if err != nil {
		http.Error(w, "Некорректный ID новости", http.StatusBadRequest)
		return
	}

	comments, err := h.DB.GetCommentsByNewsID(newsID)
	if err != nil {
		http.Error(w, "Ошибка получения комментариев", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(comments)
}
