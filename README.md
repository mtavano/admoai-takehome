# AdMoai Take Home Test - API de Anuncios

## 📋 Descripción

API REST para gestión de anuncios con sistema de TTL (Time To Live) automático. Permite crear, consultar y gestionar anuncios con expiración automática basada en tiempo.

## 🚀 Características

- **CRUD de Anuncios**: Crear, consultar y actualizar anuncios
- **Sistema TTL**: Expiración automática de anuncios basada en minutos
- **Filtros Avanzados**: Consulta por placement, status y otros criterios
- **Base de Datos SQLite**: Almacenamiento local con migraciones automáticas
- **Validaciones**: Validación de entrada con Gin y golang validator
- **Query Builder**: Uso de Squirrel para queries dinámicas y seguras

## 🛠️ Tecnologías

- **Go 1.23.4**
- **Gin** - Framework web
- **SQLite** - Base de datos
- **Squirrel** - Query builder
- **Goose** - Migraciones de base de datos
- **Golang Validator** - Validación de datos

## 📦 Instalación

### Prerrequisitos
- Go 1.23.4 o superior
- SQLite3

### Configuración
1. Clonar el repositorio:
```bash
git clone <repository-url>
cd admoai-takehome
```

2. Instalar dependencias:
```bash
go mod tidy
```

3. Configurar variables de entorno:
```bash
cp example.dev.env dev.env
```

4. Ejecutar migraciones:
```bash
make setup
```

5. Iniciar el servidor:
```bash
make run-simple
```

## 🔧 Configuración

### Variables de Entorno (`dev.env`)
```env
ENVIRONMENT=dev
API_PORT=9001
DB_DRIVER=sqlite3
DB_DSN=./data/admoai.db
```

## 📚 API Endpoints

### Base URL
```
http://localhost:9001/v1
```

### 1. Crear Anuncio
**POST** `/ads`

Crea un nuevo anuncio con TTL opcional.

**Request Body:**

```json
{
  "title": "Anuncio de Prueba",
  "image_url": "https://example.com/image.jpg",
  "placement": "homepage",
  "ttl": 30
}
```

**Campos:**
- `title` (required): Título del anuncio
- `image_url` (required): URL de la imagen (debe ser URL válida)
- `placement` (required): Ubicación del anuncio
- `ttl` (optional): Tiempo de vida en minutos (0 = sin expiración)

**Response (201):**

```json
{
  "id": "uuid-generated",
  "title": "Anuncio de Prueba",
  "imageUrl": "https://example.com/image.jpg",
  "placement": "homepage",
  "status": "active",
  "createdAt": 1640995200,
  "expiresAt": 1640997000
}
```

### 2. Obtener Anuncio por ID
**GET** `/ads/{id}`

Obtiene un anuncio específico por su ID.

**Response (200):**

```json
{
  "id": "uuid-here",
  "title": "Anuncio de Prueba",
  "imageUrl": "https://example.com/image.jpg",
  "placement": "homepage",
  "status": "active",
  "createdAt": 1640995200,
  "expiresAt": 1640997000
}
```

**Response (404):**

```json
{
  "error": "Ad not found"
}
```

### 3. Filtrar Anuncios
**GET** `/ads?placement=homepage&status=active`

Obtiene anuncios filtrados por criterios específicos.

**Query Parameters:**
- `placement` (optional): Filtrar por ubicación
- `status` (optional): Filtrar por estado

**Response (200):**

```json
{
  "ads": [
    {
      "id": "uuid-1",
      "title": "Anuncio 1",
      "imageUrl": "https://example.com/image1.jpg",
      "placement": "homepage",
      "status": "active",
      "createdAt": 1640995200,
      "expiresAt": 1640997000
    }
  ],
  "count": 1
}
```

### 4. Desactivar Anuncio
**POST** `/ads/{id}/deactivate`

Desactiva un anuncio específico.

**Response (200):**

```json
{
  "message": "Ad deactivated successfully",
  "id": "uuid-here",
  "status": "inactive"
}
```

### 5. Health Check
**GET** `/health`

Verifica el estado del servicio.

**Response (200):**

```json
{
  "status": "running"
}
```

## 🗄️ Modelo de Datos

### AdvertiseRecord

```go
type AdvertiseRecord struct {
    ID        string  `db:"id" json:"id"`
    Title     string  `db:"title" json:"title"`
    ImageURL  string  `db:"image_url" json:"imageUrl"`
    Placement string  `db:"placement" json:"placement"`
    Status    string  `db:"status" json:"status"`
    CreatedAt int64   `db:"created_at" json:"createdAt"`
    ExpiresAt *int64  `db:"expires_at" json:"expiresAt"`
}
```

### Estados del Anuncio

- `active`: Anuncio activo y visible
- `inactive`: Anuncio desactivado

## ⏰ Sistema TTL

### Comportamiento

- Los anuncios pueden tener un TTL (Time To Live) en minutos
- Si `ttl = 0` o no se especifica, el anuncio no expira
- Los anuncios expirados se filtran automáticamente de las consultas
- El campo `expiresAt` se calcula como `createdAt + (ttl * 60 segundos)`

### Filtrado Automático

Todas las consultas automáticamente descartan anuncios expirados:
- Anuncios sin TTL (`expiresAt = null`) se muestran siempre
- Anuncios con TTL válido (`expiresAt > currentTime`) se muestran
- Anuncios expirados (`expiresAt <= currentTime`) se descartan

## 🛠️ Comandos Make

```bash
# Configuración inicial
make setup

# Ejecutar servidor con hot reload
make run

# Ejecutar servidor simple
make run-simple

# Ejecutar migraciones
make migrate

# Revertir migraciones
make migrate-down

# Limpiar base de datos
make clean

# Crear nueva migración
make create-migration name=migration_name
```

## 🧪 Functional Testing

### Ejemplos de uso con curl

**Crear anuncio:**

```bash
curl -X POST http://localhost:9001/v1/ads \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Anuncio de Prueba",
    "image_url": "https://example.com/image.jpg",
    "placement": "homepage",
    "ttl": 30
  }'
```

**Obtener anuncio por ID:**

```bash
curl -X GET http://localhost:9001/v1/ads/your-uuid-here
```

**Filtrar anuncios:**

```bash
curl -X GET "http://localhost:9001/v1/ads?placement=homepage&status=active"
```

**Desactivar anuncio:**

```bash
curl -X POST http://localhost:9001/v1/ads/your-uuid-here/deactivate
```

## 📁 Estructura del Proyecto

```
admoai-takehome/
├── cmd/
│   └── server/
│       └── main.go              # Punto de entrada
├── internal/
│   ├── api/                     # Handlers de la API
│   │   ├── handle_func.go
│   │   ├── post_ads_handler.go
│   │   ├── get_ads_by_id.go
│   │   ├── get_ads_by_filters_handler.go
│   │   ├── post_deactivate_ads_handler.go
│   │   └── register_routes.go
│   └── store/                   # Capa de datos
│       ├── database.go
│       ├── implementation.go
│       ├── models.go
│       └── query/
│           ├── insert_ads.go
│           ├── select_ads.go
│           └── update_ads.go
├── migrations/                  # Migraciones de BD
│   └── 20250622182727_init_setup.go
├── dev.env                      # Variables de entorno
├── go.mod                       # Dependencias
├── Makefile                     # Comandos útiles
└── README.md                    # Este archivo
```

## 🔒 Validaciones

### Validaciones de Entrada

- **title**: No puede estar vacío
- **image_url**: No puede estar vacío y debe ser URL válida
- **placement**: No puede estar vacío
- **ttl**: Opcional, debe ser mayor a 0 si se especifica

### Validaciones de Negocio

- Al menos un filtro debe estar presente en consultas por filtros
- Los anuncios expirados se descartan automáticamente
- Los anuncios desactivados mantienen su estado

## 🚨 Códigos de Error

- **400 Bad Request**: Datos de entrada inválidos
- **404 Not Found**: Recurso no encontrado
- **500 Internal Server Error**: Error interno del servidor

## 📝 Notas de Desarrollo

- La API usa SQLite para simplicidad de desarrollo
- Todas las queries usan Squirrel para prevenir SQL injection
- Los timestamps se manejan como Unix timestamps (int64)
- El sistema es idempotente para operaciones de desactivación 