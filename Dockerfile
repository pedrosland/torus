FROM golang:latest
MAINTAINER Barak Michener <barak.michener@coreos.com>

# Set up workdir
WORKDIR /go/src/github.com/coreos/torus

# Add and install torus
ADD . .
RUN go get -d ./...
RUN go install -v github.com/coreos/torus/cmd/torus
RUN go install -v github.com/coreos/torus/cmd/torusctl
RUN go install -v github.com/coreos/torus/cmd/torusblock

# Expose the port and volume for configuration and data persistence.
VOLUME ["/data", "/plugin"]
EXPOSE 40000 4321

CMD ["./entrypoint.sh"]
