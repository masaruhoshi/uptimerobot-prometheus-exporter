# Grafana Dashboard

Also available at [grafana.com](https://grafana.com/dashboards/6282)


## Docker compose 

You can test the exporter using `docker` and `docker-compose`


#### Running 

Preparing:
```sh
$ mkdir -p ./prometheus/{data,etc} ./grafana
$ cp prometheus.yml ./prometheus/etc
```

Then:

```sh
$ docker-compose up 
```

Login on grafana, add the datasource then import the `Dashboard.json`.
You can check some screenshots [here](https://grafana.com/dashboards/6282).