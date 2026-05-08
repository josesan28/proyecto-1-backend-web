# Proyecto 1 — Backend
### Sistemas y Tecnologías Web
### José Manuel Sanchez Hernández

---
API REST para el Series Tracker. Expone endpoints JSON que consume el cliente.

Repositorio del cliente: https://github.com/josesan28/proyecto-1-cliente-web

Backend en producción: https://proyecto-1-backend-web-production.up.railway.app

---

## Tech Stack

- Lenguaje: Go
- Base de datos: PostgreSQL
- Deploy: Railway
- Dependencias externas: lib/pq, driver de PostgreSQL, godotenv, carga de variables de entorno.

---

## Requisitos para correr localmente

- Go 1.21 o superior
- PostgreSQL 14 o superior

---

## Instrucciones

**1. Clonar el repositorio**

```bash
git clone https://github.com/josesan28/proyecto-1-backend-web
cd proyecto-1-backend-web
```

**2. Configurar variables de entorno**

```bash
cp .env.example .env
```
Se comparte la estructura del .env por motivos académicos.

Editar `.env` con las credenciales de tu instancia local de PostgreSQL:

```
PORT=8080
DB_HOST=localhost
DB_PORT=5432
DB_USER=tu_usuario
DB_PASSWORD=tu_password
DB_NAME=series_tracker
DB_SSLMODE=disable
```

**3. Crear la base de datos**

Correr el siguiente comando asegurándose que se tiene PostgreSQL instalado:

```bash
createdb series_tracker
```

Las tablas se crean automáticamente al arrancar el servidor. No es necesario correr el schema manualmente.

**4. Instalar dependencias y correr**

```bash
go mod download
go run main.go
```

El servidor queda disponible en http://localhost:8080

---

## Endpoints

| Metodo | Ruta | Descripción |
|--------|------|-------------|
| GET | /series | Listar series. Acepta ?page, ?limit, ?q, ?sort, ?order |
| GET | /series/:id | Obtener una serie por ID |
| POST | /series | Crear una serie |
| PUT | /series/:id | Editar una serie |
| DELETE | /series/:id | Eliminar una serie |
| GET | /generos | Listar géneros |
| POST | /generos | Crear un género |
| GET | /series/ratings/:serie_id | Obtener ratings de una serie |
| POST | /series/ratings/:serie_id | Agregar un rating |
| PUT | /series/ratings/:serie_id/:rating_id | Editar un rating |
| DELETE | /series/ratings/:serie_id/:rating_id | Eliminar un rating |

---

## CORS

CORS (Cross-Origin Resource Sharing) es una política de seguridad del navegador que bloquea peticiones HTTP entre orígenes distintos. Como el cliente y el servidor corren en orígenes diferentes, distintos puertos en local, distintos dominios en producción, el navegador rechaza las peticiones a menos que el servidor lo permita explícitamente. Se configuró el header Access-Control-Allow-Origin: * para permitir cualquier origen y evitar problemas con esto.

---

## Challenges implementados en el backend

- Códigos HTTP correctos en toda la API: 201 al crear, 204 al eliminar, 404 si no existe, 400 en input inválido
- Validación server-side con respuestas de error en JSON descriptivas
- Paginación en GET /series con parámetros ?page y ?limit
- Búsqueda por nombre con ?q
- Ordenamiento con ?sort y ?order=asc|desc
- Sistema de ratings con tabla propia en la base de datos y endpoints REST propios
- Soporte para subida de imágenes, guardadas en /uploads, máximo 1 MB

---

## Screenshot

![alt text](/images/image.png)

---

## Reflexión

Para el backend usé Go con la biblioteca net/http, sin ningún framework externo. Al principio me costó un poco el routing manual porque ya había trabajado con frameworks como Express, pero igual me gustó porque cada petición se maneja como uno lo decida. Usar este lenguaje sin frameworks muy grandes me sirvió bastante para entender cómo se manejan los métodos HTTP manualmente y cómo se construyen queries con parámetros en PostgreSQL. En cuanto a este último, me gustó bastante utilizarlo junto a Go porque solo lo había usado con Express, entonces ahora ya entiendo cómo se puede usar con diferentes tecnologías.

La verdad es que volvería a usar Go para backends pequeños porque no tengo mucha experiencia y creo que el desarrollo se puede hacer un poco más lento por eso, entonces para backends más grandes probablemente elegiría un framework que tenga más opciones ya hechas, pero en general me gustó bastante y lo recomendaría para entender los conceptos que mencioné antes.