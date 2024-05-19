# Stage 1: Build the Go application
FROM golang:latest AS build
WORKDIR /app

# Install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy the Go module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code and build the application
COPY . ./
RUN CGO_ENABLED=0 go build -o main .

# Generate Swagger docs
RUN swag init

# Stage 2: Setup a smaller base image
FROM alpine:latest
WORKDIR /root/

# Copy the built application and Swagger docs from the build stage
COPY --from=build /app/main .
COPY --from=build /app/docs ./docs

# Expose the application port
EXPOSE 8080

# Define the entrypoint for the container
CMD ["./main"]
