FROM golang:1.26-alpine AS builder

WORKDIR /build

# Copy Go workspace
COPY go.work go.work.sum* ./
COPY packages/core/ packages/core/
COPY packages/scanner/ packages/scanner/
COPY packages/installer/ packages/installer/
COPY packages/skills/ packages/skills/
COPY packages/github/ packages/github/
COPY packages/catalog/ packages/catalog/
COPY packages/registry/ packages/registry/
COPY packages/cli/ packages/cli/

RUN cd packages/cli && \
    go mod download && \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /autodev ./main.go

# ── Runtime image ─────────────────────────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache \
    bash \
    curl \
    wget \
    git \
    ca-certificates \
    && rm -rf /var/cache/apk/*

COPY --from=builder /autodev /usr/local/bin/autodev

RUN chmod +x /usr/local/bin/autodev

WORKDIR /workspace

ENTRYPOINT ["autodev"]
CMD ["--help"]

LABEL org.opencontainers.image.title="AutoDev"
LABEL org.opencontainers.image.description="The App Store for Developers — Clone. Scan. Install. Build."
LABEL org.opencontainers.image.url="https://github.com/HEETMEHTA18/autodev"
LABEL org.opencontainers.image.source="https://github.com/HEETMEHTA18/autodev"
LABEL org.opencontainers.image.licenses="MIT"
