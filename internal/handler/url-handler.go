package handler

import (
	"design-url-shortner/internal/service"
	"net/http"

	"github.com/go-chi/chi"
)

type URLHandler struct {
	urlService *service.URLService
}

func NewURLHandler(urlService *service.URLService) *URLHandler {
	return &URLHandler{urlService: urlService}
}

func (h *URLHandler) CreateURL(w http.ResponseWriter, r *http.Request) {

}

func (h *URLHandler) GetURLStats(w http.ResponseWriter, r *http.Request) {
	shortCode := chi.URLParam(r, "shortCode")
	if shortCode == "" {
		http.Error(w, "Short Code is required", http.StatusBadRequest)
		return
	}

	// TODO
}
