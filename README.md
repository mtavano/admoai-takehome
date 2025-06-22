# AdMoai Take Home Test - Ads API

## 📋 Description

REST API for ad management with automatic TTL (Time To Live) system. Allows creating, querying and managing ads with automatic time-based expiration.

## 🚀 Features

- **CRUD Operations**: Create, read, update and delete ads
- **TTL System**: Automatic ad expiration based on minutes
- **Advanced Filters**: Query by placement, status and other criteria
- **SQLite Database**: Local storage with automatic migrations
- **Validations**: Input validation with Gin and golang validator
- **Query Builder**: Use of Squirrel for dynamic and secure queries
- **Metrics & Monitoring**: Prometheus-formatted metrics endpoint
- **Alerting**: Simulated alerts for high active ad counts

## 🛠️ Technologies

- **Go 1.23.4**
- **Gin** - Web framework
- **SQLite** - Database
- **Squirrel** - Query builder
- **Goose** - Database migrations
- **Golang Validator** - Data validation

## 📦 Installation

### Prerequisites
- Go 1.23.4 or higher
- SQLite3

### Setup
1. Clone the repository:
```bash
git clone <repository-url>
cd admoai-takehome
```

2. Install dependencies:
```bash
go mod tidy
```

3. Configure environment variables:
```bash
cp example.dev.env dev.env
```

4. Run migrations:
```bash
make setup
```

5. Start the server:
```bash
make run-simple
```

## 🚀 Deployment

### Coolify Deployment

This project includes configuration files for easy deployment on Coolify:

#### **Files Included:**
- `nixpacks.toml` - Nixpacks configuration for Coolify
- `Dockerfile` - Alternative Docker deployment
- `.dockerignore` - Optimized Docker build
- `coolify.env.example` - Environment variables template

#### **Deployment Steps:**

1. **Connect Repository to Coolify:**
   - Add your Git repository to Coolify
   - Select the main branch

2. **Configure Build Settings:**
   - **Build Pack**: Nixpacks (recommended) or Docker
   - **Port**: 9001
   - **Root Directory**: `/` (default)

3. **Set Environment Variables:**
```env
ENVIRONMENT=production
API_PORT=9001
DB_DRIVER=sqlite3
DB_DSN=/root/data/admoai.db
```

4. **Deploy:**
   - Click "Deploy" in Coolify
   - The application will build and start automatically

#### **Nixpacks Configuration:**
The `nixpacks.toml` file configures:
- **Dependencies**: Go, GCC, SQLite
- **Build Process**: Downloads dependencies and builds binary
- **Runtime**: Alpine Linux with SQLite support
- **Start Command**: `./server`

#### **Docker Alternative:**
If you prefer Docker, the `Dockerfile` provides:
- Multi-stage build for optimized image size
- Alpine Linux base for security
- SQLite runtime support
- Production-ready configuration

### Environment Variables for Production

Copy `coolify.env.example` and configure:

```env
# Application Configuration
ENVIRONMENT=production
API_PORT=9001

# Database Configuration
DB_DRIVER=sqlite3
DB_DSN=/root/data/admoai.db

# Optional: External Database
# DB_DRIVER=postgres
# DB_DSN=postgres://username:password@host:port/database?sslmode=disable

# Logging Configuration
LOG_LEVEL=info
```

### Health Check

The application includes a health check endpoint:
```
GET /health
```

Configure this in Coolify for automatic health monitoring.

## 🔧 Configuration

### Environment Variables (`dev.env`)
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

### 1. Create Ad
**POST** `/ads`

Creates a new ad with optional TTL.

**Request Body:**
```json
{
  "title": "Test Ad",
  "image_url": "https://example.com/image.jpg",
  "placement": "homepage",
  "ttl": 30
}
```

**Fields:**
- `title` (required): Ad title
- `image_url` (required): Image URL (must be valid URL)
- `placement` (required): Ad placement
- `ttl` (optional): Time to live in minutes (0 = no expiration)

**Response (201):**
```json
{
  "id": "uuid-generated",
  "title": "Test Ad",
  "imageUrl": "https://example.com/image.jpg",
  "placement": "homepage",
  "status": "active",
  "createdAt": 1640995200,
  "expiresAt": 1640997000
}
```

### 2. Get Ad by ID
**GET** `/ads/{id}`

Gets a specific ad by its ID.

**Response (200):**
```json
{
  "id": "uuid-here",
  "title": "Test Ad",
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

### 3. Filter Ads
**GET** `/ads?placement=homepage&status=active`

Gets ads filtered by specific criteria.

**Query Parameters:**
- `placement` (optional): Filter by placement
- `status` (optional): Filter by status

**Response (200):**
```json
{
  "ads": [
    {
      "id": "uuid-1",
      "title": "Ad 1",
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

### 4. Deactivate Ad
**POST** `/ads/{id}/deactivate`

Deactivates a specific ad.

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

Verifies service status.

**Response (200):**
```json
{
  "status": "running"
}
```

### Metrics (Prometheus Format)
```bash
GET /metrics
```

## 🗄️ Data Model

### AdvertiseRecord
```go
type AdvertiseRecord struct {
ID string `db:"id" json:"id"`
Title string `db:"title" json:"title"`
ImageURL string `db:"image_url" json:"imageUrl"`
Placement string `db:"placement" json:"placement"`
Status string `db:"status" json:"status"`
CreatedAt int64 `db:"created_at" json:"createdAt"`
ExpiresAt *int64 `db:"expires_at" json:"expiresAt"`
}
```

### Ad States
- `active`: Active and visible ad
- `inactive`: Deactivated ad

## ⏰ TTL System

### Behavior
- Ads can have a TTL (Time To Live) in minutes
- If `ttl = 0` or not specified, the ad doesn't expire
- Expired ads are automatically filtered out from queries
- The `expiresAt` field is calculated as `createdAt + (ttl * 60 seconds)`

### Automatic Filtering
All queries automatically discard expired ads:
- Ads without TTL (`expiresAt = null`) are always shown
- Ads with valid TTL (`expiresAt > currentTime`) are shown
- Expired ads (`expiresAt <= currentTime`) are discarded

## 🛠️ Make Commands

```bash
# Initial setup
make setup

# Run server with hot reload
make run

# Run simple server
make run-simple

# Run migrations
make migrate

# Rollback migrations
make migrate-down

# Clean database
make clean

# Create new migration
make create-migration name=migration_name
```

## 🧪 Functional Testing

### Usage examples with curl

**Create ad:**
```bash
curl -X POST http://localhost:9001/v1/ads \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Test Ad",
    "image_url": "https://example.com/image.jpg",
    "placement": "homepage",
    "ttl": 30
  }'
```

**Get ad by ID:**
```bash
curl -X GET http://localhost:9001/v1/ads/your-uuid-here
```

**Filter ads:**
```bash
curl -X GET "http://localhost:9001/v1/ads?placement=homepage&status=active"
```

**Deactivate ad:**
```bash
curl -X POST http://localhost:9001/v1/ads/your-uuid-here/deactivate
```

## 📁 Project Structure

```
admoai-takehome/
├── cmd/
│ └── server/
│ └── main.go # Entry point
├── internal/
│ ├── api/ # API handlers
│ │ ├── handle_func.go
│ │ ├── post_ads_handler.go
│ │ ├── get_ads_by_id.go
│ │ ├── get_ads_by_filters_handler.go
│ │ ├── post_deactivate_ads_handler.go
│ │ └── register_routes.go
│ └── store/ # Data layer
│ ├── database.go
│ ├── implementation.go
│ ├── models.go
│ └── query/
│ ├── insert_ads.go
│ ├── select_ads.go
│ └── update_ads.go
├── migrations/ # Database migrations
│ └── 20250622182727_init_setup.go
├── dev.env # Environment variables
├── go.mod # Dependencies
├── Makefile # Useful commands
├── nixpacks.toml # Coolify deployment config
├── Dockerfile # Docker deployment
├── .dockerignore # Docker optimization
├── coolify.env.example # Production env template
└── README.md # This file
```

## 🔒 Validations

### Input Validations
- **title**: Cannot be empty
- **image_url**: Cannot be empty and must be valid URL
- **placement**: Cannot be empty
- **ttl**: Optional, must be greater than 0 if specified

### Business Validations
- At least one filter must be present in filter queries
- Expired ads are automatically discarded
- Deactivated ads maintain their state

## 🚨 Error Codes

- **400 Bad Request**: Invalid input data
- **404 Not Found**: Resource not found
- **500 Internal Server Error**: Internal server error

## 🏗️ Implementation Details

### Runtime Expiration System

The expiration system works **at runtime** instead of using a cronjob. This architectural decision was made to:

- **Minimize synchronization issues**: Avoid inconsistencies between multiple application instances
- **Operational simplicity**: No need to manage scheduled jobs
- **Scalability**: Each request handles its own expiration logic
- **Consistency**: Data is always up to date at query time

Expiration is automatically verified in each query through SQL filters that discard expired records based on the `expires_at` field.

### Ad States: Active, Inactive and Expired

There's an important distinction between ad states:

#### **Active (`status = 'active'`)**
- Visible and functional ad
- May or may not have TTL configured
- Shown in queries if not expired

#### **Inactive (`status = 'inactive'`)**
- Manually deactivated ad by user
- State change performed through `/ads/{id}/deactivate` endpoint
- Not shown in queries regardless of TTL

#### **Expired (`expires_at` column)**
- Ad that has exceeded its configured lifetime
- The `status` remains `active`, but is automatically filtered
- Discarded in queries even though status is active

**Result for Frontend**: Both `inactive` and `expired` ads are not shown, but at data level they represent different concepts:
- **Inactive**: Manual deactivation (reversible)
- **Expired**: Automatic expiration by TTL (irreversible)

### HandleFunc as Decorator for Testing

The implementation uses a decorator pattern with `HandleFunc` that provides several advantages:

#### **Logic Abstraction**
```go
func HandleFunc(handler func(*gin.Context, *Context) (any, int, error), ctx *Context) gin.HandlerFunc {
    return func(c *gin.Context) {
        // Common error handling and response logic
        result, status, err := handler(c, ctx)
        // Serialization and response
    }
}
```

#### **Testing Benefits**
- **Separation of concerns**: Business logic is isolated from web framework
- **Direct unit testing**: Handlers can be tested without HTTP server
- **Simplified mocking**: Easy to mock dependencies (database, etc.)
- **Clear assertions**: Can directly verify return values

#### **Testing Example**
```go
func TestPostAdsHandler(t *testing.T) {
    // Arrange
    mockDB := &MockDatabase{}
    ctx := &Context{Db: mockDB}
    
    // Act
    result, status, err := PostAdsHandler(ginContext, ctx)
    
    // Assert
    assert.Equal(t, http.StatusCreated, status)
    assert.Nil(t, err)
    // Verify result...
}
```

This architecture facilitates robust unit test development and keeps code clean and maintainable.

## 📝 Development Notes

- API uses SQLite for development simplicity
- All queries use Squirrel to prevent SQL injection
- Timestamps are handled as Unix timestamps (int64)
- System is idempotent for deactivation operations

## ⚠️ Postmortem

### Identified Issues and Improvement Areas

#### **1. Data Model Inconsistency**
**Problem**: The separation between `status` (active/inactive) and `expires_at` can be confusing and error-prone.

**Impact**:
- Clients must understand two different concepts to determine if an ad is "active"
- Easy to forget applying expiration filter in new queries
- Can generate inconsistencies if filtering method is not always called

**Proposed Solution**:
- Implement a calculated field `is_visible` that combines both states
- Create a `CalculateAndSetExpired()` method that automatically updates status
- Use database triggers to maintain consistency

#### **2. Lack of Pagination**
**Problem**: Filtering endpoints don't implement pagination.

**Impact**:
- Can cause performance issues with large data volumes
- No control over response size
- Difficult to implement infinite scroll in frontend

**Proposed Solution**:
- Add `limit` and `offset` parameters to filtering endpoints
- Implement cursor-based pagination for better performance
- Add pagination metadata in responses

#### **3. Absence of Automated Tests**
**Problem**: No unit tests or integration tests implemented.

**Impact**:
- Difficult to detect regressions
- No automatic validation of changes
- Higher risk in deployments

**Proposed Solution**:
- Implement unit tests for all handlers
- Add integration tests for database
- Configure CI/CD with automatic validation

#### **4. Inconsistent Error Handling**
**Problem**: Errors don't follow a standard format and lack structured logging.

**Impact**:
- Difficult debugging in production
- Clients receive uninformative error messages
- No error traceability

**Proposed Solution**:
- Implement standard error code system
- Add structured logging with levels
- Create centralized error handling middleware

#### **5. Lack of Business Validation**
**Problem**: No business validations like TTL limits or placement validation.

**Impact**:
- Can create ads with extreme TTLs (e.g., 100 years)
- No control over valid placement values
- Difficult to maintain data consistency

**Proposed Solution**:
- Add range validations for TTL (e.g., 1 minute - 1 year)
- Implement enum of valid placements
- Create centralized business validations

#### **6. Absence of Metrics and Monitoring**
**Problem**: No performance metrics or application monitoring.

**Impact**:
- Cannot detect performance issues
- Difficult to identify bottlenecks
- No automatic alerts

**Proposed Solution**:
- Implement Prometheus metrics
- Add detailed health checks
- Configure alerts for errors and latency

#### **7. Lack of API Documentation**
**Problem**: No OpenAPI/Swagger documentation.

**Impact**:
- Difficult for other developers to integrate with API
- No automatic request/response validation
- Lack of autocomplete in IDEs

**Proposed Solution**:
- Generate OpenAPI documentation automatically
- Implement schema validation
- Add usage examples in documentation

#### **8. Hardcoded Configuration**
**Problem**: Many values are hardcoded in the code.

**Impact**:
- Difficult configuration per environment
- No flexibility for different deployments
- Difficult testing with different configurations

**Proposed Solution**:
- Move configurations to environment variables
- Implement environment-based configuration system
- Add configuration validation at startup

#### **9. Lack of Rate Limiting**
**Problem**: No protection against API abuse.

**Impact**:
- Vulnerable to DoS attacks
- A client can saturate the system
- No usage control per client

**Proposed Solution**:
- Implement rate limiting per IP/client
- Add authentication and authorization
- Configure request limits per minute

#### **10. Absence of Backup and Recovery**
**Problem**: No backup strategy for database.

**Impact**:
- Data loss in case of failure
- No disaster recovery
- Difficult data migration

**Proposed Solution**:
- Implement automatic backups
- Create recovery scripts
- Document DR procedures

### Current Strengths
- ✅ Clean architecture with separation of concerns
- ✅ Comprehensive input validation
- ✅ Flexible filtering system
- ✅ Automatic TTL expiration
- ✅ Metrics and monitoring
- ✅ Docker containerization
- ✅ Database migrations
- ✅ Comprehensive documentation

### Areas for Improvement
- 🔄 Add unit tests and integration tests
- 🔄 Implement proper logging with structured logs
- 🔄 Add rate limiting and authentication
- 🔄 Implement caching layer
- 🔄 Add database connection pooling
- 🔄 Implement proper error handling with error codes
- 🔄 Add API versioning strategy
- 🔄 Implement proper health checks with database connectivity
- 🔄 Add configuration management
- 🔄 Implement proper graceful shutdown

### Proposed Enhancements
- 📈 Add Grafana dashboards
- 📈 Implement real alerting system (email, Slack, etc.)
- 📈 Add performance monitoring
- 📈 Implement audit logging
- 📈 Add backup and recovery procedures
- 📈 Implement horizontal scaling strategy
- 📈 Add API documentation with Swagger/OpenAPI
- 📈 Implement proper CI/CD pipeline
