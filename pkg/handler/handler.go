package handler

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/google/uuid"
)

var (
	linkStorage = make(map[string]string)
	history     = make([]string, 0)
	mu          sync.Mutex
)

type ShortLinkRequest struct {
	OriginLink string `json:"origin_link"`
}

type ShortLinkResponse struct {
	ShortLink string `json:"short_link"`
}

type NewHandler struct{}

func (h *NewHandler) CreateShortLink(w http.ResponseWriter, r *http.Request) {
	var req ShortLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Неверный запрос полезной нагрузки", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(req.OriginLink, "http://") && !strings.HasPrefix(req.OriginLink, "https://") {
		http.Error(w, "Неверный формат URL", http.StatusBadRequest)
		return
	}

	shortID := uuid.New().String()[:8]

	mu.Lock()
	linkStorage[shortID] = req.OriginLink
	history = append(history, req.OriginLink)
	mu.Unlock()

	resp := ShortLinkResponse{ShortLink: "http://localhost:8080/link/" + shortID}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Не удалось создать короткую ссылку", http.StatusInternalServerError)
	}
}

func (h *NewHandler) GetOriginalLink(w http.ResponseWriter, r *http.Request) {
	shortID := strings.TrimPrefix(r.URL.Path, "/link/")
	mu.Lock()
	originalLink, exists := linkStorage[shortID]
	mu.Unlock()

	if !exists {
		http.Error(w, "Короткая ссылка не найдена", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalLink, http.StatusFound)
}

func (h *NewHandler) GetHistory(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(history); err != nil {
		http.Error(w, "Не удалось получить историю", http.StatusInternalServerError)
	}
}
