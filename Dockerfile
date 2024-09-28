# Start with the official Golang image
FROM golang:1.22.4

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files first, then the source code
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o bin/personal-backend

# Expose the port the app runs on
EXPOSE 8080

# Command to run the executable
CMD ["./bin/personal-backend"]
