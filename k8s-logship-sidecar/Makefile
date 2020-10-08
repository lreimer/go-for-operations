VERSION = v1.0.0

.PHONY: build

clean:
	@go clean

test:
	@go test -v

build:
	@go build -o k8s-logship-sidecar -ldflags="-s -w"

all: clean build test

docker:
	@docker build -t k8s-logship-sidecar:$(VERSION) .
	@docker tag k8s-logship-sidecar:$(VERSION) lreimer/k8s-logship-sidecar:$(VERSION)
	@docker push lreimer/k8s-logship-sidecar:$(VERSION)
