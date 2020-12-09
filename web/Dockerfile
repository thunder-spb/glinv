# Start from the latest golang base image
FROM golang:latest

# Maintainer Info
LABEL maintainer="<p.a.mostovoy@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go Web App
RUN go build -o web ./cmd/web/

# Expose port 8080 to the outside
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["./web"]