# Stage 1: Build the Go binary
FROM golang:1.22-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application. CGO_ENABLED=0 is important for a static binary.
# -ldflags="-w -s" strips debug information, reducing the binary size.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o ggtask-api .

# Stage 2: Create the final, minimal image
FROM alpine:latest

# Set the working directory
WORKDIR /app

# Copy only the compiled binary from the builder stage
COPY --from=builder /app/ggtask-api .

# Expose the port the API server will run on
EXPOSE 8080

# The command to run the application
CMD ["./ggtask-api"]
