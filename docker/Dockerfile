FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.work go.mod ./
COPY apps/api/go.mod ./apps/api/
RUN go mod download

COPY apps/api/ ./apps/api/
COPY packages/shared/ ./packages/shared/

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /orbit-server ./apps/api/cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates curl

WORKDIR /root/

COPY --from=builder /orbit-server .
COPY apps/api/.env.example ./.env.example

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=5s --retries=3 \
  CMD curl -f http://localhost:8080/health/ready || exit 1

CMD ["./orbit-server"]
