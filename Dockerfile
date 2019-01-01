# Donot use this Dockerfile.This is not ready yet.

# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang

# Copy the local package files to the container's workspace.
ADD . /go/src/statusok

# Build the outyet command inside the container.
# (You may fetch or manage dependencies here,
# either manually or with a tool like "godep".)
RUN go get github.com/codegangsta/cli
RUN go get github.com/influxdata/influxdb
RUN go get github.com/influxdata/platform
RUN go get github.com/mailgun/mailgun-go
RUN go get github.com/Sirupsen/logrus
RUN go install /go/src/statusok

RUN wget http://influxdb.s3.amazonaws.com/influxdb_0.9.3_amd64.deb
RUN dpkg -i influxdb_0.9.3_amd64.deb
RUN /etc/init.d/influxdb start

RUN wget https://grafanarel.s3.amazonaws.com/builds/grafana_2.1.3_amd64.deb
RUN apt-get update
RUN apt-get install -y adduser libfontconfig
RUN dpkg -i grafana_2.1.3_amd64.deb
RUN service grafana-server start

#how to connect to localhost inside ?? http://stackoverflow.com/questions/24319662/from-inside-of-a-docker-container-how-do-i-connect-to-the-localhost-of-the-mach

ENTRYPOINT /go/bin/statusok --config /go/src/statusok/config.json

# Document that the service listens
EXPOSE 80 8083 8086 7321 3000
