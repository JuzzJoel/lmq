# --- STAGE 1: Compile SvelteKit Web Static Bundle ---
FROM node:22-alpine AS frontend-builder
WORKDIR /app
COPY package.json pnpm-workspace.yaml ./
COPY frontend/package.json ./frontend/
RUN npm install -g pnpm && pnpm install
COPY frontend/ ./frontend/
RUN pnpm --filter frontend build

# --- STAGE 2: Compile High-Performance Go Executable ---
FROM golang:1.23-alpine AS backend-builder
WORKDIR /app
COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download
COPY backend/ ./backend/
# Inject the compiled SvelteKit bundle right into Go's embed workspace context
COPY --from=frontend-builder /app/frontend/build ./backend/frontend/dist
RUN cd backend && CGO_ENABLED=0 go build -ldflags="-w -s" -o main-server .

# --- STAGE 3: Final Production Execution Layer ---
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/
COPY --from=backend-builder /app/backend/main-server .
# Create target data directory path for GeoIP uploads
RUN mkdir -p data
EXPOSE 8080
CMD ["./main-server"]
