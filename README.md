# papa-shortener

Simple URL shortener service in Go with PostgreSQL and GORM.

## Quick Start

```bash
# Start services
docker compose up -d

# Create a short URL (random)
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"original_url": "https://github.com"}'

# Create a short URL (custom)
curl -X POST http://localhost:8080/api/shorten \
  -H "Content-Type: application/json" \
  -d '{"original_url": "https://github.com", "custom_url": "gh"}'

# Visit short URL
curl -I http://localhost:8080/gh
```

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/shorten` | Create short URL |
| GET | `/:code` | Redirect to original URL |
| GET | `/health` | Health check |

## Request/Response Examples

### Create Short URL (Random)

**Request:**
```json
POST /api/shorten
{
  "original_url": "https://github.com"
}
```

**Response:**
```json
{
  "short_code": "a1b2c3d4",
  "original_url": "https://github.com",
  "short_url": "a1b2c3d4",
  "is_custom": false,
  "full_short_url": "http://localhost:8080/a1b2c3d4"
}
```

### Create Short URL (Custom)

**Request:**
```json
POST /api/shorten
{
  "original_url": "https://github.com",
  "custom_url": "mygh"
}
```

**Response:**
```json
{
  "short_code": "mygh",
  "original_url": "https://github.com",
  "short_url": "mygh",
  "is_custom": true,
  "full_short_url": "http://localhost:8080/mygh"
}
```

## Validation Rules

- `original_url`: Required, must be a valid URL
- `custom_url`: Optional, must be 3-20 lowercase alphanumeric characters only

## Configuration

Edit `config.yaml` to configure the application:

```yaml
database:
  host: postgres
  port: 5432
  user: postgres
  password: postgres
  name: urlshortener
  sslmode: disable

app:
  host: 0.0.0.0
  port: 8080
  base_url: http://localhost:8080
```