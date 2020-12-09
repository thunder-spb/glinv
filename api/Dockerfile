# Start a golang base image, version 1.14
FROM golang:1.14.0-stretch AS builder

# Disable crosscompiling
ENV GO111MODULE=on \
    CGO_ENABLED=0

WORKDIR /build

# Cache modules retrieval
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code necessary to build the application
COPY . .

# Build the application
RUN go build ./cmd/api

# /dist folder containing just the files necessary for runtime
# it will be copied as the / (root) of the output image
WORKDIR /dist
RUN cp /build/api ./api
# directories/files your app needs during runtime
RUN mkdir /app

# Start with a scratch (no layers)
FROM scratch

COPY --chown=0:0 --from=builder /dist /

# Set up the app to run as a non-root user inside the /app folder
# User ID 65534 is usually user 'nobody'
# The executor of this image should still specify a user during setup
COPY --chown=65534:0 --from=builder /app /app
USER 65534
WORKDIR /app

EXPOSE 10011

ENTRYPOINT ["/api"]