package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"time"

	"github.com/magdramazz/olx-api/internal/httpx"
	"github.com/magdramazz/olx-api/internal/middleware"
)

type listing struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       string    `json:"price"`
	City        string    `json:"city"`
	CreatedAt   time.Time `json:"created_at"`
}

type ListingHandler struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewListHandler(db *sql.DB, logger *slog.Logger) *ListingHandler {
	return &ListingHandler{
		db:     db,
		logger: logger,
	}
}
func (lh ListingHandler) List(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rows, err := lh.db.QueryContext(ctx, "SELECT * FROM listings ORDER BY created_at DESC LIMIT 100")
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

func (lh ListingHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIdFromContext(ctx)
	var req listing
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		lh.logger.Error("failed to decode request body", "err", err, "request_id", requestId)
		httpx.Error(w, http.StatusBadRequest, "invalid request body", httpx.CodeInvalidID)
		return
	}
	row := lh.db.QueryRowContext(ctx, "INSERT INTO listings (id, title, description, price, city) VALUES ($1, $2, $3, $4, $5)", req.ID, req.Title, req.Description, req.Price, req.City)
	if err := row.Scan(); err != nil {
		lh.logger.Error("failed to insert listing", "err", err, "request_id", requestId)
		httpx.Error(w, http.StatusInternalServerError, "failed to insert listing", httpx.CodeInternalError)
		return
	}
	lh.logger.Info("listing created", "request_id", requestId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": req.ID})
	w.Write([]byte("{ok}"))
}

func (lh ListingHandler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	ctx := r.Context()
	requestId := middleware.RequestIdFromContext(ctx)
	_, err := lh.db.ExecContext(ctx, "DELETE FROM listings WHERE id = $1", id)
	if err != nil {
		slog.Error("delete failed", "listing_id", id, "err", err, "request_id", requestId)
		//http.Error(w, "internal error", http.StatusInternalServerError)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
