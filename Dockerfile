FROM golang:alpine as builder

COPY . /go/src/github.com/masaruhoshi/uptimerobot-prometheus-exporter
RUN apk --update add curl git make perl && \
    curl -s https://glide.sh/get | sh && \
    cd /go/src/github.com/masaruhoshi/uptimerobot-prometheus-exporter && \
    make build

FROM alpine:latest
LABEL Authors="Masaru Hoshi <https://github.com/masaruhoshi>, Felipe Santiago <https://github.com/felipesantiago>"

EXPOSE 9429

RUN apk --update add ca-certificates && \
    rm -rf /var/cache/apk/*

COPY --from=builder /go/src/github.com/masaruhoshi/uptimerobot-prometheus-exporter/uptimerobot_exporter /usr/bin/uptimerobot_exporter

ENTRYPOINT [ "/usr/bin/uptimerobot_exporter" ]