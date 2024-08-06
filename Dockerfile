# Start from a base Go Docker image with the matching Go version
FROM golang:1.22.4-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the entire project into the container
COPY . .

# Build the Go application
RUN go build -o main .

# Command to run the executable
CMD ["./main"]
