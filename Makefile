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

.PHONY: all clean dev dist $(PLATFORMS) test
