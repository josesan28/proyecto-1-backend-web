package server

import (
    "fmt"
    "log"
    "net/http"
    "time"

    "github.com/jackc/pgx/v5/pgxpool"
    "github.com/josesan28/proyecto-1-backend-web/internal/config"
    "github.com/josesan28/proyecto-1-backend-web/internal/handler"
)

type Server struct {
    cfg  *config.Config
    pool *pgxpool.Pool
    mux  *http.ServeMux
}

func New(cfg *config.Config, pool *pgxpool.Pool) *Server {
    s := &Server{
        cfg:  cfg,
        pool: pool,
        mux:  http.NewServeMux(),
    }
    s.registerRoutes()
    return s
}

func (s *Server) registerRoutes() {
    healthHandler := handler.NewHealthHandler(s.pool)

    s.mux.HandleFunc("GET /health", healthHandler.Check)

    s.mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

}

func (s *Server) Handler() http.Handler {
    return corsMiddleware(s.mux)
}

func (s *Server) Run() error {
    srv := &http.Server{
        Addr:         fmt.Sprintf(":%s", s.cfg.Port),
        Handler:      s.Handler(),
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 30 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    log.Printf("🚀 Servidor corriendo en http://localhost:%s", s.cfg.Port)
    return srv.ListenAndServe()
}

func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "*")
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusNoContent)
            return
        }

        next.ServeHTTP(w, r)
    })
}