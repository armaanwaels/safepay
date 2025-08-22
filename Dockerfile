# ---- Build stage ----
FROM golang:1.22-alpine AS builder
WORKDIR /app
RUN apk add --no-cache build-base sqlite-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=1
RUN go build -ldflags="-s -w" -o /app/safepay .

# ---- Run stage ----
FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates sqlite-libs
COPY --from=builder /app/safepay /app/safepay
COPY db/ /app/db/
ENV DB_PATH=/app/safepay.db
EXPOSE 8080
CMD ["/app/safepay"]
