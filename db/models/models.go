package models

import "time"

type Genero struct {
	ID     int    `json:"id"`
	Nombre string `json:"nombre"`
}

type Serie struct {
	ID          int       `json:"id"`
	Titulo      string    `json:"titulo"`
	Descripcion *string   `json:"descripcion"`
	Episodios   *int      `json:"episodios"`
	Anio        *int      `json:"anio"`
	ImagePath   *string   `json:"image_path"`
	CreadoA     time.Time `json:"creado_a"`
	ActualizadoA time.Time `json:"actualizado_a"`
	Generos     []Genero  `json:"generos"`
}

type SerieInput struct {
	Titulo      string  `json:"titulo"`
	Descripcion *string `json:"descripcion"`
	Episodios   *int    `json:"episodios"`
	Anio        *int    `json:"anio"`
	GeneroIDs   []int   `json:"genero_ids"`
}

type Rating struct {
	ID          int       `json:"id"`
	SerieID     int       `json:"serie_id"`
	Valoracion  int       `json:"valoracion"`
	Review      *string   `json:"review"`
	CreadoA     time.Time `json:"creado_a"`
	ActualizadoA time.Time `json:"actualizado_a"`
}

type RatingInput struct {
	Valoracion int     `json:"valoracion"`
	Review     *string `json:"review"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type PaginatedSeries struct {
	Data       []Serie `json:"data"`
	Total      int     `json:"total"`
	Page       int     `json:"page"`
	Limit      int     `json:"limit"`
	TotalPages int     `json:"total_pages"`
}