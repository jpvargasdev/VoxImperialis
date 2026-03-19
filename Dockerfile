# Stage 1: Build the Go application
FROM golang:1.25-alpine AS builder

LABEL maintainer="Juan Vargas <vargasm.jp@gmail.com>"

WORKDIR /app

# Enable multi-arch builds
ARG TARGETARCH
ENV GOARCH=$TARGETARCH

# Copy dependency files and download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go binary for the correct architecture
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$GOARCH go build -o vox-imperialis .

# Stage 2: Create a lightweight runtime image
FROM alpine:latest

WORKDIR /app

# Copy only the binary from the builder stage
COPY --from=builder /app/vox-imperialis .

# vox-imperialis reads config from env vars — mount a .env file at runtime
# or pass variables via docker run -e / docker-compose environment

ENTRYPOINT ["./vox-imperialis"]
