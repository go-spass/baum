.PHONY: build test install clean

VERSION := $(shell git describe --tags --always 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.version=$(VERSION)"

build:
	go build $(LDFLAGS) -o baum .

test:
	go test ./...

install:
	go install $(LDFLAGS) .

clean:
	rm -f baum
