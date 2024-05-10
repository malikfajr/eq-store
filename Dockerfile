# Builder stage
# Use an official Golang runtime as a parent image
FROM golang:1.22.1-alpine AS builder

# Set the working directory to /app
WORKDIR /app

# Copy the current directory contents into the container
COPY . .

# Download and cache dependencies
RUN go get

# Build
RUN go build -o main main.go



# Run stage
FROM alpine
WORKDIR /app

COPY --from=builder /app/main .

# Expose port 8080 for the container
EXPOSE 8080

# Set environment variables at runtime
ENV DB_HOST=localhost
ENV DB_PORT=5432
ENV DB_USER=postgres
ENV DB_PASSWORD=secret
ENV DB_NAME=mydb
ENV DB_PARAMS=sslmode=disable
ENV JWT_SECRET=secretly
ENV BCRYPT_SALT=8

# Set the entry point of the container to the Go app executable
CMD ["/app/main"]