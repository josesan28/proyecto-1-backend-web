package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/josesan28/proyecto-1-backend-web/db"
	"github.com/josesan28/proyecto-1-backend-web/models"
)

func GetGeneros(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, nombre FROM genero ORDER BY nombre")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error obteniendo géneros")
		return
	}
	defer rows.Close()

	generos := []models.Genero{}
	for rows.Next() {
		var g models.Genero
		if err := rows.Scan(&g.ID, &g.Nombre); err != nil {
			writeError(w, http.StatusInternalServerError, "Error leyendo géneros")
			return
		}
		generos = append(generos, g)
	}

	writeJSON(w, http.StatusOK, generos)
}

func CreateGenero(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Nombre string `json:"nombre"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Body inválido")
		return
	}

	input.Nombre = strings.TrimSpace(input.Nombre)
	if input.Nombre == "" {
		writeError(w, http.StatusBadRequest, "El nombre del género es requerido")
		return
	}

	var g models.Genero
	err := db.DB.QueryRow(
		"INSERT INTO genero (nombre) VALUES ($1) RETURNING id, nombre",
		input.Nombre,
	).Scan(&g.ID, &g.Nombre)

	if err != nil {
		writeError(w, http.StatusConflict, "El género ya existe o hubo un error al crearlo")
		return
	}

	writeJSON(w, http.StatusCreated, g)
}