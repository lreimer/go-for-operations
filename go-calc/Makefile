VERSION = v1.0.0

.PHONY: build

test:
	@go test

build:
	# omit the symbol table, debug information and the DWARF table
	@go build -o go-calc -ldflags="-s -w -X main.version=$(VERSION)"

install: test build

docker:
	@docker build -t go-calc:$(VERSION) .
