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
	mu          sync.Mutex
)

type ShortLinkRequest struct {
	OriginalLink string `json:"original_link"`
}

type ShortLinkResponse struct {
	ShortLink string `json:"short_link"`
}

func CreateShortLink(w http.ResponseWriter, r *http.Request) {
	var req ShortLinkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if !strings.HasPrefix(req.OriginalLink, "http://") && !strings.HasPrefix(req.OriginalLink, "https://") {
		http.Error(w, "Invalid URL format", http.StatusBadRequest)
		return
	}

	shortID := uuid.New().String()[:8]

	mu.Lock()
	linkStorage[shortID] = req.OriginalLink
	mu.Unlock()

	resp := ShortLinkResponse{ShortLink: "http://localhost:8080/link/" + shortID}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to create short link", http.StatusInternalServerError)
	}
}

func GetOriginalLink(w http.ResponseWriter, r *http.Request) {
	shortID := strings.TrimPrefix(r.URL.Path, "/link/")
	mu.Lock()
	originalLink, exists := linkStorage[shortID]
	mu.Unlock()

	if !exists {
		http.Error(w, "Short link not found", http.StatusNotFound)
		return
	}

	http.Redirect(w, r, originalLink, http.StatusFound)
}
