# Use official Go image
FROM golang:1.24

WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the Go app
RUN go build -o server .

EXPOSE 8080

CMD ["./server"]
