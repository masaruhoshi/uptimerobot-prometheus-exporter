FROM alpine:latest
LABEL maintainers="Masaru Hoshi <https://github.com/masaruhoshi>, Felipe Santiago <https://github.com/felipesantiago>"

EXPOSE 9429

RUN apk --update add ca-certificates && \
    rm -rf /var/cache/apk/*

COPY ./uptimerobot-exporter /usr/bin/uptimerobot-exporter

ENTRYPOINT [ "/usr/bin/uptimerobot-exporter" ]
