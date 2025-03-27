# Start from the official Go image
FROM golang:1.19-alpine

# Set the working directory in the container
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the source code into the container
COPY . .

# Build the Go application
RUN go build -o main ./cmd/server

# Expose the port the app will run on
EXPOSE 8080

# Run the Go application when the container starts
CMD ["./main"]
