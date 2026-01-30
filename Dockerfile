FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go.mod and go.sum
COPY backend/go.mod backend/go.sum ./

# Download dependencies
RUN go mod download

# Copy source
COPY backend/ ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o api ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/api ./

EXPOSE 8080

ENV ENVIRONMENT=production
ENV DB_MODE=memory

CMD ./api
