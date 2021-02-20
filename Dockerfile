# https://cloud.google.com/run/docs/quickstarts/build-and-deploy?authuser=1#shell
# Use the offical golang image to create a binary.
# This is based on Debian and sets the GOPATH to /go.
# https://hub.docker.com/_/golang
FROM golang:1.15-buster as builder

# Create and change to the app directory.
WORKDIR /app

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY run.go ./

# Build the binary.
RUN go build -mod=readonly -v -o server

# https://github.com/GoogleCloudPlatform/cloud-sdk-docker
FROM gcr.io/google.com/cloudsdktool/cloud-sdk:328.0.0
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/server /app/server
COPY *.sh ./

# Run the web service on container startup.
CMD ["/app/server"]
