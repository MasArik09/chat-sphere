FROM golang:alpine AS builder

WORKDIR /app

# Copy dependency files and download modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app statically
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api cmd/api/main.go

# Run stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates and tzdata for secure connections and timezones
RUN apk --no-cache add ca-certificates tzdata

# Copy built binary from builder
COPY --from=builder /app/api .

# Expose Go API port
EXPOSE 8080

# Run the binary
CMD ["./api"]
