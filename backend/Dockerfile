# Stage 1: Build the Go application
FROM golang:latest AS build
WORKDIR /app

# Install swag CLI
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o main .

# Generate Swagger docs
RUN swag init

# Stage 2: Setup the application in a smaller container
FROM alpine:latest
WORKDIR /root/
COPY --from=build /app/main .
COPY --from=build /app/docs ./docs
EXPOSE 8080
CMD ["./main"]
