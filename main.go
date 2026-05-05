package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"github.com/josesan28/proyecto-1-backend-web/db"
	"github.com/josesan28/proyecto-1-backend-web/handlers"
	"github.com/josesan28/proyecto-1-backend-web/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró .env, usando variables del sistema")
	}

	if err := db.Connect(); err != nil {
		log.Fatalf("Error conectando a la DB: %v", err)
	}
	defer db.DB.Close()

	mux := http.NewServeMux()

	mux.Handle("/uploads/", http.StripPrefix("/uploads/", http.FileServer(http.Dir("uploads"))))

	// Géneros
	mux.HandleFunc("/generos", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetGeneros(w, r)
		case http.MethodPost:
			handlers.CreateGenero(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Series
	mux.HandleFunc("/series", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetSeries(w, r)
		case http.MethodPost:
			handlers.CreateSerie(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/series/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetSerieByID(w, r)
		case http.MethodPut:
			handlers.UpdateSerie(w, r)
		case http.MethodDelete:
			handlers.DeleteSerie(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Ratings
	mux.HandleFunc("/series/ratings/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.GetRatings(w, r)
		case http.MethodPost:
			handlers.CreateRating(w, r)
		case http.MethodPut:
			handlers.UpdateRating(w, r)
		case http.MethodDelete:
			handlers.DeleteRating(w, r)
		default:
			http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		}
	})

	// Wrap con middleware de CORS
	handler := middleware.CORS(mux)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Servidor corriendo en http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}