VERSION = v1.0.0

.PHONY: build

clean:
	@go clean

test:
	@go test -v

build:
	# omit the symbol table, debug information and the DWARF table
	@go build -o go-calc -ldflags="-s -w -X main.version=$(VERSION)"

all: clean build test

docker:
	@docker build -t go-calc:$(VERSION) .

release:
	@goreleaser --snapshot --skip-publish --rm-dist