CREATE TABLE usuario (
    id            SERIAL PRIMARY KEY,
    username      VARCHAR(50) NOT NULL UNIQUE,
    correo        VARCHAR(255) NOT NULL UNIQUE,
    hash_password VARCHAR(255) NOT NULL,
    creado_a      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    actualizado_a TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE genero (
    id   SERIAL PRIMARY KEY,
    nombre VARCHAR(100) NOT NULL UNIQUE
);

CREATE TABLE serie (
    id             SERIAL PRIMARY KEY,
    titulo         VARCHAR(255) NOT NULL,
    descripcion    TEXT,
    episodios      INTEGER     CHECK (episodios >= 0),
    anio           INTEGER      CHECK (anio >= 1900 AND anio <= 2100),
    image_path     VARCHAR(500),  
    creado_a       TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    actualizado_a  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);

CREATE TABLE serie_genero (
    serie_id INTEGER NOT NULL REFERENCES serie(id) ON DELETE CASCADE,
    genero_id  INTEGER NOT NULL REFERENCES genero(id) ON DELETE CASCADE,
    PRIMARY KEY (serie_id, genero_id)
);

CREATE TABLE rating (
    id            SERIAL PRIMARY KEY,
    serie_id      INTEGER      NOT NULL REFERENCES serie(id) ON DELETE CASCADE,
    valoracion    SMALLINT     NOT NULL CHECK (valoracion >= 1 AND valoracion <= 10),
    review        TEXT,
    creado_a      TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    actualizado_a TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);