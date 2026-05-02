package model

import "time"

type Serie struct {
    ID          int       `json:"id"`
    Titulo      string    `json:"titulo"`
    Descripcion *string   `json:"descripcion,omitempty"`
    Episodios   *int      `json:"episodios,omitempty"`
    Anio        *int      `json:"anio,omitempty"`
    ImagePath   *string   `json:"image_url,omitempty"`
    Generos     []Genero  `json:"generos,omitempty"`
    CreadoA     time.Time `json:"creado_a"`
    ActualizadoA time.Time `json:"actualizado_a"`
}

type Genero struct {
    ID     int    `json:"id"`
    Nombre string `json:"nombre"`
}

type Rating struct {
    ID          int       `json:"id"`
    SerieID     int       `json:"serie_id"`
    Valoracion  int       `json:"valoracion"`
    Review      *string   `json:"review,omitempty"`
    CreadoA     time.Time `json:"creado_a"`
    ActualizadoA time.Time `json:"actualizado_a"`
}

type CreateSerieRequest struct {
    Titulo      string  `json:"titulo"`
    Descripcion *string `json:"descripcion"`
    Episodios   *int    `json:"episodios"`
    Anio        *int    `json:"anio"`
    GeneroIDs   []int   `json:"genero_ids"`
}

type UpdateSerieRequest struct {
    Titulo      *string `json:"titulo"`
    Descripcion *string `json:"descripcion"`
    Episodios   *int    `json:"episodios"`
    Anio        *int    `json:"anio"`
    GeneroIDs   []int   `json:"genero_ids"`
}

type CreateRatingRequest struct {
    Valoracion int     `json:"valoracion"`
    Review     *string `json:"review"`
}

type PaginatedResponse struct {
    Data       any `json:"data"`
    Page       int `json:"page"`
    Limit      int `json:"limit"`
    Total      int `json:"total"`
    TotalPages int `json:"total_pages"`
}

type ErrorResponse struct {
    Error   string `json:"error"`
    Details any    `json:"details,omitempty"`
}