package main

import (
	"github.com/oleg6k/url-shortener/internal/app"
	"net/http"
)

func main() {
	service := app.NewService()
	controller := app.NewController("http://localhost:8080", service)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			controller.PostShorting(w, r)
		} else if r.Method == http.MethodGet {
			controller.GetRedirectToOriginal(w, r)
		} else {
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	})
	http.ListenAndServe(":8080", mux)
}
