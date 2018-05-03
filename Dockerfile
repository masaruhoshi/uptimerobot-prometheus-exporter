FROM golang:alpine as builder

COPY . /go/src/github.com/masaruhoshi/uptimerobot-prometheus-exporter

RUN apk --update add curl git alpine-sdk && \
    cd /go/src/github.com/masaruhoshi/uptimerobot-prometheus-exporter && \
    go get && make

FROM alpine:latest
LABEL maintainers="Masaru Hoshi <https://github.com/masaruhoshi>, Felipe Santiago <https://github.com/felipesantiago>"

EXPOSE 9429

RUN apk --update add ca-certificates && \
    rm -rf /var/cache/apk/*

COPY --from=builder /go/src/github.com/masaruhoshi/uptimerobot-prometheus-exporter/uptimerobot-exporter /usr/bin/uptimerobot-exporter

ENTRYPOINT [ "/usr/bin/uptimerobot-exporter" ]
