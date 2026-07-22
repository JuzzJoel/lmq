# LMQ - Modern URL Shortening Platform

<p align="center">
  <img src="https://raw.githubusercontent.com/JuzzJoel/joel-olajiire/master/public/covers/lmq.jpg" alt="LMQ URL Shortener Screenshot" width="800" />
</p>

LMQ is a high-performance, modular URL shortening platform utilizing a brutalist Bauhaus aesthetic. It provides direct, highly scalable link manipulation with deeply integrated real-time analytics, dynamic visual QR code generation, and robust IP-based rate limiting.

## Architectural Layout
- **Frontend Layer:** SvelteKit 5, utilizing raw, uncompromising HTML/CSS mapped to a unified `adapter-node` SSR server cluster for immediate user-response times and perfectly hydrated layouts.
- **Backend API Layer:** Go, driven by the Chi router, connected to Supabase PostgreSQL, and accelerated by an ephemeral Redis cache to serve high-velocity link redirects safely under heavy network load.

## Features
- Scalable bulk link shortening (JSON and CSV upload).
- Configurable URL expiration timestamps.
- Hash-secured password barriers on protected links.
- GeoIP tracking (City, Region, Country tracking).
- **A/B testing**: weighted multi-destination routing per link.
- **Burn-after-reading**: self-destructing one-time links.
- **Tags/Campaigns**: organize links with labels, filterable via API.
- **CSV analytics export**: download click data for offline analysis.
- Dedicated administrative dashboard mapping click metrics in real-time.
- Uncompromising minimalist design language (100% sharp-corners, high-contrast primary colors).

---

## A/B Testing (Weighted Routing)

LMQ supports routing a single short link to multiple destinations based on configurable weights.

**Create a link with A/B routes:**
```bash
curl -X POST http://localhost:8080/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com/base",
    "routes": [
      {"url": "https://example.com/variant-a", "weight": 70},
      {"url": "https://example.com/variant-b", "weight": 30}
    ]
  }'
```

**Response** (201 Created):
```json
{
  "data": {
    "results": [{
      "token": "AbCdEf",
      "short_url": "https://lmq.name.ng/AbCdEf",
      "long_url": "https://example.com/base",
      "routes": [
        {"url": "https://example.com/variant-a", "weight": 70},
        {"url": "https://example.com/variant-b", "weight": 30}
      ],
      "created_at": "..."
    }]
  },
  "error": null
}
```

**Redirect behavior:**
- Each visit performs a weighted random selection among the configured routes.
- The `routes` array appears in all analytics responses (list and detail).
- Routes are optional; if omitted, the link redirects directly to `long_url`.

---

## Burn After Reading

A burn-after-reading (BAR) link self-destructs after the first visit. The link is atomically deleted on access, ensuring it can only be viewed once.

**Create a BAR link:**
```bash
curl -X POST http://localhost:8080/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com/sensitive",
    "burn_after_reading": true
  }'
```

**Behavior:**
- The first `GET /{token}` redirects to the target URL and atomically deletes the link.
- Any subsequent `GET /{token}` returns a 404 (link no longer exists).
- BAR can be combined with A/B routes.
- Password-protected BAR links: the link is consumed upon successful password verification (via `POST /verify-password`).
- No analytics are recorded for BAR links (the row is deleted before the analytics insert).

---

## Tags / Campaigns

Every link can carry an optional list of tags for organizing and filtering links by campaign, category, or any label.

**Create a link with tags:**
```bash
curl -X POST http://localhost:8080/api/v1/shorten \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com/page",
    "tags": ["marketing", "q1-2026", "promo"]
  }'
```

**Filter links by tag:**
```bash
curl -H "X-Admin-Token: your-token" \
  "http://localhost:8080/api/v1/analytics/links?tag=promo"
```

Tags appear in the response of create, list, and analytics detail endpoints.

---

## Bulk CSV Upload

Create up to 50 links at once by uploading a CSV file.

**Endpoint:** `POST /api/v1/shorten/csv` (multipart/form-data)

**CSV columns:**
| Column | Required | Description |
|--------|----------|-------------|
| `url` | Yes | Target URL |
| `custom_token` | No | Custom alias |
| `password` | No | Password protection |
| `expires_in` | No | Expiration in hours |
| `burn_after_reading` | No | `true` or `false` |
| `tags` | No | Comma-separated tags |

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/shorten/csv \
  -F "file=@links.csv"
```

Example `links.csv`:
```csv
url,custom_token,password,expires_in,burn_after_reading,tags
https://example.com/page1,,,24,true,campaign-a
https://example.com/page2,my-alias,secret,0,false,campaign-b
```

---

## Analytics Export

Download click analytics as CSV for offline analysis.

**Endpoint:** `GET /api/v1/analytics/export` (requires `X-Admin-Token` header)

**Optional filters:**
| Parameter | Description |
|-----------|-------------|
| `token` | Filter to a specific short link |
| `from` | Start date (e.g. `2026-01-01`) |
| `to` | End date (e.g. `2026-12-31`) |

**Example:**
```bash
curl -H "X-Admin-Token: your-token" \
  "http://localhost:8080/api/v1/analytics/export?token=AbCdEf&from=2026-01-01" \
  -o analytics.csv
```

**CSV columns:** `token`, `long_url`, `clicked_at`, `ip_address`, `city`, `region`, `country_code`, `user_agent`, `browser`, `os`, `is_mobile`, `referer`

---

## 🔥 PRODUCTION INFRASTRUCTURE TROUBLESHOOTING MATRIX

If you encounter systemic operational failures during public cloud deployments, consult this matrix:

- **Error A: 401 Unauthorized / Invalid Token**
  - **Fix:** Set `ADMIN_SECRET_HASH` in your `.env` to the SHA-256 hash of the admin password (not the plaintext password itself). The middleware compares this directly against SHA-256(incoming token). To generate the hash: `echo -n "your-admin-password" | sha256sum`.

- **Error B: 403 Forbidden**
  - **Fix:** Verify Origin matching header inside `ALLOWED_ORIGINS` string array keys in the Go backend configuration.

- **Error C: 500 Internal Server Error**
  - **Fix:** Restart background service handlers (`systemctl restart lmq-backend`) and check the database connection pooling availability on Supabase.

- **Error D: Vite/Svelte Network Error**
  - **Fix:** Verify Nginx proxy rules mapping matching paths to `127.0.0.1:8080`. Ensure `/api` is cleanly isolated from `/` front-end requests.
