# ── Stage 1: Builder ────────────────────────────────────────────────────────
FROM --platform=$BUILDPLATFORM golang:1.24-alpine3.21 AS builder

ARG TARGETARCH

WORKDIR /app

# Download dependencies before copying source (better layer caching)
COPY go.mod ./
RUN go mod download

# Copy source and build a fully static binary
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build \
    -ldflags="-w -s" \
    -o /app/server ./main.go

# ── Stage 2: Final image ─────────────────────────────────────────────────────
# distroless/static: no shell, no libc, no package manager — minimal attack surface.
# :nonroot variant pre-configures uid=65532 and receives daily CVE patches from Google.
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/server /server

USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/server"]
