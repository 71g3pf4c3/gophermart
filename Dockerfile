# Step 1: Modules caching
FROM golang:1.25-alpine AS modules

COPY go.mod go.sum /modules/

WORKDIR /modules

RUN go mod download

# Step 2: Builder
FROM golang:1.25-alpine AS builder

COPY --from=modules /go/pkg /go/pkg
COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
	go build -o /bin/app ./cmd/gophermart

# Step 3: Final
FROM scratch

COPY --from=builder /app/migrations /migrations
COPY --from=builder /bin/app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV MIGRATIONS_PATH=/migrations
CMD ["/app"]
