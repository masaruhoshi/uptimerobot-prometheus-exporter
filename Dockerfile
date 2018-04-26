FROM       golang:alpine as builder

RUN apk --no-cache add curl git make perl
RUN curl -s https://glide.sh/get | sh
COPY . /go/src/github.com/masaruhoshi/uptimerobot-prometheus-exporter
RUN cd /go/src/github.com/masaruhoshi/uptimerobot-prometheus-exporter && make release

FROM       alpine:3.4
MAINTAINER Masaru Hoshi <hoshi@hoshi.com.br>
EXPOSE     9001

RUN apk add --update ca-certificates
COPY --from=builder /go/src/github.com/masaruhoshi/uptimerobot-prometheus-exporter/release/uptimerobot_exporter-linux-amd64 /usr/local/bin/uptimerobot_exporter

ENTRYPOINT [ "uptimerobot_exporter" ]
