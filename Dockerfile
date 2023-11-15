# Base image for Go
FROM golang:1.21

# Create a directory
RUN mkdir /app

# Copy files into the app directory
ADD . /app/

# Set the working directory
WORKDIR /app

# Get dependencies
RUN go get -d

# Build the application
RUN go build -o main .

# Run the bot
CMD ["/app/main"]
