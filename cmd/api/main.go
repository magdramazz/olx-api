package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/magdramazz/olx-api/internal/config"
	database "github.com/magdramazz/olx-api/internal/db"
	"github.com/magdramazz/olx-api/internal/handlers"
)

func main() {
	cfg := config.MustLoad()
	db, err := database.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	fmt.Println("Connected to database")
	fmt.Printf("starting olx sever...")
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", handlers.List(db))

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("olx api server is listening on port %s", srv.Addr)

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatalf(" server failed:%v", err)
	}
}
