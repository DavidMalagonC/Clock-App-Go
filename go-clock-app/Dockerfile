# syntax=docker/dockerfile:1

FROM golang:1.22.5
ENV CGO_ENABLED=1

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project directory
COPY . .

# Build
RUN CGO_ENABLED=1 GOOS=linux go build -o /docker-gs-ping

# Expose the port the application runs on
EXPOSE 8080

# Run
CMD ["/docker-gs-ping"]