package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/josesan28/proyecto-1-backend-web/db"
	"github.com/josesan28/proyecto-1-backend-web/models"
)

// GET /series/{id}/ratings
func GetRatings(w http.ResponseWriter, r *http.Request) {
	serieID, err := extractID(r.URL.Path, 3)
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID de serie inválido")
		return
	}

	var exists bool
	db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM serie WHERE id=$1)", serieID).Scan(&exists)
	if !exists {
		writeError(w, http.StatusNotFound, "Serie no encontrada")
		return
	}

	rows, err := db.DB.Query(`
		SELECT id, serie_id, valoracion, review, creado_a, actualizado_a
		FROM rating WHERE serie_id = $1
		ORDER BY creado_a DESC
	`, serieID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error obteniendo ratings")
		return
	}
	defer rows.Close()

	ratings := []models.Rating{}
	for rows.Next() {
		var rt models.Rating
		rows.Scan(&rt.ID, &rt.SerieID, &rt.Valoracion, &rt.Review, &rt.CreadoA, &rt.ActualizadoA)
		ratings = append(ratings, rt)
	}

	writeJSON(w, http.StatusOK, ratings)
}

// POST /series/{id_serie}/ratings
func CreateRating(w http.ResponseWriter, r *http.Request) {
	serieID, err := extractID(r.URL.Path, 3) // ← cambiado de 2 a 3
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID de serie inválido")
		return
	}

	var exists bool
	db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM serie WHERE id=$1)", serieID).Scan(&exists)
	if !exists {
		writeError(w, http.StatusNotFound, "Serie no encontrada")
		return
	}

	var input models.RatingInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Body inválido")
		return
	}

	if input.Valoracion < 1 || input.Valoracion > 10 {
		writeError(w, http.StatusBadRequest, "La valoración debe estar entre 1 y 10")
		return
	}

	var rt models.Rating
	err = db.DB.QueryRow(`
		INSERT INTO rating (serie_id, valoracion, review)
		VALUES ($1, $2, $3)
		RETURNING id, serie_id, valoracion, review, creado_a, actualizado_a
	`, serieID, input.Valoracion, input.Review).Scan(
		&rt.ID, &rt.SerieID, &rt.Valoracion, &rt.Review, &rt.CreadoA, &rt.ActualizadoA,
	)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error creando rating")
		return
	}

	writeJSON(w, http.StatusCreated, rt)
}

// PUT /series/{id_serie}/ratings/{ratingId}
func UpdateRating(w http.ResponseWriter, r *http.Request) {
	serieID, err := extractID(r.URL.Path, 3) // ← cambiado de 2 a 3
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID de serie inválido")
		return
	}

	ratingID, err := extractID(r.URL.Path, 4) // este se mantiene en 4
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID de rating inválido")
		return
	}

	var input models.RatingInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, "Body inválido")
		return
	}

	if input.Valoracion < 1 || input.Valoracion > 10 {
		writeError(w, http.StatusBadRequest, "La valoración debe estar entre 1 y 10")
		return
	}

	var rt models.Rating
	err = db.DB.QueryRow(`
		UPDATE rating SET valoracion=$1, review=$2, actualizado_a=NOW()
		WHERE id=$3 AND serie_id=$4
		RETURNING id, serie_id, valoracion, review, creado_a, actualizado_a
	`, input.Valoracion, input.Review, ratingID, serieID).Scan(
		&rt.ID, &rt.SerieID, &rt.Valoracion, &rt.Review, &rt.CreadoA, &rt.ActualizadoA,
	)

	if err == sql.ErrNoRows {
		writeError(w, http.StatusNotFound, "Rating no encontrado")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error actualizando rating")
		return
	}

	writeJSON(w, http.StatusOK, rt)
}

// DELETE /series/{id_serie}/ratings/{ratingId}
func DeleteRating(w http.ResponseWriter, r *http.Request) {
	serieID, err := extractID(r.URL.Path, 3) // ← cambiado de 2 a 3
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID de serie inválido")
		return
	}

	ratingID, err := extractID(r.URL.Path, 4) // este se mantiene en 4
	if err != nil {
		writeError(w, http.StatusBadRequest, "ID de rating inválido")
		return
	}

	result, err := db.DB.Exec("DELETE FROM rating WHERE id=$1 AND serie_id=$2", ratingID, serieID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Error eliminando rating")
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		writeError(w, http.StatusNotFound, "Rating no encontrado")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}