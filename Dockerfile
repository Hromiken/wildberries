# Stage 1: Build
FROM golang:alpine AS builder

WORKDIR /app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -tags migrate -o /bin/app ./cmd

# Stage 2: Run
FROM alpine:latest
COPY --from=builder /bin/app /app
COPY --from=builder /app/config /config
COPY --from=builder /app/web /web

CMD ["/app"]