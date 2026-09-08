package main

import (
	"log"
	"net/http"
	"time"

	"github.com/magdramazz/olx-api/internal/config"
)

func main() {
	cfg := config.MustLoad()
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte(`{"status":"ok"}`))
		if err != nil {
			return
		}
	})

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("olx api server is listening on port %s", srv.Addr)

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatalf(" server failed:%v", err)
	}
}
