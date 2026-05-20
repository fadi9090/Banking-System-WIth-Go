# Stage 1: Build (temporary, large)
FROM golang:1.25-alpine AS builder
WORKDIR /app
RUN apk add --no-cache git
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/app

# Stage 2
FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata
RUN adduser -D -g '' appuser
WORKDIR /app
COPY --from=builder /app/main .
RUN chown -R appuser:appuser /app
USER appuser
EXPOSE 8080
CMD ["./main"]