# --- Stage 1: Build Go API ---
FROM golang:alpine AS go-builder
WORKDIR /app/api
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api/ .
RUN CGO_ENABLED=0 GOOS=linux go build -o lakoo-api ./cmd/api

# --- Stage 2: Final Image (Python + Go Binary) ---
FROM python:3.11-slim
WORKDIR /app

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    default-libmysqlclient-dev \
    pkg-config \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Copy Go Binary from Stage 1
COPY --from=go-builder /app/api/lakoo-api /app/lakoo-api
# Copy migrations folder
COPY api/migrations /app/migrations

# Setup Python AI Service
COPY ai-service/requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY ai-service/ .

# Create a startup script to run both Go API and Python AI
RUN echo '#!/bin/bash\n\
# Start AI service in background\n\
uvicorn main:app --host 0.0.0.0 --port 8000 & \n\
# Start Go API in foreground (main process)\n\
./lakoo-api\n\
' > start.sh && chmod +x start.sh

# Expose both ports (API: 8080, AI: 8000)
EXPOSE 8080 8000

# Start both services
CMD ["./start.sh"]
