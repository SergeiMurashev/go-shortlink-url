package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateShortLink(t *testing.T) {
	reqBody, _ := json.Marshal(ShortLinkRequest{OriginalLink: "http://example.com"})
	req, err := http.NewRequest("POST", "/shorten", bytes.NewBuffer(reqBody))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(CreateShortLink)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var resp ShortLinkResponse
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		t.Fatal(err)
	}

	if !strings.HasPrefix(resp.ShortLink, "http://localhost:8080/link/") {
		t.Errorf("unexpected short link: %v", resp.ShortLink)
	}
}

func TestGetOriginalLink(t *testing.T) {
	// Create a short link first
	shortID := "test1234"
	linkStorage[shortID] = "http://example.com"

	req, err := http.NewRequest("GET", "/link/"+shortID, nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(GetOriginalLink)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusFound)
	}

	location := rr.Header().Get("Location")
	if location != "http://example.com" {
		t.Errorf("unexpected redirect location: %v", location)
	}
}
