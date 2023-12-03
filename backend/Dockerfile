# Use the official Golang image to create a build artifact.
FROM golang:1.21.0 as builder

# Set the working directory outside $GOPATH to enable the support for modules.
WORKDIR /app

# Copy the go.mod and go.sum to download all dependencies.
# Adjust the paths to point to the src directory.
COPY src/go.mod ./
COPY src/go.sum ./

# Download all dependencies.
RUN go mod download

# Copy the source code into the container.
# Again, adjust the path to copy from the src directory.
COPY src/ .

# Build the application to /gopher-api.
RUN CGO_ENABLED=0 GOOS=linux go build -v -o /gopher-api .

# Use a Docker multi-stage build to create a lean production image.
FROM alpine:latest

# Copy the binary to the production image from the builder stage.
COPY --from=builder /gopher-api /gopher-api

# Run the web service on container startup.
CMD ["/gopher-api"]
