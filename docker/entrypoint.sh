#!/bin/bash
# =============================================================================
# entrypoint.sh - Budget Tracker Container Entrypoint
# =============================================================================
# This script manages two processes in the container:
#   1. Go backend API server (port 8080 internal)
#   2. nginx reverse proxy (port 3000 external)
#
# Startup sequence:
#   1. Create data directory for SQLite database
#   2. Start Go backend in background
#   3. Wait for backend to pass health check
#   4. Start nginx in foreground
#   5. Monitor both processes, exit if either dies
#
# Shutdown sequence (on SIGTERM/SIGINT):
#   1. Stop nginx (stop accepting new connections)
#   2. Wait for backend to finish in-flight requests
#   3. Exit with appropriate code
# =============================================================================

set -e

# =============================================================================
# Configuration
# =============================================================================
DATA_DIR="/data"
BACKEND_PORT="8080"
NGINX_PORT="3000"
HEALTH_CHECK_URL="http://127.0.0.1:${BACKEND_PORT}/health"
MAX_HEALTH_ATTEMPTS=30
BACKEND_SHUTDOWN_TIMEOUT=15

# Process IDs
BACKEND_PID=""
NGINX_PID=""

# =============================================================================
# Logging Functions
# =============================================================================
log_info() { echo "[INFO] $(date '+%Y-%m-%d %H:%M:%S') $1"; }
log_error() { echo "[ERROR] $(date '+%Y-%m-%d %H:%M:%S') $1" >&2; }

# =============================================================================
# Cleanup Function
# =============================================================================
# Called on SIGTERM/SIGINT for graceful shutdown
# Does NOT call exit - lets caller control exit code
cleanup() {
  log_info "Received shutdown signal, initiating graceful shutdown..."

  # Stop nginx first (stop accepting new connections)
  if [ -n "$NGINX_PID" ] && kill -0 "$NGINX_PID" 2>/dev/null; then
    log_info "Stopping nginx..."
    nginx -s quit 2>/dev/null || true
  fi

  # Wait for nginx to finish
  sleep 2

  # Stop backend with graceful shutdown timeout
  if [ -n "$BACKEND_PID" ] && kill -0 "$BACKEND_PID" 2>/dev/null; then
    log_info "Stopping backend (PID: $BACKEND_PID)..."
    kill -TERM "$BACKEND_PID" 2>/dev/null || true

    # Wait for backend to exit gracefully
    wait_count=0
    while kill -0 "$BACKEND_PID" 2>/dev/null && [ $wait_count -lt $BACKEND_SHUTDOWN_TIMEOUT ]; do
      sleep 1
      wait_count=$((wait_count + 1))
    done

    # Force kill if still running
    if kill -0 "$BACKEND_PID" 2>/dev/null; then
      log_info "Backend did not exit gracefully, forcing..."
      kill -9 "$BACKEND_PID" 2>/dev/null || true
    fi
  fi

  log_info "Shutdown complete"
}

# Trap signals for graceful shutdown
trap 'cleanup; exit 0' TERM INT

# =============================================================================
# Initialization
# =============================================================================
log_info "Starting Budget Tracker container..."

# Create data directory for SQLite database
log_info "Ensuring data directory exists: ${DATA_DIR}"
mkdir -p "${DATA_DIR}"

# Set environment variables for backend
export TURSO_LOCAL_PATH="${DATA_DIR}/budget.db"
export PORT="${BACKEND_PORT}"

# Log configuration (without sensitive values)
log_info "Configuration:"
log_info "  TURSO_MODE: ${TURSO_MODE:-local}"
log_info "  TURSO_LOCAL_PATH: ${TURSO_LOCAL_PATH}"
log_info "  Backend port: ${BACKEND_PORT}"
log_info "  Nginx port: ${NGINX_PORT}"

# =============================================================================
# Start Backend
# =============================================================================
log_info "Starting backend server..."

cd /app/backend
./server &
BACKEND_PID=$!

log_info "Backend started (PID: ${BACKEND_PID})"

# =============================================================================
# Health Check Loop
# =============================================================================
log_info "Waiting for backend to become healthy..."

attempt=1
while [ $attempt -le $MAX_HEALTH_ATTEMPTS ]; do
  if curl -sf "${HEALTH_CHECK_URL}" >/dev/null 2>&1; then
    log_info "Backend is healthy!"
    break
  fi

  # Check if backend process is still running
  if ! kill -0 "$BACKEND_PID" 2>/dev/null; then
    log_error "Backend process exited unexpectedly"
    exit 1
  fi

  if [ $attempt -eq $MAX_HEALTH_ATTEMPTS ]; then
    log_error "Backend failed to become healthy after ${MAX_HEALTH_ATTEMPTS} attempts"
    cleanup
    exit 1
  fi

  log_info "Waiting for backend... (attempt ${attempt}/${MAX_HEALTH_ATTEMPTS})"
  sleep 1
  attempt=$((attempt + 1))
done

# =============================================================================
# Start Nginx
# =============================================================================
log_info "Starting nginx on port ${NGINX_PORT}..."

nginx -g "daemon off;" &
NGINX_PID=$!

log_info "Nginx started (PID: ${NGINX_PID})"

log_info "================================================"
log_info "Budget Tracker is ready!"
log_info "  Application: http://localhost:${NGINX_PORT}"
log_info "  API: http://localhost:${NGINX_PORT}/api/"
log_info "  Health: http://localhost:${NGINX_PORT}/health"
log_info "================================================"

# =============================================================================
# Process Monitor Loop
# =============================================================================
# Monitor both processes and exit if either dies
while true; do
  # Check if backend is still running
  if ! kill -0 "$BACKEND_PID" 2>/dev/null; then
    log_error "Backend process (PID: ${BACKEND_PID}) exited unexpectedly"
    cleanup
    exit 1
  fi

  # Check if nginx is still running
  if ! kill -0 "$NGINX_PID" 2>/dev/null; then
    log_error "Nginx process (PID: ${NGINX_PID}) exited unexpectedly"
    cleanup
    exit 1
  fi

  # Sleep before next check
  sleep 5
done
