# Stage 1: Build the Go application
FROM golang:latest AS build
WORKDIR /app

# Install swag CLI for generating Swagger docs
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . ./

# Generate Swagger docs and build the application
RUN swag init && CGO_ENABLED=0 go build -o main .

# Verify Swagger docs generation
RUN test -f docs/swagger.json && test -f docs/swagger.yaml

# Stage 2: Setup the application in a smaller container
FROM alpine:latest
WORKDIR /root/

# Copy the built application binary from the build stage
COPY --from=build /app/main .

# Expose the application port
EXPOSE 8080

# Run the application binary
CMD ["./main"]
