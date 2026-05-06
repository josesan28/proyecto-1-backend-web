package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() error {
	connStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("error abriendo conexión: %w", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("error haciendo ping a la DB: %w", err)
	}

	if err = runSchema(); err != nil {
		return fmt.Errorf("error ejecutando schema: %w", err)
	}

	return nil
}

func runSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS usuario (
		id            SERIAL PRIMARY KEY,
		username      VARCHAR(50) NOT NULL UNIQUE,
		correo        VARCHAR(255) NOT NULL UNIQUE,
		hash_password VARCHAR(255) NOT NULL,
		creado_a      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
		actualizado_a TIMESTAMPTZ NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS genero (
		id     SERIAL PRIMARY KEY,
		nombre VARCHAR(100) NOT NULL UNIQUE
	);

	CREATE TABLE IF NOT EXISTS serie (
		id            SERIAL PRIMARY KEY,
		titulo        VARCHAR(255) NOT NULL,
		descripcion   TEXT,
		episodios     INTEGER      CHECK (episodios >= 0),
		anio          INTEGER      CHECK (anio >= 1900 AND anio <= 2100),
		image_path    VARCHAR(500),
		creado_a      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
		actualizado_a TIMESTAMPTZ  NOT NULL DEFAULT NOW()
	);

	CREATE TABLE IF NOT EXISTS serie_genero (
		serie_id  INTEGER NOT NULL REFERENCES serie(id) ON DELETE CASCADE,
		genero_id INTEGER NOT NULL REFERENCES genero(id) ON DELETE CASCADE,
		PRIMARY KEY (serie_id, genero_id)
	);

	CREATE TABLE IF NOT EXISTS rating (
		id            SERIAL PRIMARY KEY,
		serie_id      INTEGER      NOT NULL REFERENCES serie(id) ON DELETE CASCADE,
		valoracion    SMALLINT     NOT NULL CHECK (valoracion >= 1 AND valoracion <= 10),
		review        TEXT,
		creado_a      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
		actualizado_a TIMESTAMPTZ  NOT NULL DEFAULT NOW()
	);
	`

	_, err := DB.Exec(schema)
	return err
}