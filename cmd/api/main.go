package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/magdramazz/olx-api/internal/config"
	database "github.com/magdramazz/olx-api/internal/db"
	"github.com/magdramazz/olx-api/internal/handlers"
	"github.com/magdramazz/olx-api/internal/middleware"
)

func main() {
	cfg := config.MustLoad()
	db, err := database.Connect(cfg.DatabaseUrl)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	}))
	slog.SetDefault(logger)
	fmt.Println("Connected to database")
	fmt.Printf("starting olx sever...")
	lh := handlers.NewListHandler(db, logger)
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handlers.Health)
	mux.HandleFunc("GET /listings", lh.List)
	mux.HandleFunc("DELETE /listings/{id}", lh.DeleteListing)

	handler := middleware.RequestId(mux)
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      handler,
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
