# =========================
# Stage 1: Build Frontend
# =========================
FROM node:20 AS frontend-build

# Install pnpm
RUN corepack enable && corepack prepare pnpm@latest --activate

# Set working directory
WORKDIR /app

# Copy package files and install dependencies
# COPY app/package.json app/package-lock.json ./
# RUN npm ci --network-timeout=100000 --include=optional
COPY app/package.json app/pnpm-lock.yaml ./
RUN pnpm install --frozen-lockfile

# Copy the rest of the frontend code
COPY app/ .

# Build the frontend using SvelteKit
# RUN npm run build
RUN pnpm run build

# =========================
# Stage 2: Build Backend
# =========================
FROM golang:1.22-alpine AS backend-build

# Set working directory
WORKDIR /server

# Install git (required for some Go modules)
RUN apk add --no-cache git

# Copy Go module files and download dependencies
COPY server/go.mod server/go.sum ./
RUN go mod download

# Copy the rest of the backend code
COPY server/ .

# Copy the built frontend from the previous stage
# Adjust 'build' to your actual SvelteKit output directory if different
COPY --from=frontend-build /app/build ./app/

# Build the Go server
RUN go build -o main .

# =========================
# Stage 3: Final Image
# =========================
FROM alpine:latest

# Set working directory
WORKDIR /server

# Install necessary CA certificates
RUN apk add --no-cache ca-certificates

# Copy the built Go server and frontend from the previous stage
COPY --from=backend-build /server/main .
COPY --from=backend-build /server/app ./app/

# Expose the server port
EXPOSE 8080

# Set environment variables if needed
ENV PORT=8080

# Command to run the Go server
CMD ["./main"]
