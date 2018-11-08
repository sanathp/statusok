FROM golang:1.11

# Copy the local package files to the container's workspace.
ADD . /go/src/github.com/sanathp/statusok
WORKDIR /go/src/github.com/sanathp/statusok
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o statusok .

# Production image
FROM alpine:3.6

RUN apk add --no-cache tzdata ca-certificates gettext
COPY --from=0 /go/src/github.com/sanathp/statusok/statusok /statusok
COPY --from=0 /go/src/github.com/sanathp/statusok/config.template /config.template
COPY --from=0 /go/src/github.com/sanathp/statusok/docker-entrypoint.sh /docker-entrypoint.sh

ENTRYPOINT /docker-entrypoint.sh
