# LMQ Core — Unified High-Performance URL Shortener

An enterprise-grade, single-binary URL shortener platform. The system utilizes an ultra-efficient Go backend engine that encapsulates a compiled SvelteKit frontend enhanced with Tailwind CSS v4. Data caching and sliding-window rate limiting are handled by Upstash Redis, relational persistence by Supabase PostgreSQL, and geolocation mapping by an offline GeoIP database.

## 📁 System Architecture

```
lmq/
├── package.json
├── pnpm-workspace.yaml
├── Dockerfile
├── backend/                  # Go Core Web Service & Redirector
│   ├── main.go
│   ├── go.mod
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   ├── postgres.go
│   │   ├── redis.go
│   │   └── migrations/
│   │       └── 001_init.sql
│   ├── models/
│   │   └── models.go
│   ├── handlers/
│   │   ├── shorten.go
│   │   ├── redirect.go
│   │   └── analytics.go
│   ├── middleware/
│   │   ├── ratelimit.go
│   │   └── cors.go
│   ├── services/
│   │   ├── geoip.go
│   │   ├── analytics.go
│   │   └── token.go
│   ├── spa/
│   │   └── spa.go
│   ├── frontend/dist/        # Embedded SPA (populated at build time)
│   └── data/                 # GeoIP database directory
└── frontend/                 # SvelteKit Application
    ├── package.json
    ├── svelte.config.js
    ├── vite.config.ts
    ├── tsconfig.json
    └── src/
        ├── app.html
        ├── app.css
        ├── lib/
        │   ├── api.ts
        │   ├── types.ts
        │   └── components/
        │       ├── Seo.svelte
        │       ├── Chart.svelte
        │       ├── StatCard.svelte
        │       ├── LinkTable.svelte
        │       ├── UrlForm.svelte
        │       └── Skeleton.svelte
        └── routes/
            ├── +layout.svelte
            ├── +layout.ts
            ├── +page.svelte
            └── dashboard/
                ├── +layout.svelte
                ├── +page.svelte
                └── [token]/
                    ├── +page.ts
                    └── +page.svelte
```

---

## ⚡ Local Development Pipeline

To avoid heavy resource demands on your machine, run the development environment as two hot-reloading processes decoupled over a local network proxy.

### 1. Initialization

Ensure your global command runner and package manager pathways are securely configured:

```cmd
pnpm setup
```

Install all workspace project packages from the root directory:

```cmd
pnpm install
```

### 2. Configuration Setup (`.env`)

Create an environment file inside the `backend/` directory to handle local secrets:

```env
PORT=8080
DATABASE_URL="postgres://postgres.[PROJECT_ID]:[PASSWORD]@aws-0-[REGION].pooler.supabase.com:5432/postgres"
REDIS_URL="rediss://default:[TOKEN]@[ENDPOINT].upstash.io:6379"
GEOIP_DB_PATH="data/GeoLite2-City.mmdb"
```

### 3. Running the Dev Instance

Open two separate terminal consoles from your workspace root:

**Console A: SvelteKit Frontend Engine**

```cmd
pnpm --filter frontend dev
```

Spins up the SvelteKit application on `http://localhost:5173`. Vite's configuration automatically proxies any incoming `/api/*` or token request over to the Go port.

**Console B: Go Backend Core**

```cmd
cd backend
go run main.go
```

Launches your API endpoint manager and redirection loop on `http://localhost:8080`.

---

## 📡 Provisioning Cloud Resources

### 1. Relational Storage (Supabase)

1. Register a project cluster on [Supabase](https://supabase.com).
2. Navigate to the **Connect** options panel on your dashboard.
3. Select **Poolers**, set the connection behavior configuration to **Session Mode**, and copy the complete database connection string into your `DATABASE_URL` environment parameter.

### 2. Cache & Rate-Limiting (Upstash)

1. Launch a database instance on [Upstash](https://upstash.com) choosing **Redis**.
2. Enable **TLS enforcement** to guarantee secure network transits.
3. Extract the unified connection string (`rediss://...`) and bind it to your `REDIS_URL` target variable.

### 3. Geolocation Mapping Setup (MaxMind)

1. Log into your MaxMind workspace dashboard shown in your account view.
2. Locate the row titled GeoLite2 City (Edition ID: GeoLite2-City).
3. Click Download GZIP next to the Binary .mmdb format asset.
4. Extract the compressed file using a file extractor (like 7-Zip or WinRAR).
5. Move the internal uncompressed database file named `GeoLite2-City.mmdb` into your local project workspace directory path: `C:\Users\user\Documents\GitHub\lmq\backend\data\GeoLite2-City.mmdb`.

### 4. The New Relational Analytic Layout (001_init.sql)

The migration script executed inside Supabase now maps geolocation tracking metrics cleanly down to the city level:

```sql
-- Links master manifest table
CREATE TABLE IF NOT EXISTS links (
    id          BIGSERIAL PRIMARY KEY,
    token       VARCHAR(20) UNIQUE NOT NULL,
    long_url    TEXT NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ,
    click_count BIGINT NOT NULL DEFAULT 0
);
CREATE INDEX idx_links_token ON links(token);

-- Expanded Click events table for deep city analytics
CREATE TABLE IF NOT EXISTS click_events (
    id           BIGSERIAL PRIMARY KEY,
    link_id      BIGINT NOT NULL REFERENCES links(id) ON DELETE CASCADE,
    clicked_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address   INET,
    city         VARCHAR(100) DEFAULT 'Unknown',
    region       VARCHAR(100) DEFAULT 'Unknown',
    country_code VARCHAR(2) DEFAULT 'XX',
    user_agent   TEXT,
    browser      VARCHAR(100),
    os           VARCHAR(100),
    is_mobile    BOOLEAN DEFAULT FALSE,
    referer      TEXT
);
CREATE INDEX idx_click_events_link_id ON click_events(link_id);
CREATE INDEX idx_click_events_city ON click_events(city);
CREATE INDEX idx_click_events_clicked_at ON click_events(clicked_at);
```

> **Note:** The Go source code features an integrated safety validation layer that gracefully logs a fallback warning without crashing the boot runtime if this file is missing.

---

## 🚀 Single-Binary Multi-Stage Production Build

Production environments utilize a multi-stage compilation loop. This isolates compiler footprints, builds the Tailwind-extracted static web directories, passes them natively into Go's binary composition engine, and drops the finished asset into a highly optimized container.

### Building Locally

```cmd
docker build -t lmq .
docker run -p 8080:8080 --env-file backend/.env lmq
```

### Deploying to Render

1. Commit and push your local workspace code structure to your GitHub repository.
2. Log in to your cloud compute space at [Render](https://render.com).
3. Initiate a **Web Service** project pipeline tracking your repository.
4. Set the **Environment/Language Runtime** choice strictly to **Docker**.
5. Map your secure cloud secrets (`DATABASE_URL`, `REDIS_URL`, etc.) inside Render's advanced environment variable panel.
6. Click **Deploy**. Render will automatically build your single-binary system container.

---

## 🔒 API Reference

### `POST /api/v1/shorten`

Create a shortened URL.

**Request Body:**
```json
{
  "url": "https://example.com/very/long/path",
  "custom_token": "my-brand"
}
```

**Response `201 Created`:**
```json
{
  "data": {
    "token": "aB3xK7",
    "short_url": "https://lmq.app/aB3xK7",
    "long_url": "https://example.com/very/long/path",
    "created_at": "2026-07-18T03:00:00Z"
  },
  "error": null
}
```

### `GET /api/v1/analytics?token=aB3xK7`

Retrieve analytics for a specific link.

**Response `200 OK`:**
```json
{
  "data": {
    "token": "aB3xK7",
    "long_url": "https://example.com/very/long/path",
    "total_clicks": 1423,
    "clicks_by_day": [...],
    "countries": [...],
    "browsers": [...],
    "recent_clicks": [...]
  },
  "error": null
}
```

### `GET /api/v1/analytics/links`

List all shortened links with summary counts.

### `GET /{token}`

Redirects visitor to the original URL via **HTTP 301**. Records analytics asynchronously.

---

## ⚙️ Technical Specifications

| Component | Technology | Version |
|-----------|-----------|---------|
| Backend Router | Chi v5 | `go-chi/chi/v5` |
| Database | PostgreSQL (Supabase) | pgx/v5 |
| Cache | Redis (Upstash) | go-redis/v9 |
| GeoIP | MaxMind GeoLite2-Country | maxminddb-golang/v2 |
| User-Agent Parser | mileusna/useragent | Latest |
| Frontend Framework | SvelteKit | Latest |
| CSS Framework | Tailwind CSS | v4 |
| Charts | Chart.js | Latest |
| Language | Go 1.23+ / TypeScript | Strict |
