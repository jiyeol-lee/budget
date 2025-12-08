# =============================================================================
# Stage 1: Frontend Builder
# =============================================================================
FROM node:20-alpine AS frontend-builder

WORKDIR /build
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci --frozen-lockfile
COPY frontend/ ./
RUN npm run build

# =============================================================================
# Stage 2: Backend Builder
# =============================================================================
# Using Debian-based image because go-libsql requires glibc (Alpine uses musl)
FROM golang:1.24-bookworm AS backend-builder

# Install build dependencies for CGO
RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential \
    gcc \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /build
COPY backend/go.mod backend/go.sum ./
RUN go mod download
COPY backend/ ./
RUN go mod tidy

ENV CGO_ENABLED=1
RUN go build -ldflags="-s -w" -o /server ./cmd/server

# =============================================================================
# Stage 3: Runtime
# =============================================================================
# Using Debian slim for glibc compatibility with go-libsql
FROM debian:bookworm-slim AS runtime

# Install runtime dependencies
RUN apt-get update && apt-get install -y --no-install-recommends \
    nginx \
    ca-certificates \
    curl \
    bash \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user
RUN groupadd -g 1000 app && \
    useradd -u 1000 -g app -s /bin/sh -m app

# Create required directories with proper ownership
RUN mkdir -p /app/frontend/build /app/backend /data /var/log/nginx /run /run/nginx \
    /var/lib/nginx /var/lib/nginx/body /var/lib/nginx/proxy /var/lib/nginx/fastcgi \
    /var/lib/nginx/uwsgi /var/lib/nginx/scgi && \
    chown -R app:app /app /data /var/log/nginx /run /run/nginx /etc/nginx /var/lib/nginx

# Copy frontend build artifacts
COPY --from=frontend-builder --chown=app:app /build/build /app/frontend/build

# Copy backend binary
COPY --from=backend-builder --chown=app:app /server /app/backend/server

# Copy migration files (needed for database initialization)
COPY --chown=app:app backend/internal/repository/migrations /app/backend/internal/repository/migrations

# Copy nginx configuration
COPY --chown=app:app docker/nginx.conf /etc/nginx/nginx.conf

# Copy entrypoint script
COPY --chown=app:app docker/entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh && \
    chmod +x /app/backend/server

# Expose nginx port (non-privileged)
EXPOSE 3000

# Health check against nginx on port 3000
HEALTHCHECK --interval=30s --timeout=5s --start-period=30s --retries=3 \
    CMD curl -f http://localhost:3000/health || exit 1

WORKDIR /app

# Run as non-root user
USER app

ENTRYPOINT ["/app/entrypoint.sh"]
