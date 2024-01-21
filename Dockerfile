# Use the official Golang image as the base image
FROM golang:1.21

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files to the working directory
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the local package files to the container's working directory
COPY . .

# Build the application
RUN go build -o /app/main ./cmd

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./main"]
