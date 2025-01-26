# Use Rust slim image for building the project
FROM rust:slim AS builder

# Set the working directory
WORKDIR /app

# Install necessary build tools (for building dependencies)
RUN apt-get update && apt-get install -y build-essential curl vim && apt-get clean

# Copy only Cargo files first for caching dependencies
COPY Cargo.toml ./

# If Cargo.lock is missing, create it after copying Cargo.toml
RUN cargo fetch || (echo "Cargo.lock missing, generating it..." && cargo generate-lockfile)

# Now copy the full source code
COPY . .

# Build the project in release mode
RUN cargo build --release

# Use a distroless image for runtime
FROM gcr.io/distroless/cc

# Set the working directory
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/target/release/rss .

# Expose the application port
EXPOSE 9090

# Run the compiled binary
CMD ["./rss"]
