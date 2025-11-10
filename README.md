# Books API en Go

Este repositorio implementa una API REST sencilla para manejar libros (CRUD) escrita en Go.

Características principales:
- Servidor HTTP con endpoints para crear, listar, leer, actualizar y eliminar libros.
- Almacenamiento usando SQLite (archivo local `books.db`).

## Estructura del proyecto

La estructura principal del proyecto (resumen):

- `main.go`           : Inicializa la base de datos y arranca el servidor HTTP.
- `internal/model`    : Definición de modelos (por ejemplo `book.go`).
- `internal/service`  : Lógica de negocio.
- `internal/transport`: Handlers HTTP.
- `store`             : Acceso a la base de datos (`book_store.go`).
- `test`              : Scripts o utilidades de prueba (PowerShell en `test.ps1`).

Revisa el árbol completo en el repositorio para más detalles.

## Requisitos

- Go (recomendado >= 1.20).
- El módulo `modernc.org/sqlite` se usa como driver SQLite (ya referenciado en `main.go`).

## Cómo ejecutar (PowerShell)

1. Abrir PowerShell en la raíz del proyecto:

```powershell
# opcional: descargar dependencias y tidy
go mod tidy

# compilar
go build -o books-api

# ejecutar
./books-api
```

2. Al iniciar, `main.go` crea (si no existe) una base de datos SQLite local llamada `books.db` y arranca el servidor.

Por defecto el servidor escucha en el puerto `8000` (ver `main.go`). Los endpoints expuestos son:

- GET    /books          -> Obtener todos los libros
- POST   /books          -> Crear un libro (JSON body)
- GET    /books/{id}     -> Obtener un libro por ID
- PUT    /books/{id}     -> Actualizar un libro por ID (JSON body)
- DELETE /books/{id}     -> Eliminar un libro por ID

Ejemplo rápido con PowerShell (crear y listar):

```powershell
# Crear un libro
Invoke-RestMethod -Method Post -Uri http://localhost:8000/books -Body (@{title='El Quijote'; author='Cervantes'} | ConvertTo-Json) -ContentType 'application/json'

# Listar libros
Invoke-RestMethod -Method Get -Uri http://localhost:8000/books
```

O con curl:

```powershell
curl -X POST http://localhost:8000/books -H "Content-Type: application/json" -d '{"title":"El Quijote","author":"Cervantes"}'
curl http://localhost:8000/books
```

## Notas

- El archivo de base de datos se crea en el directorio desde donde se ejecuta la aplicación (por defecto `./books.db`).
- Si cambias el puerto o la ubicación de la DB, actualiza `main.go`.
- El proyecto contiene una estructura simple (store/service/transport) para separar responsabilidades.

## Siguientes pasos recomendados

- Añadir pruebas unitarias/integ. para `store` y `service`.
- Documentar el contrato JSON (ej. ejemplos de request/response).
- Añadir un Makefile o scripts para facilitar ejecución y pruebas.

---

Si quieres, actualizo también el README en inglés o añado ejemplos de requests más completos (con curl/HTTPie/Postman).