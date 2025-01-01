FROM golang:1.23 AS builder

RUN apt-get update -y && apt-get upgrade -y && apt-get install -y make build-essential
# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy local code to the container image.
COPY . .

# Build the binary.
RUN mkdir -p ./tmp/ && make

# Use the official Debian slim image for a production container.
# https://hub.docker.com/_/debian
FROM alpine:3.20.0 AS production
# Create a folder to store the unirest binary
RUN mkdir -p /unirest

RUN apk add --no-cache ca-certificates

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/tmp/main.exe /usr/local/bin/unirest

# Copy configuration
COPY ./unirest-sample.conf /unirest/.env

# Set workdir
WORKDIR /unirest

# Set entrypoint
ENTRYPOINT ["/usr/local/bin/unirest"]	
