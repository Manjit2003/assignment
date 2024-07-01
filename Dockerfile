# Build Stage
FROM golang:1.22 AS build

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go app and output to /app/bin/api (statically linked)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/bin/api ./cmd/api/main.go

# Verify the build output
RUN ls -la /app/bin

# Final Stage
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /root/

# Copy the built binary from the build stage
COPY --from=build /app/bin/api /root/bin/api

# Ensure the binary has executable permissions
RUN chmod +x /root/bin/api

# Copy the configuration file
COPY --from=build /app/config.yaml /root/config.yaml

# Expose the port on which the app will run
EXPOSE 8080

# Command to run the executable
CMD ["/root/bin/api"]
