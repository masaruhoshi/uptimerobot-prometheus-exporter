# Uptimerobot Prometheus Exporter
A golang prometheus exporter for UptimeRobot.

![Build Status](https://travis-ci.org/masaruhoshi/uptimerobot-prometheus-exporter.svg?branch=master)

## Building and running
The default make file behavior is to build the binary

```sh
$ make
$ export UPTIMEROBOT_API_KEY="api-key-here"
$ ./uptimerobot_exporter
```

### Docker

```sh
$ make docker
$ docker run --rm -e "UPTIMEROBOT_API_KEY=api-key-here" uptimerobot-exporter
```

## Vendoring
Package vendoring is handled with [Glide](https://github.com/Masterminds/glide).

## Environment variables
The system only depends on one single environment variable:
* `UPTIMEROBOT_API_KEY` The api key provided by UptimeRobot. For more information about
how to obtain this key, access [Uptime Robot API](https://uptimerobot.com/api).

# TODO
* Current exporter is based on [uptimerobot-go](https://github.com/uptimerobot/uptimerobot-go/).
This API is a little outdated and does not yet fully support `v2`. I created a
[fork](https://github.com/masaruhoshi/uptimerobot-go) adding some features. I have to eventually
push my changes to the main repo and use it as dependency package.
* Tests, tests, test. We need to add more tests.

For additional suggestions, *please*, use the [Issues](https://github.com/masaruhoshi/uptimerobot-prometheus-exporter/issues)
to report bugs or enhancements.
