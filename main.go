package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/josesan28/proyecto-1-backend-web/internal/config"
	"github.com/josesan28/proyecto-1-backend-web/internal/db"
	"github.com/josesan28/proyecto-1-backend-web/internal/server"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Aviso: no se cargó .env: %v", err)
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Config inválida: %v", err)
	}

	ctx := context.Background()
	pool, err := db.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("No se pudo conectar a la base de datos: %v", err)
	}
	defer pool.Close()

	log.Println("Conexión a PostgreSQL exitosa")

	srv := server.New(cfg, pool)
	if err := srv.Run(); err != nil {
		log.Fatalf("Error del servidor: %v", err)
	}
}
