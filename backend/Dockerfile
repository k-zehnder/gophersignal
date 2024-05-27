# Stage 1: Build the Go application
FROM golang:latest as builder

# Install necessary packages
RUN apt-get update && apt-get install -y \
    default-mysql-client \
    wget

# Install swag CLI for generating Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Set the working directory for the application
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Generate Swagger docs
RUN swag init

# Verify Swagger docs generation
RUN test -f docs/swagger.json && test -f docs/swagger.yaml

# Stage 2: Setup the application in a smaller container
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the application source code and Swagger docs from the build stage
COPY --from=builder /app .

# Ensure init-db.sh is executable
RUN chmod +x /app/bin/init-db.sh

# Expose port 8080
EXPOSE 8080

# Run the initialization script and start the Go application
CMD ["/app/bin/init-db.sh"]
