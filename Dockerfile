# Build stage for frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /frontend

# Copy frontend package files
COPY frontend/package.json frontend/package-lock.json* ./

# Install dependencies
RUN npm install

# Copy frontend source
COPY frontend/ ./

# Build frontend
RUN npm run build

# Build stage for backend
FROM golang:1.21-alpine AS backend-builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go.mod
COPY backend/go.mod ./go.mod

# Download dependencies
RUN go mod download

# Copy source
COPY backend/ ./

# Download and tidy dependencies
RUN go mod tidy

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy backend binary
COPY --from=backend-builder /app/api ./

# Copy frontend build
COPY --from=frontend-builder /frontend/dist ./frontend/dist

EXPOSE 8080

ENV ENVIRONMENT=production
ENV DB_MODE=memory

CMD ./api
