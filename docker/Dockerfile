# Build Stage
FROM golang:1.21 AS build

# Set the working directory
WORKDIR /app

# Copy all files to the working directory
COPY . .

# Download all dependencies
RUN go mod download

# Set CGO_ENABLED to 1
ENV CGO_ENABLED=1

# Build the Go application
RUN go build -o /app/app ./main.go

# Runtime Stage
FROM busybox AS runtime

# Copy the binary from the build stage to the runtime stage
COPY --from=build /app/app /app

# Set the entry point to execute the binary directly
ENTRYPOINT ["/app"]