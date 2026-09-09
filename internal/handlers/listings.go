package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

func List(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT * FROM listings ORDER BY created_at DESC LIMIT 100")
		if err != nil {
			log.Printf("query error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		defer func(rows *sql.Rows) {
			err := rows.Close()
			if err != nil {
				log.Printf("rows.close error: %v", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
		}(rows)
		var listings []listing
		for rows.Next() {
			var list listing
			if err := rows.Scan(&list.ID, &list.Title, &list.Description, &list.Price, &list.City, &list.CreatedAt); err != nil {
				log.Printf("rows.scan error: %v", err)
				http.Error(w, "internal error", http.StatusInternalServerError)
				return
			}
			listings = append(listings, list)
		}
		if err := rows.Err(); err != nil {
			log.Printf("rows.error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(listings)
		if err != nil {
			log.Printf("json.encode error: %v", err)
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}
}
