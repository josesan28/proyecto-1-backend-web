package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/josesan28/proyecto-1-backend-web/db"
	"github.com/josesan28/proyecto-1-backend-web/models"
)

// GET /series 
func GetSeries(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT id, titulo, descripcion, episodios, anio, image_path, creado_a, actualizado_a
		FROM serie
		ORDER BY creado_a DESC
	`)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error obteniendo series")
		return
	}
	defer rows.Close()

	series := []models.Serie{}
	for rows.Next() {
		var s models.Serie
		if err := rows.Scan(
			&s.ID, &s.Titulo, &s.Descripcion,
			&s.Episodios, &s.Anio, &s.ImagePath,
			&s.CreadoA, &s.ActualizadoA,
		); err != nil {
			writeError(w, http.StatusInternalServerError, "Error leyendo series")
			return
		}
		s.Generos = fetchGenerosDeSerie(s.ID)
		series = append(series, s)
	}

	writeJSON(w, http.StatusOK, series)
}

// GET /series/{id}
func GetSerieByID(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path, 2)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var s models.Serie
	err = db.DB.QueryRow(`
		SELECT id, titulo, descripcion, episodios, anio, image_path, creado_a, actualizado_a
		FROM serie WHERE id = $1
	`, id).Scan(
		&s.ID, &s.Titulo, &s.Descripcion,
		&s.Episodios, &s.Anio, &s.ImagePath,
		&s.CreadoA, &s.ActualizadoA,
	)

	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "Serie no encontrada")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error obteniendo serie")
		return
	}

	s.Generos = fetchGenerosDeSerie(s.ID)
	writeJSON(w, http.StatusOK, s)
}

// POST /series
func CreateSerie(w http.ResponseWriter, r *http.Request) {
	// Parsear form (imagen + campos de texto)
	if err := r.ParseMultipartForm(5 << 20); err != nil {
		// Intentar JSON si no es multipart
		if err := r.ParseForm(); err != nil {
			writeError(w, http.StatusBadRequest, "Error parseando el body")
			return
		}
	}

	titulo := strings.TrimSpace(r.FormValue("titulo"))
	if titulo == "" {
		writeError(w, http.StatusBadRequest, "El título es requerido")
		return
	}

	descripcion := nullableString(r.FormValue("descripcion"))
	episodios := nullableInt(r.FormValue("episodios"))
	anio := nullableInt(r.FormValue("anio"))

	// Validar episodios si viene
	if episodios != nil && *episodios < 0 {
		writeError(w, http.StatusBadRequest, "Los episodios no pueden ser negativos")
		return
	}

	// Manejar imagen si viene
	imagePath, err := saveUploadedImage(r, "imagen")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error guardando la imagen")
		return
	}

	// Insertar serie
	var s models.Serie
	err = db.DB.QueryRow(`
		INSERT INTO serie (titulo, descripcion, episodios, anio, image_path)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, titulo, descripcion, episodios, anio, image_path, creado_a, actualizado_a
	`, titulo, descripcion, episodios, anio, imagePath).Scan(
		&s.ID, &s.Titulo, &s.Descripcion,
		&s.Episodios, &s.Anio, &s.ImagePath,
		&s.CreadoA, &s.ActualizadoA,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error creando la serie")
		return
	}

	// Asociar géneros si vienen
	generoIDs := parseIntList(r.FormValue("genero_ids"))
	associateGeneros(s.ID, generoIDs)
	s.Generos = fetchGenerosDeSerie(s.ID)

	writeJSON(w, http.StatusCreated, s)
}

// PUT /series/{id}
func UpdateSerie(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path, 2)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	// Verificar que existe
	var exists bool
	db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM serie WHERE id=$1)", id).Scan(&exists)
	if !exists {
		writeError(w, http.StatusNotFound, "Serie no encontrada")
		return
	}

	if err := r.ParseMultipartForm(5 << 20); err != nil {
		r.ParseForm()
	}

	titulo := strings.TrimSpace(r.FormValue("titulo"))
	if titulo == "" {
		writeError(w, http.StatusBadRequest, "El título es requerido")
		return
	}

	descripcion := nullableString(r.FormValue("descripcion"))
	episodios := nullableInt(r.FormValue("episodios"))
	anio := nullableInt(r.FormValue("anio"))

	if anio != nil && (*anio < 1900 || *anio > 2100) {
		writeError(w, http.StatusBadRequest, "El año debe estar entre 1900 y 2100")
		return
	}
	if episodios != nil && *episodios < 0 {
		writeError(w, http.StatusBadRequest, "Los episodios no pueden ser negativos")
		return
	}

	// Si viene imagen nueva, guardarla, si no, mantener la existente
	imagePath, err := saveUploadedImage(r, "imagen")
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error guardando la imagen")
		return
	}

	var query string
	var queryArgs []any

	if imagePath != nil {
		query = `
			UPDATE serie SET titulo=$1, descripcion=$2, episodios=$3, anio=$4, image_path=$5, actualizado_a=NOW()
			WHERE id=$6
			RETURNING id, titulo, descripcion, episodios, anio, image_path, creado_a, actualizado_a`
		queryArgs = []any{titulo, descripcion, episodios, anio, imagePath, id}
	} else {
		query = `
			UPDATE serie SET titulo=$1, descripcion=$2, episodios=$3, anio=$4, actualizado_a=NOW()
			WHERE id=$5
			RETURNING id, titulo, descripcion, episodios, anio, image_path, creado_a, actualizado_a`
		queryArgs = []any{titulo, descripcion, episodios, anio, id}
	}

	var s models.Serie
	err = db.DB.QueryRow(query, queryArgs...).Scan(
		&s.ID, &s.Titulo, &s.Descripcion,
		&s.Episodios, &s.Anio, &s.ImagePath,
		&s.CreadoA, &s.ActualizadoA,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error actualizando la serie")
		return
	}

	// Reasociar géneros
	generoIDs := parseIntList(r.FormValue("genero_ids"))
	db.DB.Exec("DELETE FROM serie_genero WHERE serie_id=$1", s.ID)
	associateGeneros(s.ID, generoIDs)
	s.Generos = fetchGenerosDeSerie(s.ID)

	writeJSON(w, http.StatusOK, s)
}

// DELETE /series/{id}
func DeleteSerie(w http.ResponseWriter, r *http.Request) {
	id, err := extractID(r.URL.Path, 2)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID inválido")
		return
	}

	result, err := db.DB.Exec("DELETE FROM serie WHERE id=$1", id)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error eliminando la serie")
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeError(w, http.StatusNotFound, "Serie no encontrada")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Helpers privados

func fetchGenerosDeSerie(serieID int) []models.Genero {
	rows, err := db.DB.Query(`
		SELECT g.id, g.nombre FROM genero g
		JOIN serie_genero sg ON sg.genero_id = g.id
		WHERE sg.serie_id = $1
	`, serieID)
	if err != nil {
		return []models.Genero{}
	}
	defer rows.Close()

	generos := []models.Genero{}
	for rows.Next() {
		var g models.Genero
		rows.Scan(&g.ID, &g.Nombre)
		generos = append(generos, g)
	}
	return generos
}

func associateGeneros(serieID int, generoIDs []int) {
	for _, gid := range generoIDs {
		db.DB.Exec(
			"INSERT INTO serie_genero (serie_id, genero_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
			serieID, gid,
		)
	}
}

func nullableString(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}

func nullableInt(s string) *int {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	v, err := strconv.Atoi(s)
	if err != nil {
		return nil
	}
	return &v
}

func parseIntList(s string) []int {
	result := []int{}
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if v, err := strconv.Atoi(part); err == nil {
			result = append(result, v)
		}
	}
	return result
}

func decodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}