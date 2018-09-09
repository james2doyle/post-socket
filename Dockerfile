# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace
ADD . /go/src/github.com/james2doyle/post-socket

# bake in some environment variables?
# ENV APP_LOGGING true

# Set the working directory to avoid relative paths after this
WORKDIR /go/src/github.com/james2doyle/post-socket

# Fetch the dependencies
RUN go get .

# build the binary to run later
RUN go build

EXPOSE 8080 8081

# Run the command by default when the container starts
ENTRYPOINT /go/bin/post-socket
