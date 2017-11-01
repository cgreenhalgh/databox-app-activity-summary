# Databox App: Activity Monitor

A [databox](http://github.com/me-box/databox) app intended to provide
a visual dashboard/report on recent activity.

Status: just starting.

By Chris Greenhalgh

Copyright (c) The University of Nottingham, 2017

## Roadmap:

- read activity data from the 
[strava driver](https://github.com/cgreenhalgh/databox-driver-strava)

- produce a basic day-by-day summary view

- expand to other data sources...

- add automatic alerts for possible early warning signs

## Install

See the latest [databox](http://github.com/me-box/databox)
documentation. perhaps...
```
./databox-component-install cgreenhalgh/databox-app-activity-summary
```

## Develop

```
docker build -t databox-app-activity-summary -f Dockerfile.dev .
```
upload `databox-manifest.json` to [app store](http://127.0.0.1:8181).

Install app.

Find container
```
docker ps | grep strava
```

Copy in the latest code:
```
docker cp . CONTAINERID:/root/go/src/main
```
BUild/run
```
docker exec CONTAINERID:/root/go/src/main
/root/go/bin/dep ensure
GGO_ENABLED=0 GOOS=linux go build -a -tags netgo -installsuffix netgo -ldflags '-d -s -w -extldflags "-static"' -o app src/*.go
./app
```
or if using `ng serve` then
```
./app proxy
```

If you update the imports then don't forget to copy out `Gopkg.lock`/`Gopkg.toml`.

To build the app
```
cd my-app
ng build -bh /databox-app-activity-summary/ui/static/
cp dist/* ../www/
```
or to serve it live (using app proxy)
```
ng serve --host 0.0.0.0 --disable-host-check -bh /databox-app-activity-summary/ui/static/
```
