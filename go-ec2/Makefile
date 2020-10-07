VERSION = v1.0.0

.PHONY: build

clean:
	@go clean

build:
	@go build -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=`git rev-parse HEAD`"

all: clean build

release:
	@goreleaser --snapshot --skip-publish --rm-dist