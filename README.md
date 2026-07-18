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
- Dedicated administrative dashboard mapping click metrics in real-time.
- Uncompromising minimalist design language (100% sharp-corners, high-contrast primary colors).

---

## 🔥 PRODUCTION INFRASTRUCTURE TROUBLESHOOTING MATRIX

If you encounter systemic operational failures during public cloud deployments, consult this matrix:

- **Error A: 401 Unauthorized / Invalid Token**
  - **Fix:** Trim quotes out of `.env` environment strings. Ensure `ADMIN_SECRET_HASH` matches exact hash bytes natively generated.

- **Error B: 403 Forbidden**
  - **Fix:** Verify Origin matching header inside `ALLOWED_ORIGINS` string array keys in the Go backend configuration.

- **Error C: 500 Internal Server Error**
  - **Fix:** Restart background service handlers (`systemctl restart lmq-backend`) and check the database connection pooling availability on Supabase.

- **Error D: Vite/Svelte Network Error**
  - **Fix:** Verify Nginx proxy rules mapping matching paths to `127.0.0.1:8080`. Ensure `/api` is cleanly isolated from `/` front-end requests.
