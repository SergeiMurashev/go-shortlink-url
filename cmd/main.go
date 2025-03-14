package main

import (
    "go-url-shortener/pkg/handler"
    "log"
    "net/http"
)

func main() {
    h := &handler.NewHandler{}

    http.HandleFunc("/shorten", h.CreateShortLink)
    http.HandleFunc("/link/", h.GetOriginalLink)
    http.HandleFunc("/history", h.GetHistory)
    http.Handle("/", http.FileServer(http.Dir("./static")))

    log.Println("Starting server on port 8080")
    log.Println("Shorten URL: http://localhost:8080")
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal(err)
    }
}

func init() {
    log.Println("Initializing server...")
}