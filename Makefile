VERSION?=dev
COMMIT=$(shell git rev-parse HEAD | cut -c -8)

LDFLAGS=-ldflags "-X main.Version=${VERSION} -X main.Commit=${COMMIT}"
MODFLAGS=-mod=vendor

PLATFORMS:=darwin linux windows

all: dev

clean:
	rm -fr dist/

dev:
	go build ${MODFLAGS} ${LDFLAGS} -o dist/toolshed ./cmd/toolshed

dist: $(PLATFORMS)

$(PLATFORMS):
	GOOS=$@ GOARCH=amd64 go build ${MODFLAGS} ${LDFLAGS} -o dist/toolshed-$@-amd64 ./cmd/toolshed

test:
	go test ${MODFLAGS} ./...

archive: dist
	bsdtar -zcf /tmp/toolshed.tar.gz -s ,^dist/toolshed-linux-amd64,dist/toolshed, dist/toolshed-linux-amd64 Caddyfile toolshed.service

.PHONY: all clean dev dist $(PLATFORMS) test archive
