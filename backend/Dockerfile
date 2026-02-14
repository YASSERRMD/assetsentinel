# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./
RUN CGO_ENABLED=1 GOOS=linux go build -o /assetsentinel ./cmd/server

# Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates sqlite

WORKDIR /app

COPY --from=builder /assetsentinel .

EXPOSE 8080

CMD ["./assetsentinel"]
