# Build stage
FROM golang:1.23.5-bookworm as builder

ARG tmp=/tmp/app
WORKDIR ${tmp}

# Copy only the go.mod and go.sum files to leverage Docker cache
COPY go.mod .
COPY go.sum .

# Download dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the binary.
RUN go build -v -o quick-booking-api

# Runtime stage
FROM debian:bookworm-slim

ARG tmp=/tmp/app
ENV root=/app

WORKDIR ${root}

# Copy only the binary from the builder stage
COPY --from=builder ${tmp}/ ${root}/

CMD ["/app/quick-booking-api run"]
