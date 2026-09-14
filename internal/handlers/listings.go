package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/magdramazz/olx-api/internal/httpx"
	"github.com/magdramazz/olx-api/internal/middleware"
)

const maxBodyBytes = 1 << 20

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
	rows, err := lh.db.QueryContext(ctx, "SELECT id, title, description, price, city, created_at FROM listings ORDER BY created_at DESC LIMIT 100")
	if err != nil {
		log.Printf("query error: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("rows.close error: %v", err)
		}
	}()
	listings := make([]listing, 0)
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
	if err := json.NewEncoder(w).Encode(listings); err != nil {
		log.Printf("json.encode error: %v", err)
	}
}

func (lh ListingHandler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIdFromContext(ctx)
	var req listing
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxBodyBytes)).Decode(&req); err != nil {
		lh.logger.Error("failed to decode request body", "err", err, "request_id", requestId)
		httpx.Error(w, http.StatusBadRequest, "invalid request body", httpx.CodeInvalidBody)
		return
	}
	price, problem := req.validate()
	if problem != "" {
		httpx.Error(w, http.StatusBadRequest, problem, httpx.CodeValidationError)
		return
	}
	// id and created_at come from the database defaults, never from the client.
	err := lh.db.QueryRowContext(ctx,
		"INSERT INTO listings (title, description, price, city) VALUES ($1, $2, $3, $4) RETURNING id, created_at",
		req.Title, req.Description, price, req.City,
	).Scan(&req.ID, &req.CreatedAt)
	if err != nil {
		lh.logger.Error("failed to insert listing", "err", err, "request_id", requestId)
		httpx.Error(w, http.StatusInternalServerError, "failed to insert listing", httpx.CodeInternalError)
		return
	}
	req.Price = strconv.FormatInt(price, 10)
	lh.logger.Info("listing created", "listing_id", req.ID, "request_id", requestId)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(req); err != nil {
		lh.logger.Error("failed to encode response", "err", err, "request_id", requestId)
	}
}

// validate trims the client-supplied fields and returns the parsed price, or a
// message describing the first invalid field.
func (l *listing) validate() (int64, string) {
	l.Title = strings.TrimSpace(l.Title)
	l.Description = strings.TrimSpace(l.Description)
	l.City = strings.TrimSpace(l.City)
	switch {
	case l.Title == "":
		return 0, "title is required"
	case l.Description == "":
		return 0, "description is required"
	case l.City == "":
		return 0, "city is required"
	}
	price, err := strconv.ParseInt(strings.TrimSpace(l.Price), 10, 64)
	if err != nil || price < 0 {
		return 0, "price must be a non-negative integer"
	}
	return price, ""
}

func (lh ListingHandler) DeleteListing(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := middleware.RequestIdFromContext(ctx)
	id, err := uuid.Parse(r.PathValue("id"))
	if err != nil {
		httpx.Error(w, http.StatusBadRequest, "listing id must be a UUID", httpx.CodeInvalidID)
		return
	}
	res, err := lh.db.ExecContext(ctx, "DELETE FROM listings WHERE id = $1", id.String())
	if err != nil {
		slog.Error("delete failed", "listing_id", id, "err", err, "request_id", requestId)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}
	deleted, err := res.RowsAffected()
	if err != nil {
		slog.Error("delete rows affected failed", "listing_id", id, "err", err, "request_id", requestId)
		httpx.Error(w, http.StatusInternalServerError, "something went wrong", httpx.CodeInternalError)
		return
	}
	if deleted == 0 {
		httpx.Error(w, http.StatusNotFound, "listing not found", httpx.CodeNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
