package handler

import (
    "context"
    "encoding/json"
    "net/http"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
)

type HealthHandler struct {
    pool *pgxpool.Pool
}

func NewHealthHandler(pool *pgxpool.Pool) *HealthHandler {
    return &HealthHandler{pool: pool}
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
    defer cancel()

    dbStatus := "ok"
    if err := h.pool.Ping(ctx); err != nil {
        dbStatus = "error: " + err.Error()
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]any{
        "status":    "ok",
        "db":        dbStatus,
        "timestamp": time.Now().UTC(),
    })
}