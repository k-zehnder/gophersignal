# Start from the latest Golang base image
FROM golang:latest

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum from the backend directory and download dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the backend source code to the Working Directory inside the container
COPY backend/ ./

# Build the Go app
RUN go build -o main .

# Command to run the executable
CMD ["./main"]
