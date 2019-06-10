FROM golang:1.6.3

ENV STATUSOK_VERSION 0.1.1

RUN apt-get update \
    && apt-get install -y unzip \
    && wget https://github.com/sanathp/statusok/releases/download/$STATUSOK_VERSION/statusok_linux.zip \
    && unzip statusok_linux.zip \
    && mv ./statusok_linux/statusok /go/bin/StatusOk \
    && rm -rf ./statusok_linux* \
    && apt-get remove -y unzip git \
    && apt-get autoremove -y \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

VOLUME /config
COPY ./docker-entrypoint.sh /docker-entrypoint.sh
ENTRYPOINT /docker-entrypoint.sh
