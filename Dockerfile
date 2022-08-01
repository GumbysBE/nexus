FROM golang:1.15-alpine3.14

WORKDIR /opt/nexus

COPY ./ /opt/nexus

RUN set -ex \
    && apk add --no-cache --virtual .build-deps \
        make gcc musl-dev\
    && make service \
    && make service install \
    && mkdir -p /run/nexus/ \
    && apk del --no-network .build-deps

CMD ["/go/bin/nexusd", "-unix", "/run/nexus/socket", "-ws", "0.0.0.0:80", "-realm", "mbl-service-realm"]
