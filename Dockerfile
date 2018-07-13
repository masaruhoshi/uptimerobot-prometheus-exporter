FROM golang:1.9-alpine3.7 AS build

RUN apk --no-cache add git ca-certificates make bash && \
    wget -O /go/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64 && \
    chmod +x /go/bin/dep

RUN mkdir -p /go/src/github.com/masaruhoshi/uptimerobot-prometheus-exporter
WORKDIR /go/src/github.com/masaruhoshi/uptimerobot-prometheus-exporter
COPY . /go/src/github.com/masaruhoshi/uptimerobot-prometheus-exporter

RUN make build

FROM alpine:3.7

EXPOSE 9429

COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=build /go/src/github.com/masaruhoshi/uptimerobot-prometheus-exporter/bin/uptimerobot-exporter /uptimerobot-exporter

ENTRYPOINT [ "/uptimerobot-exporter" ]

RUN apk --update add ca-certificates && \
    rm -rf /var/cache/apk/*

ARG BUILD_DATE
ARG VCS_REF
ARG VERSION

LABEL maintainers="Masaru Hoshi <https://github.com/masaruhoshi>, Felipe Santiago <https://github.com/felipesantiago>" \
      org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.vcs-ref=$VCS_REF \
      org.label-schema.version=$VERSION \
      org.label-schema.vcs-url="https://github.com/masaruhoshi/uptimerobot-prometheus-exporter"
