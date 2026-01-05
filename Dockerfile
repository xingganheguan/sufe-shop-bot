# Build stage
FROM --platform=$BUILDPLATFORM golang:1.22-alpine AS builder

# Build arguments for version info
ARG VERSION=dev
ARG COMMIT=unknown
ARG BUILD_TIME=unknown
ARG TARGETOS
ARG TARGETARCH

# Install build dependencies
RUN apk add --no-cache git make ca-certificates tzdata

# Set working directory
WORKDIR /build

# Copy go mod files first for better layer caching
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Tidy dependencies (in case go.sum is missing)
RUN go mod tidy

# Build the application with version info
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s -X main.Version=${VERSION} -X main.Commit=${COMMIT} -X main.BuildTime=${BUILD_TIME}" \
    -o shopbot ./cmd/server

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 shopbot && \
    adduser -u 1000 -G shopbot -s /bin/sh -D shopbot

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/shopbot /app/shopbot

# Copy static files and templates
COPY --from=builder /build/templates /app/templates
COPY --from=builder /build/static /app/static

# Create directories for logs and data
RUN mkdir -p /app/logs /app/data && \
    chown -R shopbot:shopbot /app

# Switch to non-root user
USER shopbot

# Expose port
EXPOSE 9147

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:9147/healthz || exit 1

# Run the application
CMD ["/app/shopbot"]