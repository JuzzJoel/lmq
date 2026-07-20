# LMQ - Modern URL Shortening Platform

LMQ is a high-performance, modular URL shortening platform utilizing a brutalist Bauhaus aesthetic. It provides direct, highly scalable link manipulation with deeply integrated real-time analytics, dynamic visual QR code generation, and robust IP-based rate limiting.

## Architectural Layout
- **Frontend Layer:** SvelteKit 5, utilizing raw, uncompromising HTML/CSS mapped to a unified `adapter-node` SSR server cluster for immediate user-response times and perfectly hydrated layouts.
- **Backend API Layer:** Go, driven by the Chi router, connected to Supabase PostgreSQL, and accelerated by an ephemeral Redis cache to serve high-velocity link redirects safely under heavy network load.

## Features
- Scalable bulk link shortening.
- Configurable URL expiration timestamps.
- Hash-secured password barriers on protected links.
- GeoIP tracking (City, Region, Country tracking).
- **A/B testing**: weighted multi-destination routing per link.
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
