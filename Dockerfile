# --- Build Stage ---
# Use an official Go image as a parent image for the build environment.
# 'AS builder' names this stage, so we can refer to it later.
# alpine for a smaller base image
FROM golang:1.24-alpine AS builder

# Set ARG for build flags, allowing them to be overridden if needed.
ARG TARGETOS=linux
# Or arm64, etc., depending on your deployment target
ARG TARGETARCH=amd64 

# Set the current working dir
# All subsequent commands (COPY, RUN) 
# will be executed to /app
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies first.
# This leverages Docker's layer caching. If these files don't change,
# the dependency download step won't re-run on subsequent builds.
COPY go.mod go.sum ./

RUN go mod download
RUN go mod verify

# Install air into the container
RUN go install github.com/air-verse/air@latest

# Copy the rest of the app source
# into the container
COPY . .

# Build the Go application as a static binary.
# CGO_ENABLED=0 disables Cgo, leading to a more portable static binary.
# -ldflags="-w -s" strips debug information and symbol table, reducing binary size.
# Output the binary to /app/goseinaka-server in the builder stage.
RUN CGO_ENABLED=0 GOOS=${TARGETOS} GOARCH=${TARGETARCH} go build \
    -ldflags="-w -s" \
    -o /app/goseinaka-server \
    ./cmd/goseinaka-server/main.go

# --- Final Stage ---
# Use a minimal base image for the final application.
FROM alpine:latest
# This stage will contain the compiled application and 'air'.

# Create a non-root user and group for security.
# RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the Current Working Directory.
WORKDIR /app

# Copy only the compiled application binary from the builder stage.
COPY --from=builder /app/goseinaka-server /app/goseinaka-server

# If your application needs to read configuration files from a 'configs' directory
# that are NOT built into the binary (e.g., static assets, templates, though GoSeinaka is API only),
# you would copy them here. For GoSeinaka, this is likely not needed if all config is via env vars.
# COPY --from=builder /app/configs ./configs

# Ensure the binary is executable by the appuser.
# RUN chown appuser:appgroup /app/goseinaka-server

# Switch to the non-root user.
# USER appuser

# Expose the port the application will listen on.
# This should match the port your application is configured to use via environment variables.
# Adjust if your typical production API_PORT is different.
EXPOSE 8080 

# Command to run the application.
# This directly executes your compiled Go binary.
CMD ["/app/goseinaka-server"]
