# start with node as it is probably bigger
FROM node:8.8.1-alpine as gobuild
# FROM golang:1.8.3-alpine3.6 as gobuild

# copied from https://github.com/docker-library/golang/blob/cffcff7fce7f6b6b5c82fc8f7b3331a10590a661/1.8/alpine3.6/Dockerfile
RUN apk add --no-cache ca-certificates

ENV GOLANG_VERSION 1.8.5

# https://golang.org/issue/14851 (Go 1.8 & 1.7)
# https://golang.org/issue/17847 (Go 1.7)
COPY go/*.patch /go-alpine-patches/

RUN set -eux; \
	apk add --no-cache --virtual .build-deps \
		bash \
		gcc \
		musl-dev \
		openssl \
		go \
	; \
	export \
# set GOROOT_BOOTSTRAP such that we can actually build Go
		GOROOT_BOOTSTRAP="$(go env GOROOT)" \
# ... and set "cross-building" related vars to the installed system's values so that we create a build targeting the proper arch
# (for example, if our build host is GOARCH=amd64, but our build env/image is GOARCH=386, our build needs GOARCH=386)
		GOOS="$(go env GOOS)" \
		GOARCH="$(go env GOARCH)" \
		GO386="$(go env GO386)" \
		GOARM="$(go env GOARM)" \
		GOHOSTOS="$(go env GOHOSTOS)" \
		GOHOSTARCH="$(go env GOHOSTARCH)" \
	; \
	\
	wget -O go.tgz "https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz"; \
	echo '4949fd1a5a4954eb54dd208f2f412e720e23f32c91203116bed0387cf5d0ff2d *go.tgz' | sha256sum -c -; \
	tar -C /usr/local -xzf go.tgz; \
	rm go.tgz; \
	\
	cd /usr/local/go/src; \
	for p in /go-alpine-patches/*.patch; do \
		[ -f "$p" ] || continue; \
		patch -p2 -i "$p"; \
	done; \
	./make.bash; \
	\
	rm -rf /go-alpine-patches; \
	apk del .build-deps; \
	\
	export PATH="/usr/local/go/bin:$PATH"; \
	go version

ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
COPY go/go-wrapper /usr/local/bin/

# And python for node-gyp! plus vega

RUN apk add --no-cache linux-headers build-base \
  python2 python2-dev py-setuptools \
  cairo-dev

# Back to my own stuff...

# angular build deps
RUN npm install -g @angular/cli

RUN apk update && apk add git
RUN mkdir -p /root/go
ENV GOPATH="/root/go"
RUN go get -u github.com/golang/dep/cmd/dep

RUN mkdir -p /root/go/src/main
WORKDIR /root/go/src/main
# why another user??
# RUN addgroup -S databox && adduser -S -g databox databox

RUN mkdir /root/go/src/main/my-app
ADD my-app/package.json my-app/package-lock.json my-app/
RUN cd /root/go/src/main/my-app && npm install

ADD Gopkg.* ./
RUN $GOPATH/bin/dep ensure -vendor-only
ADD src src
RUN GGO_ENABLED=0 GOOS=linux go build -a -tags netgo -installsuffix netgo -ldflags '-d -s -w -extldflags "-static"' -o app src/*.go
ADD . .

RUN cd /root/go/src/main/my-app && ng build -bh /databox-app-activity-summary/ui/static/ && cp dist/* ../www/

FROM scratch
# COPY --from=gobuild /etc/passwd /etc/passwd
# USER databox
WORKDIR /root
COPY --from=gobuild /root/go/src/main/app .
COPY --from=gobuild /root/go/src/main/www/ /root/www/
# COPY --from=gobuild /root/go/src/main/tmpl/ /root/tmpl/
# COPY --from=gobuild /root/go/src/main/etc/ /root/etc/
COPY --from=gobuild /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
LABEL databox.type="app"
EXPOSE 8080

CMD ["./app"]
#CMD ["sleep","2147483647"]
