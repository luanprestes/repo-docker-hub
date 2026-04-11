# ── Stage 1: Builder ─────────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -o /app/webhook ./main.go

# ── Stage 2: Final image ──────────────────────────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache docker-cli

COPY --from=builder /app/webhook /webhook

EXPOSE 9000

ENTRYPOINT ["/webhook"]
