VERSION = v1.0.0

.PHONY: build

clean:
	@go clean

install:
	@go install -ldflags="-s -w -X main.version=$(VERSION) -X main.commit=`git rev-parse HEAD`"

all: clean install

release:
	@goreleaser --snapshot --skip-publish --rm-dist