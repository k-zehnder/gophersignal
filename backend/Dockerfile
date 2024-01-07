# Stage 1: Build the Go application
FROM golang:latest AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 go build -o main .

# Stage 2: Setup a smaller base image
FROM alpine:latest  
WORKDIR /root/
COPY --from=build /app/main .
COPY .env . 
EXPOSE 8080
CMD ["./main"]
