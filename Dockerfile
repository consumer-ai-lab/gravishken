# Use a multi-stage build for efficiency
FROM golang:1.23.3-bullseye AS builder

# Install required system dependencies
RUN apt-get update && apt-get install -y \
    build-essential \
    curl \
    git \
    unzip \
    && rm -rf /var/lib/apt/lists/*
# Install bun for web builds
RUN curl -fsSL https://bun.sh/install | bash

# Set working directory
WORKDIR /app

# Copy the entire project
COPY . .

# Set environment variables
ENV BUILD_MODE=PROD
ENV PATH="/root/.bun/bin:${PATH}"

# Run setup and build
RUN chmod +x ./run.sh && \
    ./run.sh setup && \
    ./run.sh build-server

# Use a smaller base image for the final stage
FROM debian:bullseye-slim

# Install runtime dependencies if needed
RUN apt-get update && apt-get install -y \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

# Copy only the necessary files from builder
COPY --from=builder /app/build/server /app/server
# COPY --from=builder /app/backend/dist /app/dist

# Set default environment variables
ENV BUILD_MODE=PROD
ENV SERVER_PORT=6201
ENV SERVER_SECURE="true"

# Expose the server port
EXPOSE 6201

# Run the server
CMD ["./server"]