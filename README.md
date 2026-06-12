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
- `custom_url`: Optional, must be at least 3 characters. Alphanumeric (case-insensitive), hyphen (`-`), and underscore (`_`) are allowed.

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

## Database Design

The application uses PostgreSQL with GORM. Tables are automatically created/migrated via GORM's `AutoMigrate` when the application starts, provided that the database itself already exists.

### Table Schema

#### 1. `urls` Table
Stores information about shortened URLs.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | SERIAL | Primary Key | Auto-incrementing identifier |
| `short_code` | VARCHAR(10) | Unique, Not Null | The shortened path/alias |
| `original` | VARCHAR(2048) | Not Null | The original destination URL |
| `custom` | BOOLEAN | Default `false` | Indicates if it is a custom short code |
| `active` | BOOLEAN | Default `true` | Link active status |
| `created_at` | TIMESTAMPTZ | Not Null | Timestamp when record was created |
| `updated_at` | TIMESTAMPTZ | Not Null | Timestamp when record was updated |
| `deleted_at` | TIMESTAMPTZ | Index | Used for GORM soft delete |

#### 2. `admins` Table
Stores administrator credentials for the dashboard.

| Column | Type | Constraints | Description |
|---|---|---|---|
| `id` | SERIAL | Primary Key | Auto-incrementing identifier |
| `username` | VARCHAR(50) | Unique, Not Null | Admin username |
| `password` | VARCHAR(255) | Not Null | Bcrypt hashed password |
| `created_at` | TIMESTAMPTZ | Not Null | Timestamp when record was created |
| `updated_at` | TIMESTAMPTZ | Not Null | Timestamp when record was updated |
| `deleted_at` | TIMESTAMPTZ | Index | Used for GORM soft delete |

---

## Database Setup & Troubleshooting

### Error: "database ... does not exist"
If you get a database connection failure stating that the database does not exist, it means:
1. **Local Setup:** The PostgreSQL server is running, but the database (e.g. `urlshortener`) hasn't been created yet. GORM's `AutoMigrate` can create tables, columns, and indexes, but **cannot create the database itself** because it requires a successful connection to a database first.
2. **Managed/Cloud Database (Render, Railway, Neon, etc.):** Managed database services usually provision a single, specific database for you (e.g., `papa-postgres`). You do not need to (and often cannot) run `CREATE DATABASE`. Instead, you should configure your application to use the existing database name provided by the platform.

---

### How to Create / Configure the Database

#### Option 1: Using Docker Compose (Recommended for local dev)
If you start the service using Docker Compose, the database will be created automatically:
```bash
docker compose up -d
```
The database name is predefined as `urlshortener` via the `POSTGRES_DB` environment variable in the `docker-compose.yml` file, which PostgreSQL initializes automatically on first startup.

#### Option 2: Using Managed Cloud Databases (Render, Railway, Supabase, etc.)
If you are deploying or connecting to a managed/cloud PostgreSQL instance:
1. Do not try to create a new database.
2. Update your `.env` or configuration variables on the hosting platform to point to the provided database name:
   * **`DB_NAME`**: Use the database name given in your dashboard (e.g., `papa-postgres`).
   * **`DB_USER`**: Use the database username given (e.g., `papa-lab-postgres`).
   * **`DB_HOST`**: Use the internal or external host provided.
   * **`DB_PASSWORD`**: Use the password provided.
3. GORM will connect to this pre-existing database and automatically create the required tables (`urls`, `admins`) on start.

#### Option 3: Manual Database Creation (Local Postgres)
If you are running a local PostgreSQL instance (outside Docker), you must create the database before launching the Go application.

1. **Using PostgreSQL CLI (`psql`):**
   When running `psql`, if you do not specify a database name using the `-d` flag, PostgreSQL defaults to connecting to a database with the same name as the user (which might fail with `database "user" does not exist`).
   
   To connect and create the database:
   ```bash
   # Connect by specifying both user (-U) and default database (-d), e.g. postgres
   psql -U postgres -d postgres
   ```
   Once connected, execute the SQL command to create the database:
   ```sql
   CREATE DATABASE urlshortener;
   ```
   Or directly from your terminal/command prompt:
   ```bash
   createdb -U postgres urlshortener
   ```

2. **Using pgAdmin or other GUI Clients:**
   * Open pgAdmin and connect to your server.
   * Right-click on **Databases** -> **Create** -> **Database...**
   * Enter `urlshortener` as the Database name.
   * Click **Save**.