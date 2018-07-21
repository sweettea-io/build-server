# Use Go image for building of binary executable.
FROM golang:1.10 AS builder

# Install dep so that dependencies can be installed.
RUN apt-get update && apt-get install -y unzip --no-install-recommends && \
    apt-get autoremove -y && apt-get clean -y && \
 	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh

# Create classic GOPATH structure.
RUN mkdir -p /go/src/github.com/sweettea-io/build-server

# Switch to project dir as new working dir.
WORKDIR /go/src/github.com/sweettea-io/build-server

# Copy files needed by dep in order to install dependencies.
COPY Gopkg.toml Gopkg.lock ./

# Install dependencies.
RUN dep ensure -vendor-only

# Copy this project into working dir.
COPY . .

# Build Go binary.
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -a -o main .

# Switch over to Docker-in-Docker image.
FROM docker:dind

# Set working dir to /root inside Docker-in-Docker image.
WORKDIR /root

# Copy Go binary built in first image over to this image.
COPY --from=builder /go/src/github.com/sweettea-io/build-server/main ./main

# Execute Go binary.
ENTRYPOINT ["./main"]