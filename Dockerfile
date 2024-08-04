FROM golang:alpine AS backend

ARG VERSION

ENV GOFLAGS="-mod=vendor"
ENV CGO_ENABLED=0

ADD . /build
WORKDIR /build

RUN apk add --no-cache --update git tzdata ca-certificates

RUN \
    version=${VERSION}-$(date +%Y%m%dT%H:%M:%S) && \
    echo "version=$version" && \
    cd app && go build -o /build/maptp -ldflags "-X 'main.Version=${version}'"

FROM scratch

COPY --from=backend /build/maptp /srv/maptp
COPY --from=backend /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=backend /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

WORKDIR /srv
ENTRYPOINT ["/srv/maptp"]